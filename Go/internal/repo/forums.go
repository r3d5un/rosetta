package repo

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/r3d5un/rosetta/Go/internal/database"
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
	Description *string `json:"description,omitzero"`
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
	DeletedAt *time.Time `json:"deletedAt,omitzero"`
}

func newForumFromRow(row data.Forum) *Forum {
	return &Forum{
		ID:          row.ID,
		OwnerID:     row.OwnerID,
		Name:        row.Name,
		Description: database.NullStringToPtr(row.Description),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		Deleted:     row.Deleted,
		DeletedAt:   database.NullTimeToPtr(row.DeletedAt),
	}
}

func (f *Forum) Row() data.Forum {
	return data.Forum{
		ID:          f.ID,
		OwnerID:     f.OwnerID,
		Name:        f.Name,
		Description: database.NewNullString(f.Description),
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
		Deleted:     f.Deleted,
		DeletedAt:   database.NewNullTime(f.DeletedAt),
	}
}

type ForumPatch struct {
	// ID is the unique identifier of a forum.
	//
	// Upon creating a new forum, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	ID uuid.UUID `json:"id"`
	// OwnerID is the unique identifier of a forum.
	OwnerID *uuid.UUID `json:"ownerId"`
	// Name is the human readable name of the forum
	Name *string `json:"name"`
	// Description contains a description about the purposes and topics of a forum.
	Description *string `json:"description,omitzero"`
}

func (f *ForumPatch) Row() data.ForumPatch {
	return data.ForumPatch{
		ID:          f.ID,
		OwnerID:     database.NewNullUUID(f.OwnerID),
		Name:        database.NewNullString(f.Name),
		Description: database.NewNullString(f.Description),
	}
}

type ForumReader interface {
	Read(context.Context, uuid.UUID, bool) (*Forum, error)
	List(context.Context, data.Filters, bool) ([]*Forum, *data.Metadata, error)
}

type ForumWriter interface {
	Create(context.Context, Forum) (*Forum, error)
	Delete(context.Context, uuid.UUID) (*Forum, error)
	Restore(context.Context, uuid.UUID) (*Forum, error)
	PermanentlyDelete(context.Context, uuid.UUID) (*Forum, error)
}

type ForumRepository struct {
	models *data.Models
}

func NewForumRepository(models *data.Models) ForumRepository {
	return ForumRepository{
		models: models,
	}
}

func (r *ForumRepository) Read(ctx context.Context, id uuid.UUID, include bool) (*Forum, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.String("id", id.String()), slog.Bool("include", include)))

	logger.LogAttrs(ctx, slog.LevelInfo, "retrieving forum")
	row, err := r.models.Forums.Select(ctx, id)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to select forum", slog.String("error", err.Error()),
		)
		return nil, err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "forum retrieved")

	return newForumFromRow(*row), nil
}

func (r *ForumRepository) List(
	ctx context.Context,
	filter data.Filters,
	include bool,
) ([]*Forum, *data.Metadata, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.Any("filters", filter), slog.Bool("include", include)))

	logger.LogAttrs(ctx, slog.LevelInfo, "retrieving forums")
	rows, metadata, err := r.models.Forums.SelectAll(ctx, filter)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to select forum", slog.String("error", err.Error()),
		)
	}
	logger = logging.LoggerFromContext(ctx).With(slog.Group(
		"parameters",
		slog.Any("filters", filter),
		slog.Any("metadata", metadata)),
		slog.Bool("include", include))
	logger.LogAttrs(ctx, slog.LevelInfo, "forums retrieved")

	forums := make([]*Forum, len(rows))
	for i, row := range rows {
		forums[i] = newForumFromRow(*row)
	}

	return forums, metadata, nil
}

func (r *ForumRepository) Create(ctx context.Context, forum Forum) (*Forum, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.Any("forum", forum)))

	logger.LogAttrs(ctx, slog.LevelInfo, "creating forum")
	row, err := r.models.Forums.Insert(ctx, forum.Row())
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to create forum", slog.String("error", err.Error()),
		)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "forum created")

	return newForumFromRow(*row), nil
}

func (r *ForumRepository) Delete(ctx context.Context, id uuid.UUID) (*Forum, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.String("id", id.String())))

	logger.LogAttrs(ctx, slog.LevelInfo, "deleting forum")
	row, err := r.models.Forums.SoftDelete(ctx, id)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to delete forum", slog.String("error", err.Error()),
		)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "forum deleted")

	return newForumFromRow(*row), nil
}

func (r *ForumRepository) Restore(ctx context.Context, id uuid.UUID) (*Forum, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.String("id", id.String())))

	logger.LogAttrs(ctx, slog.LevelInfo, "restoring forum")
	row, err := r.models.Forums.Restore(ctx, id)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to restore forum", slog.String("error", err.Error()),
		)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "forum restored")

	return newForumFromRow(*row), nil
}

func (r *ForumRepository) PermanentlyDelete(ctx context.Context, id uuid.UUID) (*Forum, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.String("id", id.String())))

	logger.LogAttrs(ctx, slog.LevelInfo, "deleting forum")
	row, err := r.models.Forums.Delete(ctx, id)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to delete forum", slog.String("error", err.Error()),
		)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "forum deleted")

	return newForumFromRow(*row), nil
}
