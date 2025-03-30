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
