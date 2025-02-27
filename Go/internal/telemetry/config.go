package telemetry

type TelemetryOutput string

const (
	OutputStdOut TelemetryOutput = "stdout"
	OutputGRPC   TelemetryOutput = "grpc"
	OutputHTTP   TelemetryOutput = "http"
)

type TelemetryConfig struct {
	ServiceName    string          `json:"serviceName"`
	ServiceVersion string          `json:"serviceVersion"`
	Output         TelemetryOutput `json:"output"`
	URL            string          `json:"url"`
	Port           int             `json:"port"`
}
