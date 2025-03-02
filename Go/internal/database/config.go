package database

import "time"

type DatabaseConfig struct {
	ConnStr         string `json:"-"`
	MaxOpenConns    int32  `json:"maxOpenConns"`
	IdleTimeMinutes int    `json:"idleTimeMinutes"`
	TimeoutSeconds  int    `json:"timeoutSeconds"`
}

func (c *DatabaseConfig) TimeoutDuration() time.Duration {
	return time.Duration(c.TimeoutSeconds) * time.Second
}

func (c *DatabaseConfig) IdleTime() time.Duration {
	return time.Duration(c.IdleTimeMinutes) * time.Minute
}
