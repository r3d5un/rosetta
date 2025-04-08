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
	// Upon creating a new thread, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	ID uuid.UUID `json:"id"`
	// ForumID is the parent forum this thread belongs to.
	ForumID uuid.UUID `json:"forumId"`
	// Title is the subject the thread is about.
	Title string `json:"title"`
	// AuthorID is the unique identifier of the author of the thread.
	AuthorID uuid.UUID `json:"authorId"`
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
	// IsLocked denotes whether a thread had been locked for changes.
	//
	// This field is ignored when updating or creating new threads.
	IsLocked bool `json:"isLocked"`
	// Deleted is a soft delete flag for a thread.
	Deleted bool `json:"deleted,omitzero"`
	// DeletedAt denotes when a thread was marked as deleted.
	//
	// This field is ignored when updating or creating new thread.
	DeletedAt sql.NullTime `json:"deletedAt,omitzero"`
}

type ThreadPatch struct {
	// ID is the unique identifier of the thread
	//
	// Upon creating a new thread, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	ID uuid.UUID `json:"id"`
	// ForumID is the parent forum this thread belongs to.
	ForumID *uuid.UUID `json:"forumId"`
	// Title is the subject the thread is about.
	Title *string `json:"title"`
	// AuthorID is the unique identifier of the author of the thread.
	AuthorID *uuid.UUID `json:"authorId"`
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
	logger.Info("thread selected", slog.Any("thread", t))

	return &t, nil
}

func (m *ThreadModel) SelectAll(
	ctx context.Context,
	filters Filters,
) ([]*Thread, *Metadata, error) {
	query := `
SELECT id, forum_id, title, author_id, created_at, updated_at, is_locked, deleted, deleted_at
FROM forum.threads
WHERE ($2::UUID IS NULL OR id = $2::UUID)
  AND ($3::UUID IS NULL OR forum_id = $3::UUID)
  AND ($4::VARCHAR(256) IS NULL OR title = $4::VARCHAR(256))
  AND ($5::UUID IS NULL OR author_id = $5::UUID)
  AND ($6::TIMESTAMP IS NULL or created_at >= $6::TIMESTAMP)
  AND ($7::TIMESTAMP IS NULL or created_at <= $7::TIMESTAMP)
  AND ($8::TIMESTAMP IS NULL or updated_at >= $8::TIMESTAMP)
  AND ($9::TIMESTAMP IS NULL or updated_at <= $9::TIMESTAMP)
  AND ($10::BOOLEAN IS NULL or is_locked = $10::BOOLEAN)
  AND ($11::BOOLEAN IS NULL or deleted = $11::BOOLEAN)
  AND ($12::TIMESTAMP IS NULL or deleted_at >= $12::TIMESTAMP)
  AND ($13::TIMESTAMP IS NULL or deleted_at <= $13::TIMESTAMP)
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
		filters.ForumID,
		filters.Title,
		filters.AuthorID,
		filters.CreatedAtFrom,
		filters.CreatedAtTo,
		filters.UpdatedAtFrom,
		filters.UpdatedAtTo,
		filters.IsLocked,
		filters.Deleted,
		filters.DeletedAtFrom,
		filters.DeletedAtTo,
		filters.LastSeen,
	)
	if err != nil {
		logger.Error("unable to perform query", slog.String("error", err.Error()))
		return nil, nil, err
	}

	threads := []*Thread{}

	for rows.Next() {
		var t Thread

		err := rows.Scan(
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
			return nil, nil, handleError(err, logger)
		}
		threads = append(threads, &t)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, handleError(err, logger)
	}
	length := len(threads)
	var metadata Metadata
	if length > 0 {
		metadata.LastSeen = threads[length-1].ID
		metadata.Next = true
	}
	metadata.ResponseLength = length

	logger.Info("thread selected", slog.Any("metadata", metadata))
	return threads, &metadata, nil
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
	logger.Info("thread selected", slog.Any("thread", t))

	return &t, nil
}

func (m *ThreadModel) Update(ctx context.Context, input ThreadPatch) (*Thread, error) {
	const query string = `
UPDATE forum.threads
SET forum_id = COALESCE($2::UUID, forum_id),
    title = COALESCE($3::VARCHAR(256), title),
    author_id = COALESCE($4::UUID, author_id)
WHERE id = $1::UUID
RETURNING id, forum_id, title, author_id, created_at, updated_at, is_locked, deleted, deleted_at;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(query)),
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
		input.ID,
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
	logger.Info("thread selected", slog.Any("forum", t))

	return &t, nil
}

func (m *ThreadModel) SoftDelete(ctx context.Context, id uuid.UUID) (*Thread, error) {
	const query string = `
UPDATE forum.threads
SET deleted    = TRUE,
    deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1::UUID
RETURNING id, forum_id, title, author_id, created_at, updated_at, is_locked, deleted, deleted_at;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(query)),
		slog.String("id", id.String()),
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
	logger.Info("thread selected", slog.Any("thread", t))

	return &t, nil
}

func (m *ThreadModel) Restore(ctx context.Context, id uuid.UUID) (*Thread, error) {
	const query string = `
UPDATE forum.threads
SET deleted    = FALSE,
    deleted_at = NULL,
    updated_at = NOW()
WHERE id = $1::UUID
RETURNING id, forum_id, title, author_id, created_at, updated_at, is_locked, deleted, deleted_at;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(query)),
		slog.String("id", id.String()),
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
	logger.Info("thread selected", slog.Any("thread", t))

	return &t, nil
}
