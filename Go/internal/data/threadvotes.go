package data

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/rosetta/Go/internal/logging"
)

// ThreadVote represents a vote for a thread.
type ThreadVote struct {
	// ThreadID is the unique identifier of the thread that was voted on.
	ThreadID uuid.UUID `json:"threadId"`
	// UserID is the unique identifier of the user which voted.
	UserID uuid.UUID `json:"userId"`
	// Vote is the value of the vote.
	Vote int8 `json:"vote"`
}

type ThreadVoteModel struct {
	DB      *pgxpool.Pool
	Timeout *time.Duration
}

func (m *ThreadVoteModel) SelectSum(ctx context.Context, filters Filters) (*int, error) {
	const query string = `
SELECT CASE
           WHEN SUM(vote) IS NULL THEN 0
           ELSE SUM(vote)
           END AS total_votes
FROM forum.thread_votes
WHERE ($1::UUID IS NULL OR thread_id = $1::UUID)
  AND ($2::UUID IS NULL OR user_id = $2::UUID);
`
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(query)),
		slog.Any("filters", filters),
		slog.Duration("timeout", *m.Timeout),
	))

	var count int

	logger.Info("performing query")
	err := m.DB.QueryRow(ctx, query, filters.ThreadID, filters.UserID).Scan(&count)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("count complete", slog.Int("count", count))

	return &count, nil
}

// Vote performs a upsert for to record any votes for any thread. If the vote is 0, the record is
// deleted.
func (m *ThreadVoteModel) Vote(ctx context.Context, vote ThreadVote) (*ThreadVote, error) {
	const query string = `
WITH input_data AS (SELECT $1::UUID     AS thread_id,
                           $2::UUID     AS user_id,
                           $3::SMALLINT AS vote),
     delete_if_zero AS (
         DELETE FROM forum.thread_votes
             WHERE thread_id = (SELECT thread_id FROM input_data)
                 AND user_id = (SELECT user_id FROM input_data)
                 AND (SELECT vote FROM input_data) = 0)
INSERT
INTO forum.thread_votes (thread_id, user_id, vote)
SELECT thread_id, user_id, vote
FROM input_data
WHERE vote != 0
ON CONFLICT (thread_id, user_id) DO UPDATE
    SET vote = EXCLUDED.vote;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(query)),
		slog.Any("vote", vote),
		slog.Duration("timeout", *m.Timeout),
	))

	logger.Info("performing query")
	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()
	result, err := m.DB.Exec(ctx, query, vote.ThreadID, vote.UserID, vote.Vote)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("query performed", slog.Any("affectedRows", result.RowsAffected()))

	return &vote, nil
}
