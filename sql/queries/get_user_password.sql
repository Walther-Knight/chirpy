-- name: GetUserPassword :one
SELECT *
FROM users
WHERE email = $1;