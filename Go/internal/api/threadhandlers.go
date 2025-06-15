package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/r3d5un/rosetta/Go/internal/repo"
	"github.com/r3d5un/rosetta/Go/internal/rest"
)

type ThreadResponse struct {
	Data repo.Thread `json:"data"`
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
