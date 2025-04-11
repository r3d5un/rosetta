package data

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
