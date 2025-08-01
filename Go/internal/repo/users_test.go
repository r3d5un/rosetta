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

	var user repo.User

	t.Run("Create", func(t *testing.T) {
		u, err := repository.UserWriter.Create(ctx, repo.UserInput{
			Name:     "Johnny Silverhand",
			Username: "samurai",
			Email:    "jsilverhand@samurai.com",
		})
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

	t.Run("Update", func(t *testing.T) {
		username := "silverhand"
		u, err := repository.UserWriter.Update(
			ctx,
			repo.UserPatch{ID: user.ID, Username: &username},
		)
		assert.NoError(t, err)
		assert.Equal(t, u.ID, user.ID)
		assert.Equal(t, u.Username, username)
	})

	t.Run("Delete", func(t *testing.T) {
		u, err := repository.UserWriter.Delete(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, u.Deleted, true)
	})

	t.Run("Restore", func(t *testing.T) {
		u, err := repository.UserWriter.Restore(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, u.Deleted, false)
	})

	t.Run("PermanentlyDelete", func(t *testing.T) {
		_, err := repository.UserWriter.PermanentlyDelete(ctx, user.ID)
		assert.NoError(t, err)
	})
}
