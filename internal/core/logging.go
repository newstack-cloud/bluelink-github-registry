package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CreateAppLogger creates a logger for application logs
// (debug, info, error etc.) that writes to the provided
// stdout and stderr targets.
// It will set the formatting to JSON if the environment
// is not "development".
func CreateAppLogger(
	stdoutTarget zapcore.WriteSyncer,
	stderrTarget zapcore.WriteSyncer,
	config *Config,
) (*zap.Logger, error) {

	zapLevel, err := zapcore.ParseLevel(config.LoggingLevel)
	if err != nil {
		return nil, err
	}

	zapConf := zap.NewDevelopmentEncoderConfig()
	if config.Environment == "production" {
		zapConf = zap.NewProductionEncoderConfig()
	}

	consoleErrors := zapcore.Lock(stderrTarget)
	consoleDebugging := zapcore.Lock(stdoutTarget)

	jsonEncoder := zapcore.NewJSONEncoder(zapConf)
	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, consoleDebugging, zapLevel),
		zapcore.NewCore(jsonEncoder, consoleErrors, zapLevel),
	)
	return zap.New(core), nil
}
