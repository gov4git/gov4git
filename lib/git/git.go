package git

import (
	"context"

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

type Worktree = git.Worktree

func CloneOrInitBranch(ctx context.Context, addr Address) (*Repository, error) {
	repo, err := CloneBranch(ctx, addr)
	if err == nil {
		return repo, nil
	}
	if err != transport.ErrEmptyRemoteRepository {
		return nil, err
	}
	if repo, err = git.Init(memory.NewStorage(), memfs.New()); err != nil {
		println("A")
		return nil, err
	}
	if _, err = repo.CreateRemote(&config.RemoteConfig{Name: Origin, URLs: []string{string(addr.Repo)}}); err != nil {
		println("B")
		return nil, err
	}
	if err = repo.CreateBranch(&config.Branch{Name: string(addr.Branch), Remote: Origin}); err != nil {
		println("C")
		return nil, err
	}
	wt, err := repo.Worktree()
	if err != nil {
		println("D")
		return nil, err
	}
	//XXX
	file, err := wt.Filesystem.Create("ok")
	if err != nil {
		panic(err)
	}
	file.Close()
	wt.Add("ok")
	h, err := wt.Commit("ok", &git.CommitOptions{})
	if err != nil {
		println("X")
		return nil, err
	}
	println(h.String())

	if err = wt.Checkout(&git.CheckoutOptions{Branch: plumbing.ReferenceName(addr.Branch), Create: true}); err != nil {
		println("E", addr.Branch, "E")
		return nil, err
	}

	g, err := repo.Head()
	if err != nil {
		println("Y")
		return nil, err
	}
	println(g.String())

	return repo, nil
}

func CloneBranch(ctx context.Context, addr Address) (*Repository, error) {
	repo, err := git.CloneContext(ctx,
		memory.NewStorage(),
		memfs.New(),
		&git.CloneOptions{
			URL:           string(addr.Repo),
			ReferenceName: plumbing.ReferenceName(addr.Branch),
		},
	)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func MustInitBare(ctx context.Context, path string) *Repository {
	repo, err := git.PlainInit(path, true)
	if err != nil {
		must.Panic(ctx, err)
	}
	return repo
}

func MustCloneOrInitBranch(ctx context.Context, addr Address) *Repository {
	repo, err := CloneOrInitBranch(ctx, addr)
	if err != nil {
		must.Panic(ctx, err)
	}
	return repo
}

func MustCloneBranch(ctx context.Context, addr Address) *Repository {
	repo, err := CloneBranch(ctx, addr)
	if err != nil {
		must.Panic(ctx, err)
	}
	return repo
}

func MustWorktree(ctx context.Context, repo *Repository) *Worktree {
	wt, err := repo.Worktree()
	if err != nil {
		must.Panic(ctx, err)
	}
	return wt
}

func MustAdd(ctx context.Context, wt *Worktree, path string) {
	if _, err := wt.Add(path); err != nil {
		must.Panic(ctx, err)
	}
}

func MustCommit(ctx context.Context, wt *Worktree, msg string) {
	if _, err := wt.Commit(msg, &git.CommitOptions{}); err != nil {
		must.Panic(ctx, err)
	}
}

func MustPush(ctx context.Context, r *Repository) {
	if err := r.PushContext(ctx, &git.PushOptions{}); err != nil {
		must.Panic(ctx, err)
	}
}
