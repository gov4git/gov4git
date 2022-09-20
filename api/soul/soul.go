package soul

import (
	"context"

	"github.com/petar/gitsoc/files"
	"github.com/petar/gitsoc/git"
	"github.com/petar/gitsoc/proto/forms"
	"github.com/petar/gitsoc/proto/layout"
)

type SoulAPI struct {
	Address forms.SoulAddress
}

func (x SoulAPI) InitBareRemote(ctx context.Context) error {
	repo := git.Local{Path: files.DirOf(ctx).Path}
	if err := repo.Init(ctx); err != nil {
		return err
	}
	var stage files.FormFiles
	XXX
	// if err := repo.Dir().WriteFormFiles(stage); err != nil {
	// 	return err
	// }
	if err := repo.Add(ctx, stage.Paths()); err != nil {
		return err
	}
	if err := repo.Commit(ctx, "initialize empty private soul repo"); err != nil {
		return err
	}
	if err := repo.RenameBranch(ctx, layout.MainBranch); err != nil {
		return err
	}
	if err := repo.AddRemoteOrigin(ctx, repoURL); err != nil {
		return err
	}
	return repo.PushToOrigin(ctx, layout.MainBranch)
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
