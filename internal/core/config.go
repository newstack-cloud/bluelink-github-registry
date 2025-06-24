package core

import "github.com/caarlos0/env/v11"

// Config holds the configuration for the github
// registry service.
type Config struct {
	Port                    int    `env:"BLUELINK_GITHUB_REGISTRY_PORT" envDefault:"8085"`
	AuthTokenHeader         string `env:"BLUELINK_GITHUB_REGISTRY_AUTH_TOKEN_HEADER" envDefault:"bluelink-gh-registry-token"`
	RegistryBaseURL         string `env:"BLUELINK_GITHUB_REGISTRY_BASE_URL"`
	PublicSigningKeysString string `env:"BLUELINK_GITHUB_REGISTRY_SIGNING_PUBLIC_KEYS"`
	HTTPClientTimeout       int    `env:"BLUELINK_GITHUB_REGISTRY_HTTP_CLIENT_TIMEOUT" envDefault:"60"`
	LoggingLevel            string `env:"BLUELINK_GITHUB_REGISTRY_LOGGING_LEVEL" envDefault:"info"`
	Environment             string `env:"BLUELINK_GITHUB_REGISTRY_ENVIRONMENT" envDefault:"production"`
	AccessLogFile           string `env:"BLUELINK_GITHUB_REGISTRY_ACCESS_LOG_FILE"`
	OutputLogFile           string `env:"BLUELINK_GITHUB_REGISTRY_OUTPUT_LOG_FILE"`
	ErrorLogFile            string `env:"BLUELINK_GITHUB_REGISTRY_ERROR_LOG_FILE"`
}

// LoadConfigFromEnv loads the application
// configuration from environment variables.
func LoadConfigFromEnv() (Config, error) {
	return env.ParseAs[Config]()
}
