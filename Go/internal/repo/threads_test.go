package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/r3d5un/rosetta/Go/internal/repo"
	"github.com/stretchr/testify/assert"
)

func TestThreadRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := repo.User{
		Name:     "Adam Smasher",
		Username: "a.smasher",
		Email:    "a.smasher@arasaka.com",
	}
	u, err := repository.UserWriter.Create(ctx, user)
	assert.NoError(t, err)

	forum := repo.Forum{
		OwnerID: u.ID,
		Name:    "Crushing Militech",
	}
	f, err := repository.ForumWriter.Create(ctx, forum)
	assert.NoError(t, err)

	thread := repo.Thread{
		AuthorID: u.ID,
		ForumID:  f.ID,
		Title:    "Johnny Boy",
	}

	t.Run("Create", func(t *testing.T) {
		createdThread, err := repository.ThreadWriter.Create(ctx, thread)
		assert.NoError(t, err)

		thread = *createdThread
	})

	t.Run("Read", func(t *testing.T) {
		readThread, err := repository.ThreadReader.Read(ctx, thread.ID, true)
		assert.NoError(t, err)
		assert.Equal(t, readThread.ID, thread.ID)
	})

	t.Run("List", func(t *testing.T) {
		listedThreads, metadata, err := repository.ThreadReader.List(
			ctx, data.Filters{PageSize: 100}, true,
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, metadata)
		assert.GreaterOrEqual(t, len(listedThreads), 1)
	})

	t.Run("Update", func(t *testing.T) {
		newTitle := "Neurochipped Johnny Boy"
		updatedThread, err := repository.ThreadWriter.Update(ctx, repo.ThreadPatch{
			ID:    thread.ID,
			Title: &newTitle,
		})
		assert.NoError(t, err)
		assert.NotEqual(t, thread, *updatedThread)
		assert.Equal(t, newTitle, updatedThread.Title)
	})

	t.Run("Delete", func(t *testing.T) {
	})

	t.Run("Restore", func(t *testing.T) {
	})

	t.Run("PermanentlyDelete", func(t *testing.T) {
	})
}
