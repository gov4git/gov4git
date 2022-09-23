package forms

import (
	"context"

	"github.com/petar/gitty/sys/form"
)

type PublicInfo struct {
	PublicRepoURL    string           `gitty:"public_repo_url"`
	PublicKeyEd25519 Ed25519PublicKey `gitty:"public_key_ed25519"`
}

func (x PublicInfo) Skeletize(ctx context.Context) any {
	return form.SkeletizeStruct(ctx, x)
}

func (x *PublicInfo) DeSkeletize(ctx context.Context, from any) error {
	return form.DeSkeletizeStruct(ctx, x, from)
}

type PrivateInfo struct {
	PrivateRepoURL string
	PrivateKey     Ed25519PrivateKey
	PublicInfo     PublicInfo
}
