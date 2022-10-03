package proto

import (
	"context"
	"crypto/ed25519"
	"encoding/json"

	"github.com/petar/gov4git/lib/form"
)

type Ed25519PublicKey = form.Bytes

type Ed25519PrivateKey = form.Bytes

func GenerateCredentials(publicURL, privateURL string) (*PrivateCredentials, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	return &PrivateCredentials{
		PrivateURL:        privateURL,
		PrivateKeyEd25519: Ed25519PrivateKey(privKey),
		PublicCredentials: PublicCredentials{
			PublicURL:        publicURL,
			PublicKeyEd25519: Ed25519PublicKey(pubKey),
		},
	}, nil
}

func Sign[Statement form.Form](ctx context.Context, priv *PrivateCredentials, stmt Statement) (*Signed[Statement], error) {
	data, err := form.EncodeForm(ctx, stmt)
	if err != nil {
		return nil, err
	}
	signature := ed25519.Sign(ed25519.PrivateKey(priv.PrivateKeyEd25519), data)
	return &Signed[Statement]{
		Statement:        stmt,
		PublicKeyEd25519: Ed25519PublicKey(priv.PublicCredentials.PublicKeyEd25519),
		Signature:        form.Bytes(signature),
		Original:         string(data),
	}, nil
}

type Signed[Statement form.Form] signed[Statement]

type signed[Statement form.Form] struct {
	Statement        Statement        `json:"-"`
	PublicKeyEd25519 Ed25519PublicKey `json:"public_key_ed25519"`
	Signature        form.Bytes       `json:"signature"`
	Original         string           `json:"original"` // signed repn
}

func (x *Signed[Statement]) Verify() bool {
	return ed25519.Verify(ed25519.PublicKey(x.PublicKeyEd25519), []byte(x.Original), x.Signature)
}

func (x *Signed[Statement]) UnmarshalJSON(d []byte) error {
	if err := json.Unmarshal(d, (*signed[Statement])(x)); err != nil {
		return err
	}
	return json.Unmarshal([]byte(x.Original), &x.Statement)
}
