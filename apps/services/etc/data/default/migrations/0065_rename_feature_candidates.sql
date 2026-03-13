-- +goose Up
ALTER TABLE "profile" RENAME COLUMN "feature_candidates" TO "feature_referrals";

-- +goose Down
ALTER TABLE "profile" RENAME COLUMN "feature_referrals" TO "feature_candidates";
