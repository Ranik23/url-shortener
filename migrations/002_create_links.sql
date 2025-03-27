-- +goose Up
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    default_link TEXT UNIQUE NOT NULL,
    shortened_link TEXT UNIQUE NOT NULL,
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE links;
