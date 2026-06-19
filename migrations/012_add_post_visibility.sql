ALTER TABLE posts
    ADD COLUMN visibility VARCHAR(10) NOT NULL DEFAULT 'public'
    CHECK (visibility IN ('public', 'friends'));

-- индекс нужен для фильтрации ленты по visibility
CREATE INDEX idx_posts_visibility ON posts(visibility);
