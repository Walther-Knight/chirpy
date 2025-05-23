-- name: GetUserFromRefreshToken :one
SELECT expires_at, user_id, revoked_at
FROM refresh_tokens
WHERE token = $1;