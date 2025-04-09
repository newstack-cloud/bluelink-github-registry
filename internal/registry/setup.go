package registry

import (
	"io"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/two-hundred/celerity-github-registry/internal/auth"
	"github.com/two-hundred/celerity-github-registry/internal/core"
)

// Setup initialises and configures the http endpoint
// handlers for the registry.
// This returns the port the server should listen on,
// the access log writer, and any error that occurred
// during setup.
func Setup(router *mux.Router) (int, io.WriteCloser, error) {
	config, err := core.LoadConfigFromEnv()
	if err != nil {
		return 0, nil, err
	}

	accessLogWriter, err := getAccessLogWriter(&config)
	if err != nil {
		return 0, nil, err
	}

	// Writes access logs to the io.Writer in the Apache Combined Log Format.
	router.Use(func(next http.Handler) http.Handler {
		return handlers.CombinedLoggingHandler(accessLogWriter, next)
	})

	// We need to serve a manifest for service discovery
	// as per the Service Discovery protoocol used by the Celerity CLI
	// and other clients.
	router.Handle(
		"/.well-known/celerity-services.json",
		GetManifestHandler(&config),
	).Methods("GET")

	// The registry protocol endpoints come under the "/plugins/" path prefix.
	protectedRouter := router.PathPrefix("/plugins").Subrouter()
	// Registry protocol endpoints should be protected by GitHub PAT authentication.
	protectedRouter.Use(
		auth.TokenMiddleware(
			config.AuthTokenHeader,
			auth.NewGitHubTokenService(),
		),
	)

	return config.Port, accessLogWriter, nil
}

func getAccessLogWriter(config *core.Config) (io.WriteCloser, error) {
	if config.AccessLogFile != "" {
		return os.OpenFile(
			config.AccessLogFile,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND,
			0600,
		)
	}
	return &noCloseStdOut{os.Stdout}, nil
}

// noCloseStdOut allows us to use the same abstraction for files we want
// to write access logs to and stdout without requiring branching logic
// to determine whether or not we should close the writer.
// This is important as we don't want to call close on os.Stdout.
type noCloseStdOut struct {
	stdout *os.File
}

func (w *noCloseStdOut) Write(p []byte) (n int, err error) {
	return w.stdout.Write(p)
}

func (w *noCloseStdOut) Close() error {
	return nil
}
