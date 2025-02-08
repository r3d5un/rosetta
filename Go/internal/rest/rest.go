package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/r3d5un/rosetta/Go/internal/logging"
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
