package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestPostModel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := models.Users.Insert(ctx, data.User{
		Name:     "delemain",
		Username: "delamain",
		Email:    "delamain@delamain.com",
	})
	assert.NoError(t, err)

	forum, err := models.Forums.Insert(ctx, data.Forum{
		OwnerID: user.ID,
		Name:    "Troublesome taxis",
	})
	assert.NoError(t, err)

	insertedThread, err := models.Threads.Insert(ctx, data.Thread{
		AuthorID: user.ID,
		ForumID:  forum.ID,
		Title:    "Rouge cars",
	})
	assert.NoError(t, err)

	newPost := data.Post{
		ThreadID: insertedThread.ID,
		ReplyTo:  uuid.NullUUID{Valid: false},
		Content:  "A rogue taxi is nearby, here are the precise coordinates",
		AuthorID: user.ID,
	}

	t.Run("Insert", func(t *testing.T) {
		insertedPost, err := models.Posts.Insert(ctx, newPost)
		assert.NoError(t, err)

		newPost = *insertedPost
	})

	t.Run("Select", func(t *testing.T) {
		selectedPost, err := models.Posts.Select(ctx, newPost.ID)
		assert.NoError(t, err)
		assert.Equal(t, newPost, *selectedPost)
	})

	t.Run("SelectAll", func(t *testing.T) {
		selectedPosts, metadata, err := models.Posts.SelectAll(ctx, data.Filters{PageSize: 25})
		assert.NoError(t, err)
		assert.NotEmpty(t, metadata.LastSeen)
		assert.GreaterOrEqual(t, len(selectedPosts), 0)
	})

	t.Run("Update", func(t *testing.T) {
		updatedContent := "A rogue taxi is nearby, here are the precise coordinates: 1.1.1.1"
		updatedPost, err := models.Posts.Update(ctx, data.PostPatch{
			ID:       newPost.ID,
			ThreadID: newPost.ThreadID,
			Content:  &updatedContent,
		})
		assert.NoError(t, err)
		assert.Equal(t, updatedPost.Content, updatedContent)
	})

	t.Run("SoftDelete", func(t *testing.T) {
		deletedPost, err := models.Posts.SoftDelete(ctx, newPost.ID)
		assert.NoError(t, err)
		assert.Equal(t, deletedPost.Deleted, true)
	})

	t.Run("Restore", func(t *testing.T) {
		deletedPost, err := models.Posts.Restore(ctx, newPost.ID)
		assert.NoError(t, err)
		assert.Equal(t, deletedPost.Deleted, false)
	})

	t.Run("Delete", func(t *testing.T) {
		deletedPost, err := models.Posts.Delete(ctx, newPost.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, deletedPost)
	})
}
