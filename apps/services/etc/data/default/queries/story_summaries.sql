-- name: GetUnsummarizedPublishedStories :many
-- Returns published story translations that have no AI summary yet.
SELECT
  st.story_id,
  RTRIM(st.locale_code) AS locale_code,
  st.title,
  st.summary,
  st.content
FROM "story_tx" st
WHERE (st.summary_ai IS NULL OR st.summary_ai = '')
  AND EXISTS (
    SELECT 1 FROM "story_publication" sp
    WHERE sp.story_id = st.story_id
      AND sp.published_at IS NOT NULL
      AND sp.deleted_at IS NULL
  )
  AND EXISTS (
    SELECT 1 FROM "story" s
    WHERE s.id = st.story_id
      AND s.deleted_at IS NULL
  )
ORDER BY st.story_id, st.locale_code
LIMIT sqlc.arg(max_items);

-- name: UpsertStorySummaryAI :exec
-- Updates the AI-generated summary for a specific story translation.
UPDATE "story_tx"
SET summary_ai = sqlc.arg(summary_ai)
WHERE story_id = sqlc.arg(story_id)
  AND locale_code = sqlc.arg(locale_code);
