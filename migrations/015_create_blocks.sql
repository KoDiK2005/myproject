CREATE TABLE blocks (
    id         SERIAL PRIMARY KEY,
    blocker_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blocked_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE (blocker_id, blocked_id),
    CHECK (blocker_id != blocked_id)
);

-- проверка "X блокировал Y" и список тех кого блокирую
CREATE INDEX idx_blocks_blocker ON blocks(blocker_id);
-- проверка с другой стороны (меня кто-то заблокировал?)
CREATE INDEX idx_blocks_blocked ON blocks(blocked_id);
