//nolint:exhaustruct,varnamelen,cyclop,funlen // table-free scenario tests with small in-memory fakes
package privacy_test

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/events"
	"github.com/eser/aya.is/services/pkg/api/business/privacy"
)

var errBoom = errors.New("boom")

// fakeRepo is an in-memory privacy.Repository.
type fakeRepo struct {
	account    *privacy.Account
	contact    *privacy.Contact
	owned      []privacy.OwnedProfile
	exports    map[string]*privacy.DataExport
	order      []string
	datasets   map[string][]string
	hidden     bool
	redacted   bool
	redactOpts privacy.DeleteOptions
	failCreate bool
	claimed    map[string]bool
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		account: &privacy.Account{UserID: "user1", ProfileID: "profile1", ProfileSlug: "jane"},
		contact: &privacy.Contact{Email: "jane@example.com", Locale: "en", Slug: "jane"},
		exports: map[string]*privacy.DataExport{},
		claimed: map[string]bool{},
		datasets: map[string][]string{
			"account": {`{"id":"user1"}`},
			"stories": {`{"id":"s1"}`, `{"id":"s2"}`},
			"empty":   {},
		},
	}
}

func (f *fakeRepo) GetAccount(_ context.Context, _ string) (*privacy.Account, error) {
	return f.account, nil
}

func (f *fakeRepo) GetContact(_ context.Context, _ string) (*privacy.Contact, error) {
	return f.contact, nil
}

func (f *fakeRepo) HideProfile(_ context.Context, _ string, _ string) (bool, error) {
	changed := !f.hidden
	f.hidden = true

	return changed, nil
}

func (f *fakeRepo) UnhideProfile(_ context.Context, _ string) (bool, error) {
	changed := f.hidden
	f.hidden = false

	return changed, nil
}

func (f *fakeRepo) ListSoleOwnedProfiles(
	_ context.Context,
	_ string,
) ([]privacy.OwnedProfile, error) {
	return f.owned, nil
}

func (f *fakeRepo) CreateExport(_ context.Context, export *privacy.DataExport) error {
	if f.failCreate {
		return errBoom
	}

	copied := *export
	f.exports[export.ID] = &copied
	f.order = append(f.order, export.ID)

	return nil
}

func (f *fakeRepo) GetLatestExport(_ context.Context, _ string) (*privacy.DataExport, error) {
	if len(f.order) == 0 {
		return nil, nil //nolint:nilnil
	}

	return f.exports[f.order[len(f.order)-1]], nil
}

func (f *fakeRepo) GetExport(_ context.Context, id string) (*privacy.DataExport, error) {
	return f.exports[id], nil
}

func (f *fakeRepo) ClaimExport(_ context.Context, id string) (bool, error) {
	export := f.exports[id]
	if export == nil ||
		(export.Status != privacy.ExportStatusPending && export.Status != privacy.ExportStatusProcessing) {
		return false, nil
	}

	export.Status = privacy.ExportStatusProcessing
	f.claimed[id] = true

	return true, nil
}

func (f *fakeRepo) MarkExportReady(
	_ context.Context,
	id string,
	storageKey string,
	sizeBytes int64,
	expiresAt time.Time,
) error {
	export := f.exports[id]
	export.Status = privacy.ExportStatusReady
	export.StorageKey = &storageKey
	export.SizeBytes = &sizeBytes
	export.ExpiresAt = &expiresAt

	return nil
}

func (f *fakeRepo) MarkExportFailed(_ context.Context, id string, message string) error {
	f.exports[id].Status = privacy.ExportStatusFailed
	f.exports[id].ErrorMessage = &message

	return nil
}

func (f *fakeRepo) MarkExportNotified(_ context.Context, id string) error {
	now := time.Now()
	f.exports[id].NotifiedAt = &now

	return nil
}

func (f *fakeRepo) ListExpiredExports(_ context.Context, _ int) ([]*privacy.DataExport, error) {
	return nil, nil
}

func (f *fakeRepo) MarkExportExpired(_ context.Context, _ string) error { return nil }

func (f *fakeRepo) ListExportsByUser(_ context.Context, _ string) ([]*privacy.DataExport, error) {
	list := make([]*privacy.DataExport, 0, len(f.exports))
	for _, export := range f.exports {
		list = append(list, export)
	}

	return list, nil
}

func (f *fakeRepo) ExportDatasetNames(_ context.Context) ([]string, error) {
	return []string{"account", "stories", "empty"}, nil
}

