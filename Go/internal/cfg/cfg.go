package cfg

import (
	"context"
	"strings"

	"github.com/r3d5un/rosetta/Go/internal/database"
	"github.com/r3d5un/rosetta/Go/internal/telemetry"
	"github.com/spf13/viper"
)

type AppCfg struct {
	Name        string                    `json:"name"`
	Version     string                    `json:"version"`
	Environemnt string                    `json:"environment"`
	Server      ServerCfg                 `json:"server"`
	Telemetry   telemetry.TelemetryConfig `json:"telemetry"`
	Database    database.DatabaseConfig   `json:"database"`
}

type ServerCfg struct {
	Port int `json:"port"`
}

func New(ctx context.Context) (*AppCfg, error) {
	// Load configurations named 'cfg.yaml' from the given paths
	viper.SetConfigName("cfg")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/rosetta")

	// Check for environment variables with the ROSETTA prefix in uppercase.
	// ROSETTA_SERVER_PORT is equivalent to server.port
	viper.AutomaticEnv()
	viper.SetEnvPrefix("rosetta")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg AppCfg
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
