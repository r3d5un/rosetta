package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/justinas/alice"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type API struct {
	mux    *http.ServeMux
	logger slog.Logger
}

func NewAPI(logger slog.Logger) *API {
	return &API{
		mux:    http.NewServeMux(),
		logger: logger,
	}
}

func (api *API) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 4000),
		Handler:      api.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(api.logger.Handler(), slog.LevelError),
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		slog.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	api.logger.Info("starting server", slog.String("addr", srv.Addr))
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}
	api.logger.Info("stopped server", "addr", srv.Addr)

	return nil
}

func (api *API) routes() http.Handler {
	api.logger.Info("creating standard middleware chain")
	standard := alice.New(
		otelhttp.NewMiddleware("rosetta"),
		api.recoverPanic,
		api.enableCORS,
		api.logRequest,
	)

	endpoints := []struct {
		Path    string
		Handler http.HandlerFunc
	}{
		{"GET /api/v1/healthcheck", api.healthcheckHandler},
		// profiling
		{"GET /debug/pprof/", http.DefaultServeMux.ServeHTTP},
		{"GET /debug/pprof/profile", http.DefaultServeMux.ServeHTTP},
		{"GET /debug/pprof/heap", http.DefaultServeMux.ServeHTTP},
	}

	api.logger.Info("registering endpoints")
	for _, d := range endpoints {
		api.logger.Info("registering endpoint", slog.String("endpoint", d.Path))
		api.mux.Handle(d.Path, d.Handler)
	}

	handler := standard.Then(api.mux)
	return handler
}
