package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/r3d5un/rosetta/Go/internal/logging"
	"github.com/r3d5un/rosetta/Go/internal/repo"
	"github.com/r3d5un/rosetta/Go/internal/rest"
)

type UserReponse struct {
	Data repo.User `json:"data"`
}

type UserListResponse struct {
	Data     []*repo.User   `json:"data"`
	Metadata *data.Metadata `json:"metadata"`
}

func (api *API) getUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := rest.ReadPathParamID(ctx, "id", r)
	if err != nil {
		rest.InvalidParameterResponse(ctx, w, r, "id", err)
		return
	}

	qs := r.URL.Query()
	include := rest.ReadQueryBoolean(qs, "include", false)

	user, err := api.repo.UserReader.Read(ctx, *id, include)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			rest.NotFoundResponse(ctx, w, r)
			return
		case errors.Is(err, context.DeadlineExceeded):
			rest.TimeoutResponse(ctx, w, r)
			return
		default:
			rest.ServerErrorResponse(w, r, err)
			return
		}
	}

	rest.RespondWithJSON(w, r, http.StatusOK, UserReponse{Data: *user}, nil)
}

func (api *API) listUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)
	filters := data.Filters{}

	qs := r.URL.Query()
	// TODO: ID
	// TODO: Name
	// TODO: Username
	// TODO: email
	// TODO: created_at_from
	// TODO: created_at_to
	// TODO: updated_at_from
	// TODO: updated_at_to
	include := rest.ReadQueryBoolean(qs, "include", false)

	logger.Info(
		"validating filters",
		slog.Any("filters", filters),
		slog.Any("queryParameters", qs),
	)
	// TODO: Validate filters

	users, metadata, err := api.repo.UserReader.List(ctx, filters, include)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			rest.TimeoutResponse(ctx, w, r)
			return
		default:
			rest.ServerErrorResponse(w, r, err)
			return
		}
	}

	rest.RespondWithJSON(
		w,
		r,
		http.StatusOK,
		UserListResponse{Data: users, Metadata: metadata},
		nil,
	)
}
