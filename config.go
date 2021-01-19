package server

import (
	"os"

	log "github.com/towl/logger"
)

// Config is the configuration structure
type Config struct {
	LogLevel   string
	Host       string
	Port       string
	WorkingDir string
	APIPattern string
}

var config = &Config{}
var logger = log.GetLoggerFromEnv("API_SERVER_", false)

func init() {
	config.loadEnv()
	logger.Info("Server config loaded successfully.")
}

func (c *Config) loadEnv() {
	c.Host = os.Getenv("API_SERVER_HOST")
	c.Port = os.Getenv("API_SERVER_PORT")
	c.WorkingDir = os.Getenv("API_SERVER_WORKING_DIR")
	c.APIPattern = os.Getenv("API_SERVER_PATTERN")
}
