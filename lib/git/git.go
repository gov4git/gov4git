package git

import (
	"context"
	"fmt"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/gov4git/gov4git/lib/must"
)

type URL string

type Branch string

const MainBranch Branch = "main"

const Origin = "origin"

type Address struct {
	Repo   URL
	Branch Branch
}

func NewAddress(repo URL, branch Branch) Address {
	return Address{Repo: repo, Branch: branch}
}

type Repository = git.Repository

type Tree = git.Worktree

func CloneOrInitBranch(ctx context.Context, addr Address) *Repository {
	repo, err := must.Try1(func() *Repository { return CloneBranch(ctx, addr) })
	if err == nil {
		return repo
	}
	if err != transport.ErrEmptyRemoteRepository {
		must.Panic(ctx, err)
	}
	repo, err = git.Init(memory.NewStorage(), memfs.New())
	must.NoError(ctx, err)

	_, err = repo.CreateRemote(&config.RemoteConfig{Name: Origin, URLs: []string{string(addr.Repo)}})
	must.NoError(ctx, err)

	err = repo.CreateBranch(&config.Branch{Name: string(addr.Branch), Remote: Origin})
	must.NoError(ctx, err)

	ChangeDefaultBranch(ctx, repo, addr.Branch)

	return repo
}

func ChangeDefaultBranch(ctx context.Context, repo *Repository, main Branch) {
	// https://github.com/hairyhenderson/gomplate/pull/1217/files#diff-06d907e05a1688ce7548c3d8b4877a01a61b3db506755db4419761dbe9fe0a5bR232
	branch := plumbing.NewBranchReferenceName(string(main))
	h := plumbing.NewSymbolicReference(plumbing.HEAD, branch)
	err := repo.Storer.SetReference(h)
	must.NoError(ctx, err)

	c, err := repo.Config()
	must.NoError(ctx, err)
	c.Init.DefaultBranch = string(main)

	err = repo.Storer.SetConfig(c)
	must.NoError(ctx, err)
}

func InitPlain(ctx context.Context, path string, isBare bool) *Repository {
	repo, err := git.PlainInit(path, isBare)
	if err != nil {
		must.Panic(ctx, err)
	}
	ChangeDefaultBranch(ctx, repo, MainBranch)
	return repo
}

func CloneBranch(ctx context.Context, addr Address) *Repository {
	repo, err := git.CloneContext(ctx,
		memory.NewStorage(),
		memfs.New(),
		&git.CloneOptions{
			URL:           string(addr.Repo),
			ReferenceName: plumbing.NewBranchReferenceName(string(addr.Branch)),
			SingleBranch:  true,
		},
	)
	must.NoError(ctx, err)

	return repo
}

func Worktree(ctx context.Context, repo *Repository) *Tree {
	wt, err := repo.Worktree()
	if err != nil {
		must.Panic(ctx, err)
	}
	return wt
}

func Add(ctx context.Context, wt *Tree, path string) {
	if _, err := wt.Add(path); err != nil {
		must.Panic(ctx, err)
	}
}

func Commit(ctx context.Context, wt *Tree, msg string) {
	if _, err := wt.Commit(msg, &git.CommitOptions{}); err != nil {
		must.Panic(ctx, err)
	}
}

func Checkout(ctx context.Context, wt *Tree, branch Branch) {
	err := wt.Checkout(
		&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(string(branch)),
		},
	)
	must.NoError(ctx, err)
}

func Push(ctx context.Context, r *Repository) {
	if err := r.PushContext(ctx, &git.PushOptions{}); err != nil {
		must.Panic(ctx, err)
	}
}

func Head(ctx context.Context, r *Repository) *plumbing.Reference {
	h, err := r.Head()
	must.NoError(ctx, err)
	return h
}

func Remotes(ctx context.Context, r *Repository) []*git.Remote {
	h, err := r.Remotes()
	must.NoError(ctx, err)
	return h
}

func Branches(ctx context.Context, r *Repository) []*plumbing.Reference {
	iter, err := r.Branches()
	must.NoError(ctx, err)
	refs := []*plumbing.Reference{}
	for {
		ref, err := iter.Next()
		if err != nil {
			break
		}
		refs = append(refs, ref)
	}
	return refs
}

func Dump(ctx context.Context, r *Repository) {
	iter, err := r.Log(&git.LogOptions{})
	if err != nil {
		fmt.Println("HEAD not found")
		return
	}
	fmt.Println("Log:")
	for {
		c, err := iter.Next()
		if err != nil {
			break
		}
		fmt.Println(c)
	}
	fmt.Println("HEAD:", Head(ctx, r))
	for _, r := range Remotes(ctx, r) {
		fmt.Println("REMOTE:", r)
	}
	for _, r := range Branches(ctx, r) {
		fmt.Println("BRANCH:", r)
	}
}
