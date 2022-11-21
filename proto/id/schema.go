package id

import (
	"github.com/google/uuid"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/lib4git/git"
)

var (
	HomeNS    = proto.RootNS.Sub("id")
	VaultNS = proto.RootNS.Sub("id")
)

var (
	PublicCredentialsNS  = HomeNS.Sub("public_credentials.json")
	PrivateCredentialsNS = VaultNS.Sub("private_credentials.json")
)

type ID string

func GenerateUniqueID() ID {
	return ID(uuid.New().String())
}

type PublicCredentials struct {
	ID               ID               `json:"id"`
	HomeAddress      git.Address      `json:"public_address"`
	PublicKeyEd25519 Ed25519PublicKey `json:"public_key_ed25519"`
}

type PrivateCredentials struct {
	VaultAddress      git.Address       `json:"private_origin"`
	PrivateKeyEd25519 Ed25519PrivateKey `json:"private_key_ed25519"`
	PublicCredentials PublicCredentials `json:"public_credentials"`
}
