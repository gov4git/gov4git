package gov

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/govproto"
)

func GetSnapshotDirLocal(local git.Local, sourceAddr proto.Address, sourceCommit string) files.Dir {
	return local.Dir().Subdir(govproto.SnapshotDir(sourceAddr, sourceCommit))
}

// SnapshotOriginLatest downloads the latest commit on a given branch at a remote source repo, and
// places it into a local community repo.
// It stages but does not commit the changes made to the local community repo.
func SnapshotOriginLatest(
	ctx context.Context,
	into git.Local,
	sourceAddr proto.Address,
) (git.Commit, error) {

	// clone source repo locally at the branch
	source, err := git.MakeLocal(ctx)
	if err != nil {
		return "", err
	}
	if err := source.CloneOrigin(ctx, sourceAddr); err != nil {
		return "", err
	}

	// get current commit
	latestCommit, err := source.HeadCommitHash(ctx)
	if err != nil {
		return "", err
	}

	// directory inside into where snapshot lives
	srcPath := govproto.SnapshotDir(sourceAddr, latestCommit)
	srcParent, _ := filepath.Split(srcPath)

	// if the community repo already has a snapshot of the source commit, remove it.
	if err := into.Dir().RemoveAll(srcPath); err != nil {
		return "", err
	}

	// prepare the parent directory of the snapshot
	if err := into.Dir().Mkdir(srcParent); err != nil {
		return "", err
	}

	// remove the .git directory from the source clone
	if err := source.Dir().RemoveAll(".git"); err != nil {
		return "", err
	}

	// move the source clone to the community clone
	if err := files.Rename(source.Dir(), into.Dir().Subdir(srcPath)); err != nil {
		return "", err
	}

	// stage changes in community clone
	if err := into.Add(ctx, []string{srcPath}); err != nil {
		return "", err
	}

	return git.Commit(latestCommit), nil
}
