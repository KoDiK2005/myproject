CREATE TABLE notifications (
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    actor_id   INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type       VARCHAR(20) NOT NULL
               CHECK (type IN ('friend_request', 'friend_accept', 'like', 'comment')),
    post_id    INT REFERENCES posts(id) ON DELETE CASCADE,
    read       BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- лента уведомлений юзера, новые сверху; и отдельно — счётчик непрочитанных
CREATE INDEX idx_notifications_user ON notifications(user_id, created_at DESC);
CREATE INDEX idx_notifications_unread ON notifications(user_id, read);
