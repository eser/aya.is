package storage

import (
	"regexp"
	"strings"
)

// privacyDataset is one file of a personal data export. Its query takes the user id as $1 and the individual
// profile id as $2 and returns one jsonb document ("doc") per row.
//
// Documents are built with to_jsonb(row) minus an explicit deny-list of secret columns, so a column added to a
// table later is exported automatically, while secrets must be named here to be withheld. The unit test
// TestPrivacyDatasetsDoNotExposeSecrets guards the deny-list; TestPrivacyDatasetsCoverEveryTable (integration)
// guards that new tables are classified.
type privacyDataset struct {
	name  string
	query string
}

// secretColumns are never exported, whatever table they appear in.
var secretColumns = []string{ //nolint:gochecknoglobals
	"auth_access_token",
	"auth_refresh_token",
	"oauth_access_token",
	"oauth_request_code_verifier",
	"oauth_request_state",
	"search_vector",
	"storage_key",
}

// docQuery builds `SELECT <row as jsonb minus secrets> FROM <table> x WHERE <where>`; extraExcluded names
// columns that are withheld for this dataset only (for example the identity of an anonymous asker).
func docQuery(table string, where string, extraExcluded ...string) string {
	excluded := make([]string, 0, len(secretColumns)+len(extraExcluded))
	excluded = append(excluded, secretColumns...)
	excluded = append(excluded, extraExcluded...)

	quoted := make([]string, 0, len(excluded))
	for _, column := range excluded {
		quoted = append(quoted, "'"+column+"'")
	}

	// The cross join references both parameters in every query: PostgreSQL cannot infer the type of a parameter
	// a statement does not mention, which would fail datasets that only need the user id or only the profile id.
	return "SELECT d.doc FROM (SELECT to_jsonb(x) - ARRAY[" + strings.Join(
		quoted,
		", ",
	) + "]::text[] AS doc " +
		"FROM " + table + " x WHERE " + where + ") d " +
		"CROSS JOIN (SELECT $1::text AS user_id, $2::text AS profile_id) params ORDER BY d.doc::text"
}

const (
	ownedLinks      = `(SELECT id FROM "profile_link" WHERE profile_id = $2 OR added_by_profile_id = $2)`
	ownedPages      = `(SELECT id FROM "profile_page" WHERE profile_id = $2 OR added_by_profile_id = $2)`
	ownedStories    = `(SELECT id FROM "story" WHERE author_profile_id = $2)`
	ownMemberships  = `(SELECT id FROM "profile_membership" WHERE member_profile_id = $2)`
	ownCandidates   = `(SELECT id FROM "profile_membership_candidate" WHERE profile_id = $2 OR referred_profile_id = $2)`
	ownSubscription = `(SELECT id FROM "bulletin_subscription" WHERE profile_id = $2)`
	ownSessions     = `(SELECT id FROM "session" WHERE logged_in_user_id = $1)`
	ownConversation = `(SELECT conversation_id FROM "mailbox_participant" WHERE profile_id = $2)`
)

