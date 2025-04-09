CREATE TABLE IF NOT EXISTS forum.post_votes
(
    post_id UUID     NOT NULL,
    user_id UUID     NOT NULL,
    vote    SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT pk_post_votes PRIMARY KEY (post_id, user_id),
    CONSTRAINT fk_post_id FOREIGN KEY (post_id) REFERENCES forum.posts (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES forum.users (id) ON DELETE CASCADE,
    CONSTRAINT chk_vote CHECK (vote IN (-1, 0, 1))
);

CREATE TABLE IF NOT EXISTS forum.thread_votes
(
    thread_id UUID     NOT NULL,
    user_id UUID     NOT NULL,
    vote    SMALLINT NOT NULL DEFAULT 0,
    CONSTRAINT pk_thread_votes PRIMARY KEY (thread_id, user_id),
    CONSTRAINT fk_thread_id FOREIGN KEY (thread_id) REFERENCES forum.posts (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES forum.users (id) ON DELETE CASCADE,
    CONSTRAINT chk_vote CHECK (vote IN (-1, 0, 1))
);

ALTER TABLE forum.threads
ADD COLUMN likes BIGINT DEFAULT 0 NOT NULL;
