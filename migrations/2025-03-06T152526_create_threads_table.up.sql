CREATE TABLE IF NOT EXISTS forum.threads
(
    id         UUID                    NOT NULL,
    forum_id   UUID                    NOT NULL,
    title      VARCHAR(128)            NOT NULL,
    author_id  UUID                    NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    is_locked  BOOL      DEFAULT FALSE NOT NULL,
    CONSTRAINT pk_threads PRIMARY KEY (id),
    CONSTRAINT foreign_key_forum_id FOREIGN KEY (forum_id) REFERENCES forum.forums (id),
    CONSTRAINT foreign_key_user_id FOREIGN KEY (author_id) REFERENCES forum.users (id)
);
