-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL
);

-- +goode Down
DROP TABLE users;