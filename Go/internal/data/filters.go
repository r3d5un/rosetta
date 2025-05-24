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
	ID            *uuid.UUID `json:"id,omitzero"`
	OwnerID       *uuid.UUID `json:"ownerId,omitzero"`
	UserID        *uuid.UUID `json:"userId,omitzero"`
	PostID        *uuid.UUID `json:"postId,omitzero"`
	ThreadID      *uuid.UUID `json:"threadId,omitzero"`
	ForumID       *uuid.UUID `json:"forumId,omitzero"`
	AuthorID      *uuid.UUID `json:"authorId,omitzero"`
	Name          *string    `json:"name,omitzero"`
	Title         *string    `json:"title,omitzero"`
	Username      *string    `json:"username,omitzero"`
	Email         *string    `json:"email,omitzero"`
	CreatedAtFrom *time.Time `json:"createdAtFrom,omitzero"`
	CreatedAtTo   *time.Time `json:"createdAtTo,omitzero"`
	UpdatedAtFrom *time.Time `json:"updatedAtFrom,omitzero"`
	UpdatedAtTo   *time.Time `json:"updatedAtTo,omitzero"`
	DeletedAtFrom *time.Time `json:"deletedAtFrom,omitzero"`
	DeletedAtTo   *time.Time `json:"deletedAtTo,omitzero"`
	Deleted       *bool      `json:"deleted,omitzero"`
	IsLocked      *bool      `json:"isLocked,omitzero"`

	OrderBy         []string  `json:"order_by,omitzero"`
	OrderBySafeList []string  `json:"order_by_safe_list,omitzero"`
	LastSeen        uuid.UUID `json:"lastSeen,omitzero"`
	PageSize        int       `json:"page_size,omitzero"`
}
