// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: create_user_password.sql

package database

import (
	"context"
)

const createUserPassword = `-- name: CreateUserPassword :exec
INSERT INTO users(hashed_password)
VALUES (
    $1
)
`

func (q *Queries) CreateUserPassword(ctx context.Context, hashedPassword string) error {
	_, err := q.db.ExecContext(ctx, createUserPassword, hashedPassword)
	return err
}
