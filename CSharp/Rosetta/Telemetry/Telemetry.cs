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

    public WebApplicationBuilder Configure(WebApplicationBuilder builder)
    {
        if (builder == null)
        {
            throw new ArgumentNullException(nameof(builder));
        }
        
        var logger = new LoggerConfiguration()
            .WriteTo.Console(new Serilog.Formatting.Json.JsonFormatter())
            .Enrich.FromLogContext()
            .CreateLogger();
        builder.Logging.AddOpenTelemetry(options =>
        {
            options .SetResourceBuilder(
                ResourceBuilder
                    .CreateDefault()
                    .AddService(_configuration.ServiceName))
                .AddConsoleExporter();
        });
        builder.Logging.ClearProviders();
        builder.Logging.AddSerilog(logger);
        
        builder.Services.AddOpenTelemetry()
              .ConfigureResource(resource => resource.AddService(_configuration.ServiceName))
              .WithTracing(tracing => tracing
                  .AddAspNetCoreInstrumentation()
                  .AddConsoleExporter())
              .WithMetrics(metrics => metrics
                  .AddAspNetCoreInstrumentation()
                  .AddConsoleExporter());
        
        builder.Logging.AddOpenTelemetry(options => options
            .SetResourceBuilder(ResourceBuilder.CreateDefault().AddService(
                serviceName: _configuration.ServiceName,
                serviceVersion: _configuration.ServiceVersion))
            .AddConsoleExporter());
        
        return builder;
    }
}