package data

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound       = errors.New("record not found")
	ErrConstraintViolation  = errors.New("constraint violation")
	ErrUniqueIndexViolation = errors.New("unique index violation")
	ErrUniqueKeyViolation   = errors.New("unique key violation")
)

type Models struct {
	// TODO: Add models
}

func NewModels(pool *pgxpool.Pool, timeout *time.Duration) Models {
	return Models{
		// TODO: Instantiate models
	}
}
