package proto

import (
	"testing"

	"github.com/petar/gitty/sys/form"
)

func TestSign(t *testing.T) {
	priv, err := GenerateKeyPair("pub", "priv")
	if err != nil {
		t.Fatal(err)
	}
	signed, err := Sign(priv, 123)
	if err != nil {
		t.Fatal(err)
	}
	if !signed.Verify() {
		t.Errorf("signature does not verify")
	}
	buf, err := form.EncodeForm(signed)
	if err != nil {
		t.Fatal(err)
	}
	var signed2 Signed[int]
	if err = form.DecodeForm(buf, &signed2); err != nil {
		t.Fatal(err)
	}
	if !signed2.Verify() {
		t.Errorf("signature does not verify")
	}
}
