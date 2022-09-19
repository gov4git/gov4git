package soul

import (
	"context"
	"path/filepath"

	"github.com/petar/gitsoc/config"
	"github.com/petar/gitsoc/files"
	"github.com/petar/gitsoc/forms"
	"github.com/petar/gitsoc/git"
	"github.com/petar/gitsoc/workspace"
)

type Soul struct {
	Workspace workspace.Workspace
	Address   forms.SoulAddress
}

func (x Soul) InitKeyGen(ctx context.Context) (*forms.PrivateInfo, error) {

	// prepare local workspace
	dir, err := x.Workspace.MakeEphemeralDir("Soul.InitKeyGen")
	if err != nil {
		return nil, err
	}
	privDir, pubDir := filepath.Join(dir, "priv"), filepath.Join(dir, "pub")

	// build private repo
	privStage := files.Files{
		files.File{Path: config.PrivateSoulInfoPath, Body: XXX},
	}
	if err = git.InitStageCommitPushToOrigin(ctx, privDir, XXXurl, privStage, XXXprivCommit); err != nil {
		return nil, err
	}

	// build public repo
	XXX
}
