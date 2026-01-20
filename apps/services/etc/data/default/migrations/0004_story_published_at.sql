-- Add published_at column to story table
ALTER TABLE "story" ADD COLUMN IF NOT EXISTS "published_at" TIMESTAMP WITH TIME ZONE;
