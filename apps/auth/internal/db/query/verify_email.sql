-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
  username,
  email,
  secret_code,
  token
) VALUES (
  $1, $2, $3, $4
) RETURNING *;
