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

	u, err := repository.UserWriter.Create(ctx, repo.UserInput{
		Name:     "delemain",
		Username: "delamain",
		Email:    "delamain@delamain.com",
	})
	assert.NoError(t, err)

	f, err := repository.ForumWriter.Create(ctx, repo.ForumInput{
		OwnerID: u.ID,
		Name:    "Troublesome taxis",
	})
	assert.NoError(t, err)

	thread, err := repository.ThreadWriter.Create(ctx, repo.ThreadInput{
		AuthorID: u.ID,
		ForumID:  f.ID,
		Title:    "Rouge cars",
	})
	assert.NoError(t, err)

	var post repo.Post

	t.Run("Create", func(t *testing.T) {
		p, err := repository.PostWriter.Create(ctx, repo.PostInput{
			ThreadID: thread.ID,
			Content:  "A rogue taxi is nearby, here are the precise coordinates",
			AuthorID: u.ID,
		})
		assert.NoError(t, err)

		post = *p
	})

	t.Run("Read", func(t *testing.T) {
		p, err := repository.PostReader.Read(ctx, f.ID, thread.ID, post.ID, true)
		assert.NoError(t, err)
		assert.Equal(t, p.ID, post.ID)
	})

	t.Run("List", func(t *testing.T) {
		posts, metadata, err := repository.PostReader.List(
			ctx, f.ID, thread.ID, data.Filters{PageSize: 100}, true,
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, metadata)
		assert.GreaterOrEqual(t, len(posts), 1)
	})

	t.Run("Update", func(t *testing.T) {
		updatedContent := "A rogue taxi is nearby, here are the precise coordinates: 1.1.1.1"
		p, err := repository.PostWriter.Update(ctx, repo.PostPatch{
			ID:       post.ID,
			ThreadID: post.ThreadID,
			Content:  &updatedContent,
		})
		assert.NoError(t, err)
		assert.NotEqual(t, thread, *p)
		assert.Equal(t, updatedContent, p.Content)
	})

	t.Run("Delete", func(t *testing.T) {
		p, err := repository.PostWriter.Delete(ctx, post.ID)
		assert.NoError(t, err)
		assert.Equal(t, p.Deleted, true)
	})

	t.Run("Restore", func(t *testing.T) {
		p, err := repository.PostWriter.Restore(ctx, post.ID)
		assert.NoError(t, err)
		assert.Equal(t, p.Deleted, false)
	})

	t.Run("PermanentlyDelete", func(t *testing.T) {
		_, err := repository.PostWriter.PermanentlyDelete(ctx, post.ID)
		assert.NoError(t, err)
	})
}
