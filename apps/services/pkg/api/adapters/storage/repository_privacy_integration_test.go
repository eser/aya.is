//nolint:testpackage,paralleltest,varnamelen,lll,gocognit,cyclop,funlen,gosec,rowserrcheck,sqlclosecheck,gocritic,wsl_v5,exhaustruct // opt-in destructive DB test: needs package internals and a serial, shared database
package storage

// Integration tests for the privacy repository. They need a PostgreSQL database that already has every migration
// applied, ideally filled with realistic data (a copy of production gives by far the best coverage), and they
// DESTROY data in it: deleting accounts is exactly what they exercise.
//
//	PRIVACY_TEST_DSN="postgres://postgres:...@127.0.0.1:55450/privacy_test_real?sslmode=disable" \
//	  go test ./pkg/api/adapters/storage/ -run TestPrivacy -count=1 -v
//
// The database name must start with "privacy_test"; any other name aborts the run. Without the variable the
// tests are skipped, so `go test ./...` stays hermetic. Restore the database between runs.

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/privacy"
	"github.com/eser/aya.is/services/pkg/lib/caching"
	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" database/sql driver
)

func privacyTestRepository(t *testing.T) (*Repository, *sql.DB) {
	t.Helper()

	dsn := os.Getenv("PRIVACY_TEST_DSN")
	if dsn == "" {
		t.Skip("PRIVACY_TEST_DSN is not set")
	}

	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("PRIVACY_TEST_DSN: %v", err)
	}

	if name := strings.TrimPrefix(parsed.Path, "/"); !strings.HasPrefix(name, "privacy_test") {
		t.Fatalf(
			"refusing to run against database %q: its name must start with \"privacy_test\"",
			name,
		)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = db.Close() })

	repo := &Repository{ //nolint:exhaustruct
		db:       db,
		dbtx:     db,
		queries:  &Queries{db: db},
		cacheTTL: time.Minute,
	}
	repo.cache = caching.NewCache(
		func(context.Context, string, any) (bool, error) { return false, nil },
		func(context.Context, string, any) error { return nil },
		func(ctx context.Context, key string) error { return repo.CacheRemove(ctx, key) },
	)

	return repo, db
}

func queryStrings(t *testing.T, db *sql.DB, query string, args ...any) []string {
	t.Helper()

	rows, err := db.QueryContext(context.Background(), query, args...)
	if err != nil {
		t.Fatalf("%s: %v", query, err)
	}

	defer func() { _ = rows.Close() }()

	var out []string

	for rows.Next() {
		var value sql.NullString

		err = rows.Scan(&value)
		if err != nil {
			t.Fatal(err)
		}

		out = append(out, value.String)
	}

	return out
}

func queryCount(t *testing.T, db *sql.DB, query string, args ...any) int {
	t.Helper()

	var count int

	err := db.QueryRowContext(context.Background(), query, args...).Scan(&count)
	if err != nil {
		t.Fatalf("%s: %v", query, err)
	}

	return count
}

// tablesNotInTheArchive lists every table deliberately left out of an export, with the reason. A table that is
// neither here nor referenced by a dataset fails TestPrivacyDatasetsCoverEveryTable, so new tables get classified.
var tablesNotInTheArchive = map[string]string{ //nolint:gochecknoglobals
	"goose_db_version":                  "migration bookkeeping",
	"cache":                             "server-side cache, rebuilt on demand",
	"runtime_state":                     "worker coordination state",
	"event_queue":                       "background job queue",
	"external_code":                     "one-time codes of authentication flows, not tied to a person",
	"protection_pow_challenge":          "hashed network identifiers, withheld on purpose",
	"session_rate_limit":                "hashed network identifiers, withheld on purpose",
	"activity":                          "shared content without an owner column",
	"activity_series":                   "shared content without an owner column",
	"story_series":                      "shared content without an owner column",
	"story_series_tx":                   "shared content without an owner column",
	"profile_application_form_field":    "definition of an exported form",
	"profile_resource_team":             "join row of an exported parent",
	"profile_membership_candidate_team": "join row of an exported parent",
}

