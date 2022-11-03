package id

import (
	"context"
	"crypto/ed25519"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
)

type Ed25519PublicKey = form.Bytes

type Ed25519PrivateKey = form.Bytes

func GenerateCredentials(public git.Address, private git.Address) (PrivateCredentials, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return PrivateCredentials{}, err
	}
	return PrivateCredentials{
		PrivateAddress:    private,
		PrivateKeyEd25519: Ed25519PrivateKey(privKey),
		PublicCredentials: PublicCredentials{
			ID:               GenerateUniqueID(),
			PublicAddress:    public,
			PublicKeyEd25519: Ed25519PublicKey(pubKey),
		},
	}, nil
}

type SignedPlaintext struct {
	Plaintext        form.Bytes       `json:"plaintext"`
	Signature        form.Bytes       `json:"signature"`
	PublicKeyEd25519 Ed25519PublicKey `json:"ed25519_public_key"`
}

func SignPlaintext(ctx context.Context, priv *PrivateCredentials, plaintext []byte) (*SignedPlaintext, error) {
	signature := ed25519.Sign(ed25519.PrivateKey(priv.PrivateKeyEd25519), plaintext)
	return &SignedPlaintext{
		Plaintext:        plaintext,
		Signature:        signature,
		PublicKeyEd25519: priv.PublicCredentials.PublicKeyEd25519,
	}, nil
}

func (signed *SignedPlaintext) Verify() bool {
	return ed25519.Verify(ed25519.PublicKey(signed.PublicKeyEd25519), signed.Plaintext, signed.Signature)
}
