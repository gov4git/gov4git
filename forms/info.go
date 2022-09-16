package forms

type PublicInfo struct {
	Form             string
	PublicRepoURL    string
	Ed25519PublicKey Bytes
}

func (x *PublicInfo) Tidy() {
	x.Form = "PublicInfo"
}

type PrivateInfo struct {
	Form              string
	PrivateRepoURL    string
	Ed25519PrivateKey Bytes
	PublicInfo        PublicInfo
}

func (x *PrivateInfo) Tidy() {
	x.Form = "PrivateInfo"
}
