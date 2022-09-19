package git

import (
	"context"

	"github.com/petar/gitsoc/config"
	"github.com/petar/gitsoc/files"
)

func InitStageCommitPushToOrigin(ctx context.Context, repoDir string, repoURL string, stage files.Files, commitMsg string) error {
	if err := Init(ctx, repoDir); err != nil {
		return err
	}
	if err := files.WriteFiles(repoDir, stage); err != nil {
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
