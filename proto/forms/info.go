package forms

type PublicInfo struct {
	PublicRepoURL string
	PublicKey     Ed25519PublicKey
}

type PrivateInfo struct {
	PrivateRepoURL string
	PrivateKey     Ed25519PrivateKey
	PublicInfo     PublicInfo
}
