package privacy

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/events"
)

const (
	maxStoredErrorLength = 500
	expiredBatchSize     = 50

	exportVisibilityTimeoutSecs = 900
)

// DefaultIDGenerator generates a unique record id.
func DefaultIDGenerator() string {
	return lib.IDsGenerateUnique()
}

// Service implements self-service privacy: hiding a profile, exporting personal data and deleting an account.
type Service struct {
	logger  *logfx.Logger
	config  *Config
	repo    Repository
	storage ObjectStorage // nil when no private bucket is configured
	queue   Enqueuer
	mailer  Mailer // nil when e-mail is not configured
	audit   *events.AuditService
	newID   func() string
	now     func() time.Time
}

// NewService creates a new privacy service. storage and mailer are optional.
func NewService(
	logger *logfx.Logger,
	config *Config,
	repo Repository,
	storage ObjectStorage,
	queue Enqueuer,
	mailer Mailer,
	audit *events.AuditService,
) *Service {
	return &Service{
		logger:  logger,
		config:  config,
		repo:    repo,
		storage: storage,
		queue:   queue,
		mailer:  mailer,
		audit:   audit,
		newID:   DefaultIDGenerator,
		now:     time.Now,
	}
}

// SetClock replaces the time source. It exists for tests, which need to move time forward without sleeping.
func (s *Service) SetClock(now func() time.Time) {
	s.now = now
}

func (s *Service) account(ctx context.Context, userID string) (*Account, error) {
	account, err := s.repo.GetAccount(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToGetAccount, err)
	}

	if account == nil {
		return nil, ErrNoAccount
	}

	if account.Redacted {
		return nil, ErrAccountDeleted
	}

	return account, nil
}

func (s *Service) recordAudit(
	ctx context.Context,
	eventType events.EventType,
	userID string,
	entityType string,
	entityID string,
	payload map[string]any,
) {
	if s.audit == nil {
		return
	}

	s.audit.Record(ctx, events.AuditParams{
		ActorID:    &userID,
		SessionID:  nil,
		Payload:    payload,
		EventType:  eventType,
		EntityType: entityType,
		EntityID:   entityID,
		ActorKind:  events.ActorUser,
	})
}

// GetState returns what the privacy settings page needs to render.
func (s *Service) GetState(ctx context.Context, userID string) (*State, error) {
	account, err := s.account(ctx, userID)
	if err != nil {
		return nil, err
	}

	latest, err := s.repo.GetLatestExport(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	return &State{
		ProfileSlug: account.ProfileSlug,
		IsHidden:    account.HiddenAt != nil,
		Export:      s.exportState(latest),
	}, nil
}

func (s *Service) exportState(latest *DataExport) ExportState {
	state := ExportState{ //nolint:exhaustruct
		Status:     ExportStatusNone,
		Available:  s.storage != nil,
		CanRequest: s.storage != nil,
	}

	if latest == nil {
		return state
	}

	now := s.now()

	state.ID = &latest.ID
	state.Status = latest.Status
	state.RequestedAt = &latest.RequestedAt
	state.ReadyAt = latest.CompletedAt
	state.ExpiresAt = latest.ExpiresAt
	state.SizeBytes = latest.SizeBytes

	// A ready archive past its expiry is as good as expired even before the cleanup worker removes it.
	if latest.Status == ExportStatusReady && latest.ExpiresAt != nil &&
		!latest.ExpiresAt.After(now) {
		state.Status = ExportStatusExpired
	}

	state.CanDownload = state.Status == ExportStatusReady

	if latest.Status == ExportStatusFailed {
		return state
	}

	next := latest.RequestedAt.Add(s.config.ExportCooldown)
	state.NextRequestAt = &next

	if state.Status == ExportStatusPending || state.Status == ExportStatusProcessing ||
		now.Before(next) {
		state.CanRequest = false
	}

	return state
}

// RequestExport queues a personal data export for the user.
func (s *Service) RequestExport(ctx context.Context, userID string) (*DataExport, error) {
	if s.storage == nil {
		return nil, ErrExportUnavailable
	}

	account, err := s.account(ctx, userID)
	if err != nil {
		return nil, err
	}

	latest, err := s.repo.GetLatestExport(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	err = s.checkCanRequest(latest)
	if err != nil {
		return nil, err
	}

	export := &DataExport{ //nolint:exhaustruct
		ID:          s.newID(),
		UserID:      userID,
		ProfileID:   account.ProfileID,
		Status:      ExportStatusPending,
		RequestedAt: s.now(),
	}

	err = s.repo.CreateExport(ctx, export)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	_, err = s.queue.Enqueue(ctx, events.QueueEnqueueParams{ //nolint:exhaustruct
		Type:    events.QueueItemTypePrivacyExport,
		Payload: map[string]any{"export_id": export.ID},
		// Collecting and uploading an archive can take minutes; a short visibility timeout would hand the same
		// request to a second worker while the first is still busy.
		VisibilityTimeoutSecs: exportVisibilityTimeoutSecs,
	})
	if err != nil {
		failErr := s.repo.MarkExportFailed(ctx, export.ID, "could not be scheduled")
		if failErr != nil {
			s.logger.ErrorContext(ctx, "failed to mark unscheduled export as failed",
				slog.String("export_id", export.ID),
				slog.Any("error", failErr))
		}

		return nil, fmt.Errorf("%w: %w", ErrFailedToEnqueue, err)
	}

	s.recordAudit(ctx, events.PrivacyExportRequested, userID, "user", userID,
		map[string]any{"export_id": export.ID})

	return export, nil
}

func (s *Service) checkCanRequest(latest *DataExport) error {
	if latest == nil || latest.Status == ExportStatusFailed {
		return nil
	}

	if latest.Status == ExportStatusPending || latest.Status == ExportStatusProcessing {
		return ErrExportInProgress
	}

	next := latest.RequestedAt.Add(s.config.ExportCooldown)
	if s.now().Before(next) {
		return &CooldownError{NextRequestAt: next}
	}

	return nil
}

// GetExportDownloadURL returns a short-lived signed URL for a finished archive owned by the user.
func (s *Service) GetExportDownloadURL(
	ctx context.Context,
	userID string,
	exportID string,
) (string, error) {
	if s.storage == nil {
		return "", ErrExportUnavailable
	}

	export, err := s.repo.GetExport(ctx, exportID)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	err = s.checkDownloadable(export, userID)
	if err != nil {
		return "", err
	}

	filename := "aya-personal-data-" + export.RequestedAt.UTC().Format("2006-01-02") + ".zip"

	url, err := s.storage.PresignedDownloadURL(
		ctx,
		*export.StorageKey,
		filename,
		s.config.DownloadURLTTL,
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	return url, nil
}

// checkDownloadable decides whether the user may download the export right now.
func (s *Service) checkDownloadable(export *DataExport, userID string) error {
	// Another user's export is reported as missing, never as forbidden, so ids cannot be probed.
	if export == nil || export.UserID != userID {
		return ErrExportNotFound
	}

	if export.Status == ExportStatusExpired ||
		(export.ExpiresAt != nil && !export.ExpiresAt.After(s.now())) {
		return ErrExportExpired
	}

	if export.Status != ExportStatusReady || export.StorageKey == nil {
		return ErrExportNotReady
	}

	return nil
}

// HideAccount hides the user's profile and signs them out. The profile comes back at their next login.
func (s *Service) HideAccount(ctx context.Context, userID string) error {
	account, err := s.account(ctx, userID)
	if err != nil {
		return err
	}

	changed, err := s.repo.HideProfile(ctx, userID, account.ProfileID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToHide, err)
	}

	if changed {
		s.recordAudit(ctx, events.PrivacyProfileHidden, userID, "profile", account.ProfileID, nil)
	}

	return nil
}

// UnhideOnLogin restores a profile its owner hid, if that is the case. Profiles an administrator removed
// are left alone.
func (s *Service) UnhideOnLogin(ctx context.Context, userID string) error {
	changed, err := s.repo.UnhideProfile(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUnhide, err)
	}

	if changed {
		s.recordAudit(ctx, events.PrivacyProfileUnhidden, userID, "user", userID, nil)
	}

	return nil
}

// DeleteAccount irreversibly erases the user's personal data and frees their username.
func (s *Service) DeleteAccount(
	ctx context.Context,
	userID string,
	confirmSlug string,
	opts DeleteOptions,
) error {
	account, err := s.account(ctx, userID)
	if err != nil {
		return err
	}

	if !strings.EqualFold(strings.TrimSpace(confirmSlug), account.ProfileSlug) {
		return ErrConfirmationMismatch
	}

	owned, err := s.repo.ListSoleOwnedProfiles(ctx, account.ProfileID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToDelete, err)
	}

	if len(owned) > 0 {
		return &OwnedProfilesError{Profiles: owned}
	}

	s.removeExportObjects(ctx, userID)

	err = s.repo.RedactAccount(ctx, userID, account.ProfileID, opts)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToDelete, err)
	}

	s.recordAudit(ctx, events.PrivacyAccountDeleted, userID, "user", userID,
		map[string]any{"delete_content": opts.DeleteContent})

	return nil
}

