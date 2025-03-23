package data

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID        uuid.UUID  `json:"id"`
	ThreadID  uuid.UUID  `json:"threadId"`
	ReplyTo   *uuid.UUID `json:"replyTo"`
	AuthorID  uuid.UUID  `json:"authorId"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	Likes     int64      `json:"likes"`
}

type PostModel struct {
	DB      *pgxpool.Pool
	Timeout *time.Duration
}
