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
	DeletedAt *time.Time `json:"deletedAt,omitzero"`
	// Likes is the sum of votes the thread has received.
	Likes int64 `json:"likes"`
	// Forum that the thread belongs to.
	Forum Forum `json:"forum,omitzero"`
	// Author of the thread.
	Author User `json:"author,omitzero"`
}

func newThreadFromRow(row data.Thread) *Thread {
	return &Thread{
		ID:        row.ID,
		ForumID:   row.ForumID,
		Title:     row.Title,
		AuthorID:  row.AuthorID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		IsLocked:  row.IsLocked,
		Deleted:   row.Deleted,
		DeletedAt: database.NullTimeToPtr(row.DeletedAt),
		Likes:     row.Likes,
	}
}

func (f *Thread) Row() data.Thread {
	return data.Thread{
		ID:        f.ID,
		ForumID:   f.ForumID,
		Title:     f.Title,
		AuthorID:  f.AuthorID,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
		IsLocked:  f.IsLocked,
		Deleted:   f.Deleted,
		DeletedAt: database.NewNullTime(f.DeletedAt),
		Likes:     f.Likes,
	}
}

type ThreadPatch struct {
	// ID is the unique identifier of the thread
	//
	// Upon creating a new thread, any existing values in this field is ignored.
	ID uuid.UUID `json:"id"`
	// ForumID is the parent forum this thread belongs to.
	ForumID *uuid.UUID `json:"forumId"`
	// Title is the subject the thread is about.
	Title *string `json:"title"`
	// AuthorID is the unique identifier of the author of the thread.
	AuthorID *uuid.UUID `json:"authorId"`
}

func (f *ThreadPatch) Row() data.ThreadPatch {
	return data.ThreadPatch{
		ID:       f.ID,
		ForumID:  database.NewNullUUID(f.ForumID),
		Title:    database.NewNullString(f.Title),
		AuthorID: database.NewNullUUID(f.AuthorID),
	}
}

type ThreadReader interface {
	Read(context.Context, uuid.UUID, bool) (*Thread, error)
	List(context.Context, data.Filters, bool) ([]*Thread, *data.Metadata, error)
}

type ThreadWriter interface {
	Create(context.Context, Thread) (*Thread, error)
	Delete(context.Context, uuid.UUID) (*Thread, error)
	Restore(context.Context, uuid.UUID) (*Thread, error)
	PermanentlyDelete(context.Context, uuid.UUID) (*Thread, error)
}

type ThreadRepository struct {
	models      *data.Models
	forumReader ForumReader
	userReader  UserReader
}

func NewThreadRepository(
	models *data.Models,
	forumReader ForumReader,
	userReader UserReader,
) ThreadRepository {
	return ThreadRepository{
		models:      models,
		forumReader: forumReader,
		userReader:  userReader,
	}
}

func (r *ThreadRepository) Read(ctx context.Context, id uuid.UUID, include bool) (*Thread, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.String("id", id.String()), slog.Bool("include", include)))

	logger.LogAttrs(ctx, slog.LevelInfo, "retrieving user")
	row, err := r.models.Threads.Select(ctx, id)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to select user", slog.String("error", err.Error()),
		)
		return nil, err
	}
	thread := newThreadFromRow(*row)
	logger.LogAttrs(ctx, slog.LevelInfo, "user retrieved")

	if !include {
		return thread, nil
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	var forumMu sync.Mutex

	wg.Add(1)
	go func() {
		defer wg.Done()
		author, err := r.userReader.Read(ctx, thread.AuthorID, false)
		if err != nil {
			errCh <- err
		}

		forumMu.Lock()
		thread.Author = *author
		forumMu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		forum, err := r.forumReader.Read(ctx, thread.ForumID, true)
		if err != nil {
			errCh <- err
		}

		forumMu.Lock()
		thread.Forum = *forum
		forumMu.Unlock()
	}()

	close(errCh)

	wg.Wait()

	return thread, nil
}
