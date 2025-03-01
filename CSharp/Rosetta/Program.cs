using Rosetta.Telemetry;

var builder = WebApplication.CreateBuilder(args);

var serviceName = builder.Configuration["ServiceName"] ?? "Rosetta";
var version = builder.Configuration["Version"] ?? "0.0.0";

var telemetry = new Telemetry(
    new Configuration(
        serviceName,
        version,
        TelemetryOutput.StdOut,
        string.Empty,
        0
    ));
builder = telemetry.Configure(builder);
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