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

type Thread struct {
	// ID is the unique identifier of the thread
	//
	// Upon creating a new forum, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	ID uuid.UUID `json:"id"`
	// ForumID is the parent forum this thread belongs to.
	ForumID uuid.UUID `json:"forumId"`
	// Title is the subject the thread is about.
	Title string `json:"title"`
	// AuthorID is the unique identifier of the author of the thread.
	AuthorID uuid.UUID `json:"authorId"`
	// CreatedAt denotes when a forum was created.
	//
	// Upon creating a new forum, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt denotes when a forum was last updated.
	//
	// Upon creating a new forum, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	UpdatedAt time.Time `json:"updatedAt"`
	// IsLocked denotes whether a thread had been locked for changes.
	//
	// This field is ignored when updating or creating new threads.
	IsLocked bool `json:"isLocked"`
	// Deleted is a soft delete flag for a forum.
	Deleted bool `json:"deleted,omitzero"`
	// DeletedAt denotes when a forum was marked as deleted.
	//
	// This field is ignored when updating or creating new forums.
	DeletedAt sql.NullTime `json:"deletedAt,omitzero"`
}

type ThreadModel struct {
	DB      *pgxpool.Pool
	Timeout *time.Duration
}

func (m *ThreadModel) Select(ctx context.Context, id uuid.UUID) (*Thread, error) {
	const query string = `
SELECT id, forum_id, title, author_id, created_at, updated_at, is_locked, deleted, deleted_at
FROM forum.threads
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
	var t Thread
	err := m.DB.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&t.ID,
		&t.ForumID,
		&t.Title,
		&t.AuthorID,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.IsLocked,
		&t.Deleted,
		&t.DeletedAt,
	)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("forum selected", slog.Any("forum", t))

	return &t, nil
}

func (m *ThreadModel) Insert(ctx context.Context, input Thread) (*Thread, error) {
	const query string = `
INSERT INTO forum.threads(forum_id, title, author_id)
VALUES($1::UUID, $2::VARCHAR(256), $3::UUID)
RETURNING id, forum_id, title, author_id, created_at, updated_at, is_locked, deleted, deleted_at;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", query),
		slog.Any("input", input),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.Info("performing query")
	var t Thread
	err := m.DB.QueryRow(
		ctx,
		query,
		input.ForumID,
		input.Title,
		input.AuthorID,
	).Scan(
		&t.ID,
		&t.ForumID,
		&t.Title,
		&t.AuthorID,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.IsLocked,
		&t.Deleted,
		&t.DeletedAt,
	)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("forum selected", slog.Any("forum", t))

	return &t, nil
}
