package git

import (
	"context"

	"github.com/petar/gitsoc/config"
	"github.com/petar/gitsoc/files"
)

func InitStageCommitPushToOrigin(ctx context.Context, repo Local, repoURL string, stage files.FormFiles, commitMsg string) error {
	if err := repo.Init(ctx); err != nil {
		return err
	}
	if err := repo.Dir().WriteFormFiles(stage); err != nil {
		return err
	}
	if err := Add(ctx, repoDir, stage.Paths()...); err != nil {
		return err
	}
	if err := Commit(ctx, repoDir, commitMsg); err != nil {
		return err
	}
	if err := RenameBranch(ctx, repoDir, config.MainBranch); err != nil {
		return err
	}
	if err := AddRemoteOrigin(ctx, repoDir, repoURL); err != nil {
		return err
	}
	return PushToOrigin(ctx, repoDir, config.MainBranch)
}
