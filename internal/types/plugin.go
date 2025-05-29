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
	SupportedProtocols []string          `json:"supportedProtocols"`
	Dependencies       map[string]string `json:"dependencies,omitempty"`
}

// PluginVersionPackage holds the information about
// the package information for a plugin version
// that is published with a plugin version release.
// This represents the structure of a plugin version package
// expected by the Celerity registry protocol.
type PluginVersionPackage struct {
	SupportedProtocols  []string              `json:"supportedProtocols"`
	OS                  string                `json:"os"`
	Arch                string                `json:"arch"`
	Filename            string                `json:"filename"`
	DownloadURL         string                `json:"downloadUrl"`
	SHASumsURL          string                `json:"shasumsUrl"`
	SHASumsSignatureURL string                `json:"shasumsSignatureUrl"`
	SHASum              string                `json:"shasum"`
	SigningKeys         *PublicGPGSigningKeys `json:"signingKeys"`
	Dependencies        map[string]string     `json:"dependencies,omitempty"`
}

// PublicGPGSigningKeys holds the information about
// the public GPG signing keys for a plugin version
// that is published with a plugin version release.
type PublicGPGSigningKeys struct {
	GPG []*PublicGPGSigningKey `json:"gpg"`
}

// PublicGPGSigningKey holds the information about
// a public GPG signing key that can be used to verify
// the authenticity of a plugin version package.
type PublicGPGSigningKey struct {
	HexKeyID string `json:"keyId"`
	// The ASCII armored public key
	// that can be used to verify the signature
	// of the plugin version package.
	PublicKey string `json:"publicKey"`
}
