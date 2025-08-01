package data

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/rosetta/Go/internal/logging"
)

type User struct {
	// ID is the unique identifier of a user.
	ID uuid.UUID `json:"id"`
	// Name is the full name of the user.
	Name string `json:"name"`
	// Username is the unique human readable name of the account.
	Username string `json:"username,omitzero"`
	// Email is the unique email beloging to a given user account.
	Email string `json:"email,omitzero"`
	// CreatedAt denotes when a user was created.
	//
	// Upon creating a new user, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt denotes when a user was last updated.
	//
	// Upon creating a new user, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	UpdatedAt time.Time `json:"updatedAt"`
	// Deleted is a soft delete flag for a user.
	Deleted bool `json:"deleted,omitzero"`
	// DeletedAt denotes when a user was last updated.
	//
	// Upon creating a new user, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	DeletedAt sql.NullTime `json:"deletedAt,omitzero"`
}

type UserInput struct {
	// Name is the full name of the user.
	Name string `json:"name"`
	// Username is the unique human readable name of the account.
	Username string `json:"username,omitzero"`
	// Email is the unique email beloging to a given user account.
	Email string `json:"email,omitzero"`
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
	// Upon creating a new user, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	DeletedAt *time.Time `json:"deletedAt,omitzero"`
}

type UserModel struct {
	DB      *pgxpool.Pool
	Timeout *time.Duration
}

func (m *UserModel) Select(ctx context.Context, id uuid.UUID) (*User, error) {
	const query string = `
SELECT id, name, username, email, created_at, updated_at, deleted, deleted_at
FROM forum.users
WHERE id = $1;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(query)),
		slog.String("id", id.String()),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.Info("performing query")
	var u User
	err := m.DB.QueryRow(
		ctx,
		query,
		id.String(),
	).Scan(
		&u.ID,
		&u.Name,
		&u.Username,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Deleted,
		&u.DeletedAt,
	)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("user selected", slog.Any("user", u))

	return &u, nil
}

func (m *UserModel) SelectAll(ctx context.Context, filters Filters) ([]*User, *Metadata, error) {
	query := `
SELECT id, name, username, email, created_at, updated_at, deleted, deleted_at
FROM forum.users
WHERE ($2::UUID IS NULL OR id = $2::UUID)
  AND ($3::VARCHAR(256) IS NULL or name = $3::VARCHAR(256))
  AND ($4::VARCHAR(256) IS NULL or username = $4::VARCHAR(256))
  AND ($5::VARCHAR(256) IS NULL or email = $5::VARCHAR(256))
  AND ($6::TIMESTAMP IS NULL or created_at >= $6::TIMESTAMP)
  AND ($7::TIMESTAMP IS NULL or created_at <= $7::TIMESTAMP)
  AND ($8::TIMESTAMP IS NULL or updated_at >= $8::TIMESTAMP)
  AND ($9::TIMESTAMP IS NULL or updated_at <= $9::TIMESTAMP)
` + CreateOrderByClause(filters.OrderBy) + `
LIMIT $1::INTEGER
`

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(query)),
		slog.Any("filters", filters),
	))

	logger.Info("performing query")
	rows, err := m.DB.Query(
		ctx,
		query,
		filters.PageSize,
		filters.ID,
		filters.Name,
		filters.Username,
		filters.Email,
		filters.CreatedAtFrom,
		filters.CreatedAtTo,
		filters.UpdatedAtFrom,
		filters.UpdatedAtTo,
	)
	if err != nil {
		logger.Error("unable to perform query", slog.String("error", err.Error()))
		return nil, nil, err
	}

	users := []*User{}

	for rows.Next() {
		var u User

		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Username,
			&u.Email,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.Deleted,
			&u.DeletedAt,
		)
		if err != nil {
			return nil, nil, handleError(err, logger)
		}
		users = append(users, &u)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, handleError(err, logger)
	}
	length := len(users)
	var metadata Metadata
	if length > 0 {
		metadata.LastSeen = users[length-1].ID
	}
	if length >= filters.PageSize {
		metadata.Next = true
	}
	metadata.ResponseLength = length

	logger.Info("users selected", slog.Any("metadata", metadata))
	return users, &metadata, nil
}

func (m *UserModel) Insert(ctx context.Context, input UserInput) (*User, error) {
	const query string = `
INSERT INTO forum.users(name, username, email)
VALUES ($1, $2, $3)
RETURNING id, name, username, email, created_at, updated_at, deleted, deleted_at;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(query)),
		slog.Any("input", input),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.Info("performing query")
	var u User
	err := m.DB.QueryRow(
		ctx,
		query,
		input.Name,
		input.Username,
		input.Email,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Username,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Deleted,
		&u.DeletedAt,
	)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("user created", slog.Any("user", u))

	return &u, nil
}

func (m *UserModel) Update(ctx context.Context, input UserPatch) (*User, error) {
	const query string = `
UPDATE forum.users
SET name       = COALESCE($2, name),
    username   = COALESCE($3, username),
    email      = COALESCE($4, email),
    deleted    = COALESCE($5, deleted),
    deleted_at = COALESCE($6, deleted_at),
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, username, email, created_at, updated_at, deleted, deleted_at;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(query)),
		slog.Any("input", input),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.Info("performing query")
	var u User
	err := m.DB.QueryRow(
		ctx,
		query,
		input.ID,
		input.Name,
		input.Username,
		input.Email,
		input.Deleted,
		input.DeletedAt,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Username,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Deleted,
		&u.DeletedAt,
	)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("user updated", slog.Any("user", u))

	return &u, nil
}

func (m *UserModel) SoftDelete(ctx context.Context, id uuid.UUID) (*User, error) {
	const query string = `
UPDATE forum.users
SET deleted    = TRUE,
    deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, username, email, created_at, updated_at, deleted, deleted_at;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(query)),
		slog.String("id", id.String()),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.Info("performing query")
	var u User
	err := m.DB.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Username,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Deleted,
		&u.DeletedAt,
	)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("user soft deleted", slog.Any("user", u))

	return &u, nil
}

func (m *UserModel) Restore(ctx context.Context, id uuid.UUID) (*User, error) {
	const query string = `
UPDATE forum.users
SET deleted    = FALSE,
    deleted_at = NULL,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, username, email, created_at, updated_at, deleted, deleted_at;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(query)),
		slog.String("id", id.String()),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.Info("performing query")
	var u User
	err := m.DB.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Username,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Deleted,
		&u.DeletedAt,
	)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("user restored", slog.Any("user", u))

	return &u, nil
}

func (m *UserModel) Delete(ctx context.Context, id uuid.UUID) (*User, error) {
	const query string = `
DELETE
FROM forum.users
WHERE id = $1
RETURNING id, name, username, email, created_at, updated_at, deleted, deleted_at;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(query)),
		slog.String("id", id.String()),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.Info("performing query")
	var u User
	err := m.DB.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Username,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Deleted,
		&u.DeletedAt,
	)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("user updated", slog.Any("user", u))

	return &u, nil
}
