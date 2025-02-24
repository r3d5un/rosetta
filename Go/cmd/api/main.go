package main

import (
	"context"
	"errors"
	"log/slog"
	_ "net/http/pprof"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/api"
	"github.com/r3d5un/rosetta/Go/internal/cfg"
	"github.com/r3d5un/rosetta/Go/internal/telemetry"
)

func main() {
	if err := run(); err != nil {
		slog.Error("an error occurred", "error", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config, err := cfg.New(ctx)
	if err != nil {
		return err
	}

	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler).With(
		slog.Group(
			"applicationInstance",
			slog.String("name", config.Name),
			slog.String("version", config.Version),
			slog.String("instanceId", uuid.New().String()),
		),
	)
	slog.SetDefault(logger)
	logger.Info("starting application", slog.Any("config", config))

	logger.Info("setting up telemetry", slog.Any("telemeteryConfig", config.Telemetry))
	shutdownTelemetry, err := telemetry.SetupTelemetry(
		ctx,
		config.Name,
		config.Version,
		config.Telemetry,
	)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, shutdownTelemetry(context.Background()))
	}()

	logger.Info("instantiating API")
	app := api.NewAPI(*logger)
	if err := app.Serve(); err != nil {
		logger.Error("unable to start server", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func otelShutdown(context context.Context) {
	panic("unimplemented")
}
