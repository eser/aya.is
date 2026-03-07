package story_date_proposals

import (
	"context"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
)

// IDGenerator is a function that generates unique IDs.
type IDGenerator func() string

// DefaultIDGenerator returns a new ULID.
func DefaultIDGenerator() string {
	return lib.IDsGenerateUnique()
}

// VoteDirection defines agree/disagree.
type VoteDirection int

const (
	VoteAgree    VoteDirection = 1
	VoteDisagree VoteDirection = -1
)

// voteFlipMultiplier is the score multiplier when flipping a vote direction.
const voteFlipMultiplier = 2

// StoryProvider retrieves activity configuration from the stories domain.
type StoryProvider interface {
	GetActivityDateConfig(ctx context.Context, storyID string) (*stories.ActivityDateConfig, error)
	FinalizeActivityDate(
		ctx context.Context,
		storyID string,
		datetimeStart time.Time,
		datetimeEnd *time.Time,
	) error
}

// AccessChecker checks membership levels between profiles.
type AccessChecker interface {
	GetMembershipKindBetween(ctx context.Context, profileID, memberProfileID string) (string, error)
}

// DateProposal represents a proposed date for an activity.
type DateProposal struct {
	ID                string     `json:"id"`
	StoryID           string     `json:"story_id"`
	ProposerProfileID string     `json:"proposer_profile_id"`
	DatetimeStart     time.Time  `json:"datetime_start"`
	DatetimeEnd       *time.Time `json:"datetime_end"`
	IsFinalized       bool       `json:"is_finalized"`
	VoteScore         int        `json:"vote_score"`
	UpvoteCount       int        `json:"upvote_count"`
	DownvoteCount     int        `json:"downvote_count"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
}

// DateProposalWithProfile extends DateProposal with proposer's profile display info.
type DateProposalWithProfile struct {
	ID                        string     `json:"id"`
	StoryID                   string     `json:"story_id"`
	ProposerProfileID         string     `json:"proposer_profile_id"`
	DatetimeStart             time.Time  `json:"datetime_start"`
	DatetimeEnd               *time.Time `json:"datetime_end"`
	IsFinalized               bool       `json:"is_finalized"`
	VoteScore                 int        `json:"vote_score"`
	UpvoteCount               int        `json:"upvote_count"`
	DownvoteCount             int        `json:"downvote_count"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 *time.Time `json:"updated_at"`
	ProposerProfileSlug       string     `json:"proposer_profile_slug"`
	ProposerProfileTitle      string     `json:"proposer_profile_title"`
	ProposerProfilePictureURI *string    `json:"proposer_profile_picture_uri"`
	ProposerProfileKind       string     `json:"proposer_profile_kind"`
	ViewerVoteDirection       int        `json:"viewer_vote_direction"`
}

// DateProposalVote represents a user's vote on a proposal.
type DateProposalVote struct {
	ID             string    `json:"id"`
	ProposalID     string    `json:"proposal_id"`
	VoterProfileID string    `json:"voter_profile_id"`
	Direction      int       `json:"direction"`
	CreatedAt      time.Time `json:"created_at"`
}

// VoteResponse is returned after a vote operation.
type VoteResponse struct {
	VoteScore           int `json:"vote_score"`
	ViewerVoteDirection int `json:"viewer_vote_direction"`
}

// voteDeltas holds the computed delta values from a vote operation.
type voteDeltas struct {
	score           int
	upvote          int
	downvote        int
	viewerDirection int
}