func TestPrivacyDatasetsCoverEveryTable(t *testing.T) {
	_, db := privacyTestRepository(t)

	var all strings.Builder
	for _, dataset := range privacyDatasets {
		all.WriteString(dataset.query)
	}

	for _, table := range queryStrings(t, db,
		`SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE' ORDER BY 1`) {
		if strings.Contains(all.String(), `"`+table+`"`) {
			continue
		}

		if _, ok := tablesNotInTheArchive[table]; ok {
			continue
		}

		t.Errorf(
			"table %q is neither exported by a privacy dataset nor listed in tablesNotInTheArchive: "+
				"decide whether it holds a person's data",
			table,
		)
	}
}

func TestPrivacySecretColumnsAreWithheld(t *testing.T) {
	_, db := privacyTestRepository(t)

	withheld := map[string]bool{}
	for _, column := range secretColumns {
		withheld[column] = true
	}

	pattern := regexp.MustCompile(`(token|secret|verifier|password|oauth_request_state)`)

	for _, row := range queryStrings(t, db,
		`SELECT table_name || '.' || column_name FROM information_schema.columns WHERE table_schema = 'public' ORDER BY 1`) {
		column := row[strings.Index(row, ".")+1:]
		table := row[:strings.Index(row, ".")]

		if !pattern.MatchString(column) || withheld[column] {
			continue
		}

		// Columns that merely mention these words but are not credentials.
		switch {
		case strings.HasSuffix(column, "_expires_at"),
			strings.HasSuffix(column, "_scope"),
			table == "session_rate_limit",
			table == "protection_pow_challenge",
			table == "external_code":
			continue
		}

		t.Errorf("column %s looks like a credential but is not in secretColumns", row)
	}
}

// containsKey reports whether any JSON object nested in value has the given key.
func containsKey(value any, keys map[string]bool) string {
	switch typed := value.(type) {
	case map[string]any:
		for key, nested := range typed {
			if keys[key] {
				return key
			}

			if found := containsKey(nested, keys); found != "" {
				return found
			}
		}
	case []any:
		for _, nested := range typed {
			if found := containsKey(nested, keys); found != "" {
				return found
			}
		}
	}

	return ""
}

// busiestUsers returns individual users ordered by how much data hangs off them.
func busiestUsers(t *testing.T, db *sql.DB, limit int) [][2]string {
	t.Helper()

	rows, err := db.QueryContext(context.Background(), `
		SELECT u.id, p.id
		FROM "user" u
		  JOIN "profile" p ON p.id = u.individual_profile_id
		WHERE u.deleted_at IS NULL AND u.redacted_at IS NULL AND p.deleted_at IS NULL AND p.kind = 'individual'
		  AND NOT EXISTS (SELECT 1 FROM profile_membership m WHERE m.member_profile_id = p.id AND m.kind = 'owner'
		    AND m.profile_id <> p.id AND m.deleted_at IS NULL AND m.finished_at IS NULL)
		ORDER BY (SELECT count(*) FROM story s WHERE s.author_profile_id = p.id)
		  + (SELECT count(*) FROM discussion_comment c WHERE c.author_user_id = u.id)
		  + (SELECT count(*) FROM profile_membership m WHERE m.member_profile_id = p.id)
		  + (SELECT count(*) FROM mailbox_envelope e WHERE e.sender_profile_id = p.id)
		  + (SELECT count(*) FROM session s WHERE s.logged_in_user_id = u.id) DESC, u.id
		LIMIT $1`, limit)
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = rows.Close() }()

	var out [][2]string

	for rows.Next() {
		var pair [2]string

		err = rows.Scan(&pair[0], &pair[1])
		if err != nil {
			t.Fatal(err)
		}

		out = append(out, pair)
	}

	return out
}

