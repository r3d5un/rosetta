package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/r3d5un/rosetta/Go/internal/data"
	"github.com/r3d5un/rosetta/Go/internal/repo"
	"github.com/r3d5un/rosetta/Go/internal/rest"
	"github.com/r3d5un/rosetta/Go/internal/validator"
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
	include := rest.ReadRequiredQueryBoolean(qs, "include", false)

	user, err := api.repo.UserReader.Read(ctx, *id, include)
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

	rest.RespondWithJSON(w, r, http.StatusOK, UserReponse{Data: *user}, nil)
}

func (api *API) listUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	v := validator.New()
	qs := r.URL.Query()
	filters := data.Filters{}

	filters.PageSize = rest.ReadRequiredQueryInt(qs, "page_size", 25, v)
	filters.ID = rest.ReadOptionalQueryUUID(qs, "id", v)
	filters.Name = rest.ReadOptionalQueryString(qs, "name")
	filters.Username = rest.ReadOptionalQueryString(qs, "username")
	filters.Email = rest.ReadOptionalQueryString(qs, "email")
	filters.CreatedAtFrom = rest.ReadOptionalQueryDate(qs, "created_at_from", v)
	filters.CreatedAtTo = rest.ReadOptionalQueryDate(qs, "created_at_to", v)
	filters.UpdatedAtFrom = rest.ReadOptionalQueryDate(qs, "updated_at_from", v)
	filters.UpdatedAtTo = rest.ReadOptionalQueryDate(qs, "updated_at_to", v)
	include := rest.ReadRequiredQueryBoolean(qs, "include", false)

	if !v.Valid() {
		rest.ValidationFailedResponse(ctx, w, r, v.Errors)
		return
	}

	users, metadata, err := api.repo.UserReader.List(ctx, filters, include)
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
		UserListResponse{Data: users, Metadata: metadata},
		nil,
	)
}

func (api *API) postUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input repo.UserInput

	err := rest.ReadJSON(r, &input)
	if err != nil {
		rest.BadRequestResponse(w, r, err, "unable to parse JSON request body")
		return
	}

	user, err := api.repo.UserWriter.Create(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrUniqueConstraintViolation):
			rest.ConstraintViolationResponse(w, r, err, "user ID already exists")
		case errors.Is(err, data.ErrCheckConstraintViolation):
			rest.ConstraintViolationResponse(w, r, err, "used failed input checks")
		case errors.Is(err, context.DeadlineExceeded):
			rest.TimeoutResponse(ctx, w, r)
		default:
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	rest.RespondWithJSON(w, r, http.StatusOK, UserReponse{Data: *user}, nil)
}

func (api *API) patchUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input repo.UserPatch

	err := rest.ReadJSON(r, &input)
	if err != nil {
		rest.BadRequestResponse(w, r, err, "unable to parse JSON request body")
		return
	}

	user, err := api.repo.UserWriter.Update(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrUniqueConstraintViolation):
			rest.ConstraintViolationResponse(w, r, err, "user ID already exists")
		case errors.Is(err, data.ErrCheckConstraintViolation):
			rest.ConstraintViolationResponse(w, r, err, "used failed input checks")
		case errors.Is(err, context.DeadlineExceeded):
			rest.TimeoutResponse(ctx, w, r)
		default:
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	rest.RespondWithJSON(w, r, http.StatusOK, UserReponse{Data: *user}, nil)
}
