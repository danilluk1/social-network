-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
  username,
  email,
  secret_code,
  token
) VALUES (
  $1, $2, $3, $4
) RETURNING *;


-- name: VerifyEmail :one
UPDATE verify_emails
SET is_used = TRUE
WHERE
  (secret_code = $1 OR token = $2)
  AND NOW() < expires_at
  AND is_used = false
RETURNING *;
