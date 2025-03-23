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
	Forums  ForumModel
	Users   UserModel
	Threads ThreadModel
	Posts   PostModel
}

func NewModels(pool *pgxpool.Pool, timeout *time.Duration) Models {
	return Models{
		Forums:  ForumModel{DB: pool, Timeout: timeout},
		Users:   UserModel{DB: pool, Timeout: timeout},
		Threads: ThreadModel{DB: pool, Timeout: timeout},
		Posts:   PostModel{DB: pool, Timeout: timeout},
	}
}
