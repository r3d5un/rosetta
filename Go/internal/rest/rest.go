package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/logging"
)

var (
	ErrPathParamID = errors.New("path parameter is invalid")
)

const (
	notFoundMsg string = "resource not found"
	timeoutMsg  string = "the server took to long to respond"
)

type ErrorMessage struct {
	Message any `json:"message"`
}

func ErrorResponse(
	w http.ResponseWriter, r *http.Request, status int, message any,
) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.Info("writing error response", slog.Int("status", status), slog.Any("message", message))
	RespondWithJSON(w, r, status, ErrorMessage{Message: message}, nil)
}

func LogError(r *http.Request, err error) {
	logging.LoggerFromContext(r.Context()).Error(
		"an error occurred",
		slog.String("request_method", r.Method),
		slog.String("request_url", r.URL.String()),
		slog.String("error", err.Error()),
	)
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	LogError(r, err)
	const serverErrorMsg string = "the server encountered a problem and could not process your request"
	ErrorResponse(w, r, http.StatusInternalServerError, serverErrorMsg)
}

func InvalidParameterResponse(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	param string,
	err error,
) {
	logger := logging.LoggerFromContext(ctx)
	logger.LogAttrs(ctx, slog.LevelInfo, "parameter invalid", slog.String("error", err.Error()))

	ErrorResponse(w, r, http.StatusNotFound, fmt.Sprintf("%s is not a valid parameter", param))
}

func NotFoundResponse(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(ctx)
	logger.LogAttrs(ctx, slog.LevelInfo, notFoundMsg)
	ErrorResponse(w, r, http.StatusNotFound, notFoundMsg)
}

func TimeoutResponse(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context())
	logger.LogAttrs(ctx, slog.LevelInfo, timeoutMsg)
	ErrorResponse(w, r, http.StatusRequestTimeout, timeoutMsg)
}

func RespondWithJSON(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	data any,
	headers http.Header,
) {
	logger := logging.LoggerFromContext(r.Context())

	logger.Info("marshalling data")
	js, err := json.Marshal(data)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}

	js = append(js, '\n')

	logger.Info("adding headers")
	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	logger.Info("writing response")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(js); err != nil {
		ServerErrorResponse(w, r, err)
	}
}

func ReadPathParamID(ctx context.Context, key string, r *http.Request) (*uuid.UUID, error) {
	logger := logging.LoggerFromContext(ctx)

	pathValue := r.PathValue(key)
	if pathValue == "" {
		logger.LogAttrs(ctx, slog.LevelInfo, "path parameter empty", slog.String("key", key))
		return nil, ErrPathParamID
	}

	id, err := uuid.Parse(pathValue)
	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelInfo,
			"unable to parse path parameter UUID",
			slog.String("key", key),
			slog.String("value", pathValue),
		)
		return nil, ErrPathParamID
	}

	return &id, err
}

func ReadQueryBoolean(
	qs url.Values,
	key string,
	defaultValue bool,
) bool {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return defaultValue
	}
	return b
}
