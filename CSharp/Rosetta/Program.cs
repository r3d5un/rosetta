using Microsoft.EntityFrameworkCore;
using Rosetta.Data;
using Rosetta.Telemetry;

var builder = WebApplication.CreateBuilder(args);

var version = builder.Configuration["Version"] ?? "0.0.0";

var telemetryOptions = builder.Configuration.GetSection("Telemetry").Get<TelemetryOptions>();
if (telemetryOptions != null)
{
    var telemetry = new Telemetry(telemetryOptions);
    telemetry.ConfigureBuilder(builder);
}
builder.Services.AddDbContextPool<RosettaDbContext>(options =>
{
    options.UseNpgsql(builder.Configuration.GetConnectionString("Database"));
});
builder.Services.AddOpenApi();
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSingleton(new Version(version));

var app = builder.Build();

using (var scope = app.Services.CreateScope())
{
    var logger = scope.ServiceProvider.GetRequiredService<ILogger<Program>>();
    var lifetime = scope.ServiceProvider.GetRequiredService<IHostApplicationLifetime>();
    
    logger.LogInformation("Starting up");
    
    logger.LogInformation("Testing database connection");
    var dbContext = scope.ServiceProvider.GetRequiredService<RosettaDbContext>();
    if (!dbContext.Database.CanConnect())
    {
        logger.LogInformation("Cannot connect to database. Shutting down.");
        lifetime.StopApplication();
        return;
    }
}

app.MapOpenApi();
app.MapControllers();

app.UseHttpsRedirection();

app.Run();

public record Version(string version);