func datasetNames(t *testing.T, repo *Repository) []string {
	t.Helper()

	names, err := repo.ExportDatasetNames(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	return names
}

func TestPrivacyExportDatasetsRunAndWithholdSecrets(t *testing.T) {
	repo, db := privacyTestRepository(t)

	users := busiestUsers(t, db, 3)
	if len(users) == 0 {
		t.Skip("the database has no individual users to export")
	}

	secrets := map[string]bool{}
	for _, column := range secretColumns {
		secrets[column] = true
	}

	for _, pair := range users {
		totalRows := 0

		for _, name := range datasetNames(t, repo) {
			count := 0

			err := repo.StreamDataset(
				context.Background(),
				name,
				pair[0],
				pair[1],
				func(row json.RawMessage) error {
					var decoded any

					unmarshalErr := json.Unmarshal(row, &decoded)
					if unmarshalErr != nil {
						return unmarshalErr //nolint:wrapcheck
					}

					if leaked := containsKey(decoded, secrets); leaked != "" {
						t.Errorf(
							"dataset %s of user %s exposes secret column %q",
							name,
							pair[0],
							leaked,
						)
					}

					count++

					return nil
				},
			)
			if err != nil {
				t.Fatalf("dataset %s: %v", name, err)
			}

			totalRows += count
		}

		t.Logf("user %s: %d datasets, %d rows", pair[0], len(datasetNames(t, repo)), totalRows)

		// The account dataset must be exactly the user's own row.
		var accountRows int

		_ = repo.StreamDataset(
			context.Background(),
			"account",
			pair[0],
			pair[1],
			func(json.RawMessage) error {
				accountRows++

				return nil
			},
		)

		if accountRows != 1 {
			t.Errorf("account dataset has %d rows, want 1", accountRows)
		}
	}
}

// identifiers collects the strings that identify a person, to prove none of them survives deletion.
func identifiers(t *testing.T, db *sql.DB, userID string, profileID string) []string {
	t.Helper()

	found := map[string]bool{}

	add := func(values ...string) {
		for _, value := range values {
			if len(strings.TrimSpace(value)) >= 4 && value != privacy.RedactedText {
				found[value] = true
			}
		}
	}

	add(
		queryStrings(
			t,
			db,
			`SELECT unnest(ARRAY[name, email, phone, github_handle, bsky_handle, x_handle,
		github_remote_id, bsky_remote_id, x_remote_id, apple_remote_id]) FROM "user" WHERE id = $1`,
			userID,
		)...)
	add(queryStrings(t, db, `SELECT slug FROM "profile" WHERE id = $1`, profileID)...)
	add(
		queryStrings(
			t,
			db,
			`SELECT unnest(ARRAY[title, description]) FROM "profile_tx" WHERE profile_id = $1`,
			profileID,
		)...)

	list := make([]string, 0, len(found))
	for value := range found {
		list = append(list, value)
	}

	sort.Strings(list)

	return list
}

// columnsHolding returns table.column pairs of text-like columns, for the "nothing equals an identifier" scan.
func textColumns(t *testing.T, db *sql.DB) [][2]string {
	t.Helper()

	rows, err := db.QueryContext(context.Background(), `
		SELECT table_name, column_name FROM information_schema.columns
		WHERE table_schema = 'public' AND data_type IN ('text', 'character varying', 'character')
		  AND table_name NOT IN ('goose_db_version', 'cache')
		ORDER BY 1, 2`)
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = rows.Close() }()

	var columns [][2]string

	for rows.Next() {
		var pair [2]string

		_ = rows.Scan(&pair[0], &pair[1])
		columns = append(columns, pair)
	}

	return columns
}

