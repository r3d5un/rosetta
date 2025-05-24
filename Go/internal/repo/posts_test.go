package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/r3d5un/rosetta/Go/internal/repo"
	"github.com/stretchr/testify/assert"
)

func TestPostRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := repo.User{
		Name:     "delemain",
		Username: "delamain",
		Email:    "delamain@delamain.com",
	}
	u, err := repository.UserWriter.Create(ctx, user)
	assert.NoError(t, err)

	forum := repo.Forum{
		OwnerID: u.ID,
		Name:    "Troublesome taxis",
	}
	f, err := repository.ForumWriter.Create(ctx, forum)
	assert.NoError(t, err)

	thread := repo.Thread{
		AuthorID: u.ID,
		ForumID:  f.ID,
		Title:    "Rouge cars",
	}
	newThread, err := repository.ThreadWriter.Create(ctx, thread)
	assert.NoError(t, err)

	post := repo.Post{
		ThreadID: newThread.ID,
		Content:  "A rogue taxi is nearby, here are the precise coordinates",
		AuthorID: u.ID,
	}

	t.Run("Create", func(t *testing.T) {
		p, err := repository.PostWriter.Create(ctx, post)
		assert.NoError(t, err)

		post = *p
	})

	t.Run("Read", func(t *testing.T) {
		p, err := repository.PostReader.Read(ctx, post.ID, true)
		assert.NoError(t, err)
		assert.Equal(t, p.ID, post.ID)
	})

	t.Run("List", func(t *testing.T) {
		posts, metadata, err := repository.PostReader.List(
			ctx, data.Filters{PageSize: 100}, true,
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, metadata)
		assert.GreaterOrEqual(t, len(posts), 1)
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
