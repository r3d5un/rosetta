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
