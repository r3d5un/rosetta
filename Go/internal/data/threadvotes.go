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
	// ThreadID is the unique identifier of the post that was voted on.
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

func (m *ThreadVoteModel) SelectCount(ctx context.Context, filters Filters) (*int, error) {
	const query string = `
SELECT COUNT(*)
FROM forum.thread_votes
WHERE ($2::UUID IS NULL OR thread_id = $2::UUID)
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
	err := m.DB.QueryRow(ctx, query, filters.ThreadID, filters.UserID).Scan(&count)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("count complete", slog.Int("count", count))

	return &count, nil
}
