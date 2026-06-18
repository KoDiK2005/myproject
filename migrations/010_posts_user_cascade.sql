-- удаление пользователя падало с foreign key violation, т.к. posts.user_id
-- не каскадировался (в отличие от comments и likes)
ALTER TABLE posts DROP CONSTRAINT posts_user_id_fkey;
ALTER TABLE posts ADD CONSTRAINT posts_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
