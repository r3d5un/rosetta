package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/r3d5un/rosetta/Go/internal/repo"
	"github.com/r3d5un/rosetta/Go/internal/rest"
)

type ForumResponse struct {
	Data repo.Forum `json:"data"`
}

func (api *API) getForumHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := rest.ReadPathParamID(ctx, "id", r)
	if err != nil {
		rest.InvalidParameterResponse(ctx, w, r, "id", err)
		return
	}

	qs := r.URL.Query()
	include := rest.ReadQueryBoolean(qs, "include", false)

	forum, err := api.repo.ForumReader.Read(ctx, *id, include)
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

	rest.RespondWithJSON(w, r, http.StatusOK, ForumResponse{Data: *forum}, nil)
}
