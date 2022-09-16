package git

import (
	"context"

	"github.com/petar/gitsoc/files"
)

var (
	MainBranch = "main"
)

func InitPopulatePushToOrigin(ctx context.Context, repoDir string, repoURL string, stage files.Dir, commitMsg string) error {
	if err := Init(ctx, repoDir); err != nil {
		return err
	}
	if err := files.WriteDir(stage, repoDir); err != nil {
		return err
	}
	// XXX: stage files
	if err := Commit(ctx, repoDir, commitMsg); err != nil {
		return err
	}
	if err := RenameBranch(ctx, repoDir, MainBranch); err != nil {
		return err
	}
	if err := AddRemoteOrigin(ctx, repoDir, repoURL); err != nil {
		return err
	}
	return PushToOrigin(ctx, repoDir, MainBranch)
}
