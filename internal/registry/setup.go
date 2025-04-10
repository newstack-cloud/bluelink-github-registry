package registry

import (
	"io"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/two-hundred/celerity-github-registry/internal/core"
	"github.com/two-hundred/celerity-github-registry/internal/plugins"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type registryDependencies struct {
	pluginService plugins.Service
}

type dependenciesRetriever func(
	config *core.Config,
	logger *zap.Logger,
) *registryDependencies

// Setup initialises and configures the http endpoint
// handlers for the registry.
// This returns the port the server should listen on,
// the access log writer, and any error that occurred
// during setup.
func Setup(
	router *mux.Router,
	getDeps dependenciesRetriever,
) (int, io.WriteCloser, error) {
	config, err := core.LoadConfigFromEnv()
	if err != nil {
		return 0, nil, err
	}

	accessLogWriter, err := getAccessLogWriter(&config)
	if err != nil {
		return 0, nil, err
	}

	outWriter, errWriter, err := getAppLogWriters(&config)
	if err != nil {
		return 0, nil, err
	}

	appLogger, err := core.CreateAppLogger(
		zapcore.AddSync(outWriter),
		zapcore.AddSync(errWriter),
		&config,
	)
	if err != nil {
		return 0, nil, err
	}

	deps := getDeps(&config, appLogger)

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
	// The auth token provided by the client will be passed through to make requests
	// to the underlying repositories, if those requests fail due to auth issues,
	// then those errors will be returned to the client.
	protocolRouter := router.PathPrefix("/plugins/").Subrouter()

	protocolRouter.Handle(
		"/{organisation}/{plugin}/versions",
		GetPluginVersionsHandler(&config, appLogger, deps.pluginService),
	).Methods("GET")

	protocolRouter.Handle(
		"/{organisation}/{plugin}/{version}/package/{os}/{arch}",
		GetPluginPackageHandler(&config, appLogger, deps.pluginService),
	).Methods("GET")

	return config.Port, accessLogWriter, nil
}

func getAccessLogWriter(config *core.Config) (io.WriteCloser, error) {
	if config.AccessLogFile != "" {
		return getLogWriter(config.AccessLogFile)
	}
	return &noCloseOutput{os.Stdout}, nil
}

func getLogWriter(logFile string) (io.WriteCloser, error) {
	return os.OpenFile(
		logFile,
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		0600,
	)
}

func getAppLogWriters(config *core.Config) (io.WriteCloser, io.WriteCloser, error) {
	var outWriter io.WriteCloser
	outWriter = &noCloseOutput{os.Stdout}
	var err error
	if config.OutputLogFile != "" {
		outWriter, err = getLogWriter(config.OutputLogFile)
		if err != nil {
			return nil, nil, err
		}
	}

	var errWriter io.WriteCloser
	errWriter = &noCloseOutput{os.Stderr}
	if config.ErrorLogFile != "" {
		errWriter, err = getLogWriter(config.ErrorLogFile)
		if err != nil {
			return nil, nil, err
		}
	}

	return outWriter, errWriter, nil
}

// noCloseOutput allows us to use the same abstraction for files we want
// to write access logs to and stdout without requiring branching logic
// to determine whether or not we should close the writer.
// This is important as we don't want to call close on os.Stdout.
type noCloseOutput struct {
	stdout *os.File
}

func (w *noCloseOutput) Write(p []byte) (n int, err error) {
	return w.stdout.Write(p)
}

func (w *noCloseOutput) Close() error {
	return nil
}
