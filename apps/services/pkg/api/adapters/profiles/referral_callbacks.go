package profiles

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/mailbox"
	profilesbiz "github.com/eser/aya.is/services/pkg/api/business/profiles"
)

// NewReferralAutoAccepter returns a callback that automatically creates a membership
// when a profile_join invitation is accepted — completing the referral flow.
func NewReferralAutoAccepter(
	profileService *profilesbiz.Service,
	logger *logfx.Logger,
) mailbox.OnEnvelopeAcceptedFunc {
	return func(ctx context.Context, envelope *mailbox.Envelope) {
		if envelope.Kind != mailbox.KindInvitation {
			return
		}

		props, ok := parseProfileJoinProps(envelope, logger)
		if !ok {
			return
		}

		// Create membership for the accepted profile.
		recipientProfileID := envelope.TargetProfileID

		err := profileService.CreateProfileMembership(
			ctx,
			"", // system-initiated, no user ID
			props.ProfileID,
			&recipientProfileID,
			string(profilesbiz.MembershipKindMember),
		)
		if err != nil {
			logger.ErrorContext(ctx, "Failed to create membership from referral acceptance",
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

		logger.InfoContext(ctx, "Referral invitation accepted, membership created",
			slog.String("referral_id", props.ReferralID),
			slog.String("profile_id", props.ProfileID),
			slog.String("recipient_profile_id", recipientProfileID))
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

		props, ok := parseProfileJoinProps(envelope, logger)
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
	envelope *mailbox.Envelope,
	logger *logfx.Logger,
) (mailbox.InvitationProperties, bool) {
	propsJSON, marshalErr := json.Marshal(envelope.Properties)
	if marshalErr != nil {
		logger.WarnContext(context.Background(), "Failed to marshal envelope properties",
			slog.String("envelope_id", envelope.ID),
			slog.String("error", marshalErr.Error()))

		return mailbox.InvitationProperties{}, false
	}

	var props mailbox.InvitationProperties

	unmarshalErr := json.Unmarshal(propsJSON, &props)
	if unmarshalErr != nil {
		logger.WarnContext(context.Background(), "Failed to unmarshal invitation properties",
			slog.String("envelope_id", envelope.ID),
			slog.String("error", unmarshalErr.Error()))

		return mailbox.InvitationProperties{}, false
	}

	if props.InvitationKind != mailbox.InvitationKindProfileJoin {
		return mailbox.InvitationProperties{}, false
	}

	if props.ReferralID == "" || props.ProfileID == "" {
		logger.WarnContext(
			context.Background(),
			"Profile join invitation missing referral_id or profile_id",
			slog.String("envelope_id", envelope.ID),
		)

		return mailbox.InvitationProperties{}, false
	}

	return props, true
}
