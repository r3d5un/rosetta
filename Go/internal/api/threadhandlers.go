package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/r3d5un/rosetta/Go/internal/repo"
	"github.com/r3d5un/rosetta/Go/internal/rest"
	"github.com/r3d5un/rosetta/Go/internal/validator"
)

type ThreadResponse struct {
	Data repo.Thread `json:"data"`
}

type ThreadListResponse struct {
	Data     []*repo.Thread `json:"data"`
	Metadata *data.Metadata `json:"metadata"`
}

func (api *API) getThreadHandler(w http.ResponseWriter, r *http.Request) {
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

	qs := r.URL.Query()
	include := rest.ReadRequiredQueryBoolean(qs, "include", false)

	forum, err := api.repo.ThreadReader.Read(ctx, *forumID, *threadID, include)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			rest.NotFoundResponse(ctx, w, r)
		case errors.Is(err, context.DeadlineExceeded):
			rest.TimeoutResponse(ctx, w, r)
		default:
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	rest.RespondWithJSON(w, r, http.StatusOK, ThreadResponse{Data: *forum}, nil)
}

func (api *API) listThreadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	v := validator.New()
	qs := r.URL.Query()
	filters := data.Filters{}
	var err error

	filters.ForumID, err = rest.ReadPathParamID(ctx, "forum_id", r)
	if err != nil {
		rest.InvalidParameterResponse(ctx, w, r, "forum_id", err)
		return
	}

	filters.PageSize = rest.ReadRequiredQueryInt(qs, "page_size", 25, v)
	filters.ID = rest.ReadOptionalQueryUUID(qs, "id", v)
	filters.AuthorID = rest.ReadOptionalQueryUUID(qs, "author_id", v)
	filters.CreatedAtFrom = rest.ReadOptionalQueryDate(qs, "created_at_from", v)
	filters.CreatedAtTo = rest.ReadOptionalQueryDate(qs, "created_at_to", v)
	filters.UpdatedAtFrom = rest.ReadOptionalQueryDate(qs, "updated_at_from", v)
	filters.UpdatedAtTo = rest.ReadOptionalQueryDate(qs, "updated_at_to", v)
	filters.Deleted = rest.ReadOptionalQueryBoolean(qs, "deleted")
	filters.DeletedAtFrom = rest.ReadOptionalQueryDate(qs, "deleted_at_from", v)
	filters.DeletedAtTo = rest.ReadOptionalQueryDate(qs, "deleted_at_to", v)
	filters.LastSeen = *rest.ReadRequiredQueryUUID(qs, "deleted_at_to", v, uuid.MustParse("00000000-0000-0000-0000-000000000000"))
	include := rest.ReadRequiredQueryBoolean(qs, "include", false)

	if !v.Valid() {
		rest.ValidationFailedResponse(ctx, w, r, v.Errors)
		return
	}

	threads, metadata, err := api.repo.ThreadReader.List(ctx, filters, include)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			rest.TimeoutResponse(ctx, w, r)
		default:
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	rest.RespondWithJSON(
		w,
		r,
		http.StatusOK,
		ThreadListResponse{Data: threads, Metadata: metadata},
		nil,
	)
}

func (api *API) postThreadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	forumID, err := rest.ReadPathParamID(ctx, "forum_id", r)
	if err != nil {
		rest.InvalidParameterResponse(ctx, w, r, "forum_id", err)
		return
	}

	var input repo.ThreadInput

	err = rest.ReadJSON(r, &input)
	if err != nil {
		rest.BadRequestResponse(w, r, err, "unable to parse JSON request body")
		return
	}
	if *forumID != input.ForumID {
		rest.BadRequestResponse(w, r, err, "response body thread ID does not match path parameter")
		return
	}

	forum, err := api.repo.ThreadWriter.Create(ctx, input)
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

	rest.RespondWithJSON(w, r, http.StatusOK, ThreadResponse{Data: *forum}, nil)
}
