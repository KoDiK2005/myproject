CREATE TABLE IF NOT EXISTS comments (
    id         SERIAL PRIMARY KEY,
    post_id    INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body       TEXT    NOT NULL CHECK (char_length(body) >= 1 AND char_length(body) <= 1000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
