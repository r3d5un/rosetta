DROP TRIGGER IF EXISTS trigger_post_vote_sum_on_insert ON forum.post_votes;
DROP TRIGGER IF EXISTS trigger_post_vote_sum_on_update ON forum.post_votes;
DROP TRIGGER IF EXISTS trigger_post_vote_sum_on_delete ON forum.post_votes;

DROP TRIGGER IF EXISTS trigger_thread_vote_sum_on_insert ON forum.thread_votes;
DROP TRIGGER IF EXISTS trigger_thread_vote_sum_on_update ON forum.thread_votes;
DROP TRIGGER IF EXISTS trigger_thread_vote_sum_on_delete ON forum.thread_votes;
