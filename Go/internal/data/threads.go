package data

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Threads struct {
	ID        uuid.UUID `json:"id"`
	ForumID   uuid.UUID `json:"forumId"`
	Title     string    `json:"title"`
	AuthorID  uuid.UUID `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	IsLocked  bool      `json:"isLocked"`
}

type ThreadModel struct {
	DB      *pgxpool.Pool
	Timeout *time.Duration
}
