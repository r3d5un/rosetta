package data

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Forum struct {
	ID          uuid.UUID      `json:"id"`
	OwnerID     uuid.UUID      `json:"ownerId"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

type ForumModel struct {
	DB      *pgxpool.Pool
	Timeout *time.Duration
}
