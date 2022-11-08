package id

import (
	"context"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
)

func Init(
	ctx context.Context,
	publicAddr git.Address,
	privateAddr git.Address,
) git.Change[PrivateCredentials] {
	public := git.CloneOrInitBranch(ctx, publicAddr)
	private := git.CloneOrInitBranch(ctx, privateAddr)
	publicTree := git.Worktree(ctx, public)
	privateTree := git.Worktree(ctx, private)

	privChg := InitLocal(ctx, publicAddr, privateAddr, publicTree, privateTree)

	git.Push(ctx, private)
	git.Push(ctx, public)
	return privChg
}

func InitLocal(
	ctx context.Context,
	publicAddr git.Address,
	privateAddr git.Address,
	publicTree *git.Tree,
	privateTree *git.Tree,
) git.Change[PrivateCredentials] {
	privChg := InitPrivate(ctx, privateTree, publicAddr, privateAddr)
	pubChg := InitPublic(ctx, publicTree, privChg.Result.PublicCredentials)
	git.Commit(ctx, privateTree, privChg.Msg)
	git.Commit(ctx, publicTree, pubChg.Msg)
	return privChg
}

func InitPrivate(ctx context.Context, priv *git.Tree, publicAddr git.Address, privateAddr git.Address) git.Change[PrivateCredentials] {
	if _, err := priv.Filesystem.Stat(PrivateCredentialsNS.Path()); err == nil {
		must.Errorf(ctx, "private credentials file already exists")
	}
	cred, err := GenerateCredentials(publicAddr, privateAddr)
	if err != nil {
		must.Panic(ctx, err)
	}
	git.ToFileStage(ctx, priv, PrivateCredentialsNS.Path(), cred)
	return git.Change[PrivateCredentials]{
		Result: cred,
		Msg:    "Initialized private credentials.",
	}
}

func InitPublic(ctx context.Context, pub *git.Tree, cred PublicCredentials) git.ChangeNoResult {
	if _, err := pub.Filesystem.Stat(PublicCredentialsNS.Path()); err == nil {
		must.Errorf(ctx, "public credentials file already exists")
	}
	git.ToFileStage(ctx, pub, PublicCredentialsNS.Path(), cred)
	return git.ChangeNoResult{
		Msg: "Initialized public credentials.",
	}
}
