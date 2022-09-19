package forms

import "github.com/petar/gitsoc/form"

type PublicInfo struct {
	Form             string
	PublicRepoURL    string
	Ed25519PublicKey form.Bytes
}

func (x *PublicInfo) Tidy() {
	x.Form = "PublicInfo"
}

type PrivateInfo struct {
	Form              string
	PrivateRepoURL    string
	Ed25519PrivateKey form.Bytes
	PublicInfo        PublicInfo
}

func (x *PrivateInfo) Tidy() {
	x.Form = "PrivateInfo"
}
