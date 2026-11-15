//nolint:lll // the SQL statements read best unwrapped, one per step
package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/eser/aya.is/services/pkg/api/business/privacy"
)

// redactionStep is one statement of the account deletion transaction.
type redactionStep struct {
	label string
	query string
	// when restricts the step to one of the two content modes; empty means "always".
	when redactionMode
	args []any
}

type redactionMode string

const (
	modeAlways     redactionMode = ""
	modeKeepOnly   redactionMode = "keep"
	modeDeleteOnly redactionMode = "delete"
)

// redactionSteps lists, in execution order, everything that happens when an account is deleted.
//
// The strategy keeps foreign keys intact and removes the *personal data* rather than the rows: rows that
// are pure personal state (sessions, votes, reactions, subscriptions, OAuth links) are deleted; rows that
// other people's data hangs off (memberships, stories, comments, conversations) are soft-deleted or handed
// to the shared "deleted user" placeholder, and every free-text field of the departing user is overwritten
// with the redaction marker. oldSlug is the username being freed and profileSlug the new, unusable slug of the deleted profile.
//
//nolint:funlen,maintidx // one long, flat list is the clearest form for an erasure procedure that is reviewed line by line
func redactionSteps(
	userID string,
	profileID string,
	oldSlug string,
	profileSlug string,
) []redactionStep {
	const (
		tombUser    = privacy.TombstoneUserID
		tombProfile = privacy.TombstoneProfileID
		red         = privacy.RedactedText
	)

	ownMemberships := `(SELECT id FROM "profile_membership" WHERE member_profile_id = $1)`
	ownCandidates := `(SELECT id FROM "profile_membership_candidate" WHERE profile_id = $1 OR referred_profile_id = $1)`
	ownLinks := `(SELECT id FROM "profile_link" WHERE profile_id = $1)`
	ownStories := `(SELECT id FROM "story" WHERE author_profile_id = $1)`
	ownPages := `(SELECT id FROM "profile_page" WHERE profile_id = $1)`

	steps := []redactionStep{
		// ---- the shared placeholder that inherits kept content --------------------------------------------
		{
			"ensure placeholder profile", `INSERT INTO "profile" (id, slug, kind, default_locale, approved_at, redacted_at)
			VALUES ($1, $2, 'individual', 'en', NOW(), NOW()) ON CONFLICT (id) DO NOTHING`,
			modeAlways,
			[]any{tombProfile, privacy.TombstoneProfileSlug},
		},
		{
			"ensure placeholder profile text", `INSERT INTO "profile_tx" (profile_id, locale_code, title, description)
			VALUES ($1, 'en', $2, '') ON CONFLICT (profile_id, locale_code) DO NOTHING`,
			modeAlways,
			[]any{tombProfile, privacy.TombstoneTitle},
		},
		{
			"ensure placeholder user", `INSERT INTO "user" (id, kind, name, individual_profile_id, redacted_at)
			VALUES ($1, 'individual', $2, $3, NOW()) ON CONFLICT (id) DO NOTHING`,
			modeAlways,
			[]any{tombUser, privacy.TombstoneTitle, tombProfile},
		},

		// ---- secrets and pure personal state: deleted -----------------------------------------------------
		{
			"sessions",
			`DELETE FROM "session" WHERE logged_in_user_id = $1`,
			modeAlways,
			[]any{userID},
		},
		{
			"comment votes",
			`DELETE FROM "discussion_comment_vote" WHERE user_id = $1`,
			modeAlways,
			[]any{userID},
		},
		{
			"profile question votes",
			`DELETE FROM "profile_question_vote" WHERE user_id = $1`,
			modeAlways,
			[]any{userID},
		},
		{
			"date proposal votes",
			`DELETE FROM "story_date_proposal_vote" WHERE voter_profile_id = $1`,
			modeAlways,
			[]any{profileID},
		},
		{
			"story interactions",
			`DELETE FROM "story_interaction" WHERE profile_id = $1`,
			modeAlways,
			[]any{profileID},
		},
		{
			"mailbox reactions",
			`DELETE FROM "mailbox_reaction" WHERE profile_id = $1`,
			modeAlways,
			[]any{profileID},
		},
		{
			"activity attendance",
			`DELETE FROM "activity_attendance" WHERE profile_id = $1`,
			modeAlways,
			[]any{profileID},
		},
		{
			"candidate votes cast", `DELETE FROM "profile_membership_candidate_vote" WHERE voter_membership_id IN ` + ownMemberships,
			modeAlways,
			[]any{profileID},
		},
		{"bulletin log", `DELETE FROM "bulletin_log" WHERE subscription_id IN
			(SELECT id FROM "bulletin_subscription" WHERE profile_id = $1)`, modeAlways, []any{profileID}},
		{
			"bulletin subscriptions",
			`DELETE FROM "bulletin_subscription" WHERE profile_id = $1`,
			modeAlways,
			[]any{profileID},
		},
		{
			"link imports",
			`DELETE FROM "profile_link_import" WHERE profile_link_id IN ` + ownLinks,
			modeAlways,
			[]any{profileID},
		},
		{
			"link translations",
			`DELETE FROM "profile_link_tx" WHERE profile_link_id IN ` + ownLinks,
			modeAlways,
			[]any{profileID},
		},
		{
			"links (and their OAuth tokens)",
			`DELETE FROM "profile_link" WHERE profile_id = $1`,
			modeAlways,
			[]any{profileID},
		},
		{
			"custom domains",
			`DELETE FROM "profile_custom_domain" WHERE profile_id = $1`,
			modeAlways,
			[]any{profileID},
		},
		{
			"export requests",
			`DELETE FROM "profile_data_export" WHERE user_id = $1`,
			modeAlways,
			[]any{userID},
		},

		// ---- applications and candidacies about this person: deleted with everything attached ---------------
		{
			"candidate responses",
			`DELETE FROM "profile_candidate_response" WHERE candidate_id IN ` + ownCandidates,
			modeAlways,
			[]any{profileID},
		},
		{
			"candidate teams",
			`DELETE FROM "profile_membership_candidate_team" WHERE candidate_id IN ` + ownCandidates,
			modeAlways,
			[]any{profileID},
		},
		{
			"candidate votes",
			`DELETE FROM "profile_membership_candidate_vote" WHERE candidate_id IN ` + ownCandidates,
			modeAlways,
			[]any{profileID},
		},
		{
			"candidates",
			`DELETE FROM "profile_membership_candidate" WHERE profile_id = $1 OR referred_profile_id = $1`,
			modeAlways,
			[]any{profileID},
		},
		{
			"application forms",
			`DELETE FROM "profile_application_form" WHERE profile_id = $1`,
			modeAlways,
			[]any{profileID},
		},

		// ---- memberships, teams, resources: soft-deleted and scrubbed -----------------------------------------
		{
			"membership teams", `UPDATE "profile_membership_team" SET deleted_at = COALESCE(deleted_at, NOW())
			WHERE profile_membership_id IN (SELECT id FROM "profile_membership" WHERE member_profile_id = $1 OR profile_id = $1)`,
			modeAlways,
			[]any{profileID},
		},
		{"memberships", `UPDATE "profile_membership" SET deleted_at = COALESCE(deleted_at, NOW()),
			finished_at = COALESCE(finished_at, NOW()), properties = NULL
			WHERE member_profile_id = $1 OR profile_id = $1`, modeAlways, []any{profileID}},
		{
			"teams",
			`UPDATE "profile_team" SET deleted_at = COALESCE(deleted_at, NOW()) WHERE profile_id = $1`,
			modeAlways,
			[]any{profileID},
		},
		{
			"resources",
			`UPDATE "profile_resource" SET deleted_at = COALESCE(deleted_at, NOW()), title = $2, description = $2,
			url = $2, remote_id = NULL, public_id = NULL, properties = NULL WHERE profile_id = $1`,
			modeAlways,
			[]any{profileID, red},
		},
		{
			"things the user added to other profiles (links)", `UPDATE "profile_link" SET added_by_profile_id = $2 WHERE added_by_profile_id = $1`,
			modeAlways,
			[]any{profileID, tombProfile},
		},
		{
			"things the user added to other profiles (resources)", `UPDATE "profile_resource" SET added_by_profile_id = $2 WHERE added_by_profile_id = $1`,
			modeAlways,
			[]any{profileID, tombProfile},
		},
		{
			"things the user added to other profiles (pages)",
			`UPDATE "profile_page" SET added_by_profile_id = $2
			WHERE added_by_profile_id = $1 AND profile_id <> $1`,
			modeAlways,
			[]any{profileID, tombProfile},
		},

		// ---- pages ---------------------------------------------------------------------------------------------
		{
			"pages: erase text", `UPDATE "profile_page_tx" SET title = $2, summary = $2, content = $2 WHERE profile_page_id IN ` + ownPages,
			modeDeleteOnly,
			[]any{profileID, red},
		},
		{
			"pages: remove",
			`UPDATE "profile_page" SET deleted_at = COALESCE(deleted_at, NOW()), cover_picture_uri = NULL
			WHERE profile_id = $1`,
			modeDeleteOnly,
			[]any{profileID},
		},
		{"pages: hand over", `UPDATE "profile_page" SET profile_id = $2, added_by_profile_id = $2,
			slug = slug || '-' || lower(right(id, 8)) WHERE profile_id = $1`, modeKeepOnly, []any{profileID, tombProfile}},

		// ---- stories -------------------------------------------------------------------------------------------
		{
			"stories: erase text",
			`UPDATE "story_tx" SET title = $2, summary = $2, content = $2, summary_ai = NULL
			WHERE story_id IN ` + ownStories,
			modeDeleteOnly,
			[]any{profileID, red},
		},
		{
			"stories: remove",
			`UPDATE "story" SET deleted_at = COALESCE(deleted_at, NOW()), story_picture_uri = NULL,
			properties = NULL, remote_id = NULL WHERE author_profile_id = $1`,
			modeDeleteOnly,
			[]any{profileID},
		},
		{
			"publications: hand over",
			`UPDATE "story_publication" sp SET profile_id = $2 WHERE sp.profile_id = $1
			AND NOT EXISTS (SELECT 1 FROM "story_publication" o WHERE o.story_id = sp.story_id
				AND o.profile_id = $2 AND o.kind = sp.kind)`,
			modeKeepOnly,
			[]any{profileID, tombProfile},
		},
		{
			"publications: remove",
			`DELETE FROM "story_publication" WHERE profile_id = $1`,
			modeAlways,
			[]any{profileID},
		},
		{
			"stories: change author", `UPDATE "story" SET author_profile_id = $2 WHERE author_profile_id = $1`,
			modeAlways,
			[]any{profileID, tombProfile},
		},
		{
			"date proposals", `UPDATE "story_date_proposal" SET proposer_profile_id = $2 WHERE proposer_profile_id = $1`,
			modeAlways,
			[]any{profileID, tombProfile},
		},

		// ---- discussions, questions ----------------------------------------------------------------------------
		{
			"threads",
			`UPDATE "discussion_thread" SET profile_id = $2 WHERE profile_id = $1`,
			modeAlways,
			[]any{profileID, tombProfile},
		},
		{
			"comments: erase text",
			`UPDATE "discussion_comment" SET content = $2, deleted_at = COALESCE(deleted_at, NOW())
			WHERE author_user_id = $1`,
			modeDeleteOnly,
			[]any{userID, red},
		},
		{
			"comments: change author", `UPDATE "discussion_comment" SET author_user_id = $2 WHERE author_user_id = $1`,
			modeAlways,
			[]any{userID, tombUser},
		},
		{
			"profile questions asked: erase text",
			`UPDATE "profile_question" SET content = $2, deleted_at = COALESCE(deleted_at, NOW())
			WHERE author_user_id = $1`,
			modeDeleteOnly,
			[]any{userID, red},
		},
		{
			"profile questions asked: change author", `UPDATE "profile_question" SET author_user_id = $2 WHERE author_user_id = $1`,
			modeAlways,
			[]any{userID, tombUser},
		},
		{
			"profile questions answered by the user",
			`UPDATE "profile_question" SET answer_content = $3, answer_uri = NULL, answered_by = $2
			WHERE answered_by = $1 OR profile_id = $1`,
			modeAlways,
			[]any{profileID, tombProfile, red},
		},

		// ---- private messages: anonymised ----------------------------------------------------------------------
		{
			"messages sent: erase text", `UPDATE "mailbox_envelope" SET message = $2 WHERE sender_profile_id = $1`,
			modeDeleteOnly,
			[]any{profileID, red},
		},
		{
			"messages sent: change sender (profile)", `UPDATE "mailbox_envelope" SET sender_profile_id = $2 WHERE sender_profile_id = $1`,
			modeAlways,
			[]any{profileID, tombProfile},
		},
		{
			"messages sent: change sender (user)", `UPDATE "mailbox_envelope" SET sender_user_id = $2 WHERE sender_user_id = $1`,
			modeAlways,
			[]any{userID, tombUser},
		},
		{
			"conversations: change creator", `UPDATE "mailbox_conversation" SET created_by_profile_id = $2 WHERE created_by_profile_id = $1`,
			modeAlways,
			[]any{profileID, tombProfile},
		},
		{
			"conversations: change participant", `UPDATE "mailbox_participant" mp SET profile_id = $2 WHERE mp.profile_id = $1
			AND NOT EXISTS (SELECT 1 FROM "mailbox_participant" o WHERE o.conversation_id = mp.conversation_id AND o.profile_id = $2)`,
			modeAlways,
			[]any{profileID, tombProfile},
		},
		{
			"conversations: leave",
			`DELETE FROM "mailbox_participant" WHERE profile_id = $1`,
			modeAlways,
			[]any{profileID},
		},

		// ---- points, audit trail -----------------------------------------------------------------------------
		{
			"points received: erase descriptions", `UPDATE "profile_point_transaction" SET description = $2 WHERE target_profile_id = $1`,
			modeAlways,
			[]any{profileID, red},
		},
		{
			"points given: change origin", `UPDATE "profile_point_transaction" SET origin_profile_id = $2 WHERE origin_profile_id = $1`,
			modeAlways,
			[]any{profileID, tombProfile},
		},
		{
			"pending awards", `UPDATE "profile_point_pending_award" SET description = $2, metadata = NULL WHERE target_profile_id = $1`,
			modeAlways,
			[]any{profileID, red},
		},
		{"audit trail payloads", `UPDATE "event_audit" SET payload = '{}'::jsonb
			WHERE actor_id IN ($1, $2) OR entity_id IN ($1, $2)`, modeAlways, []any{userID, profileID}},
		// Some events (profile visits) use the username as the entity id.
		{
			"audit trail username references",
			`UPDATE "event_audit" SET entity_id = $2, payload = '{}'::jsonb
			WHERE entity_id = $1 AND $1 <> $2`,
			modeAlways,
			[]any{oldSlug, profileID},
		},

		// ---- the account itself --------------------------------------------------------------------------------
		{
			"profile text", `UPDATE "profile_tx" SET title = $2, description = $2, properties = NULL WHERE profile_id = $1`,
			modeAlways,
			[]any{profileID, red},
		},
		{
			"profile (rename, which frees the username)",
			`UPDATE "profile" SET slug = $2, profile_picture_uri = NULL, pronouns = NULL,
			properties = NULL, hidden_at = NULL, deleted_at = COALESCE(deleted_at, NOW()), redacted_at = NOW(), updated_at = NOW()
			WHERE id = $1`,
			modeAlways,
			[]any{profileID, profileSlug},
		},
		{"user", `UPDATE "user" SET name = $2,
			email = CASE WHEN email IS NULL THEN NULL ELSE $2 END,
			phone = CASE WHEN phone IS NULL THEN NULL ELSE $2 END,
			github_handle = CASE WHEN github_handle IS NULL THEN NULL ELSE $2 END,
			bsky_handle = CASE WHEN bsky_handle IS NULL THEN NULL ELSE $2 END,
			x_handle = CASE WHEN x_handle IS NULL THEN NULL ELSE $2 END,
			github_remote_id = NULL, bsky_remote_id = NULL, x_remote_id = NULL, apple_remote_id = NULL,
			profile_picture_uri = NULL,
			deleted_at = COALESCE(deleted_at, NOW()), redacted_at = NOW(), updated_at = NOW()
			WHERE id = $1`, modeAlways, []any{userID, red}},
	}

	// The old Q&A tables exist under their legacy_ names in production; steps for a missing table are skipped.
	for _, tables := range [][2]string{{"question", "question_vote"}, {"legacy_question", "legacy_question_vote"}} {
		question, vote := `"`+tables[0]+`"`, `"`+tables[1]+`"`

		steps = append(
			steps,
			redactionStep{
				tables[1] + " votes",
				`DELETE FROM ` + vote + ` WHERE user_id = $1`,
				modeAlways,
				[]any{userID},
			},
			redactionStep{
				tables[0] + " asked: erase text",
				`UPDATE ` + question + ` SET content = $2,
				deleted_at = COALESCE(deleted_at, NOW()) WHERE user_id = $1`,
				modeDeleteOnly,
				[]any{userID, red},
			},
			redactionStep{
				tables[0] + " asked: change author", `UPDATE ` + question + ` SET user_id = $2 WHERE user_id = $1`,
				modeAlways,
				[]any{userID, tombUser},
			},
			redactionStep{
				tables[0] + " answered by the user",
				`UPDATE ` + question + ` SET answer_content = $2, answer_uri = NULL
				WHERE profile_id = $1 AND answered_at IS NOT NULL`,
				modeAlways,
				[]any{profileID, red},
			},
		)
	}

	return steps
}