func (f *fakeRepo) StreamDataset(
	_ context.Context,
	name string,
	_ string,
	_ string,
	fn func(row json.RawMessage) error,
) error {
	for _, row := range f.datasets[name] {
		err := fn(json.RawMessage(row))
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *fakeRepo) RedactAccount(
	_ context.Context,
	_ string,
	_ string,
	opts privacy.DeleteOptions,
) error {
	f.redacted = true
	f.redactOpts = opts

	return nil
}

type fakeStorage struct {
	objects map[string][]byte
	removed []string
}

func (s *fakeStorage) PutObject(_ context.Context, key string, _ string, body io.ReadSeeker) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return err //nolint:wrapcheck
	}

	if s.objects == nil {
		s.objects = map[string][]byte{}
	}

	s.objects[key] = data

	return nil
}

func (s *fakeStorage) PresignedDownloadURL(
	_ context.Context,
	key string,
	filename string,
	_ time.Duration,
) (string, error) {
	return "https://signed.example/" + key + "?name=" + filename, nil
}

func (s *fakeStorage) RemoveObject(_ context.Context, key string) error {
	s.removed = append(s.removed, key)

	return nil
}

type fakeQueue struct{ items []events.QueueEnqueueParams }

func (q *fakeQueue) Enqueue(_ context.Context, params events.QueueEnqueueParams) (string, error) {
	q.items = append(q.items, params)

	return "item", nil
}

type fakeMailer struct{ sent []string }

func (m *fakeMailer) SendExportReady(
	_ context.Context,
	to string,
	_ string,
	url string,
	_ int,
) error {
	m.sent = append(m.sent, to+" "+url)

	return nil
}

func newService(
	repo *fakeRepo,
	storage privacy.ObjectStorage,
	queue *fakeQueue,
	mailer privacy.Mailer,
	now *time.Time,
) *privacy.Service {
	config := &privacy.Config{
		ExportCooldown:  30 * 24 * time.Hour,
		ExportRetention: 7 * 24 * time.Hour,
		DownloadURLTTL:  10 * time.Minute,
		ExportMaxBytes:  1 << 20,
		FrontendURI:     "https://aya.is",
	}

	service := privacy.NewService(logfx.NewLogger(), config, repo, storage, queue, mailer, nil)
	service.SetClock(func() time.Time { return *now })

	return service
}

func TestRequestExportQueuesWorkAndEnforcesCooldown(t *testing.T) {
	t.Parallel()

	repo, queue, now := newFakeRepo(), &fakeQueue{}, time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC)
	service := newService(repo, &fakeStorage{}, queue, nil, &now)

	export, err := service.RequestExport(context.Background(), "user1")
	if err != nil {
		t.Fatalf("first request: %v", err)
	}

	if len(queue.items) != 1 || queue.items[0].Type != events.QueueItemTypePrivacyExport ||
		queue.items[0].Payload["export_id"] != export.ID {
		t.Fatalf("export was not queued correctly: %+v", queue.items)
	}

	// While it is pending, another request is refused.
	_, err = service.RequestExport(context.Background(), "user1")
	if !errors.Is(err, privacy.ErrExportInProgress) {
		t.Fatalf("want ErrExportInProgress, got %v", err)
	}

	// Once finished, the 30-day cooldown applies and reports when the next request is allowed.
	repo.exports[export.ID].Status = privacy.ExportStatusReady
	now = now.Add(5 * 24 * time.Hour)

	_, err = service.RequestExport(context.Background(), "user1")

	var cooldown *privacy.CooldownError
	if !errors.As(err, &cooldown) {
		t.Fatalf("want CooldownError, got %v", err)
	}

	if want := export.RequestedAt.Add(30 * 24 * time.Hour); !cooldown.NextRequestAt.Equal(want) {
		t.Fatalf("next request at = %v, want %v", cooldown.NextRequestAt, want)
	}

	// After the cooldown a new request is accepted.
	now = now.Add(26 * 24 * time.Hour)

	_, err = service.RequestExport(context.Background(), "user1")
	if err != nil {
		t.Fatalf("request after cooldown: %v", err)
	}
}

