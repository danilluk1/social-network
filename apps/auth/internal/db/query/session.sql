-- name: CreateSession :one
INSERT INTO sessions (
  id,
  username,
  refresh_token,
  user_agent,
  client_ip,
  is_blocked,
  expires_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: UpdateSession :one
UPDATE sessions
SET
  refresh_token = COALESCE(sqlc.narg(refresh_token), refresh_token),
  user_agent = COALESCE(sqlc.narg(user_agent), user_agent),
  client_ip = COALESCE(sqlc.narg(client_ip), client_ip),
  is_blocked = COALESCE(sqlc.narg(is_blocked), is_blocked),
  expires_at = COALESCE(sqlc.narg(expires_at), expires_at)
WHERE
  id = sqlc.arg(id)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;
