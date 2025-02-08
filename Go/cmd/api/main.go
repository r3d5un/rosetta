package main

import (
	"log/slog"
	_ "net/http/pprof"
	"os"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/api"
)

const version = "0.0.1"

func main() {
	if err := run(); err != nil {
		slog.Error("an error occurred", "error", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler).With(
		slog.Group(
			"applicationInstance",
			slog.String("version", version),
			slog.String("instanceId", uuid.New().String()),
		),
	)
	slog.SetDefault(logger)

	app := api.NewAPI(*logger)
	if err := app.Serve(); err != nil {
		logger.Error("unable to start server", slog.String("error", err.Error()))
		return err
	}

	return nil
}
