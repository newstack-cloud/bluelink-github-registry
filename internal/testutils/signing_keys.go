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

	fmt.Println("length of parsed keys:", len(signingKeysInternal.Keys))
	for _, key := range signingKeysInternal.Keys {
		fmt.Println("length of public key:", len(key.PublicKey))
		fmt.Println("public bits and pieces:",
			strings.Replace(
				strings.Replace(
					key.PublicKey,
					"-----BEGIN PGP PUBLIC KEY BLOCK-----",
					"SOME-HEADER",
					1,
				),
				"-----END PGP PUBLIC KEY BLOCK-----",
				"SOME-FOOTER",
				1,
			),
		)
		fmt.Println(
			"contains ASCII armor header?",
			strings.Contains(key.PublicKey, "-----BEGIN PGP PUBLIC KEY BLOCK-----"),
		)
		fmt.Println(
			"contains ASCII armor footer?",
			strings.Contains(key.PublicKey, "-----END PGP PUBLIC KEY BLOCK-----"),
		)
		fmt.Println(
			"contains literal new line characters that have not been unescaped?",
			strings.Contains(key.PublicKey, "\\n"),
		)
		keyID, err := signingkeys.ExtractHexKeyID(testKey)
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

var (
	testKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----

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
a7KeKz6lT3XUTk5yfQN3eVZ5
=gaJ5
-----END PGP PUBLIC KEY BLOCK-----`
)
