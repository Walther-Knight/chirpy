-- +goose Up
ALTER TABLE users
ADD is_chirpy_red BOOLEAN DEFAULT false;

-- +goose Down
DROP is_chirpy_red FROM users;