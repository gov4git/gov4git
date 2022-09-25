package proto

import (
	"crypto/ed25519"
	"encoding/json"

	"github.com/petar/gitty/sys/form"
)

type Ed25519PublicKey = form.Bytes

type Ed25519PrivateKey = form.Bytes

func GenerateKeyPair(publicURL, privateURL string) (*PrivateInfo, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	return &PrivateInfo{
		PrivateRepoURL: privateURL,
		PrivateKey:     Ed25519PrivateKey(privKey),
		PublicInfo: PublicInfo{
			PublicRepoURL: publicURL,
			PublicKey:     Ed25519PublicKey(pubKey),
		},
	}, nil
}

func Sign[Statement form.Form](priv *PrivateInfo, stmt Statement) (*Signed[Statement], error) {
	data, err := form.EncodeForm(stmt)
	if err != nil {
		return nil, err
	}
	signature := ed25519.Sign(ed25519.PrivateKey(priv.PrivateKey), data)
	return &Signed[Statement]{
		Statement: stmt,
		PublicKey: Ed25519PublicKey(priv.PublicInfo.PublicKey),
		Signature: form.Bytes(signature),
		Original:  string(data),
	}, nil
}

type Signed[Statement form.Form] signed[Statement]

type signed[Statement form.Form] struct {
	Statement Statement `json:"-"`
	PublicKey Ed25519PublicKey
	Signature form.Bytes
	Original  string // encoded statement
}

func (x *Signed[Statement]) Verify() bool {
	return ed25519.Verify(ed25519.PublicKey(x.PublicKey), []byte(x.Original), x.Signature)
}

func (x *Signed[Statement]) UnmarshalJSON(d []byte) error {
	if err := json.Unmarshal(d, (*signed[Statement])(x)); err != nil {
		return err
	}
	var stmt Statement
	if err := json.Unmarshal([]byte(x.Original), &stmt); err != nil {
		return err
	}
	x.Statement = stmt
	return nil
}
