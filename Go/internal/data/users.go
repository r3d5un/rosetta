package data

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/rosetta/Go/internal/logging"
)

type User struct {
	// ID is the unique identifier of a user.
	//
	// Upon creating a new user, any existing values in this field is ignored. The database handles
	// setting the value upon insertion.
	ID uuid.UUID `json:"id"`
	// Name is the full name of the user.
	Name string `json:"name"`
	// Username is the unique human readable name of the account.
	Username string `json:"username,omitempty"`
	// Email is the unique email beloging to a given user account.
	Email string `json:"email,omitempty"`
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
}

type UserModel struct {
	DB      *pgxpool.Pool
	Timeout *time.Duration
}

func (m *UserModel) Select(ctx context.Context, id uuid.UUID) (*User, error) {
	const query string = `
SELECT id, name, username, email, created_at, updated_at
FROM forum.users
WHERE id = $1;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", query),
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
	)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("user selected", slog.Any("user", u))

	return &u, nil
}

func (m *UserModel) SelectAll(ctx context.Context, filters Filters) ([]*User, *Metadata, error) {
	query := `
SELECT id, name, username, email, created_at, updated_at
FROM forum.users
WHERE ($2 IS NULL OR id = $2)
  AND ($3 IS NULL or name = $3)
  AND ($4 IS NULL or username = $4)
  AND ($5 IS NULL or email = $5)
  AND ($6 IS NULL or created_at >= $6)
  AND ($7 IS NULL or created_at <= $7)
  AND ($8 IS NULL or updated_at >= $8)
  AND ($9 IS NULL or updated_at <= $9)
` + CreateOrderByClause(filters.OrderBy) + `
LIMIT $1
`

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("statement", logging.MinifySQL(query)),
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
		var user User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, nil, handleError(err, logger)
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, handleError(err, logger)
	}
	length := len(users)
	var metadata Metadata
	if length > 0 {
		metadata.LastSeen = users[length-1].ID
		metadata.Next = true
	}
	metadata.ResponseLength = length

	logger.Info("users selected", slog.Any("metadata", metadata))
	return users, &metadata, nil
}

func (m *UserModel) Insert(ctx context.Context, input User) (*User, error) {
	const query string = `
INSERT INTO forum.users(name, username, email)
VALUES ($1, $2, $3)
RETURNING id, name, username, email, created_at, updated_at;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", query),
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
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, username, email, created_at, updated_at;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", query),
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
	)
	if err != nil {
		return nil, handleError(err, logger)
	}
	logger.Info("user updated", slog.Any("user", u))

	return &u, nil
}
