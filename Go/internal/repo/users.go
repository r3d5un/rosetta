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
	Update(context.Context, UserPatch) (*User, error)
	Delete(context.Context, uuid.UUID) (*User, error)
	Restore(context.Context, uuid.UUID) (*User, error)
	PermanentlyDelete(context.Context, uuid.UUID) (*User, error)
}

type UserRepository struct {
	models *data.Models
}

func NewUserRepository(models *data.Models) UserRepository {
	return UserRepository{models: models}
}

func (r *UserRepository) Read(ctx context.Context, id uuid.UUID, include bool) (*User, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.String("id", id.String()), slog.Bool("include", include)))

	logger.LogAttrs(ctx, slog.LevelInfo, "retrieving user")
	row, err := r.models.Users.Select(ctx, id)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to select user", slog.String("error", err.Error()),
		)
		return nil, err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "user retrieved")

	return newUserFromRow(*row), nil
}

func (r *UserRepository) List(
	ctx context.Context,
	filter data.Filters,
	include bool,
) ([]*User, *data.Metadata, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.Any("filters", filter), slog.Bool("include", include)))

	logger.LogAttrs(ctx, slog.LevelInfo, "retrieving users")
	rows, metadata, err := r.models.Users.SelectAll(ctx, filter)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to select user", slog.String("error", err.Error()),
		)
	}
	logger = logging.LoggerFromContext(ctx).With(slog.Group(
		"parameters",
		slog.Any("filters", filter),
		slog.Any("metadata", metadata)),
		slog.Bool("include", include))
	logger.LogAttrs(ctx, slog.LevelInfo, "users retrieved")

	users := make([]*User, len(rows))
	for i, row := range rows {
		users[i] = newUserFromRow(*row)
	}

	return users, metadata, nil
}

func (r *UserRepository) Update(ctx context.Context, patch UserPatch) (*User, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.Any("patch", patch)))

	logger.LogAttrs(ctx, slog.LevelInfo, "updating user")
	row, err := r.models.Users.Update(ctx, patch.Row())
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to update user", slog.String("error", err.Error()),
		)
		return nil, err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "user updated")

	return newUserFromRow(*row), nil
}

func (r *UserRepository) Create(ctx context.Context, user User) (*User, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.Any("user", user)))

	logger.LogAttrs(ctx, slog.LevelInfo, "creating user")
	row, err := r.models.Users.Insert(ctx, user.Row())
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to create user", slog.String("error", err.Error()),
		)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "user created", slog.Any("row", row))

	return newUserFromRow(*row), nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) (*User, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.String("id", id.String())))

	logger.LogAttrs(ctx, slog.LevelInfo, "deleting user")
	row, err := r.models.Users.SoftDelete(ctx, id)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to delete user", slog.String("error", err.Error()),
		)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "user deleted")

	return newUserFromRow(*row), nil
}

func (r *UserRepository) Restore(ctx context.Context, id uuid.UUID) (*User, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.String("id", id.String())))

	logger.LogAttrs(ctx, slog.LevelInfo, "restoring user")
	row, err := r.models.Users.Restore(ctx, id)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to restore user", slog.String("error", err.Error()),
		)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "user restored")

	return newUserFromRow(*row), nil
}

func (r *UserRepository) PermanentlyDelete(ctx context.Context, id uuid.UUID) (*User, error) {
	logger := logging.LoggerFromContext(ctx).
		With(slog.Group("parameters", slog.String("id", id.String())))

	logger.LogAttrs(ctx, slog.LevelInfo, "deleting user")
	row, err := r.models.Users.Delete(ctx, id)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to delete user", slog.String("error", err.Error()),
		)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "user deleted")

	return newUserFromRow(*row), nil
}
