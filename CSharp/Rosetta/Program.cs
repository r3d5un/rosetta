using Serilog;


var logger = new LoggerConfiguration()
    .WriteTo.Console(new Serilog.Formatting.Json.JsonFormatter())
    .Enrich.FromLogContext()
    .CreateLogger();


var builder = WebApplication.CreateBuilder(args);
builder.Logging.ClearProviders();
builder.Logging.AddSerilog(logger);
builder.Services.AddOpenApi();
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
var version = builder.Configuration["Version"] ?? "0.0.0";
builder.Services.AddSingleton(new Version(version));
logger.ForContext("Version", version).Information("Starting up");

var app = builder.Build();

app.MapOpenApi();
app.MapControllers();

app.UseHttpsRedirection();

app.Run();

public record Version(string version);