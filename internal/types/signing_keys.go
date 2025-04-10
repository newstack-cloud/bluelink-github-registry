package types

// IntermediarySigningKeys is a struct that represents an intermediate
// representation of signing keys provided in the signing keys environment variable.
type IntermediarySigningKeys struct {
	Keys []*IntermediarySigningKey `json:"keys"`
}

// IntermediarySigningKey is a struct that represents an individual signing key
// in the intermediate representation of signing keys.
type IntermediarySigningKey struct {
	PublicKey string `json:"publicKey"`
}
