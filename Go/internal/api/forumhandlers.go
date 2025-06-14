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

type ForumResponse struct {
	Data repo.Forum `json:"data"`
}

type ForumListResponse struct {
	Data     []*repo.Forum  `json:"data"`
	Metadata *data.Metadata `json:"metadata"`
}

func (api *API) getForumHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := rest.ReadPathParamID(ctx, "id", r)
	if err != nil {
		rest.InvalidParameterResponse(ctx, w, r, "id", err)
		return
	}

	qs := r.URL.Query()
	include := rest.ReadRequiredQueryBoolean(qs, "include", false)

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

func (api *API) listForumHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	v := validator.New()
	qs := r.URL.Query()
	filters := data.Filters{}

	filters.PageSize = rest.ReadRequiredQueryInt(qs, "page_size", 25, v)
	filters.ID = rest.ReadOptionalQueryUUID(qs, "id", v)
	filters.OwnerID = rest.ReadOptionalQueryUUID(qs, "owner_id", v)
	filters.Name = rest.ReadOptionalQueryString(qs, "name")
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

	forums, metadata, err := api.repo.ForumReader.List(ctx, filters, include)
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
		ForumListResponse{Data: forums, Metadata: metadata},
		nil,
	)
}

func (api *API) postForumHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input repo.ForumInput

	err := rest.ReadJSON(r, &input)
	if err != nil {
		rest.BadRequestResponse(w, r, err, "unable to parse JSON request body")
		return
	}

	forum, err := api.repo.ForumWriter.Create(ctx, input)
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

	rest.RespondWithJSON(w, r, http.StatusOK, ForumResponse{Data: *forum}, nil)
}

func (api *API) patchForumHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input repo.ForumPatch

	err := rest.ReadJSON(r, &input)
	if err != nil {
		rest.BadRequestResponse(w, r, err, "unable to parse JSON request body")
		return
	}

	forum, err := api.repo.ForumWriter.Update(ctx, input)
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

	rest.RespondWithJSON(w, r, http.StatusOK, ForumResponse{Data: *forum}, nil)
}

func (api *API) deleteForumHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := rest.ReadPathParamID(ctx, "id", r)
	if err != nil {
		rest.InvalidParameterResponse(ctx, w, r, "id", err)
		return
	}

	forum, err := api.repo.ForumWriter.Delete(ctx, *id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			rest.NotFoundResponse(ctx, w, r)
		case errors.Is(err, context.DeadlineExceeded):
			rest.TimeoutResponse(ctx, w, r)
		case errors.Is(err, data.ErrForeignKeyConstraintViolation):
			rest.ConstraintViolationResponse(w, r, err, "forum referenced by other resources")
		default:
			rest.ServerErrorResponse(w, r, err)
		}
		return
	}

	rest.RespondWithJSON(w, r, http.StatusOK, ForumResponse{Data: *forum}, nil)
}

func (api *API) restoreForumHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := rest.ReadPathParamID(ctx, "id", r)
	if err != nil {
		rest.InvalidParameterResponse(ctx, w, r, "id", err)
		return
	}

	forum, err := api.repo.ForumWriter.Restore(ctx, *id)
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
