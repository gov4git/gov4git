package idproto

import (
	"github.com/google/uuid"
	"github.com/gov4git/gov4git/lib/git"
)

type ID string

func GenerateUniqueID() ID {
	return ID(uuid.New().String())
}

type PublicCredentials struct {
	ID               ID               `json:"id"`
	PublicOrigin     git.Origin       `json:"public_origin"`
	PublicKeyEd25519 Ed25519PublicKey `json:"public_key_ed25519"`
}

type PrivateCredentials struct {
	PrivateOrigin     git.Origin        `json:"private_origin"`
	PrivateKeyEd25519 Ed25519PrivateKey `json:"private_key_ed25519"`
	PublicCredentials PublicCredentials `json:"public_credentials"`
}
