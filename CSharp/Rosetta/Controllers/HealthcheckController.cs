using Microsoft.AspNetCore.Mvc;

namespace Rosetta.Controllers;

[Route("/api/v1/healthcheck"), ApiController]
public class HealthcheckController : Controller
{
    private readonly ILogger<HealthcheckController> _logger;
    private readonly Version _version;

    public HealthcheckController(ILogger<HealthcheckController> logger, Version version)
    {
        _logger = logger;
        _version = version;
    }

    [HttpGet]
    public Task<ActionResult> Get()
    {
        _logger.LogInformation("Rosetta Version {Version}", _version);
        return Task.FromResult<ActionResult>(Ok());
    }
}
