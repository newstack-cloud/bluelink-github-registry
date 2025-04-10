package registry

import (
	"errors"
	"net/http"

	"github.com/two-hundred/celerity-github-registry/internal/httputils"
	"github.com/two-hundred/celerity-github-registry/internal/plugins"
	"go.uber.org/zap"
)

func handlePluginError(
	w http.ResponseWriter,
	err error,
	logger *zap.Logger,
) {
	if errors.Is(err, plugins.ErrRepoNotFound) {
		httputils.HTTPError(
			w,
			http.StatusNotFound,
			"Plugin repository not found",
		)
		return
	}

	if errors.Is(err, plugins.ErrUnauthorised) {
		httputils.HTTPError(
			w,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	if errors.Is(err, plugins.ErrForbidden) {
		httputils.HTTPError(
			w,
			http.StatusForbidden,
			"Forbidden",
		)
		return
	}

	logger.Error(
		"Error retrieving plugin version information",
		zap.Error(err),
	)
	httputils.HTTPError(
		w,
		http.StatusInternalServerError,
		"An unexpected error occurred",
	)
}
