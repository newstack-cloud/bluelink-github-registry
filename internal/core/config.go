package core

import "github.com/caarlos0/env/v11"

// Config holds the configuration for the github
// registry service.
type Config struct {
	Port            int    `env:"CELERITY_GITHUB_REGISTRY_PORT" envDefault:"8085"`
	AccessLogFile   string `env:"CELERITY_GITHUB_REGISTRY_ACCESS_LOG_FILE"`
	AuthTokenHeader string `env:"CELERITY_GITHUB_REGISTRY_AUTH_TOKEN_HEADER" envDefault:"celerity-gh-registry-token"`
	RegistryBaseURL string `env:"CELERITY_GITHUB_REGISTRY_BASE_URL"`
}

// LoadConfigFromEnv loads the application
// configuration from environment variables.
func LoadConfigFromEnv() (Config, error) {
	return env.ParseAs[Config]()
}
