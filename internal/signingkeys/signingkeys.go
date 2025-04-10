package signingkeys

import (
	"strings"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

// ExtractHexKeyID extracts the key ID from an armored GPG public key
// and returns it as a hexadecimal string.
func ExtractHexKeyID(armoredPublicKey string) (string, error) {
	key, err := crypto.NewKeyFromArmored(armoredPublicKey)
	if err != nil {
		return "", err
	}

	return strings.ToUpper(key.GetHexKeyID()), nil
}
