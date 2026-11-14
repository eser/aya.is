package privacy

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/events"
)

const hoursPerDay = 24

// archiveNotes explain to the reader what the archive is and, just as important, what it leaves out.
var archiveNotes = []string{ //nolint:gochecknoglobals
	"This archive contains the personal data the platform stores about your account.",
	"Every file is a JSON array; every element is a database row exactly as stored.",
	"Authentication tokens, OAuth secrets and hashed network identifiers are never included.",
	"Content other people wrote to you (for example private messages) is limited to what was addressed to you.",
}

type manifestDataset struct {
	Name string `json:"name"`
	File string `json:"file"`
	Rows int    `json:"rows"`
}

type manifest struct {
	GeneratedAt time.Time         `json:"generated_at"`
	UserID      string            `json:"user_id"`
	ProfileID   string            `json:"profile_id"`
	ProfileSlug string            `json:"profile_slug"`
	Notes       []string          `json:"notes"`
	Datasets    []manifestDataset `json:"datasets"`
}

// limitedWriter fails once more than limit bytes have been written, so a runaway export cannot exhaust memory.
type limitedWriter struct {
	writer    io.Writer
	remaining int64
}

func (l *limitedWriter) Write(data []byte) (int, error) {
	l.remaining -= int64(len(data))
	if l.remaining < 0 {
		return 0, ErrExportTooLarge
	}

	return l.writer.Write(data) //nolint:wrapcheck
}

