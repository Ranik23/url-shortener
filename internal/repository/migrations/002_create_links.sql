-- +goose Up
CREATE TABLE links (
    ID SERIAL PRIMARY KEY,
    default_link TEXT UNIQUE NOT NULL,
    shortened_link TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE links;