func TestFailedExportDoesNotBlockARetry(t *testing.T) {
	t.Parallel()

	repo, queue, now := newFakeRepo(), &fakeQueue{}, time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC)
	service := newService(repo, &fakeStorage{}, queue, nil, &now)

	export, err := service.RequestExport(context.Background(), "user1")
	if err != nil {
		t.Fatal(err)
	}

	repo.exports[export.ID].Status = privacy.ExportStatusFailed

	_, err = service.RequestExport(context.Background(), "user1")
	if err != nil {
		t.Fatalf("a failed export must not count against the cooldown: %v", err)
	}
}

func TestRequestExportUnavailableWithoutStorage(t *testing.T) {
	t.Parallel()

	now := time.Now()
	service := newService(newFakeRepo(), nil, &fakeQueue{}, nil, &now)

	_, err := service.RequestExport(context.Background(), "user1")
	if !errors.Is(err, privacy.ErrExportUnavailable) {
		t.Fatalf("want ErrExportUnavailable, got %v", err)
	}

	state, err := service.GetState(context.Background(), "user1")
	if err != nil || state.Export.Available || state.Export.CanRequest {
		t.Fatalf("state must report the feature as unavailable: %+v err=%v", state, err)
	}
}

func TestProcessExportBuildsArchiveAndNotifies(t *testing.T) {
	t.Parallel()

	repo, storage, mailer, now := newFakeRepo(), &fakeStorage{}, &fakeMailer{}, time.Date(
		2026,
		1,
		10,
		12,
		0,
		0,
		0,
		time.UTC,
	)
	service := newService(repo, storage, &fakeQueue{}, mailer, &now)

	export, err := service.RequestExport(context.Background(), "user1")
	if err != nil {
		t.Fatal(err)
	}

	err = service.ProcessExport(context.Background(), export.ID, false)
	if err != nil {
		t.Fatalf("ProcessExport: %v", err)
	}

	done := repo.exports[export.ID]
	if done.Status != privacy.ExportStatusReady || done.StorageKey == nil || done.ExpiresAt == nil {
		t.Fatalf("export not marked ready: %+v", done)
	}

	if want := now.Add(7 * 24 * time.Hour); !done.ExpiresAt.Equal(want) {
		t.Fatalf("expires at %v, want %v", done.ExpiresAt, want)
	}

	archive := storage.objects[*done.StorageKey]
	if len(archive) == 0 {
		t.Fatal("no archive was stored")
	}

	files := readZip(t, archive)

	for _, name := range []string{"account.json", "stories.json", "empty.json", "manifest.json"} {
		if _, ok := files[name]; !ok {
			t.Fatalf("archive misses %s (has %v)", name, keys(files))
		}
	}

	var stories []map[string]string

	err = json.Unmarshal(files["stories.json"], &stories)
	if err != nil || len(stories) != 2 {
		t.Fatalf("stories.json: %v %s", err, files["stories.json"])
	}

	var empty []any

	err = json.Unmarshal(files["empty.json"], &empty)
	if err != nil || len(empty) != 0 {
		t.Fatalf("an empty dataset must be an empty JSON array: %v %s", err, files["empty.json"])
	}

	var doc struct {
		Datasets []struct {
			Name string `json:"name"`
			Rows int    `json:"rows"`
		} `json:"datasets"`
		UserID string `json:"user_id"`
	}

	err = json.Unmarshal(files["manifest.json"], &doc)
	if err != nil || doc.UserID != "user1" || len(doc.Datasets) != 3 || doc.Datasets[1].Rows != 2 {
		t.Fatalf("manifest: %v %+v", err, doc)
	}

	if len(mailer.sent) != 1 ||
		mailer.sent[0] != "jane@example.com https://aya.is/en/jane/settings/privacy" {
		t.Fatalf("notification: %v", mailer.sent)
	}

	if done.NotifiedAt == nil {
		t.Fatal("the notification must be recorded")
	}

	// Processing the same request again must not build a second archive.
	storage.objects = map[string][]byte{}

	err = service.ProcessExport(context.Background(), export.ID, false)
	if err != nil || len(storage.objects) != 0 {
		t.Fatalf(
			"a finished export must not be rebuilt: err=%v objects=%d",
			err,
			len(storage.objects),
		)
	}
}

func TestProcessExportDoesNotEmailRedactedAddresses(t *testing.T) {
	t.Parallel()

	repo, storage, mailer, now := newFakeRepo(), &fakeStorage{}, &fakeMailer{}, time.Now()
	repo.contact.Email = privacy.RedactedText
	service := newService(repo, storage, &fakeQueue{}, mailer, &now)

	export, _ := service.RequestExport(context.Background(), "user1")

	err := service.ProcessExport(context.Background(), export.ID, false)
	if err != nil || len(mailer.sent) != 0 {
		t.Fatalf("err=%v sent=%v", err, mailer.sent)
	}
}