// privacyDatasets is the ordered list of files that make up an export archive.
var privacyDatasets = []privacyDataset{ //nolint:gochecknoglobals
	{"account", docQuery(`"user"`, `x.id = $1`)},
	{"profile", docQuery(`"profile"`, `x.id = $2`)},
	{"profile_translations", docQuery(`"profile_tx"`, `x.profile_id = $2`)},
	{
		"profile_links",
		docQuery(`"profile_link"`, `x.profile_id = $2 OR x.added_by_profile_id = $2`),
	},
	{
		"profile_link_translations",
		docQuery(`"profile_link_tx"`, `x.profile_link_id IN `+ownedLinks),
	},
	{"profile_link_imports", docQuery(`"profile_link_import"`, `x.profile_link_id IN `+ownedLinks)},
	{
		"profile_pages",
		docQuery(`"profile_page"`, `x.profile_id = $2 OR x.added_by_profile_id = $2`),
	},
	{
		"profile_page_translations",
		docQuery(`"profile_page_tx"`, `x.profile_page_id IN `+ownedPages),
	},
	{
		"profile_resources",
		docQuery(`"profile_resource"`, `x.profile_id = $2 OR x.added_by_profile_id = $2`),
	},
	{"profile_custom_domains", docQuery(`"profile_custom_domain"`, `x.profile_id = $2`)},
	{"profile_teams", docQuery(`"profile_team"`, `x.profile_id = $2`)},
	{"profile_application_forms", docQuery(`"profile_application_form"`, `x.profile_id = $2`)},
	{
		"memberships",
		docQuery(`"profile_membership"`, `x.member_profile_id = $2 OR x.profile_id = $2`),
	},
	{
		"membership_teams",
		docQuery(`"profile_membership_team"`, `x.profile_membership_id IN `+ownMemberships),
	},
	{
		"membership_candidates",
		docQuery(
			`"profile_membership_candidate"`,
			`x.profile_id = $2 OR x.referred_profile_id = $2`,
		),
	},
	{
		"membership_candidate_responses",
		docQuery(`"profile_candidate_response"`, `x.candidate_id IN `+ownCandidates),
	},
	{
		"membership_candidate_votes",
		docQuery(`"profile_membership_candidate_vote"`, `x.voter_membership_id IN `+ownMemberships),
	},
	{"stories", docQuery(`"story"`, `x.author_profile_id = $2`)},
	{"story_translations", docQuery(`"story_tx"`, `x.story_id IN `+ownedStories)},
	{
		"story_publications",
		docQuery(`"story_publication"`, `x.profile_id = $2 OR x.story_id IN `+ownedStories),
	},
	{"story_interactions", docQuery(`"story_interaction"`, `x.profile_id = $2`)},
	{"story_date_proposals", docQuery(`"story_date_proposal"`, `x.proposer_profile_id = $2`)},
	{
		"story_date_proposal_votes",
		docQuery(`"story_date_proposal_vote"`, `x.voter_profile_id = $2`),
	},
	{"discussion_threads", docQuery(`"discussion_thread"`, `x.profile_id = $2`)},
	{"discussion_comments", docQuery(`"discussion_comment"`, `x.author_user_id = $1`)},
	{"discussion_comment_votes", docQuery(`"discussion_comment_vote"`, `x.user_id = $1`)},
	{
		"mailbox_conversations",
		docQuery(
			`"mailbox_conversation"`,
			`x.created_by_profile_id = $2 OR x.id IN `+ownConversation,
		),
	},
	{"mailbox_participations", docQuery(`"mailbox_participant"`, `x.profile_id = $2`)},
	{
		"mailbox_envelopes",
		docQuery(
			`"mailbox_envelope"`,
			`x.sender_profile_id = $2 OR x.sender_user_id = $1 OR x.target_profile_id = $2`,
		),
	},
	{"mailbox_reactions", docQuery(`"mailbox_reaction"`, `x.profile_id = $2`)},
	{"profile_questions_asked", docQuery(`"profile_question"`, `x.author_user_id = $1`)},
	{
		"profile_questions_received",
		docQuery(`"profile_question"`, `x.profile_id = $2`, "author_user_id"),
	},
	{"profile_question_votes", docQuery(`"profile_question_vote"`, `x.user_id = $1`)},
	{"activity_attendance", docQuery(`"activity_attendance"`, `x.profile_id = $2`)},
	{"bulletin_subscriptions", docQuery(`"bulletin_subscription"`, `x.profile_id = $2`)},
	{"bulletin_log", docQuery(`"bulletin_log"`, `x.subscription_id IN `+ownSubscription)},
	{
		"points_transactions",
		docQuery(
			`"profile_point_transaction"`,
			`x.target_profile_id = $2 OR x.origin_profile_id = $2`,
		),
	},
	{
		"points_pending_awards",
		docQuery(`"profile_point_pending_award"`, `x.target_profile_id = $2`),
	},
	{"sessions", docQuery(`"session"`, `x.logged_in_user_id = $1`)},
	{"session_preferences", docQuery(`"session_preference"`, `x.session_id IN `+ownSessions)},
	{
		"audit_events",
		docQuery(`"event_audit"`, `x.actor_id IN ($1, $2) OR x.entity_id IN ($1, $2)`),
	},
	{"data_exports", docQuery(`"profile_data_export"`, `x.user_id = $1`)},
}

// questionDatasets returns the datasets of one questions table. Production still carries the old Q&A tables under
// their legacy_ names while a freshly migrated database has them unprefixed; each variant is exported when its
// table exists and skipped when it does not (see tablesOf).
func questionDatasets(prefix string, questionTable string, voteTable string) []privacyDataset {
	return []privacyDataset{
		{prefix + "questions_asked", docQuery(`"`+questionTable+`"`, `x.user_id = $1`)},
		{
			prefix + "questions_received",
			docQuery(`"`+questionTable+`"`, `x.profile_id = $2`, "user_id"),
		},
		{prefix + "question_votes", docQuery(`"`+voteTable+`"`, `x.user_id = $1`)},
	}
}

func init() { //nolint:gochecknoinits
	privacyDatasets = append(privacyDatasets, questionDatasets("", "question", "question_vote")...)
	privacyDatasets = append(
		privacyDatasets,
		questionDatasets("legacy_", "legacy_question", "legacy_question_vote")...)
}

var quotedIdentifier = regexp.MustCompile(`"([a-z_]+)"`)

// tablesOf returns the tables a query reads, which are the only identifiers this package double-quotes.
func tablesOf(query string) []string {
	matches := quotedIdentifier.FindAllStringSubmatch(query, -1)
	tables := make([]string, 0, len(matches))

	for _, match := range matches {
		tables = append(tables, match[1])
	}

	return tables
}

// allTablesExist reports whether every table the query reads is present.
func allTablesExist(query string, existing map[string]bool) bool {
	for _, table := range tablesOf(query) {
		if !existing[table] {
			return false
		}
	}

	return true
}

// privacyDatasetsByName indexes privacyDatasets.
func privacyDatasetsByName() map[string]privacyDataset {
	index := make(map[string]privacyDataset, len(privacyDatasets))
	for _, dataset := range privacyDatasets {
		index[dataset.name] = dataset
	}

	return index
}
