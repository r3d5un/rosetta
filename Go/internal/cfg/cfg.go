package cfg

import (
	"context"
	"log/slog"
	"strings"

	"github.com/r3d5un/rosetta/Go/internal/database"
	"github.com/r3d5un/rosetta/Go/internal/logging"
	"github.com/r3d5un/rosetta/Go/internal/telemetry"
	"github.com/spf13/viper"
)

type AppCfg struct {
	Name             string                    `json:"name"`
	Version          string                    `json:"version"`
	Environemnt      string                    `json:"environment"`
	Server           ServerCfg                 `json:"server"`
	TelemetryEnabled bool                      `json:"telemetryEnabled"`
	Telemetry        telemetry.TelemetryConfig `json:"telemetry"`
	Database         database.DatabaseConfig   `json:"database"`
}

type ServerCfg struct {
	Port int `json:"port"`
}

func New(ctx context.Context) (*AppCfg, error) {
	logger := logging.LoggerFromContext(ctx)

	// Load configurations named 'cfg.yaml' from the given paths
	logger.LogAttrs(ctx, slog.LevelInfo, "setting configuration file paths")
	viper.SetConfigName("cfg")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/rosetta")

	// Check for environment variables with the ROSETTA prefix in uppercase.
	// ROSETTA_SERVER_PORT is equivalent to server.port
	logger.LogAttrs(ctx, slog.LevelInfo, "reading environment variables")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("rosetta")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	logger.LogAttrs(ctx, slog.LevelInfo, "loading configuration")
	err := viper.ReadInConfig()
	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"unable to load configuration file",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	var cfg AppCfg
	err = viper.Unmarshal(&cfg)
	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"unmarshalling configuration",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &cfg, nil
}
