-- +goose Up
ALTER TABLE users
ADD hashed_password TEXT NOT NULL DEFAULT 'unset';

-- +goose Down
DROP hashed_password FROM users;