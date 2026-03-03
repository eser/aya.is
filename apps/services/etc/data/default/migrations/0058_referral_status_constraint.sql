-- +goose Up
-- Add CHECK constraint to enforce valid referral status values.
ALTER TABLE "profile_membership_referral"
  ADD CONSTRAINT "profile_membership_referral_status_check"
    CHECK (status IN ('voting', 'frozen', 'reference_rejected', 'invitation_pending_response', 'invitation_accepted', 'invitation_rejected'));

-- +goose Down
ALTER TABLE "profile_membership_referral"
  DROP CONSTRAINT IF EXISTS "profile_membership_referral_status_check";
