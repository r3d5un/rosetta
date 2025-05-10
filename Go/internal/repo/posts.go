package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/r3d5un/rosetta/Go/internal/database"
)

type Post struct {
	// ID is the unique identifier of the post
	//
	// Upon creating a new post, any existing values in this field is ignored.
	ID uuid.UUID `json:"id"`
	// ThreadID is the ID of the parent thread.
	ThreadID uuid.UUID `json:"threadId"`
	// ReplyTo is the ID of which this post is a reply to.
	ReplyTo uuid.NullUUID `json:"replyTo"`
	// AuthorID is the unique identifier of the author of the post.
	AuthorID uuid.UUID `json:"authorId"`
	// Content is the actual text content of a post
	Content string `json:"content"`
	// CreatedAt denotes when a post was created.
	//
	// Upon creating a new post, any existing values in this field is ignored.
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt denotes when a post was last updated.
	//
	// Upon creating a new post, any existing values in this field is ignored.
	UpdatedAt time.Time `json:"updatedAt"`
	// Likes is the total number of votes a post has received
	Likes int64 `json:"likes"`
	// Deleted is a soft delete flag for a post.
	Deleted bool `json:"deleted,omitzero"`
	// DeletedAt denotes when a post was marked as deleted.
	//
	// This field is ignored when updating or creating new post.
	DeletedAt *time.Time `json:"deletedAt,omitzero"`
}

func newPostFromRow(row data.Post) *Post {
	return &Post{
		ID:        row.ID,
		ThreadID:  row.ThreadID,
		ReplyTo:   row.ReplyTo,
		AuthorID:  row.AuthorID,
		Content:   row.Content,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		Likes:     row.Likes,
		Deleted:   row.Deleted,
		DeletedAt: database.NullTimeToPtr(row.DeletedAt),
	}
}

func (p *Post) Row() data.Post {
	return data.Post{
		ID:        p.ID,
		ThreadID:  p.ThreadID,
		ReplyTo:   p.ReplyTo,
		AuthorID:  p.AuthorID,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Likes:     p.Likes,
		Deleted:   p.Deleted,
		DeletedAt: database.NewNullTime(p.DeletedAt),
	}
}

type PostPatch struct {
	// ID is the unique identifier of the post
	ID uuid.UUID `json:"id"`
	// ThreadID is the ID of the parent thread.
	ThreadID uuid.UUID `json:"threadId"`
	// Content is the actual text content of a post
	Content *string `json:"content"`
}

func (p *PostPatch) Row() data.PostPatch {
	return data.PostPatch{
		ID:       p.ID,
		ThreadID: p.ThreadID,
		Content:  database.NewNullString(p.Content),
	}
}

type PostReader interface {
	Read(context.Context, uuid.UUID, bool) (*Post, error)
	List(context.Context, data.Filters, bool) ([]*Post, *data.Metadata, error)
}

type PostWriter interface {
	Create(context.Context, Post) (*Post, error)
	Delete(context.Context, uuid.UUID) (*Post, error)
	Restore(context.Context, uuid.UUID) (*Post, error)
	PermanentlyDelete(context.Context, uuid.UUID) (*Post, error)
}

type PostRepository struct {
	models       *data.Models
	threadReader ThreadReader
	userReader   UserReader
}

func NewPostRepository(
	models *data.Models,
	threadReader ThreadReader,
	userReader UserReader,
) PostRepository {
	return PostRepository{
		models:       models,
		threadReader: threadReader,
		userReader:   userReader,
	}
}
