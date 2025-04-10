package testutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/two-hundred/celerity-github-registry/internal/signingkeys"
	"github.com/two-hundred/celerity-github-registry/internal/types"
)

type SigningKeysInfo struct {
	Expected *types.PublicGPGSigningKeys
	Input    string
}

// GetSigningKeysFromEnv retrieves the signing keys from the environment variable
// CELERITY_GITHUB_REGISTRY_SIGNING_PUBLIC_KEYS and unmarshals them into a
// PublicGPGSigningKeys struct.
// This returns an expected output set of GPG public keys that
// should be produced and the input string.
func GetSigningKeysFromEnv() (*SigningKeysInfo, error) {
	signingKeysSerialised := os.Getenv("CELERITY_GITHUB_REGISTRY_SIGNING_PUBLIC_KEYS")
	if signingKeysSerialised == "" {
		return nil, errors.New("no signing keys found in env")
	}

	fmt.Println("Serialised stuff here:", signingKeysSerialised)
	fmt.Println(
		"Contains ASCII armor header?",
		strings.Contains(signingKeysSerialised, "-----BEGIN PGP PUBLIC KEY BLOCK-----"),
	)
	fmt.Println(
		"Contains ASCII armor footer?",
		strings.Contains(signingKeysSerialised, "-----END PGP PUBLIC KEY BLOCK-----"),
	)

	signingKeysInternal := &types.IntermediarySigningKeys{}
	err := json.Unmarshal([]byte(signingKeysSerialised), signingKeysInternal)
	if err != nil {
		return nil, err
	}

	expected := &types.PublicGPGSigningKeys{
		GPG: []*types.PublicGPGSigningKey{},
	}

	for _, key := range signingKeysInternal.Keys {
		keyID, err := signingkeys.ExtractHexKeyID(key.PublicKey)
		if err != nil {
			return nil, err
		}

		expected.GPG = append(expected.GPG, &types.PublicGPGSigningKey{
			HexKeyID:  keyID,
			PublicKey: key.PublicKey,
		})
	}

	return &SigningKeysInfo{
		Expected: expected,
		Input:    signingKeysSerialised,
	}, nil
}
