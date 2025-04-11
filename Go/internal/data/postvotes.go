package data

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/rosetta/Go/internal/logging"
)

// PostVote represents a vote for a post.
type PostVote struct {
	// PostID is the unique identifier of the post that was voted on.
	PostID uuid.UUID `json:"postId"`
	// UserID is the unique identifier of the user which voted.
	UserID uuid.UUID `json:"userId"`
	// Vote is the value of the vote.
	Vote int8 `json:"vote"`
}

type PostVoteModel struct {
	DB      *pgxpool.Pool
	Timeout *time.Duration
}

func (m *PostVoteModel) SelectCount(ctx context.Context, filters Filters) (*int, error) {
	const query string = `
SELECT CASE
           WHEN SUM(vote) IS NULL THEN 0
           ELSE SUM(vote)
           END AS total_votes
FROM forum.post_votes
WHERE ($1::UUID IS NULL OR post_id = $1::UUID)
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
	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()
	err := m.DB.QueryRow(ctx, query, filters.PostID, filters.UserID).Scan(&count)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("count complete", slog.Int("count", count))

	return &count, nil
}

// Vote performs a upsert for to record any votes for any post. If the vote is 0, the record is
// deleted.
func (m *PostVoteModel) Vote(ctx context.Context, vote PostVote) (*PostVote, error) {
	const query string = `
WITH input_data AS (SELECT $1::UUID     AS post_id,
                           $2::UUID     AS user_id,
                           $3::SMALLINT AS vote),
     delete_if_zero AS (
         DELETE FROM forum.post_votes
             WHERE post_id = (SELECT post_id FROM input_data)
                 AND user_id = (SELECT user_id FROM input_data)
                 AND (SELECT vote FROM input_data) = 0)
INSERT
INTO forum.post_votes (post_id, user_id, vote)
SELECT post_id, user_id, vote
FROM input_data
WHERE vote != 0
ON CONFLICT (post_id, user_id) DO UPDATE
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
	result, err := m.DB.Exec(ctx, query, vote.PostID, vote.UserID, vote.Vote)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("query performed", slog.Any("affectedRows", result.RowsAffected()))

	return &vote, nil
}
