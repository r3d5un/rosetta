DROP TRIGGER IF EXISTS trigger_post_vote_sum_on_insert ON forum.post_votes;
DROP TRIGGER IF EXISTS trigger_post_vote_sum_on_update ON forum.post_votes;
DROP TRIGGER IF EXISTS trigger_post_vote_sum_on_delete ON forum.post_votes;

CREATE TRIGGER trigger_post_vote_sum_update_on_insert
    AFTER INSERT
    ON forum.post_votes
    FOR EACH ROW
EXECUTE FUNCTION update_post_likes();

CREATE TRIGGER trigger_post_vote_sum_update_on_update
    AFTER UPDATE
    ON forum.post_votes
    FOR EACH ROW
EXECUTE FUNCTION update_post_likes();

CREATE TRIGGER trigger_post_vote_sum_update_on_delete
    AFTER DELETE
    ON forum.post_votes
    FOR EACH ROW
EXECUTE FUNCTION update_post_likes();

DROP TRIGGER IF EXISTS trigger_thread_vote_sum_on_insert ON forum.thread_votes;
DROP TRIGGER IF EXISTS trigger_thread_vote_sum_on_update ON forum.thread_votes;
DROP TRIGGER IF EXISTS trigger_thread_vote_sum_on_delete ON forum.thread_votes;

CREATE TRIGGER trigger_thread_vote_sum_update_on_insert
    AFTER INSERT
    ON forum.thread_votes
    FOR EACH ROW
EXECUTE FUNCTION update_thread_likes();

CREATE TRIGGER trigger_thread_vote_sum_update_on_update
    AFTER UPDATE
    ON forum.thread_votes
    FOR EACH ROW
EXECUTE FUNCTION update_thread_likes();

CREATE TRIGGER trigger_thread_vote_sum_update_on_delete
    AFTER DELETE
    ON forum.thread_votes
    FOR EACH ROW
EXECUTE FUNCTION update_thread_likes();
