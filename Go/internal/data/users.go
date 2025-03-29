package data

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/rosetta/Go/internal/logging"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			logger.Error("unable to perform query", slog.String("error", err.Error()))
			return nil, err
		}
	}
	logger.Info("query returned user")

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
			logger.Error("unable to scan query result", slog.String("error", err.Error()))
			return nil, nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		logger.Error("unable to scan query result", slog.String("error", err.Error()))
		return nil, nil, err
	}
	length := len(users)
	var metadata Metadata
	if length > 0 {
		metadata.LastSeen = users[length-1].ID
		metadata.Next = true
	}
	metadata.ResponseLength = length

	logger.Info("query returned users", slog.Any("metadata", metadata))
	return users, &metadata, nil
}
