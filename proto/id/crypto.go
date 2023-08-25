package id

import (
	"context"
	"crypto/ed25519"

	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/must"
)

type Ed25519PublicKey = form.Bytes

type Ed25519PrivateKey = form.Bytes

func GenerateCredentials() (PrivateCredentials, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return PrivateCredentials{}, err
	}
	return PrivateCredentials{
		PrivateKeyEd25519: Ed25519PrivateKey(privKey),
		PublicCredentials: PublicCredentials{
			ID:               Ed25519PubKeyToID(pubKey),
			PublicKeyEd25519: Ed25519PublicKey(pubKey),
		},
	}, nil
}

// signing

type Signed[V form.Form] struct {
	Value            V                `json:"value"`
	Plaintext        form.Bytes       `json:"plaintext"`
	Signature        form.Bytes       `json:"signature"`
	PublicKeyEd25519 Ed25519PublicKey `json:"ed25519_public_key"`
}

func (signed *Signed[V]) Verify() bool {
	// XXX: also verify Value encodes to Plaintext
	return ed25519.Verify(ed25519.PublicKey(signed.PublicKeyEd25519), signed.Plaintext, signed.Signature)
}

func SignBytes(ctx context.Context, priv PrivateCredentials, plaintext []byte) (signature []byte, pubKey []byte) {
	signature = ed25519.Sign(ed25519.PrivateKey(priv.PrivateKeyEd25519), plaintext)
	pubKey = priv.PublicCredentials.PublicKeyEd25519
	return
}

func Sign[V form.Form](ctx context.Context, priv PrivateCredentials, value V) Signed[V] {
	data, err := form.EncodeBytes(ctx, value)
	must.NoError(ctx, err)
	signature, pubKey := SignBytes(ctx, priv, data)
	return Signed[V]{Value: value, Plaintext: data, Signature: signature, PublicKeyEd25519: pubKey}
}
