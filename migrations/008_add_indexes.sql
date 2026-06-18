-- Ускоряем выборку постов по автору
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);

-- Ускоряем выборку комментариев к посту
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);

-- Ускоряем COUNT лайков у поста
CREATE INDEX IF NOT EXISTS idx_likes_post_id ON likes(post_id);
