package id

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/lib4git/form"
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

func GenerateRandomID() ID {
	const w = 512 / 8 // 512 bits, measured in bytes
	buf := make([]byte, w)
	rand.Read(buf)
	return ID(form.BytesHashForFilename(buf))
}

type PublicCredentials struct {
	ID               ID               `json:"id"`
	PublicKeyEd25519 Ed25519PublicKey `json:"public_key_ed25519"`
}

func Ed25519PubKeyToID(pubKey ed25519.PublicKey) ID {
	return ID(form.BytesHashForFilename(pubKey))
}

func (x PublicCredentials) IsValid() bool {
	return form.BytesHashForFilename(x.PublicKeyEd25519) == string(x.ID)
}

type PrivateCredentials struct {
	PrivateKeyEd25519 Ed25519PrivateKey `json:"private_key_ed25519"`
	PublicCredentials PublicCredentials `json:"public_credentials"`
}

func (x PrivateCredentials) IsValid() bool {
	return x.PublicCredentials.IsValid()
}
