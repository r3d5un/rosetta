CREATE OR REPLACE FUNCTION update_post_likes()
    RETURNS TRIGGER AS
$$
BEGIN

    UPDATE forum.posts
    SET likes = (SELECT COALESCE(SUM(vote), 0)
                 FROM forum.post_votes
                 WHERE post_id = NEW.post_id)
    WHERE id = NEW.post_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_thread_likes()
    RETURNS TRIGGER AS
$$
BEGIN

    UPDATE forum.threads
    SET likes = (SELECT COALESCE(SUM(vote), 0)
                 FROM forum.thread_votes
                 WHERE thread_id = NEW.thread_id)
    WHERE id = NEW.thread_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
