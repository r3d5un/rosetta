package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/r3d5un/rosetta/Go/internal/repo"
	"github.com/r3d5un/rosetta/Go/internal/rest"
)

type PostResponse struct {
	Data repo.Post `json:"data"`
}

type PostPostRequestBody struct {
	// ReplyTo is the ID of which this post is a reply to.
	ReplyTo *uuid.UUID `json:"replyTo"`
	// AuthorID is the unique identifier of the author of the post.
	AuthorID uuid.UUID `json:"authorId"`
	// Content is the actual text content of a post
	Content string `json:"content"`
}

func (api *API) postPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	forumID, err := rest.ReadPathParamID(ctx, "forum_id", r)
	if err != nil {
		rest.InvalidParameterResponse(ctx, w, r, "forum_id", err)
		return
	}

	threadID, err := rest.ReadPathParamID(ctx, "thread_id", r)
	if err != nil {
		rest.InvalidParameterResponse(ctx, w, r, "thread_id", err)
		return
	}

	var body PostPostRequestBody

	err = rest.ReadJSON(r, &body)
	if err != nil {
		rest.BadRequestResponse(w, r, err, "unable to parse JSON request body")
		return
	}

	post, err := api.repo.PostWriter.Create(ctx, repo.PostInput{
		ForumID:  *forumID,
		ThreadID: *threadID,
		ReplyTo:  body.ReplyTo,
		AuthorID: body.AuthorID,
		Content:  body.Content,
	})
	if err != nil {
		switch {
		case errors.Is(err, data.ErrUniqueConstraintViolation):
			rest.ConstraintViolationResponse(w, r, err, "forum ID already exists")
		case errors.Is(err, data.ErrCheckConstraintViolation):
			rest.ConstraintViolationResponse(w, r, err, "used failed input checks")
		case errors.Is(err, context.DeadlineExceeded):
			rest.TimeoutResponse(ctx, w, r)
		default:
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	rest.RespondWithJSON(w, r, http.StatusOK, PostResponse{Data: *post}, nil)
}