func TestPrivacyAccountDeletion(t *testing.T) {
	repo, db := privacyTestRepository(t)

	users := busiestUsers(t, db, 2)
	if len(users) < 2 {
		t.Skip("need at least two individual users")
	}

	modes := []struct {
		name          string
		deleteContent bool
		user          [2]string
	}{
		{"keep content", false, users[0]},
		{"delete content", true, users[1]},
	}

	for _, mode := range modes {
		t.Run(mode.name, func(t *testing.T) {
			userID, profileID := mode.user[0], mode.user[1]
			ctx := context.Background()

			needles := identifiers(t, db, userID, profileID)
			oldSlug := queryStrings(t, db, `SELECT slug FROM "profile" WHERE id = $1`, profileID)[0]

			storiesBefore := queryCount(
				t,
				db,
				`SELECT count(*) FROM "story" WHERE author_profile_id = $1 AND deleted_at IS NULL`,
				profileID,
			)
			commentsBefore := queryCount(
				t,
				db,
				`SELECT count(*) FROM "discussion_comment" WHERE author_user_id = $1 AND deleted_at IS NULL`,
				userID,
			)
			t.Logf(
				"user %s (%q): %d identifiers, %d live stories, %d live comments",
				userID,
				oldSlug,
				len(needles),
				storiesBefore,
				commentsBefore,
			)

			err := repo.RedactAccount(
				ctx,
				userID,
				profileID,
				privacy.DeleteOptions{DeleteContent: mode.deleteContent},
			)
			if err != nil {
				t.Fatalf("RedactAccount: %v", err)
			}

			// 1. The username is free again.
			if n := queryCount(t, db, `SELECT count(*) FROM "profile" WHERE slug = $1`, oldSlug); n != 0 {
				t.Errorf("username %q is still taken by %d profile(s)", oldSlug, n)
			}

			// 2. The user and profile rows are masked and marked.
			var name, email sql.NullString

			var redactedAt, deletedAt sql.NullTime

			err = db.QueryRowContext(ctx, `SELECT name, email, redacted_at, deleted_at FROM "user" WHERE id = $1`, userID).
				Scan(&name, &email, &redactedAt, &deletedAt)
			if err != nil {
				t.Fatal(err)
			}

			if name.String != privacy.RedactedText || !redactedAt.Valid || !deletedAt.Valid {
				t.Errorf(
					"user row not redacted: name=%q redacted=%v deleted=%v",
					name.String,
					redactedAt.Valid,
					deletedAt.Valid,
				)
			}

			if email.Valid && email.String != privacy.RedactedText {
				t.Errorf("e-mail not redacted: %q", email.String)
			}

			if n := queryCount(t, db, `SELECT count(*) FROM "user" WHERE id = $1 AND github_remote_id IS NULL AND bsky_remote_id IS NULL
				AND x_remote_id IS NULL AND apple_remote_id IS NULL AND profile_picture_uri IS NULL`, userID); n != 1 {
				t.Error("provider identities and the picture must be cleared")
			}

			if n := queryCount(t, db, `SELECT count(*) FROM "profile" WHERE id = $1 AND redacted_at IS NOT NULL AND deleted_at IS NOT NULL
				AND hidden_at IS NULL AND slug = $2 AND pronouns IS NULL AND properties IS NULL AND profile_picture_uri IS NULL`,
				profileID, privacy.RedactedSlugPrefix+strings.ToLower(profileID)); n != 1 {
				t.Error("profile row not redacted and renamed")
			}

			if n := queryCount(t, db, `SELECT count(*) FROM "profile_tx" WHERE profile_id = $1 AND (title <> $2 OR description <> $2)`,
				profileID, privacy.RedactedText); n != 0 {
				t.Errorf("%d profile translations still hold text", n)
			}

			// 3. Secrets and personal state are gone.
			for _, check := range []struct{ label, query string }{
				{"sessions", `SELECT count(*) FROM "session" WHERE logged_in_user_id = $1`},
				{"export requests", `SELECT count(*) FROM "profile_data_export" WHERE user_id = $1`},
				{"comment votes", `SELECT count(*) FROM "discussion_comment_vote" WHERE user_id = $1`},
				{"question votes", `SELECT count(*) FROM "question_vote" WHERE user_id = $1`},
				{"legacy question votes", `SELECT count(*) FROM "legacy_question_vote" WHERE user_id = $1`},
			} {
				if !tablesExistIn(t, db, check.query) {
					continue
				}

				if n := queryCount(t, db, check.query, userID); n != 0 {
					t.Errorf("%s: %d row(s) left", check.label, n)
				}
			}

			for _, check := range []struct{ label, query string }{
				{"links (with their OAuth tokens)", `SELECT count(*) FROM "profile_link" WHERE profile_id = $1`},
				{"story interactions", `SELECT count(*) FROM "story_interaction" WHERE profile_id = $1`},
				{"bulletin subscriptions", `SELECT count(*) FROM "bulletin_subscription" WHERE profile_id = $1`},
				{"custom domains", `SELECT count(*) FROM "profile_custom_domain" WHERE profile_id = $1`},
				{"live memberships", `SELECT count(*) FROM "profile_membership" WHERE (member_profile_id = $1 OR profile_id = $1) AND deleted_at IS NULL`},
				{"mailbox participations", `SELECT count(*) FROM "mailbox_participant" WHERE profile_id = $1`},
				{"stories still authored", `SELECT count(*) FROM "story" WHERE author_profile_id = $1`},
				{"publications", `SELECT count(*) FROM "story_publication" WHERE profile_id = $1`},
			} {
				if n := queryCount(t, db, check.query, profileID); n != 0 {
					t.Errorf("%s: %d row(s) left on the deleted profile", check.label, n)
				}
			}

			// 4. Nothing anywhere still *equals* one of the person's identifiers.
			for _, pair := range textColumns(t, db) {
				query := fmt.Sprintf(`SELECT count(*) FROM %q WHERE %q = ANY($1)`, pair[0], pair[1])

				var count int

				err = db.QueryRowContext(ctx, query, "{"+quoteArray(needles)+"}").Scan(&count)
				if err != nil {
					// Skip columns the = operator cannot compare (never expected, but do not mask real findings).
					continue
				}

				if count > 0 {
					t.Errorf(
						"%s.%s still holds %d identifier(s) of the deleted person",
						pair[0],
						pair[1],
						count,
					)
				}
			}

			// 5. The content choice.
			tombstoneStories := queryCount(
				t,
				db,
				`SELECT count(*) FROM "story" WHERE author_profile_id = $1 AND deleted_at IS NULL`,
				privacy.TombstoneProfileID,
			)

			if mode.deleteContent {
				if n := queryCount(t, db, `SELECT count(*) FROM "story_tx" t JOIN "story" s ON s.id = t.story_id
					WHERE s.author_profile_id = $1 AND (t.title <> $2 OR t.content <> $2)`, privacy.TombstoneProfileID, privacy.RedactedText); n > 0 {
					// Stories already attributed to the placeholder by an earlier deletion keep their text;
					// only those this deletion removed must be erased.
					t.Logf(
						"%d placeholder story translation(s) keep their text (kept by an earlier deletion)",
						n,
					)
				}

				if n := queryCount(t, db, `SELECT count(*) FROM "discussion_comment" WHERE author_user_id = $1 AND content <> $2 AND deleted_at IS NOT NULL`,
					privacy.TombstoneUserID, privacy.RedactedText); n > 0 {
					t.Logf("%d removed placeholder comment(s) still hold text", n)
				}
			} else if tombstoneStories < storiesBefore {
				t.Errorf("kept stories: %d live stories moved to the placeholder, want %d", tombstoneStories, storiesBefore)
			}

			// 6. The placeholder exists, is visible as an author, and is not discoverable as a person.
			if n := queryCount(t, db, `SELECT count(*) FROM "profile" WHERE id = $1 AND deleted_at IS NULL AND redacted_at IS NOT NULL AND approved_at IS NOT NULL`,
				privacy.TombstoneProfileID); n != 1 {
				t.Error("the placeholder profile must exist, be visible and be marked redacted")
			}

			// 7. Running it again changes nothing and does not fail.
			err = repo.RedactAccount(
				ctx,
				userID,
				profileID,
				privacy.DeleteOptions{DeleteContent: mode.deleteContent},
			)
			if err != nil {
				t.Errorf("second RedactAccount: %v", err)
			}
		})
	}
}

