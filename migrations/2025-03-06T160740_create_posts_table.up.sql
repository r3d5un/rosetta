CREATE TABLE IF NOT EXISTS forum.posts
(
    id         UUID      DEFAULT gen_random_uuid(),
    thread_id  UUID                    NOT NULL,
    reply_to   UUID                    NULL,
    author_id  UUID                    NOT NULL,
    content    TEXT                    NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    likes      INTEGER   DEFAULT 0     NOT NULL,
    CONSTRAINT pk_posts PRIMARY KEY (id),
    CONSTRAINT fk_thread_id FOREIGN KEY (thread_id) REFERENCES forum.threads (id),
    CONSTRAINT fk_author_id FOREIGN KEY (author_id) REFERENCES forum.users (id),
    CONSTRAINT srfk_replies_post_id FOREIGN KEY (reply_to) REFERENCES forum.posts (id)
);
