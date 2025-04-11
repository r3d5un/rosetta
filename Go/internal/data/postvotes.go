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
SELECT COUNT(*)
FROM forum.post_votes
WHERE ($2::UUID IS NULL OR post_id = $2::UUID)
  AND ($3::UUID IS NULL OR user_id = $2::UUID);
`
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", query),
		slog.Any("filters", filters),
		slog.Duration("timeout", *m.Timeout),
	))

	var count int

	logger.Info("performing query")
	err := m.DB.QueryRow(ctx, query, filters.PostID, filters.UserID).Scan(&count)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("count complete", slog.Int("count", count))

	return &count, nil
}
