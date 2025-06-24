package registry

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/newstack-cloud/bluelink-github-registry/internal/core"
	"github.com/newstack-cloud/bluelink-github-registry/internal/httputils"
	"github.com/newstack-cloud/bluelink-github-registry/internal/plugins"
	"go.uber.org/zap"
)

func GetPluginVersionsHandler(
	config *core.Config,
	logger *zap.Logger,
	pluginService plugins.Service,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			token := req.Header.Get(config.AuthTokenHeader)
			if strings.TrimSpace(token) == "" {
				httputils.HTTPError(
					w,
					http.StatusUnauthorized,
					"Unauthorized",
				)
				return
			}

			params := mux.Vars(req)
			organisation := params["organisation"]
			plugin := params["plugin"]

			pluginVersions, err := pluginService.ListVersions(
				req.Context(),
				organisation,
				plugin,
				token,
			)
			if err != nil {
				handlePluginError(
					w,
					err,
					logger,
				)
				return
			}

			respBytes, err := json.Marshal(pluginVersions)
			if err != nil {
				logger.Error(
					"Error marshalling plugin version information",
					zap.Error(err),
				)
				httputils.HTTPError(
					w,
					http.StatusInternalServerError,
					"An unexpected error occurred",
				)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(respBytes)
		},
	)
}
