-- name: ListProfileTeams :many
SELECT * FROM "profile_team"
WHERE profile_id = sqlc.arg(profile_id) AND deleted_at IS NULL
ORDER BY name ASC;

-- name: GetProfileTeamByID :one
SELECT * FROM "profile_team"
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;

-- name: CreateProfileTeam :one
INSERT INTO "profile_team" (id, profile_id, name, description)
VALUES (sqlc.arg(id), sqlc.arg(profile_id), sqlc.arg(name), sqlc.narg(description))
RETURNING *;

-- name: UpdateProfileTeam :execrows
UPDATE "profile_team"
SET name = sqlc.arg(name), description = sqlc.narg(description)
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;

-- name: DeleteProfileTeam :execrows
UPDATE "profile_team"
SET deleted_at = NOW()
WHERE id = sqlc.arg(id) AND deleted_at IS NULL;

-- name: CountProfileTeamMembers :one
SELECT COUNT(*) FROM "profile_membership_team"
WHERE profile_team_id = sqlc.arg(profile_team_id) AND deleted_at IS NULL;

-- name: ListMembershipTeams :many
SELECT pt.* FROM "profile_team" pt
JOIN "profile_membership_team" pmt ON pmt.profile_team_id = pt.id AND pmt.deleted_at IS NULL
WHERE pmt.profile_membership_id = sqlc.arg(profile_membership_id) AND pt.deleted_at IS NULL
ORDER BY pt.name ASC;

-- name: SetMembershipTeams_Delete :execrows
UPDATE "profile_membership_team"
SET deleted_at = NOW()
WHERE profile_membership_id = sqlc.arg(profile_membership_id) AND deleted_at IS NULL;

-- name: SetMembershipTeams_Insert :one
INSERT INTO "profile_membership_team" (id, profile_membership_id, profile_team_id)
VALUES (sqlc.arg(id), sqlc.arg(profile_membership_id), sqlc.arg(profile_team_id))
RETURNING *;

-- name: ListProfileTeamsWithMemberCount :many
SELECT pt.*, COUNT(pmt.id) AS member_count
FROM "profile_team" pt
LEFT JOIN "profile_membership_team" pmt ON pmt.profile_team_id = pt.id AND pmt.deleted_at IS NULL
WHERE pt.profile_id = sqlc.arg(profile_id) AND pt.deleted_at IS NULL
GROUP BY pt.id
ORDER BY pt.name ASC;
