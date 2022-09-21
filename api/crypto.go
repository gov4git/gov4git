package soul

import (
	"crypto/ed25519"

	"github.com/petar/gitsoc/proto/forms"
	"github.com/petar/gitsoc/sys/form"
)

func GenerateKeyPair() (*forms.PrivateInfo, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}
	return &forms.PrivateInfo{
		Ed25519PrivateKey: form.Bytes(privKey),
		PublicInfo: forms.PublicInfo{
			Ed25519PublicKey: form.Bytes(pubKey),
		},
	}, nil
}

func Sign(priv forms.PrivateInfo, form form.Form) (*forms.Signed, error) {
	XXX
	return &forms.Signed{
		Signature: XXX,
		Msg:       XXX,
	}, nil
}
