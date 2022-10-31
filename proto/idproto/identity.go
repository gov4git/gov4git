package idproto

import "github.com/google/uuid"

type IdentityConfig struct {
	PublicURL  string `json:"public_url"`
	PrivateURL string `json:"private_url"`
}

type ID string

func GenerateUniqueID() ID {
	return ID(uuid.New().String())
}

type PublicCredentials struct {
	ID               ID               `json:"id"`
	PublicURL        string           `json:"public_url"`
	PublicKeyEd25519 Ed25519PublicKey `json:"public_key_ed25519"`
}

type PrivateCredentials struct {
	PrivateURL        string            `json:"private_url"`
	PrivateKeyEd25519 Ed25519PrivateKey `json:"private_key_ed25519"`
	PublicCredentials PublicCredentials `json:"public_credentials"`
}
