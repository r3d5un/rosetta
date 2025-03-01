namespace Rosetta.Telemetry;

public class Configuration(
    string serviceName,
    string serviceVersion,
    TelemetryOutput telemetryOutput,
    string url,
    int port)
{
    public string ServiceName { get; set; } = serviceName;
    public string ServiceVersion { get; set; } = serviceVersion;
    public TelemetryOutput TelemetryOutput { get; set; } = telemetryOutput;
    public string Url { get; set; } = url;
    public int Port { get; set; } = port;
}

public enum TelemetryOutput
{
    StdOut,
    Grpc,
    Http
}