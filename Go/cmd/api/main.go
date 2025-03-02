package main

import (
	"context"
	"errors"
	"log/slog"
	_ "net/http/pprof"
	"os"
	"os/signal"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/r3d5un/rosetta/Go/internal/api"
	"github.com/r3d5un/rosetta/Go/internal/cfg"
	"github.com/r3d5un/rosetta/Go/internal/telemetry"
)

func main() {
	if err := run(); err != nil {
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

	shutdownTelemetry, err := telemetry.SetupTelemetry(ctx, config.Telemetry)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, shutdownTelemetry(context.Background()))
	}()
	logger := slog.Default()
	logger.Info("starting application", slog.Any("config", config))

	logger.Info("instantiating API")
	app, err := api.NewAPI(ctx, *config)
	if err != nil {
		logger.Error("unable to start API", slog.String("error", err.Error()))
		return err
	}
	if err := app.Serve(); err != nil {
		logger.Error("unable to start server", slog.String("error", err.Error()))
		return err
	}

	logger.Info("shutting down...")
	return nil
}
