package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/r3d5un/rosetta/Go/internal/data"
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

	t.Run("Read", func(t *testing.T) {
		u, err := repository.UserReader.Read(ctx, user.ID, true)
		assert.NoError(t, err)
		assert.Equal(t, u.ID, user.ID)
	})

	t.Run("List", func(t *testing.T) {
		u, metadata, err := repository.UserReader.List(ctx, data.Filters{PageSize: 100}, true)
		assert.NoError(t, err)
		assert.NotEmpty(t, metadata)
		assert.GreaterOrEqual(t, len(u), 1)
	})
}
