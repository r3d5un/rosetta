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

type Forum struct {
	// ID is the unique identifier of a forum.
	//
	// Upon creating a new forum, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	ID uuid.UUID `json:"id"`
	// OwnerID is the unique identifier of a forum.
	OwnerID uuid.UUID `json:"ownerId"`
	// Name is the human readable name of the forum
	Name string `json:"name"`
	// Description contains a description about the purposes and topics of a forum.
	Description sql.NullString `json:"description,omitzero"`
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
	// Deleted is a soft delete flag for a forum.
	Deleted bool `json:"deleted,omitzero"`
	// DeletedAt denotes when a forum was marked as deleted.
	//
	// This field is ignored when updating or creating new forums.
	DeletedAt sql.NullTime `json:"deletedAt,omitzero"`
}

type ForumModel struct {
	DB      *pgxpool.Pool
	Timeout *time.Duration
}

func (m *ForumModel) Select(ctx context.Context, id uuid.UUID) (*Forum, error) {
	const query string = `
SELECT id, owner_id, name, description, created_at, updated_at, deleted, deleted_at
FROM forum.forums
WHERE id = $1;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", query),
		slog.String("id", id.String()),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.Info("performing query")
	var f Forum
	err := m.DB.QueryRow(
		ctx,
		query,
		id.String(),
	).Scan(
		&f.ID,
		&f.OwnerID,
		&f.Name,
		&f.Description,
		&f.CreatedAt,
		&f.UpdatedAt,
		&f.Deleted,
		&f.DeletedAt,
	)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("forum selected", slog.Any("forum", f))

	return &f, nil
}

func (m *ForumModel) SelectAll(ctx context.Context, filters Filters) ([]*Forum, *Metadata, error) {
	query := `
SELECT id, owner_id, name, description, created_at, updated_at, deleted, deleted_at
FROM forum.forums
WHERE ($2::UUID IS NULL OR id = $2::UUID)
  AND ($3::UUID IS NULL OR owner_id = $3::UUID)
  AND ($4::VARCHAR(256) IS NULL or name = $4::VARCHAR(256))
  AND ($5::TIMESTAMP IS NULL or created_at >= $5::TIMESTAMP)
  AND ($6::TIMESTAMP IS NULL or created_at <= $6::TIMESTAMP)
  AND ($7::TIMESTAMP IS NULL or updated_at >= $7::TIMESTAMP)
  AND ($8::TIMESTAMP IS NULL or updated_at <= $8::TIMESTAMP)
  AND ($9::BOOLEAN IS NULL or deleted = $9::BOOLEAN)
  AND ($10::TIMESTAMP IS NULL or deleted_at >= $10::TIMESTAMP)
  AND ($11::TIMESTAMP IS NULL or deleted_at <= $11::TIMESTAMP)
` + CreateOrderByClause(filters.OrderBy) + `
LIMIT $1::INTEGER
`

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("statement", logging.MinifySQL(query)),
		slog.Any("filters", filters),
	))

	logger.Info("performing query")
	rows, err := m.DB.Query(
		ctx,
		query,
		filters.PageSize,
		filters.ID,
		filters.OwnerID,
		filters.Name,
		filters.CreatedAtFrom,
		filters.CreatedAtTo,
		filters.UpdatedAtFrom,
		filters.UpdatedAtTo,
		filters.Deleted,
		filters.DeletedAtFrom,
		filters.DeletedAtTo,
	)
	if err != nil {
		logger.Error("unable to perform query", slog.String("error", err.Error()))
		return nil, nil, err
	}

	forums := []*Forum{}

	for rows.Next() {
		var f Forum

		err := rows.Scan(
			&f.ID,
			&f.OwnerID,
			&f.Name,
			&f.Description,
			&f.CreatedAt,
			&f.UpdatedAt,
			&f.Deleted,
			&f.DeletedAt,
		)
		if err != nil {
			return nil, nil, handleError(err, logger)
		}
		forums = append(forums, &f)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, handleError(err, logger)
	}
	length := len(forums)
	var metadata Metadata
	if length > 0 {
		metadata.LastSeen = forums[length-1].ID
		metadata.Next = true
	}
	metadata.ResponseLength = length

	logger.Info("forums selected", slog.Any("metadata", metadata))
	return forums, &metadata, nil
}
