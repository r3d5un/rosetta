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
	Forum *Forum `json:"forum,omitzero"`
	// Author of the thread.
	Author *User `json:"author,omitzero"`
	// Votes is the sum of votes the thread has received
	Votes *int `json:"votes,omitzero"`
	// PostCount is the number of posts within a thread
	PostCount *int `json:"post_count,omitzero"`
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

type ThreadInput struct {
	// ForumID is the parent forum this thread belongs to.
	ForumID uuid.UUID `json:"forumId"`
	// Title is the subject the thread is about.
	Title string `json:"title"`
	// AuthorID is the unique identifier of the author of the thread.
	AuthorID uuid.UUID `json:"authorId"`
}

func (f *ThreadInput) Row() data.ThreadInput {
	return data.ThreadInput{
		ForumID:  f.ForumID,
		Title:    f.Title,
		AuthorID: f.AuthorID,
	}
}

type ThreadPatch struct {
	// ID is the unique identifier of the thread
	//
	// Upon creating a new thread, any existing values in this field is ignored.
	ID uuid.UUID `json:"id"`
	// ForumID is the parent forum this thread belongs to.
	ForumID uuid.UUID `json:"forumId"`
	// Title is the subject the thread is about.
	Title *string `json:"title"`
	// AuthorID is the unique identifier of the author of the thread.
	AuthorID *uuid.UUID `json:"authorId"`
}

func (f *ThreadPatch) Row() data.ThreadPatch {
	return data.ThreadPatch{
		ID:       f.ID,
		ForumID:  f.ForumID,
		Title:    database.NewNullString(f.Title),
		AuthorID: database.NewNullUUID(f.AuthorID),
	}
}

type ThreadReader interface {
	Read(context.Context, uuid.UUID, uuid.UUID, bool) (*Thread, error)
	List(context.Context, data.Filters, bool) ([]*Thread, *data.Metadata, error)
}

type ThreadWriter interface {
	Create(context.Context, ThreadInput) (*Thread, error)
	Update(context.Context, ThreadPatch) (*Thread, error)
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

func (r *ThreadRepository) Read(
	ctx context.Context,
	forumID uuid.UUID,
	threadID uuid.UUID,
	include bool,
) (*Thread, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group(
			"parameters",
			slog.String("forumId", forumID.String()),
			slog.String("threadId", threadID.String()),
			slog.Bool("include", include)))

	logger.LogAttrs(ctx, slog.LevelInfo, "retrieving thread")
	row, err := r.models.Threads.Select(ctx, forumID, threadID)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to select thread", slog.String("error", err.Error()),
		)
		return nil, err
	}
	thread := newThreadFromRow(*row)
	logger.LogAttrs(ctx, slog.LevelInfo, "thread retrieved")

	if !include {
		return thread, nil
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 4)
	var threadMu sync.Mutex

	wg.Add(1)
	go func() {
		defer wg.Done()
		author, err := r.userReader.Read(ctx, thread.AuthorID, false)
		if err != nil {
			errCh <- err
			return
		}

		threadMu.Lock()
		thread.Author = author
		threadMu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		forum, err := r.forumReader.Read(ctx, thread.ForumID, true)
		if err != nil {
			errCh <- err
			return
		}

		threadMu.Lock()
		thread.Forum = forum
		threadMu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		votes, err := r.models.ThreadVotes.SelectSum(
			ctx,
			data.Filters{ThreadID: &thread.ID},
		)
		if err != nil {
			errCh <- err
			return
		}

		threadMu.Lock()
		thread.Votes = votes
		threadMu.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := r.models.Posts.SelectCount(ctx, data.Filters{ThreadID: &thread.ID})
		if err != nil {
			errCh <- err
			return
		}

		threadMu.Lock()
		thread.PostCount = count
		threadMu.Unlock()
	}()

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			logger.Error("unable to include all data", slog.String("error", err.Error()))
		}
	}

	return thread, nil
}

