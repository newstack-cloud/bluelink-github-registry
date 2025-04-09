package registry

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/two-hundred/celerity-github-registry/internal/core"
	"github.com/two-hundred/celerity-github-registry/internal/httputils"
)

// Manifest is the JSON manifest
// for service discovery that provides thecontent for the
// /.well-known/celerity-services.json document.
type Manifest struct {
	ProviderV1    string            `json:"provider.v1"`
	TransformerV1 string            `json:"transformer.v1"`
	AuthV1        *AuthManifestInfo `json:"auth.v1"`
}

// AuthManifestInfo is the authentication portion
// of the manifest.
// For this registry implementation, only the API Key
// authentication method is supported.
type AuthManifestInfo struct {
	APIKeyHeader string `json:"apiKeyHeader"`
}

func GetManifestHandler(config *core.Config) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			pluginsBaseURL := fmt.Sprintf("%s/plugins", config.RegistryBaseURL)
			manifest := &Manifest{
				ProviderV1:    pluginsBaseURL,
				TransformerV1: pluginsBaseURL,
				AuthV1: &AuthManifestInfo{
					APIKeyHeader: config.AuthTokenHeader,
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
