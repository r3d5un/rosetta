package data

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	// ID is the unique identifier of the thread
	//
	// Upon creating a new thread, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	ID uuid.UUID `json:"id"`
	// ThreadID is the ID of the parent thread.
	ThreadID uuid.UUID `json:"threadId"`
	// ReplyTo is the ID of which this post is a reply to.
	ReplyTo *uuid.UUID `json:"replyTo"`
	// AuthorID is the unique identifier of the author of the thread.
	AuthorID uuid.UUID `json:"authorId"`
	// Context is the actual text content of a post
	Content string `json:"content"`
	// CreatedAt denotes when a thread was created.
	//
	// Upon creating a new thread, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt denotes when a thread was last updated.
	//
	// Upon creating a new thread, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	UpdatedAt time.Time `json:"updatedAt"`
	// Likes is the total number of votes a post has received
	Likes int64 `json:"likes"`
	// Deleted is a soft delete flag for a thread.
	Deleted bool `json:"deleted,omitzero"`
	// DeletedAt denotes when a thread was marked as deleted.
	//
	// This field is ignored when updating or creating new thread.
	DeletedAt sql.NullTime `json:"deletedAt,omitzero"`
}

type PostModel struct {
	DB      *pgxpool.Pool
	Timeout *time.Duration
}
