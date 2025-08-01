package repo

import (
	"context"
	"log/slog"
	"sync"
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
	// Owner is the user which owns the forum.
	Owner *User `json:"owner,omitzero"`
	// ThreadCount is the number of threads within the forum
	ThreadCount *int `json:"threadCount,omitzero"`
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

type ForumInput struct {
	// OwnerID is the unique identifier of a forum.
	OwnerID uuid.UUID `json:"ownerId"`
	// Name is the human readable name of the forum
	Name string `json:"name"`
	// Description contains a description about the purposes and topics of a forum.
	Description *string `json:"description,omitzero"`
}

func (f *ForumInput) Row() data.ForumInput {
	return data.ForumInput{
		OwnerID:     f.OwnerID,
		Name:        f.Name,
		Description: database.NewNullString(f.Description),
	}
}

type ForumPatch struct {
	// ID is the unique identifier of a forum.
	//
	// Upon creating a new forum, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	ID uuid.UUID `json:"id"`
	// OwnerID is the unique identifier of a forum.
	OwnerID *uuid.UUID `json:"ownerId,omitzero"`
	// Name is the human readable name of the forum
	Name *string `json:"name,omitzero"`
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
	Create(context.Context, ForumInput) (*Forum, error)
	Delete(context.Context, uuid.UUID) (*Forum, error)
	Update(context.Context, ForumPatch) (*Forum, error)
	Restore(context.Context, uuid.UUID) (*Forum, error)
	PermanentlyDelete(context.Context, uuid.UUID) (*Forum, error)
}

type ForumRepository struct {
	models     *data.Models
	userReader UserReader
}

func NewForumRepository(models *data.Models, userReader UserReader) ForumRepository {
	return ForumRepository{
		models:     models,
		userReader: userReader,
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
	forum := newForumFromRow(*row)
	logger.LogAttrs(ctx, slog.LevelInfo, "forum retrieved")

	if !include {
		return forum, nil
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	var forumMu sync.Mutex

	wg.Add(1)
	go func() {
		defer wg.Done()
		owner, err := r.userReader.Read(ctx, forum.OwnerID, false)
		if err != nil {
			errCh <- err
			return
		}

		forumMu.Lock()
		forum.Owner = owner
		forumMu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := r.models.Threads.SelectCount(ctx, data.Filters{ForumID: &forum.ID})
		if err != nil {
			errCh <- err
			return
		}

		forumMu.Lock()
		forum.ThreadCount = count
		forumMu.Unlock()
	}()

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			logger.Error("unable to include all data", slog.String("error", err.Error()))
		}
	}

	return forum, nil
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
	var wg sync.WaitGroup
	errCh := make(chan error, len(rows)*2)
	var forumsMu sync.Mutex

	for i, row := range rows {
		forums[i] = newForumFromRow(*row)

		if !include {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			owner, err := r.userReader.Read(ctx, forums[i].OwnerID, false)
			if err != nil {
				errCh <- err
				return
			}

			forumsMu.Lock()
			forums[i].Owner = owner
			forumsMu.Unlock()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			count, err := r.models.Threads.SelectCount(ctx, data.Filters{ForumID: &forums[i].ID})
			if err != nil {
				errCh <- err
				return
			}

			forumsMu.Lock()
			forums[i].ThreadCount = count
			forumsMu.Unlock()
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			logger.Error("unable to include all data", slog.String("error", err.Error()))
		}
	}

	return forums, metadata, nil
}

func (r *ForumRepository) Create(ctx context.Context, input ForumInput) (*Forum, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.Any("forum", input)))

	logger.LogAttrs(ctx, slog.LevelInfo, "creating forum")
	row, err := r.models.Forums.Insert(ctx, input.Row())
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to create forum", slog.String("error", err.Error()),
		)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "forum created")

	return newForumFromRow(*row), nil
}

func (r *ForumRepository) Update(ctx context.Context, patch ForumPatch) (*Forum, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.Any("patch", patch)))

	logger.LogAttrs(ctx, slog.LevelInfo, "updating forum")
	row, err := r.models.Forums.Update(ctx, patch.Row())
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to update forum", slog.String("error", err.Error()),
		)
		return nil, err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "forum updated")

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
