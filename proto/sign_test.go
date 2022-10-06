package proto

import (
	"context"
	"testing"
)

func TestSignPlaintext(t *testing.T) {
	ctx := context.Background()
	priv, err := GenerateCredentials("pub", "priv")
	if err != nil {
		t.Fatal(err)
	}
	_ = priv
	_ = ctx
	// TODO: add test
}