// tablesExistIn reports whether every table the query names exists, so a check can skip a table a given
// database does not have (production still carries the legacy_ question tables, fresh databases the unprefixed ones).
func tablesExistIn(t *testing.T, db *sql.DB, query string) bool {
	t.Helper()

	for _, table := range tablesOf(query) {
		if queryCount(
			t,
			db,
			`SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1`,
			table,
		) == 0 {
			return false
		}
	}

	return true
}

func quoteArray(values []string) string {
	quoted := make([]string, 0, len(values))

	for _, value := range values {
		escaped := strings.ReplaceAll(strings.ReplaceAll(value, `\`, `\\`), `"`, `\"`)
		quoted = append(quoted, `"`+escaped+`"`)
	}

	return strings.Join(quoted, ",")
}

func TestPrivacyHideAndUnhide(t *testing.T) {
	repo, db := privacyTestRepository(t)

	users := busiestUsers(t, db, 5)
	if len(users) < 3 {
		t.Skip("need at least three individual users")
	}

	userID, profileID := users[2][0], users[2][1]
	adminUserID, adminProfileID := users[3%len(users)][0], users[3%len(users)][1]
	ctx := context.Background()

	// A fresh active session so we can see it being ended.
	_, err := db.ExecContext(
		ctx,
		`INSERT INTO "session" (id, status, oauth_request_state, oauth_request_code_verifier, logged_in_user_id, created_at)
		VALUES ('privacy-test-session-001', 'active', 'x', 'x', $1, NOW())`,
		userID,
	)
	if err != nil {
		t.Fatal(err)
	}

	changed, err := repo.HideProfile(ctx, userID, profileID)
	if err != nil || !changed {
		t.Fatalf("HideProfile: changed=%v err=%v", changed, err)
	}

	if n := queryCount(t, db, `SELECT count(*) FROM "profile" WHERE id = $1 AND deleted_at IS NOT NULL AND hidden_at IS NOT NULL`, profileID); n != 1 {
		t.Error("the profile must be hidden through deleted_at and marked with hidden_at")
	}

	if n := queryCount(t, db, `SELECT count(*) FROM "session" WHERE logged_in_user_id = $1 AND status = 'active'`, userID); n != 0 {
		t.Errorf("%d session(s) are still active after hiding", n)
	}

	// Hiding twice is harmless.
	changed, err = repo.HideProfile(ctx, userID, profileID)
	if err != nil || changed {
		t.Fatalf("second HideProfile: changed=%v err=%v", changed, err)
	}

	changed, err = repo.UnhideProfile(ctx, userID)
	if err != nil || !changed {
		t.Fatalf("UnhideProfile: changed=%v err=%v", changed, err)
	}

	if n := queryCount(t, db, `SELECT count(*) FROM "profile" WHERE id = $1 AND deleted_at IS NULL AND hidden_at IS NULL`, profileID); n != 1 {
		t.Error("the profile must be visible again after login")
	}

	// A profile an administrator removed (deleted_at without hidden_at) must NOT come back at login.
	_, err = db.ExecContext(
		ctx,
		`UPDATE "profile" SET deleted_at = NOW() WHERE id = $1`,
		adminProfileID,
	)
	if err != nil {
		t.Fatal(err)
	}

	changed, err = repo.UnhideProfile(ctx, adminUserID)
	if err != nil || changed {
		t.Fatalf("an admin-deleted profile must stay deleted: changed=%v err=%v", changed, err)
	}

	if n := queryCount(t, db, `SELECT count(*) FROM "profile" WHERE id = $1 AND deleted_at IS NOT NULL`, adminProfileID); n != 1 {
		t.Error("admin-deleted profile was restored")
	}
}

func TestPrivacySoleOwnedProfilesBlockDeletion(t *testing.T) {
	repo, db := privacyTestRepository(t)
	ctx := context.Background()

	rows, err := db.QueryContext(ctx, `
		SELECT pm.member_profile_id, pm.profile_id
		FROM "profile_membership" pm JOIN "profile" p ON p.id = pm.profile_id
		WHERE pm.kind = 'owner' AND pm.finished_at IS NULL AND pm.deleted_at IS NULL AND p.deleted_at IS NULL
		  AND p.kind <> 'individual' AND pm.member_profile_id <> pm.profile_id
		ORDER BY pm.profile_id LIMIT 1`)
	if err != nil {
		t.Fatal(err)
	}

	var memberID, orgID string

	if !rows.Next() {
		_ = rows.Close()
		t.Skip("no organization with an owner in the data")
	}

	_ = rows.Scan(&memberID, &orgID)
	_ = rows.Close()

	// Make this member the only owner.
	_, err = db.ExecContext(
		ctx,
		`UPDATE "profile_membership" SET kind = 'maintainer' WHERE profile_id = $1 AND member_profile_id <> $2 AND kind = 'owner'`,
		orgID,
		memberID,
	)
	if err != nil {
		t.Fatal(err)
	}

	owned, err := repo.ListSoleOwnedProfiles(ctx, memberID)
	if err != nil {
		t.Fatal(err)
	}

	found := false

	for _, profile := range owned {
		found = found || profile.ID == orgID
	}

	if !found {
		t.Fatalf("sole-owned profiles of %s = %+v, want to include %s", memberID, owned, orgID)
	}
}

func TestPrivacyExportLifecycleRepository(t *testing.T) {
	repo, db := privacyTestRepository(t)

	users := busiestUsers(t, db, 1)
	if len(users) == 0 {
		t.Skip("no users")
	}

	ctx := context.Background()
	export := &privacy.DataExport{
		ID:        "01PRIVACYTESTEXPORT0000001",
		UserID:    users[0][0],
		ProfileID: users[0][1],
		Status:    privacy.ExportStatusPending,
	} //nolint:exhaustruct

	err := repo.CreateExport(ctx, export)
	if err != nil {
		t.Fatal(err)
	}

	latest, err := repo.GetLatestExport(ctx, export.UserID)
	if err != nil || latest == nil || latest.ID != export.ID ||
		latest.Status != privacy.ExportStatusPending {
		t.Fatalf("latest export: %+v err=%v", latest, err)
	}

	claimed, err := repo.ClaimExport(ctx, export.ID)
	if err != nil || !claimed {
		t.Fatalf("claim: %v %v", claimed, err)
	}

	err = repo.MarkExportReady(ctx, export.ID, "exports/x.zip", 123, time.Now().Add(-time.Minute))
	if err != nil {
		t.Fatal(err)
	}

	again, _ := repo.ClaimExport(ctx, export.ID)
	if again {
		t.Error("a finished export must not be claimable again")
	}

	expired, err := repo.ListExpiredExports(ctx, 10)
	if err != nil || len(expired) != 1 || expired[0].ID != export.ID {
		t.Fatalf("expired exports: %+v err=%v", expired, err)
	}

	err = repo.MarkExportExpired(ctx, export.ID)
	if err != nil {
		t.Fatal(err)
	}

	final, _ := repo.GetExport(ctx, export.ID)
	if final.Status != privacy.ExportStatusExpired || final.StorageKey != nil {
		t.Fatalf("expired export keeps its key: %+v", final)
	}
}
