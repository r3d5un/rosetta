package data

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Models struct {
	Forums      ForumModel
	Users       UserModel
	Threads     ThreadModel
	ThreadVotes ThreadVoteModel
	Posts       PostModel
	PostVotes   PostVoteModel
}

func NewModels(pool *pgxpool.Pool, timeout *time.Duration) Models {
	return Models{
		Forums:      ForumModel{DB: pool, Timeout: timeout},
		Users:       UserModel{DB: pool, Timeout: timeout},
		Threads:     ThreadModel{DB: pool, Timeout: timeout},
		ThreadVotes: ThreadVoteModel{DB: pool, Timeout: timeout},
		Posts:       PostModel{DB: pool, Timeout: timeout},
		PostVotes:   PostVoteModel{DB: pool, Timeout: timeout},
	}
}
