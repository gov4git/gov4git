package soul

import (
	"context"

	"github.com/petar/gitsoc/forms"
	"github.com/petar/gitsoc/git"
	"github.com/petar/gitsoc/workspace"
)

type Soul struct {
	Workspace *workspace.Workspace
	Address   forms.SoulAddress
}

func (x Soul) Init(ctx context.Context) (*forms.PrivateInfo, error) {
	dir, err := x.Workspace.MakeEphemeralDir("Soul.Init")
	if err != nil {
		return nil, err
	}
	XXX
	if err := git.Init(ctx, dir); err != nil {
		return nil, err
	}
	XXX
}
