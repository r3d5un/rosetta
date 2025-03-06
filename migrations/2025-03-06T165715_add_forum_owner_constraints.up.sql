ALTER TABLE forum.forums
ADD CONSTRAINT fk_owner_id FOREIGN KEY (owner_id) REFERENCES forum.users (id);
