-- +goose Up

-- Rename tables
ALTER TABLE "event_series" RENAME TO "activity_series";
ALTER TABLE "event" RENAME TO "activity";
ALTER TABLE "event_attendance" RENAME TO "activity_attendance";

-- Rename columns
ALTER TABLE "activity_series" RENAME COLUMN "event_picture_uri" TO "activity_picture_uri";
ALTER TABLE "activity" RENAME COLUMN "event_picture_uri" TO "activity_picture_uri";
ALTER TABLE "activity_attendance" RENAME COLUMN "event_id" TO "activity_id";

-- Rename constraints
ALTER TABLE "activity_series" RENAME CONSTRAINT "event_series_slug_unique" TO "activity_series_slug_unique";
ALTER TABLE "activity" RENAME CONSTRAINT "event_slug_unique" TO "activity_slug_unique";
ALTER TABLE "activity" RENAME CONSTRAINT "event_series_id_fk" TO "activity_series_id_fk";
ALTER TABLE "activity_attendance" RENAME CONSTRAINT "event_attendance_event_id_fk" TO "activity_attendance_activity_id_fk";
ALTER TABLE "activity_attendance" RENAME CONSTRAINT "event_attendance_profile_id_fk" TO "activity_attendance_profile_id_fk";
ALTER TABLE "activity_attendance" RENAME CONSTRAINT "event_attendance_event_id_profile_id_unique" TO "activity_attendance_activity_id_profile_id_unique";

-- +goose Down

-- Reverse constraint renames
ALTER TABLE "activity_attendance" RENAME CONSTRAINT "activity_attendance_activity_id_profile_id_unique" TO "event_attendance_event_id_profile_id_unique";
ALTER TABLE "activity_attendance" RENAME CONSTRAINT "activity_attendance_profile_id_fk" TO "event_attendance_profile_id_fk";
ALTER TABLE "activity_attendance" RENAME CONSTRAINT "activity_attendance_activity_id_fk" TO "event_attendance_event_id_fk";
ALTER TABLE "activity" RENAME CONSTRAINT "activity_series_id_fk" TO "event_series_id_fk";
ALTER TABLE "activity" RENAME CONSTRAINT "activity_slug_unique" TO "event_slug_unique";
ALTER TABLE "activity_series" RENAME CONSTRAINT "activity_series_slug_unique" TO "event_series_slug_unique";

-- Reverse column renames
ALTER TABLE "activity_attendance" RENAME COLUMN "activity_id" TO "event_id";
ALTER TABLE "activity" RENAME COLUMN "activity_picture_uri" TO "event_picture_uri";
ALTER TABLE "activity_series" RENAME COLUMN "activity_picture_uri" TO "event_picture_uri";

-- Reverse table renames
ALTER TABLE "activity_attendance" RENAME TO "event_attendance";
ALTER TABLE "activity" RENAME TO "event";
ALTER TABLE "activity_series" RENAME TO "event_series";
