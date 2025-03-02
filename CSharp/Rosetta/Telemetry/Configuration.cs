namespace Rosetta.Telemetry;

public enum TelemetryOutput
{
    StdOut,
    Grpc,
    Http
}

public class TelemetryOptions
{
    public required string ServiceName { get; set; }
    public required string ServiceVersion { get; set; }
    public required string Url { get; set; }
    public required int Port { get; set; }
    public TelemetryOutput TelemetryOutput { get; set; }
}