// RedactAccount erases the personal data of a user in a single transaction: either all of it happens or none.
func (r *Repository) RedactAccount(
	ctx context.Context,
	userID string,
	profileID string,
	opts privacy.DeleteOptions,
) error {
	account, err := r.GetAccount(ctx, userID)
	if err != nil {
		return err
	}

	if account == nil {
		return fmt.Errorf("%w: user %s", ErrPrivacyQuery, userID)
	}

	domains, err := r.customDomainsOfProfile(ctx, profileID)
	if err != nil {
		return err
	}

	mode := modeKeepOnly
	if opts.DeleteContent {
		mode = modeDeleteOnly
	}

	slug := privacy.RedactedSlugPrefix + strings.ToLower(profileID)

	err = r.runRedaction(ctx, redactionSteps(userID, profileID, account.ProfileSlug, slug), mode)
	if err != nil {
		return err
	}

	r.invalidateProfileCaches(ctx, userID, account.ProfileSlug, slug)

	for _, domain := range domains {
		_ = r.cache.Invalidate(ctx, "custom_domain_by_domain:"+domain)
	}

	return nil
}

// runRedaction executes the steps that apply to the chosen content mode in one transaction: all or nothing.
func (r *Repository) runRedaction(
	ctx context.Context,
	steps []redactionStep,
	mode redactionMode,
) error {
	existing, err := r.existingTables(ctx)
	if err != nil {
		return err
	}

	dbTx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPrivacyTransaction, err)
	}

	defer func() { _ = dbTx.Rollback() }()

	for _, step := range steps {
		if (step.when != modeAlways && step.when != mode) || !allTablesExist(step.query, existing) {
			continue
		}

		_, execErr := dbTx.ExecContext(ctx, step.query, step.args...)
		if execErr != nil {
			return fmt.Errorf("%w: step %q: %w", ErrPrivacyTransaction, step.label, execErr)
		}
	}

	err = dbTx.Commit()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPrivacyTransaction, err)
	}

	return nil
}

func (r *Repository) customDomainsOfProfile(
	ctx context.Context,
	profileID string,
) ([]string, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT domain FROM "profile_custom_domain" WHERE profile_id = $1`,
		profileID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPrivacyQuery, err)
	}

	defer func() { _ = rows.Close() }()

	domains := make([]string, 0)

	for rows.Next() {
		var domain sql.NullString

		scanErr := rows.Scan(&domain)
		if scanErr != nil {
			return nil, fmt.Errorf("%w: %w", ErrPrivacyQuery, scanErr)
		}

		if domain.Valid {
			domains = append(domains, domain.String)
		}
	}

	return domains, rows.Err() //nolint:wrapcheck
}
