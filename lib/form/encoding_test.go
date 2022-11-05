package form

import (
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"strings"
	"testing"
)

func TestEncoding(t *testing.T) {
	h := sha256.New()
	if _, err := h.Write([]byte("x")); err != nil {
		panic(err)
	}
	enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	x := strings.ToLower(enc.EncodeToString(h.Sum(nil)))
	fmt.Println(x)
}
