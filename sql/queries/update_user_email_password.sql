-- name: UpdateUserPasswordEmail :one
UPDATE users
SET hashed_password = $1, email = $2, updated_at = $3
RETURNING *;