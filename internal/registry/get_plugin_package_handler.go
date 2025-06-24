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

func GetPluginPackageHandler(
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
			version := params["version"]
			os := params["os"]
			arch := params["arch"]

			packageInfo, err := pluginService.GetPackageInfo(
				req.Context(),
				&plugins.PackageInfoParams{
					Organisation: organisation,
					Plugin:       plugin,
					Version:      version,
					OS:           os,
					Arch:         arch,
				},
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

			respBytes, err := json.Marshal(packageInfo)
			if err != nil {
				logger.Error(
					"Error marshalling plugin version package information",
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
