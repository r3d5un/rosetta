package data

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/rosetta/Go/internal/logging"
)

type Post struct {
	// ID is the unique identifier of the post
	//
	// Upon creating a new post, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	ID uuid.UUID `json:"id"`
	// ThreadID is the ID of the parent thread.
	ThreadID uuid.UUID `json:"threadId"`
	// ReplyTo is the ID of which this post is a reply to.
	ReplyTo *uuid.UUID `json:"replyTo"`
	// AuthorID is the unique identifier of the author of the post.
	AuthorID uuid.UUID `json:"authorId"`
	// Context is the actual text content of a post
	Content string `json:"content"`
	// CreatedAt denotes when a post was created.
	//
	// Upon creating a new post, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt denotes when a post was last updated.
	//
	// Upon creating a new post, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	UpdatedAt time.Time `json:"updatedAt"`
	// Likes is the total number of votes a post has received
	Likes int64 `json:"likes"`
	// Deleted is a soft delete flag for a post.
	Deleted bool `json:"deleted,omitzero"`
	// DeletedAt denotes when a post was marked as deleted.
	//
	// This field is ignored when updating or creating new post.
	DeletedAt sql.NullTime `json:"deletedAt,omitzero"`
}

type PostModel struct {
	DB      *pgxpool.Pool
	Timeout *time.Duration
}

func (m *PostModel) Select(ctx context.Context, id uuid.UUID) (*Post, error) {
	const query string = `
SELECT id,
       thread_id,
       reply_to,
       author_id,
       content,
       created_at,
       updated_at,
       likes,
       deleted,
       deleted_at
FROM forum.posts
WHERE id = $1::UUID;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", query),
		slog.Any("id", id),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.Info("performing query")
	var p Post
	err := m.DB.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&p.ID,
		&p.ThreadID,
		&p.ReplyTo,
		&p.AuthorID,
		&p.Content,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.Likes,
		&p.Deleted,
		&p.DeletedAt,
	)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("post selected", slog.Any("post", p))

	return &p, nil
}

func (m *PostModel) SelectAll(ctx context.Context, filters Filters) ([]*Post, *Metadata, error) {
	query := `
SELECT id,
       thread_id,
       reply_to,
       author_id,
       content,
       created_at,
       updated_at,
       likes,
       deleted,
       deleted_at
FROM forum.posts
WHERE ($2::UUID IS NULL OR id = $2::UUID)
  AND ($3::UUID IS NULL OR thread_id = $3::UUID)
  AND ($4::VARCHAR(256) IS NULL OR title = $4::VARCHAR(256))
  AND ($5::UUID IS NULL OR author_id = $5::UUID)
  AND ($6::TIMESTAMP IS NULL or created_at >= $6::TIMESTAMP)
  AND ($7::TIMESTAMP IS NULL or created_at <= $7::TIMESTAMP)
  AND ($8::TIMESTAMP IS NULL or updated_at >= $8::TIMESTAMP)
  AND ($9::TIMESTAMP IS NULL or updated_at <= $9::TIMESTAMP)
  AND ($10::BOOLEAN IS NULL or deleted = $10::BOOLEAN)
  AND ($11::TIMESTAMP IS NULL or deleted_at >= $11::TIMESTAMP)
  AND ($12::TIMESTAMP IS NULL or deleted_at <= $12::TIMESTAMP)
  AND id > $14::UUID
` + CreateOrderByClause(filters.OrderBy) + `
LIMIT $1::INTEGER;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", query),
		slog.Any("filters", filters),
		slog.Duration("timeout", *m.Timeout),
	))

	logger.Info("performing query")
	rows, err := m.DB.Query(
		ctx,
		query,
		filters.PageSize,
		filters.ID,
		filters.ThreadID,
		filters.Title,
		filters.AuthorID,
		filters.CreatedAtFrom,
		filters.CreatedAtTo,
		filters.UpdatedAtFrom,
		filters.UpdatedAtTo,
		filters.Deleted,
		filters.DeletedAtFrom,
		filters.DeletedAtTo,
		filters.LastSeen,
	)
	if err != nil {
		logger.Error("unable to perform query", slog.String("error", err.Error()))
		return nil, nil, err
	}

	posts := []*Post{}

	for rows.Next() {
		var p Post

		err := rows.Scan(
			&p.ID,
			&p.ThreadID,
			&p.ReplyTo,
			&p.AuthorID,
			&p.Content,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.Likes,
			&p.Deleted,
			&p.DeletedAt,
		)
		if err != nil {
			return nil, nil, handleError(err, logger)
		}
		posts = append(posts, &p)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, handleError(err, logger)
	}
	length := len(posts)
	var metadata Metadata
	if length > 0 {
		metadata.LastSeen = posts[length-1].ID
		metadata.Next = true
	}
	metadata.ResponseLength = length

	logger.Info("posts selected", slog.Any("metadata", metadata))
	return posts, &metadata, nil
}
