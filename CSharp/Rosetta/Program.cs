using Rosetta.Telemetry;

var builder = WebApplication.CreateBuilder(args);

var version = builder.Configuration["Version"] ?? "0.0.0";

var telemetryOptions = builder.Configuration.GetSection("Telemetry").Get<TelemetryOptions>();
if (telemetryOptions != null)
{
    var telemetry = new Telemetry(telemetryOptions);
    telemetry.ConfigureBuilder(builder);
}

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