package data_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestForumModel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := models.Users.Insert(ctx, data.User{
		Name:     "Saburo Arasaka",
		Username: "s.arasaka",
		Email:    "s.arasaka@arasaka.com",
	})
	assert.NoError(t, err)

	newForum := data.Forum{
		OwnerID: user.ID,
		Name:    "Crushing Militech",
	}

	t.Run("Insert", func(t *testing.T) {
		f, err := models.Forums.Insert(ctx, newForum)
		assert.NoError(t, err)

		if f.ID == uuid.MustParse("00000000-0000-0000-0000-000000000000") {
			t.Errorf("forum ID zero valued: %s\n", f.ID.String())
			return
		}

		newForum = *f
	})

	t.Run("Select", func(t *testing.T) {
		f, err := models.Forums.Select(ctx, newForum.ID)
		assert.NoError(t, err)
		if !assert.Equal(t, newForum, *f) {
			t.Error("inserted and selected forums do not match")
			return
		}
	})

	t.Run("SelectAll", func(t *testing.T) {
		forums, metadata, err := models.Forums.SelectAll(ctx, data.Filters{PageSize: 100})
		assert.NoError(t, err)
		assert.NotEmpty(t, forums)
		assert.NotEmpty(t, metadata)
		if !assert.Equal(t, forums[len(forums)-1].ID, metadata.LastSeen) {
			t.Errorf(
				"last seen ID %s does not match the last selected user ID %s\n",
				metadata.LastSeen,
				forums[len(forums)-1].ID,
			)
		}
	})

	t.Run("Update", func(t *testing.T) {
		newName := "Surviving Militech"
		updatedForum, err := models.Forums.Update(ctx, data.ForumPatch{
			ID:   newForum.ID,
			Name: sql.NullString{Valid: true, String: newName},
		})
		assert.NoError(t, err)
		assert.NotEqual(t, newForum, *updatedForum)
		assert.Equal(t, newName, updatedForum.Name)
	})

	t.Run("SoftDelete", func(t *testing.T) {
		deletedForum, err := models.Forums.SoftDelete(ctx, newForum.ID)
		assert.NoError(t, err)
		assert.Equal(t, deletedForum.Deleted, true)
	})

	t.Run("Restore", func(t *testing.T) {
		restoreForum, err := models.Forums.Restore(ctx, newForum.ID)
		assert.NoError(t, err)
		assert.Equal(t, restoreForum.Deleted, false)
	})

	t.Run("Delete", func(t *testing.T) {
		deletedForum, err := models.Forums.Delete(ctx, newForum.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, deletedForum)
	})
}