func (r *ThreadRepository) List(
	ctx context.Context,
	filter data.Filters,
	include bool,
) ([]*Thread, *data.Metadata, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.Any("filters", filter), slog.Bool("include", include)))

	logger.LogAttrs(ctx, slog.LevelInfo, "retrieving threads")
	rows, metadata, err := r.models.Threads.SelectAll(ctx, filter)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to select thread", slog.String("error", err.Error()),
		)
	}
	logger = logging.LoggerFromContext(ctx).With(slog.Group(
		"parameters",
		slog.Any("filters", filter),
		slog.Any("metadata", metadata)),
		slog.Bool("include", include))
	logger.LogAttrs(ctx, slog.LevelInfo, "threads retrieved")

	threads := make([]*Thread, len(rows))
	var wg sync.WaitGroup
	errCh := make(chan error, len(rows)*4)
	var threadsMu sync.Mutex

	for i, row := range rows {
		threads[i] = newThreadFromRow(*row)

		if !include {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			author, err := r.userReader.Read(ctx, threads[i].AuthorID, false)
			if err != nil {
				errCh <- err
				return
			}

			threadsMu.Lock()
			threads[i].Author = author
			threadsMu.Unlock()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			forum, err := r.forumReader.Read(ctx, threads[i].ForumID, true)
			if err != nil {
				errCh <- err
				return
			}

			threadsMu.Lock()
			threads[i].Forum = forum
			threadsMu.Unlock()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			votes, err := r.models.ThreadVotes.SelectSum(
				ctx,
				data.Filters{ThreadID: &threads[i].ID},
			)
			if err != nil {
				errCh <- err
				return
			}

			threadsMu.Lock()
			threads[i].Votes = votes
			threadsMu.Unlock()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			count, err := r.models.Posts.SelectCount(
				ctx,
				data.Filters{ThreadID: &threads[i].ID},
			)
			if err != nil {
				errCh <- err
				return
			}

			threadsMu.Lock()
			threads[i].PostCount = count
			threadsMu.Unlock()
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			logger.Error("unable to include all data", slog.String("error", err.Error()))
		}
	}

	return threads, metadata, nil
}

func (r *ThreadRepository) Create(ctx context.Context, input ThreadInput) (*Thread, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.Any("thread", input)))

	logger.LogAttrs(ctx, slog.LevelInfo, "creating thread")
	row, err := r.models.Threads.Insert(ctx, input.Row())
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to create thread", slog.String("error", err.Error()),
		)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "thread created")

	return newThreadFromRow(*row), nil
}

func (r *ThreadRepository) Update(ctx context.Context, patch ThreadPatch) (*Thread, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.Any("patch", patch)))

	logger.LogAttrs(ctx, slog.LevelInfo, "updating thread")
	row, err := r.models.Threads.Update(ctx, patch.Row())
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to update thread", slog.String("error", err.Error()),
		)
		return nil, err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "thread updated")

	return newThreadFromRow(*row), nil
}

func (r *ThreadRepository) Delete(ctx context.Context, id uuid.UUID) (*Thread, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.String("id", id.String())))

	logger.LogAttrs(ctx, slog.LevelInfo, "deleting thread")
	row, err := r.models.Threads.SoftDelete(ctx, id)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to delete thread", slog.String("error", err.Error()),
		)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "thread deleted")

	return newThreadFromRow(*row), nil
}

func (r *ThreadRepository) Restore(ctx context.Context, id uuid.UUID) (*Thread, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.String("id", id.String())))

	logger.LogAttrs(ctx, slog.LevelInfo, "restoring thread")
	row, err := r.models.Threads.Restore(ctx, id)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to restore thread", slog.String("error", err.Error()),
		)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "thread restored")

	return newThreadFromRow(*row), nil
}

func (r *ThreadRepository) PermanentlyDelete(ctx context.Context, id uuid.UUID) (*Thread, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.String("id", id.String())))

	logger.LogAttrs(ctx, slog.LevelInfo, "deleting thread")
	row, err := r.models.Threads.Delete(ctx, id)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to delete thread", slog.String("error", err.Error()),
		)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "thread deleted")

	return newThreadFromRow(*row), nil
}