func TestProcessExportFailsPermanentlyForDeletedAccounts(t *testing.T) {
	t.Parallel()

	repo, storage, now := newFakeRepo(), &fakeStorage{}, time.Now()
	service := newService(repo, storage, &fakeQueue{}, nil, &now)

	export, _ := service.RequestExport(context.Background(), "user1")
	repo.account.Redacted = true

	err := service.ProcessExport(context.Background(), export.ID, false)
	if err != nil {
		t.Fatalf("a permanent failure must not be retried: %v", err)
	}

	if repo.exports[export.ID].Status != privacy.ExportStatusFailed || len(storage.objects) != 0 {
		t.Fatalf("export should be failed with nothing stored: %+v", repo.exports[export.ID])
	}
}

func TestProcessExportTooLargeIsPermanent(t *testing.T) {
	t.Parallel()

	repo, storage, now := newFakeRepo(), &fakeStorage{}, time.Now()
	// The limit applies to the archive as stored, so the payload must not compress away: random bytes do not.
	noise := make([]byte, 2<<20)
	_, _ = rand.Read(noise)
	repo.datasets["stories"] = []string{`{"blob":"` + hex.EncodeToString(noise) + `"}`}
	service := newService(repo, storage, &fakeQueue{}, nil, &now)

	export, _ := service.RequestExport(context.Background(), "user1")

	err := service.ProcessExport(context.Background(), export.ID, false)
	if err != nil {
		t.Fatalf("too large is not retryable: %v", err)
	}

	if repo.exports[export.ID].Status != privacy.ExportStatusFailed {
		t.Fatalf("status = %s", repo.exports[export.ID].Status)
	}
}

func TestDownloadURLOnlyForOwnerAndReadyAndUnexpired(t *testing.T) {
	t.Parallel()

	repo, storage, now := newFakeRepo(), &fakeStorage{}, time.Date(
		2026,
		1,
		10,
		12,
		0,
		0,
		0,
		time.UTC,
	)
	service := newService(repo, storage, &fakeQueue{}, nil, &now)

	export, _ := service.RequestExport(context.Background(), "user1")

	_, err := service.GetExportDownloadURL(context.Background(), "user1", export.ID)
	if !errors.Is(err, privacy.ErrExportNotReady) {
		t.Fatalf("pending export: %v", err)
	}

	err = service.ProcessExport(context.Background(), export.ID, false)
	if err != nil {
		t.Fatal(err)
	}

	url, err := service.GetExportDownloadURL(context.Background(), "user1", export.ID)
	if err != nil || url == "" {
		t.Fatalf("ready export: url=%q err=%v", url, err)
	}

	// Someone else's export looks exactly like a missing one.
	_, err = service.GetExportDownloadURL(context.Background(), "someone-else", export.ID)
	if !errors.Is(err, privacy.ErrExportNotFound) {
		t.Fatalf("other user's export: %v", err)
	}

	now = now.Add(8 * 24 * time.Hour)

	_, err = service.GetExportDownloadURL(context.Background(), "user1", export.ID)
	if !errors.Is(err, privacy.ErrExportExpired) {
		t.Fatalf("expired export: %v", err)
	}

	state, _ := service.GetState(context.Background(), "user1")
	if state.Export.Status != privacy.ExportStatusExpired || state.Export.CanDownload {
		t.Fatalf(
			"an unexpired-looking but past-expiry export must read as expired: %+v",
			state.Export,
		)
	}
}

func TestProcessExportRetriesTransientFailuresUntilTheLastAttempt(t *testing.T) {
	t.Parallel()

	repo, now := newFakeRepo(), time.Now()
	failing := &failingStorage{}
	service := newService(repo, failing, &fakeQueue{}, nil, &now)

	export, _ := service.RequestExport(context.Background(), "user1")

	err := service.ProcessExport(context.Background(), export.ID, false)
	if err == nil || repo.exports[export.ID].Status != privacy.ExportStatusProcessing {
		t.Fatalf(
			"a transient failure must be returned for retry and keep the export processing: err=%v status=%s",
			err,
			repo.exports[export.ID].Status,
		)
	}

	err = service.ProcessExport(context.Background(), export.ID, true)
	if err == nil || repo.exports[export.ID].Status != privacy.ExportStatusFailed {
		t.Fatalf(
			"the final attempt must mark the export failed: err=%v status=%s",
			err,
			repo.exports[export.ID].Status,
		)
	}
}

