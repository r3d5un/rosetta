package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/r3d5un/rosetta/Go/internal/repo"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := repo.User{
		Name:     "Johnny Silverhand",
		Username: "samurai",
		Email:    "jsilverhand@samurai.com",
	}

	t.Run("Create", func(t *testing.T) {
		u, err := repository.UserWriter.Create(ctx, user)
		assert.NoError(t, err)

		user = *u
	})
}
