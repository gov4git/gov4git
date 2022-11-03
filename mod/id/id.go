package id

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
	PublicAddress    git.Address      `json:"public_address"`
	PublicKeyEd25519 Ed25519PublicKey `json:"public_key_ed25519"`
}

type PrivateCredentials struct {
	PrivateAddress    git.Address       `json:"private_origin"`
	PrivateKeyEd25519 Ed25519PrivateKey `json:"private_key_ed25519"`
	PublicCredentials PublicCredentials `json:"public_credentials"`
}
