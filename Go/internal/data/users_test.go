package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestUserModel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user data.User

	t.Run("Insert", func(t *testing.T) {
		u, err := models.Users.Insert(ctx, data.UserInput{
			Name:     "Johnny Silverhand",
			Username: "samurai",
			Email:    "jsilverhand@samurai.com",
		})
		assert.NoError(t, err)

		if u.ID == uuid.MustParse("00000000-0000-0000-0000-000000000000") {
			t.Errorf("user ID zero valued: %s\n", u.ID.String())
			return
		}

		user = *u
	})

	t.Run("Select", func(t *testing.T) {
		u, err := models.Users.Select(ctx, user.ID)
		assert.NoError(t, err)
		if !assert.Equal(t, user, *u) {
			t.Error("inserted and selected users do not match")
			return
		}
	})

	t.Run("SelectAll", func(t *testing.T) {
		users, metadata, err := models.Users.SelectAll(ctx, data.Filters{PageSize: 100})
		assert.NoError(t, err)
		assert.NotEmpty(t, users)
		assert.NotEmpty(t, metadata)
		if !assert.Equal(t, users[len(users)-1].ID, metadata.LastSeen) {
			t.Errorf(
				"last seen ID %s does not match the last selected user ID %s\n",
				metadata.LastSeen,
				users[len(users)-1].ID,
			)
		}
	})

	t.Run("Update", func(t *testing.T) {
		newName := "Silverhand"
		updatedUser, err := models.Users.Update(ctx, data.UserPatch{
			ID:   user.ID,
			Name: &newName,
		})
		assert.NoError(t, err)
		assert.NotEqual(t, user, *updatedUser)
		assert.Equal(t, newName, updatedUser.Name)
	})

	t.Run("SoftDelete", func(t *testing.T) {
		deletedUser, err := models.Users.SoftDelete(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, deletedUser.Deleted, true)
	})

	t.Run("Restore", func(t *testing.T) {
		restoredUser, err := models.Users.Restore(ctx, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, restoredUser.Deleted, false)
	})

	t.Run("Delete", func(t *testing.T) {
		deletedUser, err := models.Users.Delete(ctx, user.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, deletedUser)
	})
}
