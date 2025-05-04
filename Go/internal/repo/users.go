package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/r3d5un/rosetta/Go/internal/database"
)

type User struct {
	// ID is the unique identifier of a user.
	//
	// Upon creating a new user, any existing values in this field is ignored.
	ID uuid.UUID `json:"id"`
	// Name is the full name of the user.
	Name string `json:"name"`
	// Username is the unique human readable name of the account.
	Username string `json:"username,omitzero"`
	// Email is the unique email beloging to a given user account.
	Email string `json:"email,omitzero"`
	// CreatedAt denotes when a user was created.
	//
	// Upon creating a new user, any existing values in this field is ignored.
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt denotes when a user was last updated.
	//
	// Upon creating a new user, any existing values in this field is ignored.
	UpdatedAt time.Time `json:"updatedAt"`
	// Deleted is a soft delete flag for a user.
	Deleted bool `json:"deleted,omitzero"`
	// DeletedAt denotes when a user was last updated.
	//
	// Upon creating a new user, any existing values in this field is ignored.
	DeletedAt *time.Time `json:"deletedAt,omitzero"`
}

func newUserFromRow(row data.User) *User {
	return &User{
		ID:        row.ID,
		Name:      row.Name,
		Username:  row.Username,
		Email:     row.Email,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		Deleted:   row.Deleted,
		DeletedAt: database.NullTimeToPtr(row.DeletedAt),
	}
}

func (f *User) Row() data.User {
	return data.User{
		ID:        f.ID,
		Name:      f.Name,
		Username:  f.Username,
		Email:     f.Email,
		CreatedAt: f.CreatedAt,
		UpdatedAt: f.UpdatedAt,
		Deleted:   f.Deleted,
		DeletedAt: database.NewNullTime(f.DeletedAt),
	}
}

type UserPatch struct {
	// ID is the unique identifier of a user.
	ID uuid.UUID `json:"id"`
	// Name is the full name of the user.
	//
	// If populated, will update the name of the user.
	Name *string `json:"name"`
	// Username is the unique human readable name of the account.
	//
	// If populated, will update the username of the user.
	Username *string `json:"username,omitempty"`
	// Email is the unique email beloging to a given user account.
	//
	// If populated, will update the username of the user.
	Email *string `json:"email,omitempty"`
	// Deleted is a soft delete flag for a user.
	Deleted *bool `json:"deleted,omitzero"`
	// DeletedAt denotes when a user was last updated.
	//
	// Upon creating a new user, any existing values in this field is ignored.
	DeletedAt *time.Time `json:"deletedAt,omitzero"`
}

func (u *UserPatch) Row() data.UserPatch {
	return data.UserPatch{
		ID:        u.ID,
		Name:      u.Name,
		Username:  u.Username,
		Email:     u.Email,
		Deleted:   u.Deleted,
		DeletedAt: u.DeletedAt,
	}
}

type UserReader interface {
	Read(context.Context, uuid.UUID, bool) (*User, error)
	List(context.Context, data.Filters, bool) ([]*User, *data.Metadata, error)
}

type UserWriter interface {
	Create(context.Context, User) (*User, error)
}

type UserRepository struct {
	models *data.Models
}
