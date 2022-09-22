package forms

import (
	"encoding/json"

	"github.com/petar/gitty/proto/forms"
	"github.com/petar/gitty/sys/form"
)

type Ed25519PublicKey = form.Bytes

type Ed25519PrivateKey = form.Bytes

func Sign[Statement form.Form](priv *forms.PrivateInfo) {
	XXX
}

type Signed[Statement form.Form] signed[Statement]

type signed[Statement form.Form] struct {
	PublicKey Ed25519PublicKey
	Signature form.Bytes
	Statement Statement
	Original  []byte `json:"-"`
}

func (x *Signed[Statement]) Verify() error {
	XXX
}

func (x *Signed[Statement]) UnmarshalJSON(d []byte) error {
	if err := json.Unmarshal(d, (*signed[Statement])(x)); err != nil {
		return err
	}
	x.Original = d
	return nil
}
