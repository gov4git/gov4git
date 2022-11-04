package id

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
	"github.com/gov4git/gov4git/mod/runtime"
)

func (m PrivateMod) Init(ctx context.Context) runtime.Change[PrivateCredentials] {
	public := git.CloneOrInitBranch(ctx, m.Public)
	private := git.CloneOrInitBranch(ctx, m.Private)
	privateWt := git.Tree(ctx, private)
	publicWt := git.Tree(ctx, public)
	privChg := m.InitPrivate(ctx, privateWt)
	pubChg := m.InitPublic(ctx, publicWt, privChg.Result.PublicCredentials)
	git.Commit(ctx, privateWt, privChg.Msg)
	git.Commit(ctx, publicWt, pubChg.Msg)
	git.Push(ctx, private)
	git.Push(ctx, public)

	return privChg
}

func (m PrivateMod) InitPrivate(ctx context.Context, wt *git.Worktree) runtime.Change[PrivateCredentials] {
	filepath := m.Subpath(PrivateCredentialsFilebase)
	if _, err := wt.Filesystem.Stat(filepath); err == nil {
		must.Panic(ctx, fmt.Errorf("private credentials file already exists"))
	}
	cred, err := GenerateCredentials(m.Public, m.Private)
	if err != nil {
		must.Panic(ctx, err)
	}
	form.MustEncodeToFile(ctx, wt.Filesystem, filepath, cred)
	git.Add(ctx, wt, filepath)
	return runtime.Change[PrivateCredentials]{
		Result: cred,
		Msg:    "Initialized private credentials.",
	}
}

func (m PrivateMod) InitPublic(ctx context.Context, wt *git.Worktree, cred PublicCredentials) runtime.Change[struct{}] {
	filepath := m.Subpath(PublicCredentialsFilebase)
	if _, err := wt.Filesystem.Stat(filepath); err == nil {
		must.Panic(ctx, fmt.Errorf("public credentials file already exists"))
	}
	form.MustEncodeToFile(ctx, wt.Filesystem, filepath, cred)
	git.Add(ctx, wt, filepath)
	return runtime.Change[struct{}]{
		Msg: "Initialized public credentials.",
	}
}
