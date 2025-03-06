CREATE TABLE IF NOT EXISTS forum.users
(
    id         UUID      DEFAULT gen_random_uuid(),
    name       VARCHAR(256)            NOT NULL,
    username   VARCHAR(256)            NOT NULL,
    email      VARCHAR(256)            NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    CONSTRAINT pk_users PRIMARY KEY (id),
    CONSTRAINT email_unique_constraint UNIQUE (email),
    CONSTRAINT username_unique_constraint UNIQUE (username)
);
