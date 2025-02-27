package telemetry

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/r3d5un/rosetta/Go/internal/cfg"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// SetupTelemetry performs the setup of metrics, tracing and logging, registering
// each as globally available to the application.
func SetupTelemetry(
	ctx context.Context,
	serviceName string,
	serviceVersion string,
	config cfg.TelemetryCfg,
) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	resource := newResource(serviceName, serviceVersion)

	propagator := newPropagator()
	otel.SetTextMapPropagator(propagator)

	traceProvider, err := newTracerProvider(ctx, resource, config)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, traceProvider.Shutdown)
	otel.SetTracerProvider(traceProvider)

	meterProvider, err := newMeterProvider(ctx, resource, config)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	loggerProvider, err := newLoggerProvider(ctx, resource, config)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)
	logger := otelslog.NewLogger(serviceName).With(
		slog.Group(
			"applicationInstance",
			slog.String("name", serviceName),
			slog.String("version", serviceVersion),
			slog.String("instanceId", uuid.New().String()),
		),
	)
	slog.SetDefault(logger)

	return shutdown, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newLoggerProvider(
	ctx context.Context,
	res *resource.Resource,
	config cfg.TelemetryCfg,
) (*log.LoggerProvider, error) {
	var processor *log.BatchProcessor
	var exporter log.Exporter
	var err error

	switch config.Output {
	case cfg.HTTP:
		exporter, err = otlploghttp.New(
			ctx,
			otlploghttp.WithEndpoint(fmt.Sprintf("%s:%d", config.URL, config.Port)),
		)
	case cfg.GRPC:
		exporter, err = otlploggrpc.New(
			ctx,
			otlploggrpc.WithEndpoint(fmt.Sprintf("%s:%d", config.URL, config.Port)),
		)
	default:
		exporter, err = stdoutlog.New()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create OLTP log exporter: %w", err)
	}
	processor = log.NewBatchProcessor(exporter)

	return log.NewLoggerProvider(log.WithProcessor(processor), log.WithResource(res)), nil
}

func newMeterProvider(
	ctx context.Context,
	res *resource.Resource,
	config cfg.TelemetryCfg,
) (*metric.MeterProvider, error) {
	var mp *metric.MeterProvider
	var exporter metric.Exporter
	var err error

	switch config.Output {
	case cfg.HTTP:
		exporter, err = otlpmetrichttp.New(
			ctx,
			otlpmetrichttp.WithEndpoint(fmt.Sprintf("%s:%d", config.URL, config.Port)),
		)
	case cfg.GRPC:
		exporter, err = otlpmetricgrpc.New(
			ctx,
			otlpmetricgrpc.WithEndpoint(fmt.Sprintf("%s:%d", config.URL, config.Port)),
		)
	default:
		exporter, err = stdoutmetric.New()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP metric exporter: %w", err)
	}
	mp = metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(mp)

	return mp, nil
}

func newTracerProvider(
	ctx context.Context,
	res *resource.Resource,
	config cfg.TelemetryCfg,
) (*trace.TracerProvider, error) {
	var tp *trace.TracerProvider
	var exporter trace.SpanExporter
	var err error

	switch config.Output {
	case cfg.HTTP:
		exporter, err = otlptracehttp.New(
			ctx,
			otlptracehttp.WithEndpoint(fmt.Sprintf("%s:%d", config.URL, config.Port)),
		)
	case cfg.GRPC:
		exporter, err = otlptracegrpc.New(
			ctx,
			otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%d", config.URL, config.Port)),
		)
	default:
		exporter, err = stdouttrace.New()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	tp = trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}

func newResource(serviceName string, serviceVersion string) *resource.Resource {
	hostName, _ := os.Hostname()

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(serviceVersion),
		semconv.HostName(hostName),
	)
}
