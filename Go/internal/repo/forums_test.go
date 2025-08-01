package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/r3d5un/rosetta/Go/internal/repo"
	"github.com/stretchr/testify/assert"
)

func TestForumRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u, err := repository.UserWriter.Create(ctx, repo.UserInput{
		Name:     "Saburo Arasaka",
		Username: "s.arasaka",
		Email:    "s.arasaka@arasaka.com",
	})
	assert.NoError(t, err)

	var forum repo.Forum

	t.Run("Create", func(t *testing.T) {
		f, err := repository.ForumWriter.Create(ctx, repo.ForumInput{
			OwnerID: u.ID,
			Name:    "Crushing Militech",
		})
		assert.NoError(t, err)

		forum = *f
	})

	t.Run("Read", func(t *testing.T) {
		f, err := repository.ForumReader.Read(ctx, forum.ID, true)
		assert.NoError(t, err)
		assert.Equal(t, f.ID, forum.ID)
	})

	t.Run("List", func(t *testing.T) {
		f, metadata, err := repository.ForumReader.List(ctx, data.Filters{PageSize: 100}, true)
		assert.NoError(t, err)
		assert.NotEmpty(t, metadata)
		assert.GreaterOrEqual(t, len(f), 1)
	})

	t.Run("Update", func(t *testing.T) {
		forumName := "Surviving Militech"
		f, err := repository.ForumWriter.Update(
			ctx,
			repo.ForumPatch{ID: forum.ID, Name: &forumName},
		)
		assert.NoError(t, err)
		assert.Equal(t, f.ID, forum.ID)
		assert.Equal(t, f.Name, forumName)
	})

	t.Run("Delete", func(t *testing.T) {
		f, err := repository.ForumWriter.Delete(ctx, forum.ID)
		assert.NoError(t, err)
		assert.Equal(t, f.Deleted, true)
	})

	t.Run("Restore", func(t *testing.T) {
		f, err := repository.ForumWriter.Restore(ctx, forum.ID)
		assert.NoError(t, err)
		assert.Equal(t, f.Deleted, false)
	})

	t.Run("PermanentlyDelete", func(t *testing.T) {
		_, err := repository.ForumWriter.PermanentlyDelete(ctx, forum.ID)
		assert.NoError(t, err)
	})
}
