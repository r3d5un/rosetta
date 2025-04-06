package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestThreadModel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := models.Users.Insert(ctx, data.User{
		Name:     "Adam Smasher",
		Username: "a.smasher",
		Email:    "a.smasher@arasaka.com",
	})
	assert.NoError(t, err)

	forum, err := models.Forums.Insert(ctx, data.Forum{
		OwnerID: user.ID,
		Name:    "Crushing Militech",
	})
	assert.NoError(t, err)

	newThread := data.Thread{
		AuthorID: user.ID,
		ForumID:  forum.ID,
		Title:    "Johnny Boy",
	}

	t.Run("Insert", func(t *testing.T) {
		insertedThread, err := models.Threads.Insert(ctx, newThread)
		assert.NoError(t, err)

		if insertedThread.ID == uuid.MustParse("00000000-0000-0000-0000-000000000000") {
			t.Errorf("thread ID zero valued: %s\n", insertedThread.ID.String())
			return
		}

		newThread = *insertedThread
	})

	t.Run("Select", func(t *testing.T) {
		f, err := models.Threads.Select(ctx, newThread.ID)
		assert.NoError(t, err)
		if !assert.Equal(t, newThread, *f) {
			t.Error("inserted and selected thread do not match")
			return
		}
	})

	t.Run("SelectAll", func(t *testing.T) {
	})

	t.Run("Update", func(t *testing.T) {
	})

	t.Run("SoftDelete", func(t *testing.T) {
	})

	t.Run("Restore", func(t *testing.T) {
	})

	t.Run("Delete", func(t *testing.T) {
	})
}
