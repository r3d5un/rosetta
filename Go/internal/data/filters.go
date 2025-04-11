package data

import (
	"time"

	"github.com/google/uuid"
)

type Metadata struct {
	LastSeen       uuid.UUID `json:"lastSeen,omitzero"`
	Next           bool      `json:"next"`
	ResponseLength int       `json:"responseLength"`
}

type Filters struct {
	ID            *uuid.UUID `json:"id"`
	OwnerID       *uuid.UUID `json:"ownerId"`
	UserID        *uuid.UUID `json:"userId"`
	PostID        *uuid.UUID `json:"postId"`
	ThreadID      *uuid.UUID `json:"threadId"`
	ForumID       *uuid.UUID `json:"forumId"`
	AuthorID      *uuid.UUID `json:"authorId"`
	Name          *string    `json:"name"`
	Title         *string    `json:"title"`
	Username      *string    `json:"username,omitempty"`
	Email         *string    `json:"email,omitempty"`
	CreatedAtFrom *time.Time `json:"createdAtFrom"`
	CreatedAtTo   *time.Time `json:"createdAtTo"`
	UpdatedAtFrom *time.Time `json:"updatedAtFrom"`
	UpdatedAtTo   *time.Time `json:"updatedAtTo"`
	DeletedAtFrom *time.Time `json:"deletedAtFrom"`
	DeletedAtTo   *time.Time `json:"deletedAtTo"`
	Deleted       *bool      `json:"deleted"`
	IsLocked      *bool      `json:"isLocked"`

	OrderBy         []string  `json:"order_by,omitempty"`
	OrderBySafeList []string  `json:"order_by_safe_list,omitempty"`
	LastSeen        uuid.UUID `json:"lastSeen"`
	PageSize        int       `json:"page_size,omitempty"`
}
