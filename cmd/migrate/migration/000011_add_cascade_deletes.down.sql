-- revert posts
ALTER TABLE posts DROP CONSTRAINT fk_user;
ALTER TABLE posts ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id);

-- revert comments fk_post
ALTER TABLE comments DROP CONSTRAINT fk_post;
ALTER TABLE comments ADD CONSTRAINT fk_post FOREIGN KEY (post_id) REFERENCES posts(id);

-- revert comments fk_user
ALTER TABLE comments DROP CONSTRAINT fk_user;
ALTER TABLE comments ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id);

-- revert user_invitations
ALTER TABLE user_invitations DROP CONSTRAINT IF EXISTS fk_user;
