package repo

import (
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/r3d5un/rosetta/Go/internal/database"
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
