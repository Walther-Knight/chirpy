-- name: CreateUserPassword :exec
INSERT INTO users(hashed_password)
VALUES (
    $1
);