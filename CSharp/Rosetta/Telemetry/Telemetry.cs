using OpenTelemetry.Exporter;
using OpenTelemetry.Logs;
using OpenTelemetry.Metrics;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using Serilog;

namespace Rosetta.Telemetry;


public class Telemetry(TelemetryOptions configuration)
{
    public void ConfigureBuilder(WebApplicationBuilder builder)
    {
        SetupLogger(builder);
        SetupTelemetry(builder);
    }

    private void SetupLogger (WebApplicationBuilder builder)
    {
        var logger = new LoggerConfiguration()
            .WriteTo.Console(new Serilog.Formatting.Json.JsonFormatter())
            .Enrich.FromLogContext()
            .CreateLogger();

        builder.Logging.AddOpenTelemetry(options =>
        {
            options.SetResourceBuilder(
                ResourceBuilder
                    .CreateDefault()
                    .AddService(
                        serviceName: configuration.ServiceName,
                        serviceVersion: configuration.ServiceVersion));

            switch (configuration.TelemetryOutput)
            {
                case TelemetryOutput.Grpc:
                    options.AddOtlpExporter(exporterOptions =>
                    {
                        exporterOptions.Protocol = OtlpExportProtocol.Grpc;
                        exporterOptions.Endpoint = new Uri($"{configuration.Url}:{configuration.Port}");
                    });
                    break;
                case TelemetryOutput.Http:
                    options.AddOtlpExporter(exporterOptions =>
                    {
                        exporterOptions.Protocol = OtlpExportProtocol.HttpProtobuf;
                        exporterOptions.Endpoint = new Uri($"{configuration.Url}:{configuration.Port}");
                    });
                    break;
                case TelemetryOutput.StdOut:
                default:
                    options.AddConsoleExporter();
                    break;
            }
        });
        builder.Logging.AddSerilog(logger);
    }

    private void SetupTelemetry(WebApplicationBuilder builder)
    {
        builder.Services.AddOpenTelemetry()
            .ConfigureResource(resource => resource.AddService(configuration.ServiceName))
            .WithTracing(tracing =>
            {
                tracing.AddAspNetCoreInstrumentation();
                switch (configuration.TelemetryOutput)
                {
                    case TelemetryOutput.Grpc:
                        tracing.AddOtlpExporter(options =>
                        {
                            options.Protocol = OtlpExportProtocol.Grpc;
                            options.Endpoint = new Uri($"{configuration.Url}:{configuration.Port}");
                        });
                        break;
                    case TelemetryOutput.Http:
                        tracing.AddOtlpExporter(options =>
                        {
                            options.Protocol = OtlpExportProtocol.HttpProtobuf;
                            options.Endpoint = new Uri($"{configuration.Url}:{configuration.Port}");
                        });
                        break;
                    case TelemetryOutput.StdOut:
                    default:
                        tracing.AddConsoleExporter();
                        break;
                }
            })
            .WithMetrics(metrics =>
            {
                metrics.AddAspNetCoreInstrumentation();
                switch (configuration.TelemetryOutput)
                {
                    case TelemetryOutput.Grpc:
                        metrics.AddOtlpExporter(options =>
                        {
                            options.Protocol = OtlpExportProtocol.Grpc;
                            options.Endpoint = new Uri($"{configuration.Url}:{configuration.Port}");
                        });
                        break;
                    case TelemetryOutput.Http:
                        metrics.AddOtlpExporter(options =>
                        {
                            options.Protocol = OtlpExportProtocol.HttpProtobuf;
                            options.Endpoint = new Uri($"{configuration.Url}:{configuration.Port}");
                        });
                        break;
                    case TelemetryOutput.StdOut:
                    default:
                        metrics.AddConsoleExporter();
                        break;
                }
            });
    }
}