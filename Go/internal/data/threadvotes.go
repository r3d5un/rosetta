package data

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
