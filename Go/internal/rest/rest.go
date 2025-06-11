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
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/logging"
	"github.com/r3d5un/rosetta/Go/internal/validator"
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

func ValidationFailedResponse(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	validationErrors map[string]string,
) {
	logger := logging.LoggerFromContext(r.Context())
	logger.LogAttrs(ctx, slog.LevelInfo, timeoutMsg, slog.Any("validationErrors", validationErrors))
	ErrorResponse(
		w,
		r,
		http.StatusUnprocessableEntity,
		fmt.Sprintf("filter validation failed: %v", validationErrors),
	)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error, msg string) {
	logger := logging.LoggerFromContext(r.Context())

	logger.Info("bad request", slog.String("error", err.Error()), slog.String("message", msg))
	ErrorResponse(w, r, http.StatusBadRequest, msg)
}

func ConstraintViolationResponse(w http.ResponseWriter, r *http.Request, err error, msg string) {
	logger := logging.LoggerFromContext(r.Context())

	logger.Info(
		"a constraint violation occurred",
		slog.String("error", err.Error()),
		slog.String("message", msg),
	)
	ErrorResponse(w, r, http.StatusConflict, msg)
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

func ReadRequiredQueryBoolean(
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

func ReadOptionalQueryBoolean(qs url.Values, key string) *bool {
	s := qs.Get(key)
	if s == "" {
		return nil
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return nil
	}
	return &b
}

func ReadRequiredQueryInt(qs url.Values, key string, defaultVal int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultVal
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultVal
	}

	return i
}

func ReadRequiredQueryUUID(
	qs url.Values,
	key string,
	v *validator.Validator,
	defaultVal uuid.UUID,
) *uuid.UUID {
	s := qs.Get(key)

	if s == "" {
		return &defaultVal
	}

	id, err := uuid.Parse(s)
	if err != nil {
		v.AddError(key, fmt.Sprintf("unable to parse value: %s", err.Error()))
	}

	return &id
}

func ReadOptionalQueryUUID(qs url.Values, key string, v *validator.Validator) *uuid.UUID {
	s := qs.Get(key)

	if s == "" {
		return nil
	}

	id, err := uuid.Parse(s)
	if err != nil {
		v.AddError(key, fmt.Sprintf("unable to parse value: %s", err.Error()))
	}

	return &id
}

func ReadOptionalQueryString(qs url.Values, key string) *string {
	s := qs.Get(key)

	if s == "" {
		return nil
	}

	return &s
}

func ReadOptionalQueryDate(qs url.Values, key string, v *validator.Validator) *time.Time {
	s := qs.Get(key)
	if s == "" {
		return nil
	}

	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05",
	}

	for _, format := range formats {
		if date, err := time.Parse(format, s); err == nil {
			return &date
		}
	}

	v.AddError(key, fmt.Sprintf("not a valid date format, accepting %s", formats))

	return nil
}

func ReadJSON(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	return nil
}
