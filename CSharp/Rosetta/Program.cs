using System.Globalization;
using OpenTelemetry.Logs;
using OpenTelemetry.Metrics;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using Serilog;


var logger = new LoggerConfiguration()
    .WriteTo.Console(new Serilog.Formatting.Json.JsonFormatter())
    .Enrich.FromLogContext()
    .CreateLogger();


var builder = WebApplication.CreateBuilder(args);

var serviceName = builder.Configuration["ServiceName"] ?? "Rosetta";
var version = builder.Configuration["Version"] ?? "0.0.0";

logger.ForContext("Version", version).Information("Starting up");
builder.Logging.AddOpenTelemetry(options =>
{
    options .SetResourceBuilder(
        ResourceBuilder.CreateDefault().AddService(serviceName)).AddConsoleExporter();
});
builder.Services.AddOpenTelemetry()
      .ConfigureResource(resource => resource.AddService(serviceName))
      .WithTracing(tracing => tracing
          .AddAspNetCoreInstrumentation()
          .AddConsoleExporter())
      .WithMetrics(metrics => metrics
          .AddAspNetCoreInstrumentation()
          .AddConsoleExporter());
builder.Logging.AddOpenTelemetry(options => options
    .SetResourceBuilder(ResourceBuilder.CreateDefault().AddService(
        serviceName: serviceName,
        serviceVersion: version))
    .AddConsoleExporter());
builder.Logging.ClearProviders();
builder.Logging.AddSerilog(logger);
builder.Services.AddOpenApi();
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSingleton(new Version(version));

var app = builder.Build();

app.MapOpenApi();
app.MapControllers();

app.UseHttpsRedirection();

app.Run();

public record Version(string version);