CREATE TABLE friendships (
    id           SERIAL PRIMARY KEY,
    requester_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    addressee_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status       VARCHAR(10) NOT NULL DEFAULT 'pending'
                 CHECK (status IN ('pending', 'accepted', 'rejected')),
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),

    -- нельзя отправить дружбу дважды в одну сторону
    UNIQUE (requester_id, addressee_id),
    -- нельзя добавить самого себя
    CHECK (requester_id != addressee_id)
);

-- ищем по addressee_id (входящие заявки, принятые дружбы)
CREATE INDEX idx_friendships_addressee ON friendships(addressee_id, status);
-- ищем по requester_id (исходящие заявки, принятые дружбы)
CREATE INDEX idx_friendships_requester ON friendships(requester_id, status);
