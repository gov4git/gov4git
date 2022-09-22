package soul

import (
	"context"

	"github.com/petar/gitty/proto/forms"
	"github.com/petar/gitty/proto/layout"
	. "github.com/petar/gitty/sys/base"
	. "github.com/petar/gitty/sys/files"
	"github.com/petar/gitty/sys/git"
)

type SoulAPI struct {
	Address forms.SoulAddress
}

func (x SoulAPI) InitKeygen(ctx context.Context) ContextError {
	repo := git.LocalFromDir(DirOf(ctx).Subdir("private"))
	// clone or init repo
	if err := repo.CloneOrInitBranch(ctx, x.Address.PrivateURL, layout.MainBranch); err != nil {
		return DoneErr(ctx, err)
	}
	// check if key files already exist
	XXX
	// generate keys
	stage := FormFiles{
		FormFile{Path: XXX, Form: XXX},
	}
	XXX
	return DoneOk(ctx)
}

// func (x SoulAPI) CloneKeyGen(ctx context.Context) (*forms.PrivateInfo, error) {

// 	// create workspace
// 	ctxDir := files.DirOf(ctx)
// 	privRepo, pubRepo := git.Local{Path: ctxDir.Abs("priv")}, git.Local{Path: ctxDir.Abs("priv")}

// 	// build private repo
// 	privFiles := files.FormFiles{
// 		files.FormFile{Path: config.PrivateSoulInfoPath, Form: XXX},
// 	}

// 	if err := privRepo.

// 	// if err = git.InitStageCommitPushToOrigin(ctx, privDir, XXXurl, privStage, XXXprivCommit); err != nil {
// 	// 	return nil, err
// 	// }

// 	// build public repo
// 	XXX
// }

// func (x SoulAPI) CheckoutSendChannel(ctx context.Context, to XXX, topic string) error {
// 	XXX
// }

// func (x SoulAPI) SyncCheckoutReceiveChannel(ctx context.Context, from XXX, topic string) error {
// 	XXX
// }
