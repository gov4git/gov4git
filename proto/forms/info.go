package forms

type PublicInfo struct {
	Form          string
	PublicRepoURL string
	PublicKey     Ed25519PublicKey
}

func (x *PublicInfo) Tidy() {
	x.Form = "PublicInfo"
}

type PrivateInfo struct {
	Form           string
	PrivateRepoURL string
	PrivateKey     Ed25519PrivateKey
	PublicInfo     PublicInfo
}

func (x *PrivateInfo) Tidy() {
	x.Form = "PrivateInfo"
}
