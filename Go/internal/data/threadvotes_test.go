package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestThreadVoteModel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := models.Users.Insert(ctx, data.UserInput{
		Name:     "Hanako Arasaka",
		Username: "h.arasaka",
		Email:    "h.arasaka@arasaka.com",
	})
	assert.NoError(t, err)

	forum, err := models.Forums.Insert(ctx, data.ForumInput{
		OwnerID: user.ID,
		Name:    "Night City Players",
	})
	assert.NoError(t, err)

	insertedThread, err := models.Threads.Insert(ctx, data.ThreadInput{
		AuthorID: user.ID,
		ForumID:  forum.ID,
		Title:    "About V",
	})
	assert.NoError(t, err)
	t.Log(insertedThread)

	newVote := data.ThreadVote{
		ThreadID: insertedThread.ID,
		UserID:   user.ID,
		Vote:     1,
	}

	t.Run("Vote", func(t *testing.T) {
		vote, err := models.ThreadVotes.Vote(ctx, newVote)
		assert.NoError(t, err)
		assert.NotEqual(t, newVote, vote)
	})

	t.Run("SelectSum", func(t *testing.T) {
		count, err := models.ThreadVotes.SelectSum(ctx, data.Filters{
			ThreadID: &insertedThread.ID,
			UserID:   &user.ID,
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, *count, 0, "thread vote less than 0")
	})
}
