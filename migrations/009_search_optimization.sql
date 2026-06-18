-- pg_trgm позволяет индексировать ILIKE запросы
-- без него ILIKE '%query%' = full table scan
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_posts_title_trgm ON posts USING GIN(title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_posts_body_trgm  ON posts USING GIN(body  gin_trgm_ops);
