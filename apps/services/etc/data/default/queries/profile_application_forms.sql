-- name: GetActiveApplicationForm :one
SELECT * FROM "profile_application_form"
WHERE profile_id = sqlc.arg(profile_id)
  AND is_active = TRUE;

-- name: CreateApplicationForm :one
INSERT INTO "profile_application_form" (
  id, profile_id, preset_key, is_active, responses_visibility, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(profile_id),
  sqlc.narg(preset_key),
  TRUE,
  sqlc.arg(responses_visibility),
  NOW()
) RETURNING *;

-- name: UpdateApplicationForm :exec
UPDATE "profile_application_form"
SET preset_key = sqlc.narg(preset_key),
    responses_visibility = sqlc.arg(responses_visibility),
    updated_at = NOW()
WHERE id = sqlc.arg(id);

-- name: DeactivateApplicationForms :exec
UPDATE "profile_application_form"
SET is_active = FALSE,
    updated_at = NOW()
WHERE profile_id = sqlc.arg(profile_id)
  AND is_active = TRUE;

-- name: ListApplicationFormFields :many
SELECT * FROM "profile_application_form_field"
WHERE form_id = sqlc.arg(form_id)
ORDER BY sort_order ASC;

-- name: CreateApplicationFormField :one
INSERT INTO "profile_application_form_field" (
  id, form_id, label, field_type, is_required, sort_order, placeholder, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(form_id),
  sqlc.arg(label),
  sqlc.arg(field_type),
  sqlc.arg(is_required),
  sqlc.arg(sort_order),
  sqlc.narg(placeholder),
  NOW()
) RETURNING *;

-- name: DeleteApplicationFormFields :exec
DELETE FROM "profile_application_form_field"
WHERE form_id = sqlc.arg(form_id);

-- name: CreateCandidateResponse :one
INSERT INTO "profile_candidate_response" (
  id, candidate_id, form_field_id, value, created_at
) VALUES (
  sqlc.arg(id),
  sqlc.arg(candidate_id),
  sqlc.arg(form_field_id),
  sqlc.arg(value),
  NOW()
) RETURNING *;

-- name: ListCandidateResponses :many
SELECT
  cr.*,
  ff.label AS field_label,
  ff.field_type,
  ff.is_required,
  ff.sort_order
FROM "profile_candidate_response" cr
  INNER JOIN "profile_application_form_field" ff ON ff.id = cr.form_field_id
WHERE cr.candidate_id = sqlc.arg(candidate_id)
ORDER BY ff.sort_order ASC;

-- name: GetApplicationFormByProfileID :one
SELECT f.*, p.feature_applications
FROM "profile_application_form" f
  INNER JOIN "profile" p ON p.id = f.profile_id
WHERE f.profile_id = sqlc.arg(profile_id)
  AND f.is_active = TRUE;
