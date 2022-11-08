package id

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
)

func Init(ctx context.Context, m PrivateMod) git.Change[PrivateCredentials] {
	public := git.CloneOrInitBranch(ctx, m.Public)
	private := git.CloneOrInitBranch(ctx, m.Private)
	privateWt := git.Worktree(ctx, private)
	publicWt := git.Worktree(ctx, public)
	privChg := InitPrivate(ctx, m, privateWt)
	pubChg := InitPublic(ctx, m, publicWt, privChg.Result.PublicCredentials)
	git.Commit(ctx, privateWt, privChg.Msg)
	git.Commit(ctx, publicWt, pubChg.Msg)
	git.Push(ctx, private)
	git.Push(ctx, public)
	return privChg
}

func InitPrivate(ctx context.Context, m PrivateMod, priv *git.Tree) git.Change[PrivateCredentials] {
	fileNS := m.Sub(PrivateCredentialsFilebase)
	if _, err := priv.Filesystem.Stat(fileNS.Path()); err == nil {
		must.Panic(ctx, fmt.Errorf("private credentials file already exists"))
	}
	cred, err := GenerateCredentials(m.Public, m.Private)
	if err != nil {
		must.Panic(ctx, err)
	}
	form.ToFile(ctx, priv.Filesystem, fileNS.Path(), cred)
	git.Add(ctx, priv, fileNS.Path())
	return git.Change[PrivateCredentials]{
		Result: cred,
		Msg:    "Initialized private credentials.",
	}
}

func InitPublic(ctx context.Context, m PrivateMod, pub *git.Tree, cred PublicCredentials) git.ChangeNoResult {
	fileNS := m.Sub(PublicCredentialsFilebase)
	if _, err := pub.Filesystem.Stat(fileNS.Path()); err == nil {
		must.Panic(ctx, fmt.Errorf("public credentials file already exists"))
	}
	form.ToFile(ctx, pub.Filesystem, fileNS.Path(), cred)
	git.Add(ctx, pub, fileNS.Path())
	return git.ChangeNoResult{
		Msg: "Initialized public credentials.",
	}
}
