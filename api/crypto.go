package soul

import "crypto/ed25519"

func GenerateKeyPair() {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
}
