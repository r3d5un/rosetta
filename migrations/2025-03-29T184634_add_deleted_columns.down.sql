ALTER TABLE forum.users
    DROP COLUMN IF EXISTS deleted,
    DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE forum.forums
    DROP COLUMN IF EXISTS deleted,
    DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE forum.threads
    DROP COLUMN IF EXISTS deleted,
    DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE forum.posts
    DROP COLUMN IF EXISTS deleted,
    DROP COLUMN IF EXISTS deleted_at;
