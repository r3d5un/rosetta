ALTER TABLE forum.threads
    ALTER COLUMN id SET DEFAULT gen_random_uuid();