// removeExportObjects deletes the archives of a departing user. Failures are logged: an orphaned object
// still expires through the bucket lifecycle rule, and must not stop the deletion.
func (s *Service) removeExportObjects(ctx context.Context, userID string) {
	if s.storage == nil {
		return
	}

	exports, err := s.repo.ListExportsByUser(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to list exports before account deletion",
			slog.String("user_id", userID), slog.Any("error", err))

		return
	}

	for _, export := range exports {
		if export.StorageKey == nil {
			continue
		}

		removeErr := s.storage.RemoveObject(ctx, *export.StorageKey)
		if removeErr != nil {
			s.logger.ErrorContext(ctx, "failed to remove export archive",
				slog.String("export_id", export.ID), slog.Any("error", removeErr))
		}
	}
}

// CleanupExpiredExports removes archives past their retention period and returns how many it handled.
func (s *Service) CleanupExpiredExports(ctx context.Context) (int, error) {
	if s.storage == nil {
		return 0, nil
	}

	expired, err := s.repo.ListExpiredExports(ctx, expiredBatchSize)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrFailedToExport, err)
	}

	handled := 0

	for _, export := range expired {
		if export.StorageKey != nil {
			removeErr := s.storage.RemoveObject(ctx, *export.StorageKey)
			if removeErr != nil {
				s.logger.ErrorContext(ctx, "failed to remove expired export archive",
					slog.String("export_id", export.ID), slog.Any("error", removeErr))

				continue
			}
		}

		markErr := s.repo.MarkExportExpired(ctx, export.ID)
		if markErr != nil {
			s.logger.ErrorContext(ctx, "failed to mark export as expired",
				slog.String("export_id", export.ID), slog.Any("error", markErr))

			continue
		}

		handled++
	}

	return handled, nil
}

// truncateError keeps stored failure messages short and free of anything that looks like a stack dump.
func truncateError(err error) string {
	message := err.Error()
	if len(message) > maxStoredErrorLength {
		message = message[:maxStoredErrorLength]
	}

	return message
}

// IsFinalAttempt reports whether a failing queue item will not be retried again.
func IsFinalAttempt(item *events.QueueItem) bool {
	return item.RetryCount >= item.MaxRetries
}
