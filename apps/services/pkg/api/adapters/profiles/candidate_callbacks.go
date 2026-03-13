package profiles

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/mailbox"
	profilesbiz "github.com/eser/aya.is/services/pkg/api/business/profiles"
)

// NewCandidateAutoAccepter returns a callback that ensures a membership exists
// (creating or upgrading as needed) and merges teams when a profile_join
// invitation is accepted — completing the candidate flow.
func NewCandidateAutoAccepter(
	profileService *profilesbiz.Service,
	logger *logfx.Logger,
) mailbox.OnEnvelopeAcceptedFunc {
	return func(ctx context.Context, envelope *mailbox.Envelope) {
		if envelope.Kind != mailbox.KindInvitation {
			return
		}

		props, ok := parseProfileJoinProps(ctx, envelope, logger)
		if !ok {
			return
		}

		recipientProfileID := envelope.TargetProfileID

		// Ensure membership exists at member+ level and merge teams.
		membershipID, err := profileService.EnsureMembershipFromCandidateInternal(
			ctx,
			props.ProfileID,
			recipientProfileID,
			props.CandidateID,
		)
		if err != nil {
			logger.ErrorContext(ctx, "Failed to ensure membership from candidate acceptance",
				slog.String("candidate_id", props.CandidateID),
				slog.String("profile_id", props.ProfileID),
				slog.String("recipient_profile_id", recipientProfileID),
				slog.String("error", err.Error()))

			return
		}

		// Update candidate status to accepted.
		statusErr := profileService.UpdateCandidateStatusInternal(
			ctx,
			props.CandidateID,
			props.ProfileID,
			profilesbiz.CandidateStatusInvitationAccepted,
		)
		if statusErr != nil {
			logger.ErrorContext(ctx, "Failed to update candidate status to accepted",
				slog.String("candidate_id", props.CandidateID),
				slog.String("error", statusErr.Error()))

			return
		}

		logger.InfoContext(ctx, "Candidate invitation accepted, membership ensured",
			slog.String("candidate_id", props.CandidateID),
			slog.String("profile_id", props.ProfileID),
			slog.String("recipient_profile_id", recipientProfileID),
			slog.String("membership_id", membershipID))
	}
}

// NewCandidateAutoRejecter returns a callback that updates candidate status to invitation_rejected
// when a profile_join invitation is rejected.
func NewCandidateAutoRejecter(
	profileService *profilesbiz.Service,
	logger *logfx.Logger,
) mailbox.OnEnvelopeRejectedFunc {
	return func(ctx context.Context, envelope *mailbox.Envelope) {
		if envelope.Kind != mailbox.KindInvitation {
			return
		}

		props, ok := parseProfileJoinProps(ctx, envelope, logger)
		if !ok {
			return
		}

		err := profileService.UpdateCandidateStatusInternal(
			ctx,
			props.CandidateID,
			props.ProfileID,
			profilesbiz.CandidateStatusInvitationRejected,
		)
		if err != nil {
			logger.ErrorContext(ctx, "Failed to update candidate status to rejected",
				slog.String("candidate_id", props.CandidateID),
				slog.String("error", err.Error()))

			return
		}

		logger.InfoContext(ctx, "Candidate invitation rejected, status updated",
			slog.String("candidate_id", props.CandidateID),
			slog.String("profile_id", props.ProfileID))
	}
}

// parseProfileJoinProps extracts InvitationProperties from an envelope and checks
// whether it's a profile_join invitation. Returns false if not applicable.
func parseProfileJoinProps(
	ctx context.Context,
	envelope *mailbox.Envelope,
	logger *logfx.Logger,
) (mailbox.InvitationProperties, bool) {
	emptyProps := mailbox.InvitationProperties{
		InvitationKind:   "",
		TelegramChatID:   0,
		GroupProfileSlug: "",
		GroupName:        "",
		InviteLink:       nil,
		CandidateID:      "",
		ProfileID:        "",
		ProfileSlug:      "",
	}

	propsJSON, marshalErr := json.Marshal(envelope.Properties)
	if marshalErr != nil {
		logger.WarnContext(ctx, "Failed to marshal envelope properties",
			slog.String("envelope_id", envelope.ID),
			slog.String("error", marshalErr.Error()))

		return emptyProps, false
	}

	var props mailbox.InvitationProperties

	unmarshalErr := json.Unmarshal(propsJSON, &props)
	if unmarshalErr != nil {
		logger.WarnContext(ctx, "Failed to unmarshal invitation properties",
			slog.String("envelope_id", envelope.ID),
			slog.String("error", unmarshalErr.Error()))

		return emptyProps, false
	}

	if props.InvitationKind != mailbox.InvitationKindProfileJoin {
		return emptyProps, false
	}

	if props.CandidateID == "" || props.ProfileID == "" {
		logger.WarnContext(
			ctx,
			"Profile join invitation missing candidate_id or profile_id",
			slog.String("envelope_id", envelope.ID),
		)

		return emptyProps, false
	}

	return props, true
}
