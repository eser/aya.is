package profiles

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/mailbox"
	profilesbiz "github.com/eser/aya.is/services/pkg/api/business/profiles"
)

// NewReferralAutoAccepter returns a callback that ensures a membership exists
// (creating or upgrading as needed) and merges teams when a profile_join
// invitation is accepted — completing the referral flow.
func NewReferralAutoAccepter(
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
		membershipID, err := profileService.EnsureMembershipFromReferralInternal(
			ctx,
			props.ProfileID,
			recipientProfileID,
			props.ReferralID,
		)
		if err != nil {
			logger.ErrorContext(ctx, "Failed to ensure membership from referral acceptance",
				slog.String("referral_id", props.ReferralID),
				slog.String("profile_id", props.ProfileID),
				slog.String("recipient_profile_id", recipientProfileID),
				slog.String("error", err.Error()))

			return
		}

		// Update referral status to accepted.
		statusErr := profileService.UpdateReferralStatusInternal(
			ctx,
			props.ReferralID,
			props.ProfileID,
			profilesbiz.ReferralStatusInvitationAccepted,
		)
		if statusErr != nil {
			logger.ErrorContext(ctx, "Failed to update referral status to accepted",
				slog.String("referral_id", props.ReferralID),
				slog.String("error", statusErr.Error()))

			return
		}

		logger.InfoContext(ctx, "Referral invitation accepted, membership ensured",
			slog.String("referral_id", props.ReferralID),
			slog.String("profile_id", props.ProfileID),
			slog.String("recipient_profile_id", recipientProfileID),
			slog.String("membership_id", membershipID))
	}
}

// NewReferralAutoRejecter returns a callback that updates referral status to invitation_rejected
// when a profile_join invitation is rejected.
func NewReferralAutoRejecter(
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

		err := profileService.UpdateReferralStatusInternal(
			ctx,
			props.ReferralID,
			props.ProfileID,
			profilesbiz.ReferralStatusInvitationRejected,
		)
		if err != nil {
			logger.ErrorContext(ctx, "Failed to update referral status to rejected",
				slog.String("referral_id", props.ReferralID),
				slog.String("error", err.Error()))

			return
		}

		logger.InfoContext(ctx, "Referral invitation rejected, status updated",
			slog.String("referral_id", props.ReferralID),
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
		ReferralID:       "",
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

	if props.ReferralID == "" || props.ProfileID == "" {
		logger.WarnContext(
			ctx,
			"Profile join invitation missing referral_id or profile_id",
			slog.String("envelope_id", envelope.ID),
		)

		return emptyProps, false
	}

	return props, true
}
