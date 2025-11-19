-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    description TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,  -- Changed to TIMESTAMP
    feed_id UUID NOT NULL,
    FOREIGN KEY (feed_id)
        REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;
