package forms

import "testing"

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
}
