package git

import (
	"context"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

type URL string

type Branch string

type Address struct {
	Repo   URL
	Branch Branch
}

type Repository = git.Repository

type Worktree = git.Worktree

func CloneBranch(ctx context.Context, addr Address) (*git.Repository, error) {
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
