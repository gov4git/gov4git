package id

import (
	"github.com/google/uuid"
	"github.com/gov4git/gov4git/proto"
)

var (
	PublicNS  = proto.RootNS.Sub("id")
	PrivateNS = proto.RootNS.Sub("id")
)

var (
	PublicCredentialsNS  = PublicNS.Sub("public_credentials.json")
	PrivateCredentialsNS = PrivateNS.Sub("private_credentials.json")
)

type ID string

func GenerateUniqueID() ID {
	return ID(uuid.New().String())
}

type PublicCredentials struct {
	ID               ID               `json:"id"`
	PublicKeyEd25519 Ed25519PublicKey `json:"public_key_ed25519"`
}

type PrivateCredentials struct {
	PrivateKeyEd25519 Ed25519PrivateKey `json:"private_key_ed25519"`
	PublicCredentials PublicCredentials `json:"public_credentials"`
}
