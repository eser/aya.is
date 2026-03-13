-- +goose Up

-- Rename referral tables to candidate
ALTER TABLE "profile_membership_referral" RENAME TO "profile_membership_candidate";
ALTER TABLE "profile_membership_referral_team" RENAME TO "profile_membership_candidate_team";
ALTER TABLE "profile_membership_referral_vote" RENAME TO "profile_membership_candidate_vote";

-- Rename foreign key columns
ALTER TABLE "profile_membership_candidate_team"
  RENAME COLUMN "profile_membership_referral_id" TO "candidate_id";

ALTER TABLE "profile_membership_candidate_vote"
  RENAME COLUMN "profile_membership_referral_id" TO "candidate_id";

-- Add source column (referral or application)
ALTER TABLE "profile_membership_candidate"
  ADD COLUMN "source" TEXT NOT NULL DEFAULT 'referral';

ALTER TABLE "profile_membership_candidate"
  ADD CONSTRAINT "profile_membership_candidate_source_check"
    CHECK (source IN ('referral', 'application'));

-- Add applicant_message for free-text motivation (used by applications)
ALTER TABLE "profile_membership_candidate"
  ADD COLUMN "applicant_message" TEXT;

-- Make referrer_membership_id nullable (NULL for self-applications)
ALTER TABLE "profile_membership_candidate"
  ALTER COLUMN "referrer_membership_id" DROP NOT NULL;

-- Update status CHECK constraint (drop old, add new with same values + application_accepted)
ALTER TABLE "profile_membership_candidate"
  DROP CONSTRAINT IF EXISTS "profile_membership_referral_status_check";

ALTER TABLE "profile_membership_candidate"
  ADD CONSTRAINT "profile_membership_candidate_status_check"
    CHECK (status IN (
      'voting', 'frozen', 'reference_rejected',
      'invitation_pending_response', 'invitation_accepted', 'invitation_rejected',
      'application_accepted'
    ));

-- Add feature toggles to profile
ALTER TABLE "profile" ADD COLUMN "feature_candidates" TEXT NOT NULL DEFAULT 'disabled';
ALTER TABLE "profile" ADD COLUMN "feature_applications" TEXT NOT NULL DEFAULT 'disabled';

-- Migrate existing orgs that had referrals: enable feature_candidates
UPDATE "profile"
SET "feature_candidates" = 'public'
WHERE "id" IN (
  SELECT DISTINCT "profile_id" FROM "profile_membership_candidate"
);

-- +goose Down

-- Remove feature toggles
ALTER TABLE "profile" DROP COLUMN IF EXISTS "feature_applications";
ALTER TABLE "profile" DROP COLUMN IF EXISTS "feature_candidates";

-- Restore status constraint
ALTER TABLE "profile_membership_candidate"
  DROP CONSTRAINT IF EXISTS "profile_membership_candidate_status_check";

ALTER TABLE "profile_membership_candidate"
  ADD CONSTRAINT "profile_membership_referral_status_check"
    CHECK (status IN (
      'voting', 'frozen', 'reference_rejected',
      'invitation_pending_response', 'invitation_accepted', 'invitation_rejected'
    ));

-- Make referrer_membership_id NOT NULL again
ALTER TABLE "profile_membership_candidate"
  ALTER COLUMN "referrer_membership_id" SET NOT NULL;

-- Remove new columns
ALTER TABLE "profile_membership_candidate" DROP COLUMN IF EXISTS "applicant_message";

ALTER TABLE "profile_membership_candidate"
  DROP CONSTRAINT IF EXISTS "profile_membership_candidate_source_check";

ALTER TABLE "profile_membership_candidate" DROP COLUMN IF EXISTS "source";

-- Rename columns back
ALTER TABLE "profile_membership_candidate_vote"
  RENAME COLUMN "candidate_id" TO "profile_membership_referral_id";

ALTER TABLE "profile_membership_candidate_team"
  RENAME COLUMN "candidate_id" TO "profile_membership_referral_id";

-- Rename tables back
ALTER TABLE "profile_membership_candidate_vote" RENAME TO "profile_membership_referral_vote";
ALTER TABLE "profile_membership_candidate_team" RENAME TO "profile_membership_referral_team";
ALTER TABLE "profile_membership_candidate" RENAME TO "profile_membership_referral";
