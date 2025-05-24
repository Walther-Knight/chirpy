-- name: GetUserFromID :one
SELECT id, email, is_chirpy_red
FROM users
WHERE id = $1;