-- +goose Up
ALTER TABLE "event_queue" RENAME TO "queue";
ALTER INDEX "event_queue_poll_idx" RENAME TO "queue_poll_idx";
ALTER INDEX "event_queue_type_idx" RENAME TO "queue_type_idx";

-- +goose Down
ALTER TABLE "queue" RENAME TO "event_queue";
ALTER INDEX "queue_poll_idx" RENAME TO "event_queue_poll_idx";
ALTER INDEX "queue_type_idx" RENAME TO "event_queue_type_idx";
