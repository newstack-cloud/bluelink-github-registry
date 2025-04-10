package types

// PluginVersions holds the information about the plugin versions
// that are available for a given plugin.
type PluginVersions struct {
	Versions []*PluginVersion `json:"versions"`
}

// PluginVersion holds information about a plugin version
// that is mostly useful for allowing the client to determine
// the correct plugin version to use for the current protocol and platform.
type PluginVersion struct {
	Version            string                   `json:"version"`
	SupportedProtocols []string                 `json:"supportedProtocols"`
	SupportedPlatforms []*PluginVersionPlatform `json:"supportedPlatforms"`
}

// PluginVersionPlatform holds the information about
// the supported OS and architectures for a plugin version.
type PluginVersionPlatform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

// PluginRegistryInfo holds the information about
// the registry information for a plugin version
// that is published with a plugin version release.
type PluginRegistryInfo struct {
	SupportedProtocols []string `json:"supportedProtocols"`
}