type failingStorage struct{ fakeStorage }

func (s *failingStorage) PutObject(_ context.Context, _ string, _ string, _ io.ReadSeeker) error {
	return errBoom
}

func TestDeleteAccountRequiresMatchingUsernameAndNoSoleOwnedProfiles(t *testing.T) {
	t.Parallel()

	repo, storage, now := newFakeRepo(), &fakeStorage{}, time.Now()
	service := newService(repo, storage, &fakeQueue{}, nil, &now)

	err := service.DeleteAccount(context.Background(), "user1", "not-jane", privacy.DeleteOptions{})
	if !errors.Is(err, privacy.ErrConfirmationMismatch) || repo.redacted {
		t.Fatalf("wrong confirmation: err=%v redacted=%v", err, repo.redacted)
	}

	repo.owned = []privacy.OwnedProfile{
		{ID: "org1", Slug: "acme", Kind: "organization", Title: "Acme"},
	}

	err = service.DeleteAccount(context.Background(), "user1", "jane", privacy.DeleteOptions{})

	var owned *privacy.OwnedProfilesError
	if !errors.As(err, &owned) || len(owned.Profiles) != 1 || repo.redacted {
		t.Fatalf("sole owner: err=%v redacted=%v", err, repo.redacted)
	}

	repo.owned = nil

	err = service.DeleteAccount(
		context.Background(),
		"user1",
		"  JANE ",
		privacy.DeleteOptions{DeleteContent: true},
	)
	if err != nil || !repo.redacted || !repo.redactOpts.DeleteContent {
		t.Fatalf("delete: err=%v redacted=%v opts=%+v", err, repo.redacted, repo.redactOpts)
	}
}

func TestDeleteAccountRemovesStoredArchives(t *testing.T) {
	t.Parallel()

	repo, storage, now := newFakeRepo(), &fakeStorage{}, time.Now()
	service := newService(repo, storage, &fakeQueue{}, nil, &now)

	export, _ := service.RequestExport(context.Background(), "user1")
	_ = service.ProcessExport(context.Background(), export.ID, false)

	key := *repo.exports[export.ID].StorageKey

	err := service.DeleteAccount(context.Background(), "user1", "jane", privacy.DeleteOptions{})
	if err != nil || len(storage.removed) != 1 || storage.removed[0] != key {
		t.Fatalf(
			"archives must be removed with the account: err=%v removed=%v",
			err,
			storage.removed,
		)
	}
}

func TestHideAndUnhide(t *testing.T) {
	t.Parallel()

	repo, now := newFakeRepo(), time.Now()
	service := newService(repo, &fakeStorage{}, &fakeQueue{}, nil, &now)

	err := service.HideAccount(context.Background(), "user1")
	if err != nil || !repo.hidden {
		t.Fatalf("hide: err=%v hidden=%v", err, repo.hidden)
	}

	err = service.UnhideOnLogin(context.Background(), "user1")
	if err != nil || repo.hidden {
		t.Fatalf("unhide: err=%v hidden=%v", err, repo.hidden)
	}
}

func TestDeletedAccountCannotUsePrivacyActions(t *testing.T) {
	t.Parallel()

	repo, now := newFakeRepo(), time.Now()
	repo.account.Redacted = true
	service := newService(repo, &fakeStorage{}, &fakeQueue{}, nil, &now)

	err := service.HideAccount(context.Background(), "user1")
	if !errors.Is(err, privacy.ErrAccountDeleted) {
		t.Fatalf("hide: %v", err)
	}

	_, err = service.RequestExport(context.Background(), "user1")
	if !errors.Is(err, privacy.ErrAccountDeleted) {
		t.Fatalf("export: %v", err)
	}
}

func readZip(t *testing.T, data []byte) map[string][]byte {
	t.Helper()

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("not a zip: %v", err)
	}

	files := map[string][]byte{}

	for _, file := range reader.File {
		handle, openErr := file.Open()
		if openErr != nil {
			t.Fatal(openErr)
		}

		content, readErr := io.ReadAll(handle)
		_ = handle.Close()

		if readErr != nil {
			t.Fatal(readErr)
		}

		files[file.Name] = content
	}

	return files
}

func keys(files map[string][]byte) []string {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}

	return names
}
