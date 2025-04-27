package repo

import (
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/r3d5un/rosetta/Go/internal/database"
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
