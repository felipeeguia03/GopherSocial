-- posts: al borrar user, borrar sus posts
ALTER TABLE posts DROP CONSTRAINT fk_user;
ALTER TABLE posts ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- comments: al borrar post, borrar sus comentarios
ALTER TABLE comments DROP CONSTRAINT fk_post;
ALTER TABLE comments ADD CONSTRAINT fk_post FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE;

-- comments: al borrar user, borrar sus comentarios
ALTER TABLE comments DROP CONSTRAINT fk_user;
ALTER TABLE comments ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- user_invitations: al borrar user, borrar sus invitaciones
ALTER TABLE user_invitations DROP CONSTRAINT IF EXISTS fk_user;
ALTER TABLE user_invitations ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
