package git

import (
	"context"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/gov4git/gov4git/lib/must"
)

type URL string

type Branch string

type Address struct {
	Repo   URL
	Branch Branch
}

type Repository = git.Repository

type Worktree = git.Worktree

func CloneBranch(ctx context.Context, addr Address) (*Repository, error) {
	fs := memfs.New()
	storer := memory.NewStorage()
	repo, err := git.CloneContext(ctx, storer, fs,
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
	if _, err := wt.Commit(msg, nil); err != nil {
		must.Panic(ctx, err)
	}
}

func MustPush(ctx context.Context, r *Repository) {
	if err := r.PushContext(ctx, nil); err != nil {
		must.Panic(ctx, err)
	}
}