// ProcessExport builds and stores the archive of one export request. The queue handler calls it. When the work
// fails for good it marks the request as failed; a transient failure is returned so the queue retries it.
func (s *Service) ProcessExport(ctx context.Context, exportID string, finalAttempt bool) error {
	export, err := s.repo.GetExport(ctx, exportID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	if export == nil {
		s.logger.WarnContext(
			ctx,
			"export request no longer exists",
			slog.String("export_id", exportID),
		)

		return nil
	}

	claimed, err := s.repo.ClaimExport(ctx, exportID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	if !claimed {
		return nil
	}

	runErr := s.runExport(ctx, export)
	if runErr == nil {
		return nil
	}

	permanent := isPermanentExportError(runErr)

	if permanent || finalAttempt {
		markErr := s.repo.MarkExportFailed(ctx, exportID, truncateError(runErr))
		if markErr != nil {
			s.logger.ErrorContext(ctx, "failed to mark export as failed",
				slog.String("export_id", exportID), slog.Any("error", markErr))
		}
	}

	if permanent {
		s.logger.WarnContext(ctx, "export failed permanently",
			slog.String("export_id", exportID), slog.Any("error", runErr))

		return nil
	}

	return runErr
}

// isPermanentExportError tells failures that retrying cannot fix from transient ones.
func isPermanentExportError(err error) bool {
	return errors.Is(err, ErrAccountDeleted) ||
		errors.Is(err, ErrNoAccount) ||
		errors.Is(err, ErrExportTooLarge) ||
		errors.Is(err, ErrExportUnavailable)
}

func (s *Service) runExport(ctx context.Context, export *DataExport) error {
	if s.storage == nil {
		return ErrExportUnavailable
	}

	account, err := s.repo.GetAccount(ctx, export.UserID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToGetAccount, err)
	}

	if account == nil {
		return ErrNoAccount
	}

	if account.Redacted {
		return ErrAccountDeleted
	}

	var archive bytes.Buffer

	err = s.writeArchive(
		ctx,
		account,
		&limitedWriter{writer: &archive, remaining: s.config.ExportMaxBytes},
	)
	if err != nil {
		return err
	}

	key := "exports/" + export.UserID + "/" + export.ID + ".zip"

	err = s.storage.PutObject(ctx, key, exportContentType, bytes.NewReader(archive.Bytes()))
	if err != nil {
		return fmt.Errorf("%w: storing the archive: %w", ErrFailedToExport, err)
	}

	err = s.repo.MarkExportReady(
		ctx,
		export.ID,
		key,
		int64(archive.Len()),
		s.now().Add(s.config.ExportRetention),
	)
	if err != nil {
		// The archive exists but is unreachable; remove it so nothing personal is left orphaned.
		_ = s.storage.RemoveObject(ctx, key)

		return fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	s.recordAudit(ctx, events.PrivacyExportCompleted, export.UserID, "user", export.UserID,
		map[string]any{"export_id": export.ID, "size_bytes": archive.Len()})

	s.notifyExportReady(ctx, export)

	return nil
}

// writeArchive streams every dataset into a zip, followed by a manifest that lists what was written.
func (s *Service) writeArchive(ctx context.Context, account *Account, out io.Writer) error {
	zipWriter := zip.NewWriter(out)

	datasets := make([]manifestDataset, 0)

	names, err := s.repo.ExportDatasetNames(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	for _, name := range names {
		file := name + ".json"

		rows, err := s.writeDataset(ctx, zipWriter, file, name, account)
		if err != nil {
			return err
		}

		datasets = append(datasets, manifestDataset{Name: name, File: file, Rows: rows})
	}

	doc := manifest{
		GeneratedAt: s.now().UTC(),
		UserID:      account.UserID,
		ProfileID:   account.ProfileID,
		ProfileSlug: account.ProfileSlug,
		Notes:       archiveNotes,
		Datasets:    datasets,
	}

	manifestEntry, err := zipWriter.Create("manifest.json")
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	encoder := json.NewEncoder(manifestEntry)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(doc)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	err = zipWriter.Close()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	return nil
}

func (s *Service) writeDataset(
	ctx context.Context,
	zipWriter *zip.Writer,
	file string,
	name string,
	account *Account,
) (int, error) {
	entry, err := zipWriter.Create(file)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	_, err = io.WriteString(entry, "[")
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	rows := 0

	err = s.repo.StreamDataset(ctx, name, account.UserID, account.ProfileID,
		func(row json.RawMessage) error {
			var pretty bytes.Buffer

			indentErr := json.Indent(&pretty, row, "  ", "  ")
			if indentErr != nil {
				return fmt.Errorf("dataset %s: %w", name, indentErr)
			}

			separator := "\n  "
			if rows > 0 {
				separator = ",\n  "
			}

			_, writeErr := io.WriteString(entry, separator+pretty.String())

			rows++

			return writeErr //nolint:wrapcheck
		})
	if err != nil {
		if errors.Is(err, ErrExportTooLarge) {
			return 0, ErrExportTooLarge
		}

		return 0, fmt.Errorf("%w: dataset %s: %w", ErrFailedToExport, name, err)
	}

	closing := "]\n"
	if rows > 0 {
		closing = "\n]\n"
	}

	_, err = io.WriteString(entry, closing)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	return rows, nil
}

// notifyExportReady e-mails the user. It never fails the export: the archive is downloadable from the
// settings page either way.
func (s *Service) notifyExportReady(ctx context.Context, export *DataExport) {
	if s.mailer == nil || s.config.FrontendURI == "" {
		return
	}

	contact, err := s.repo.GetContact(ctx, export.UserID)
	if err != nil || contact == nil || contact.Email == "" || contact.Email == RedactedText {
		return
	}

	settingsURL := strings.TrimRight(s.config.FrontendURI, "/") + "/" + contact.Locale + "/" +
		contact.Slug + "/settings/privacy"

	err = s.mailer.SendExportReady(
		ctx,
		contact.Email,
		contact.Locale,
		settingsURL,
		int(s.config.ExportRetention.Hours()/hoursPerDay),
	)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to send the export-ready e-mail",
			slog.String("export_id", export.ID), slog.Any("error", err))

		return
	}

	markErr := s.repo.MarkExportNotified(ctx, export.ID)
	if markErr != nil {
		s.logger.ErrorContext(ctx, "failed to record the export notification",
			slog.String("export_id", export.ID), slog.Any("error", markErr))
	}
}
