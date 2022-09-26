package proto

import (
	"context"
	"testing"

	"github.com/petar/gitty/sys/form"
)

func TestSign(t *testing.T) {
	ctx := context.Background()
	priv, err := GenerateCredentials("pub", "priv")
	if err != nil {
		t.Fatal(err)
	}
	signed, err := Sign(ctx, priv, 123)
	if err != nil {
		t.Fatal(err)
	}
	if !signed.Verify() {
		t.Errorf("signature does not verify")
	}
	buf, err := form.EncodeForm(ctx, signed)
	if err != nil {
		t.Fatal(err)
	}
	var signed2 Signed[int]
	if err = form.DecodeForm(ctx, buf, &signed2); err != nil {
		t.Fatal(err)
	}
	if !signed2.Verify() {
		t.Errorf("signature does not verify")
	}
}
