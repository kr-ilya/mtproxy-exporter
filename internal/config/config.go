package config

import (
	"flag"
	"os"
	"sync"
	"time"
)

var (
	once sync.Once
	cfg  *Config
)

// Config holds the configuration for the exporter
type Config struct {
	ListenAddress  string
	MTProxyURL     string
	MTProxyTimeout time.Duration
	LogLevel       string
}

// Load loads configuration from environment variables and command-line flags
// This function is safe to call multiple times, but will only parse flags once
func Load() *Config {
	once.Do(func() {
		cfg = &Config{}

		// Define flags
		flag.StringVar(&cfg.ListenAddress, "listen.address", getEnv("LISTEN_ADDRESS", ":9330"),
			"Address to listen on for metrics and health endpoints")
		flag.StringVar(&cfg.MTProxyURL, "mtproxy.url", getEnv("MTPROXY_URL", "http://localhost:8888"),
			"URL of the MTProxy stats endpoint")
		flag.DurationVar(&cfg.MTProxyTimeout, "mtproxy.timeout", getDurationEnv("MTPROXY_TIMEOUT", 10*time.Second),
			"Timeout for MTProxy HTTP requests")
		flag.StringVar(&cfg.LogLevel, "log.level", getEnv("LOG_LEVEL", "info"),
			"Log level (debug, info, warn, error)")

		flag.Parse()
	})

	return cfg
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDurationEnv gets a duration environment variable or returns a default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
