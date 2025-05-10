package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestPostVoteModel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := models.Users.Insert(ctx, data.User{
		Name:     "Morgan Blackhand",
		Username: "blackhand",
		Email:    "blackhand@afterlife.com",
	})
	assert.NoError(t, err)

	forum, err := models.Forums.Insert(ctx, data.Forum{
		OwnerID: user.ID,
		Name:    "Contracts",
	})
	assert.NoError(t, err)

	insertedThread, err := models.Threads.Insert(ctx, data.Thread{
		AuthorID: user.ID,
		ForumID:  forum.ID,
		Title:    "Bounty: Adam Smasher",
	})
	assert.NoError(t, err)

	newPost := data.Post{
		ThreadID: insertedThread.ID,
		ReplyTo:  uuid.NullUUID{Valid: false},
		Content:  "Adam Smasher located at Arasaka reginal office. Moving to apprehend.",
		AuthorID: user.ID,
	}

	insertedPost, err := models.Posts.Insert(ctx, newPost)
	assert.NoError(t, err)

	newPost = *insertedPost

	newVote := data.PostVote{
		PostID: insertedPost.ID,
		UserID: user.ID,
		Vote:   1,
	}

	t.Run("Vote", func(t *testing.T) {
		vote, err := models.PostVotes.Vote(ctx, newVote)
		assert.NoError(t, err)
		assert.NotEqual(t, newVote, vote)
	})

	t.Run("SelectSum", func(t *testing.T) {
		count, err := models.PostVotes.SelectSum(ctx, data.Filters{
			PostID: &insertedPost.ID,
			UserID: &user.ID,
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, *count, 0, "post vote less than 0")
	})
}
