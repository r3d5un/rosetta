namespace Rosetta.Telemetry;

using Microsoft.AspNetCore.Builder;
using OpenTelemetry.Logs;
using OpenTelemetry.Metrics;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using Serilog;

public class Telemetry
{
    private readonly Configuration _configuration;

    public Telemetry(Configuration configuration)
    {
        _configuration = configuration;
    }

    public void ConfigureBuilder(WebApplicationBuilder builder)
    {
        SetupLogger(builder);
        SetupTracing(builder);
        SetupMetrics(builder);
    }

    private void SetupLogger (WebApplicationBuilder builder)
    {
        var logger = new LoggerConfiguration()
            .WriteTo.Console(new Serilog.Formatting.Json.JsonFormatter())
            .Enrich.FromLogContext()
            .CreateLogger();

        switch (_configuration.TelemetryOutput)
        {
            case TelemetryOutput.Grpc:
            case TelemetryOutput.Http:
            case TelemetryOutput.StdOut:
            default:
                builder.Logging.AddOpenTelemetry(options =>
                {
                    options.SetResourceBuilder(
                        ResourceBuilder
                            .CreateDefault()
                            .AddService(
                                serviceName: _configuration.ServiceName,
                                serviceVersion: _configuration.ServiceVersion)
                    ).AddConsoleExporter();
                });
                break;
        }
    builder.Logging.ClearProviders();
    builder.Logging.AddSerilog(logger);
    }

    private void SetupTracing(WebApplicationBuilder builder)
    {
        switch (_configuration.TelemetryOutput)
        {
            case TelemetryOutput.Grpc:
            case TelemetryOutput.Http:
            case TelemetryOutput.StdOut:
            default:
                builder.Services.AddOpenTelemetry()
                      .ConfigureResource(resource => resource.AddService(_configuration.ServiceName))
                      .WithTracing(tracing => tracing
                          .AddAspNetCoreInstrumentation()
                          .AddConsoleExporter())
                      .WithMetrics(metrics => metrics
                          .AddAspNetCoreInstrumentation()
                          .AddConsoleExporter());
                break;
        }
    }

    private void SetupMetrics(WebApplicationBuilder builder)
    {
        switch (_configuration.TelemetryOutput)
        {
            case TelemetryOutput.Grpc:
            case TelemetryOutput.Http:
            case TelemetryOutput.StdOut:
            default:
                builder.Logging.AddOpenTelemetry(options => options
                    .SetResourceBuilder(ResourceBuilder.CreateDefault().AddService(
                        serviceName: _configuration.ServiceName,
                        serviceVersion: _configuration.ServiceVersion))
                    .AddConsoleExporter());
                break;
        }
    }
}