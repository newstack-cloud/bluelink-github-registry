package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/go-github/v70/github"
	"github.com/newstack-cloud/bluelink-github-registry/internal/httputils"
	"github.com/newstack-cloud/bluelink-github-registry/internal/signingkeys"
	"github.com/newstack-cloud/bluelink-github-registry/internal/types"
)

// ExtractPluginVersions extracts the plugin versions from the GitHub releases
// and returns them in a format that is compatible with the
// Bluelink registry protocol.
func ExtractPluginVersions(
	ctx context.Context,
	repository string,
	releases []*github.RepositoryRelease,
	client httputils.Client,
	token string,
) (*types.PluginVersions, error) {
	versions := []*types.PluginVersion{}
	for _, release := range releases {
		if !validTagPattern.MatchString(release.GetTagName()) {
			// Ignore releases that are not semantic versions prefixed with "v".
			continue
		}

		registryInfoURL := getRegistryInfoURL(release.Assets)
		registryInfo, err := getRegistryInfoFromURL(ctx, client, registryInfoURL, token)
		if err != nil {
			return nil, err
		}

		supportedPlatforms, err := extractSupportedPlatforms(repository, release)
		if err != nil {
			return nil, err
		}

		versions = append(versions, &types.PluginVersion{
			Version:            versionFromTag(release.GetTagName()),
			SupportedProtocols: registryInfo.SupportedProtocols,
			SupportedPlatforms: supportedPlatforms,
		})
	}

	return &types.PluginVersions{
		Versions: versions,
	}, nil
}

var (
	// A regex pattern that matches semantic versioning.
	// It matches versions like 1.0.0, 1.0.0-alpha, 1.0.0-beta, etc.
	// We prefix it with "v" as per expected in plugin releases for the registry.
	// The regexp is taken from the https://semver.org/ docs and modified to
	// include the "v" prefix.
	validTagPattern = regexp.MustCompile(
		`^v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`,
	)
)

func versionFromTag(tag string) string {
	// Tags are expected to be in the format "vX.Y.Z"
	// where X, Y, and Z are integers.
	return strings.TrimPrefix(tag, "v")
}

var (
	// A regex pattern that matches an archive name that dismisses
	// the plugin name and version to extract the os and architecture.
	platformPattern = regexp.MustCompile(
		`^.+_(linux|windows|darwin|freebsd)_(amd64|arm64|arm|386).zip$`,
	)
)

func extractSupportedPlatforms(
	repository string,
	release *github.RepositoryRelease,
) ([]*types.PluginVersionPlatform, error) {
	platforms := []*types.PluginVersionPlatform{}

	for _, asset := range release.Assets {
		isReleaseArchive, err := checkMatchesArchiveFileName(
			repository,
			release.GetTagName(),
			asset.GetName(),
		)
		if err != nil {
			return nil, err
		}

		if isReleaseArchive {
			matches := platformPattern.FindStringSubmatch(asset.GetName())
			if len(matches) == 3 {
				platforms = append(platforms, &types.PluginVersionPlatform{
					OS:   matches[1],
					Arch: matches[2],
				})
			}
		}
	}

	return platforms, nil
}

func checkMatchesArchiveFileName(
	repository string,
	tagName string,
	assetName string,
) (bool, error) {
	// The release must have at least one asset in the form:
	// <repo-name>_<version>_<platform>_<arch>.zip
	return regexp.MatchString(
		fmt.Sprintf(
			`^%s_%s_(linux|windows|darwin|freebsd)_(amd64|arm64|arm|386).zip$`,
			repository,
			strings.TrimPrefix(tagName, "v"),
		),
		assetName,
	)
}

func getRegistryInfoURL(
	assets []*github.ReleaseAsset,
) string {
	for _, asset := range assets {
		if strings.HasSuffix(asset.GetName(), "_registry_info.json") {
			return asset.GetURL()
		}
	}

	return ""
}

func getRegistryInfoFromURL(
	ctx context.Context,
	client httputils.Client,
	url string,
	token string,
) (*types.PluginRegistryInfo, error) {
	respBodyBytes, err := downloadFromGitHub(
		ctx,
		client,
		url,
		token,
	)
	if err != nil {
		return nil, err
	}

	var registryInfo types.PluginRegistryInfo
	err = json.Unmarshal(respBodyBytes, &registryInfo)
	if err != nil {
		return nil, err
	}

	return &registryInfo, nil
}

// FindPluginRepo searches for a plugin repository in the list of
// repositories based on the organisation and plugin name.
// It returns the first matching repository found or nil if none is found.
func FindPluginRepo(
	repositories []*github.Repository,
	organisation string,
	pluginName string,
) *github.Repository {
	pluginRepo := (*github.Repository)(nil)
	i := 0
	for pluginRepo == nil && i < len(repositories) {
		repo := repositories[i]
		candidateProviderRepo := RepoName(pluginName, "provider")
		candidateTransformerRepo := RepoName(pluginName, "transformer")

		if repo.GetName() == candidateProviderRepo ||
			repo.GetName() == candidateTransformerRepo {
			pluginRepo = repo
		}

		i += 1
	}

	return pluginRepo
}

