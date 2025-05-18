package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/r3d5un/rosetta/Go/internal/repo"
	"github.com/stretchr/testify/assert"
)

func TestForumRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := repo.User{
		Name:     "Saburo Arasaka",
		Username: "s.arasaka",
		Email:    "s.arasaka@arasaka.com",
	}
	u, err := repository.UserWriter.Create(ctx, user)
	assert.NoError(t, err)

	forum := repo.Forum{
		OwnerID: u.ID,
		Name:    "Crushing Militech",
	}

	t.Run("Create", func(t *testing.T) {
		f, err := repository.ForumWriter.Create(ctx, forum)
		assert.NoError(t, err)

		forum = *f
	})

	t.Run("Read", func(t *testing.T) {
		f, err := repository.ForumReader.Read(ctx, forum.ID, true)
		assert.NoError(t, err)
		assert.Equal(t, f.ID, forum.ID)
	})

	t.Run("List", func(t *testing.T) {
	})

	t.Run("Update", func(t *testing.T) {
	})

	t.Run("Delete", func(t *testing.T) {
	})

	t.Run("Restore", func(t *testing.T) {
	})

	t.Run("PermanentlyDelete", func(t *testing.T) {
	})
}
