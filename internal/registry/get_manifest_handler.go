package registry

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/newstack-cloud/bluelink-github-registry/internal/core"
	"github.com/newstack-cloud/bluelink-github-registry/internal/httputils"
)

// Manifest is the JSON manifest
// for service discovery that provides thecontent for the
// /.well-known/bluelink-services.json document.
type Manifest struct {
	ProviderV1    *PluginTypeManifestInfo `json:"provider.v1"`
	TransformerV1 *PluginTypeManifestInfo `json:"transformer.v1"`
	AuthV1        *AuthManifestInfo       `json:"auth.v1"`
}

// PluginTypeManifestInfo is the plugin type
// portion of the manifest.
type PluginTypeManifestInfo struct {
	Endpoint string `json:"endpoint"`
	// The content type that the client should specify in the `Accept` header
	// when downloading the plugin and related artifacts.
	DownloadAcceptContentType string `json:"downloadAcceptContentType"`
}

// AuthManifestInfo is the authentication portion
// of the manifest.
// For this registry implementation, only the API Key
// authentication method is supported.
type AuthManifestInfo struct {
	APIKeyHeader string `json:"apiKeyHeader"`
	DownloadAuth string `json:"downloadAuth,omitempty"`
}

const (
	// DownloadContentType is the content type
	// that the client should specify in the `Accept` header
	// when downloading the plugin and related artifacts.
	// This is required for a GitHub registry due to the way
	// GitHub serves files from private repositories.
	DownloadContentType = "application/octet-stream"
)

func GetManifestHandler(config *core.Config) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			pluginsBaseURL := fmt.Sprintf("%s/plugins", config.RegistryBaseURL)
			manifest := &Manifest{
				ProviderV1: &PluginTypeManifestInfo{
					Endpoint:                  pluginsBaseURL,
					DownloadAcceptContentType: DownloadContentType,
				},
				TransformerV1: &PluginTypeManifestInfo{
					Endpoint:                  pluginsBaseURL,
					DownloadAcceptContentType: DownloadContentType,
				},
				AuthV1: &AuthManifestInfo{
					APIKeyHeader: config.AuthTokenHeader,
					// GitHub expects the `Authorization: Bearer <token>` header
					// to be set when downloading artifacts from private repositories.
					DownloadAuth: "bearer",
				},
			}

			manifestBytes, err := json.Marshal(manifest)
			if err != nil {
				httputils.HTTPError(
					w,
					http.StatusInternalServerError,
					"Failed to marshal manifest",
				)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(manifestBytes)
		},
	)
}