// RepoName generates the repository name for a plugin
// based on the plugin name and type.
func RepoName(pluginName string, pluginType string) string {
	return fmt.Sprintf(
		"bluelink-%s-%s",
		pluginType,
		pluginName,
	)
}

// ExtractPluginVersionPackageParams holds the parameters needed to
// extract the plugin version package from a GitHub release.
type ExtractPluginVersionPackageParams struct {
	Repository            string
	Release               *github.RepositoryRelease
	Version               string
	OS                    string
	Arch                  string
	SigningKeysSerialised string
}

func ExtractPluginVersionPackage(
	ctx context.Context,
	params *ExtractPluginVersionPackageParams,
	client httputils.Client,
	token string,
) (*types.PluginVersionPackage, error) {
	pluginPackage := &types.PluginVersionPackage{
		OS:   params.OS,
		Arch: params.Arch,
	}

	registryInfoURL := getRegistryInfoURL(params.Release.Assets)
	registryInfo, err := getRegistryInfoFromURL(
		ctx,
		client,
		registryInfoURL,
		token,
	)
	if err != nil {
		return nil, err
	}
	pluginPackage.SupportedProtocols = registryInfo.SupportedProtocols
	pluginPackage.Dependencies = registryInfo.Dependencies

	attachReleaseFileInfo(
		params.Repository,
		params.Version,
		params.Release,
		pluginPackage,
	)

	signingKeys, err := prepareSigningKeysForPackageInfo(
		params.SigningKeysSerialised,
	)
	if err != nil {
		return nil, err
	}
	pluginPackage.SigningKeys = signingKeys

	shasum, err := getSHASumFromURL(
		ctx,
		client,
		pluginPackage.SHASumsURL,
		pluginPackage.Filename,
		token,
	)
	if err != nil {
		return nil, err
	}
	pluginPackage.SHASum = shasum

	return pluginPackage, nil
}

func prepareSigningKeysForPackageInfo(
	signingKeysSerialised string,
) (*types.PublicGPGSigningKeys, error) {
	if signingKeysSerialised == "" {
		return nil, errors.New("empty signing keys provided")
	}

	var signingKeys types.IntermediarySigningKeys
	err := json.Unmarshal(
		[]byte(signingKeysSerialised),
		&signingKeys,
	)
	if err != nil {
		return nil, err
	}

	publicSigningKeys := &types.PublicGPGSigningKeys{}

	for _, key := range signingKeys.Keys {
		publicKey := &types.PublicGPGSigningKey{
			PublicKey: key.PublicKey,
		}
		hexKeyID, err := signingkeys.ExtractHexKeyID(key.PublicKey)
		if err != nil {
			return nil, err
		}
		publicKey.HexKeyID = hexKeyID
		publicSigningKeys.GPG = append(publicSigningKeys.GPG, publicKey)
	}

	return publicSigningKeys, nil
}

func attachReleaseFileInfo(
	repository string,
	version string,
	release *github.RepositoryRelease,
	versionPackage *types.PluginVersionPackage,
) {
	archive := fmt.Sprintf(
		"%s_%s_%s_%s.zip",
		repository,
		version,
		versionPackage.OS,
		versionPackage.Arch,
	)
	shasumsFile := fmt.Sprintf(
		"%s_%s_SHA256SUMS",
		repository,
		version,
	)
	shasumsSignatureFile := fmt.Sprintf(
		"%s_%s_SHA256SUMS.sig",
		repository,
		version,
	)

	for _, asset := range release.Assets {
		if asset.GetName() == archive {
			versionPackage.Filename = archive
			versionPackage.DownloadURL = asset.GetURL()
		}

		if asset.GetName() == shasumsFile {
			versionPackage.SHASumsURL = asset.GetURL()
		}

		if asset.GetName() == shasumsSignatureFile {
			versionPackage.SHASumsSignatureURL = asset.GetURL()
		}
	}
}

func getSHASumFromURL(
	ctx context.Context,
	client httputils.Client,
	shasumsURL string,
	archiveFilename string,
	token string,
) (string, error) {
	shasumBytes, err := downloadFromGitHub(
		ctx,
		client,
		shasumsURL,
		token,
	)
	if err != nil {
		return "", err
	}

	shasumLines := strings.Split(string(shasumBytes), "\n")
	for _, line := range shasumLines {
		if strings.Contains(line, archiveFilename) {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				return parts[0], nil
			}
		}
	}

	return "", fmt.Errorf(
		"failed to find shasum for archive file: %s",
		archiveFilename,
	)
}

func downloadFromGitHub(
	ctx context.Context,
	client httputils.Client,
	url string,
	token string,
) ([]byte, error) {
	if url == "" {
		return nil, nil
	}

	req, err := createGitHubDownloadRequest(
		ctx,
		url,
		token,
	)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"failed to fetch from url %q: status code: %s",
			url,
			resp.Status,
		)
	}

	return io.ReadAll(resp.Body)
}

func createGitHubDownloadRequest(
	ctx context.Context,
	url string,
	token string,
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	return req, nil
}
