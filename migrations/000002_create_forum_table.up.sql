CREATE TABLE IF NOT EXISTS forum.forums
(
    id          UUID      DEFAULT gen_random_uuid(),
    name        VARCHAR(256)            NOT NULL,
    description TEXT                    NULL,
    created_at  TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at  TIMESTAMP DEFAULT NOW() NOT NULL,
    CONSTRAINT pk_forums PRIMARY KEY (id)
);

