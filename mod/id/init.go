package id

import (
	"context"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
)

func Init(ctx context.Context, publicAddr git.Address, privateAddr git.Address) git.Change[PrivateCredentials] {
	public := git.CloneOrInitBranch(ctx, publicAddr)
	private := git.CloneOrInitBranch(ctx, privateAddr)
	privateWt := git.Worktree(ctx, private)
	publicWt := git.Worktree(ctx, public)
	privChg := InitPrivate(ctx, privateWt, publicAddr, privateAddr)
	pubChg := InitPublic(ctx, publicWt, privChg.Result.PublicCredentials)
	git.Commit(ctx, privateWt, privChg.Msg)
	git.Commit(ctx, publicWt, pubChg.Msg)
	git.Push(ctx, private)
	git.Push(ctx, public)
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
