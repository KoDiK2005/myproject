CREATE TABLE IF NOT EXISTS messages (
    id          SERIAL PRIMARY KEY,
    sender_id   INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content     TEXT NOT NULL CHECK (length(content) > 0 AND length(content) <= 5000),
    read_at     TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    CHECK (sender_id != receiver_id)
);

-- по паре юзеров — чтобы грузить историю конкретного чата
CREATE INDEX idx_messages_pair ON messages(
    LEAST(sender_id, receiver_id),
    GREATEST(sender_id, receiver_id),
    created_at DESC
);

-- непрочитанные сообщения конкретного юзера
CREATE INDEX idx_messages_receiver_unread ON messages(receiver_id, read_at) WHERE read_at IS NULL;
