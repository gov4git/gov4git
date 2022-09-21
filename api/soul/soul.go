package soul

import (
	"context"

	"github.com/petar/gitsoc/files"
	"github.com/petar/gitsoc/git"
	"github.com/petar/gitsoc/proto/forms"
)

type SoulAPI struct {
	Address forms.SoulAddress
}

func (x SoulAPI) InitKeygen(ctx context.Context) error {
	repo := git.Local{Path: files.DirOf(ctx).Path}
	if err := repo.Init(ctx); err != nil {
		return err
	}
	var stage files.FormFiles
	XXX
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
