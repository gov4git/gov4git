package git

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type RepoBranch struct {
	RepoURL string
	Branch  string
}

type WorkspaceConfig struct {
	RootPath string // local path, root of workspace
}

func (x WorkspaceConfig) Init() error {
	return x.MkDir("")
}

// MkDir makes a directory at a path, relative to the workspace root.
func (x WorkspaceConfig) MkDir(path string) error {
	return os.MkdirAll(filepath.Join(x.RootPath, path), 0755)
}

func (x WorkspaceConfig) AbsPath(path string) string {
	return filepath.Join(x.RootPath, path)
}

type Manager struct {
	space  WorkspaceConfig
	lk     sync.Mutex
	clones map[RepoBranch]*BranchClone
}

func NewManager(w WorkspaceConfig) (*Manager, error) {
	return &Manager{clones: map[RepoBranch]*BranchClone{}}, nil
}

// Provision may return a repo that is not recently synced with the remote.
func (x *Manager) Provision(ctx context.Context, repoBranch RepoBranch) (*BranchClone, error) {
	x.lk.Lock()
	defer x.lk.Unlock()
	if bc, present := x.clones[repoBranch]; present {
		return bc, nil
	}

	// compute workspace directory
	url2Path, err := URL2Path(repoBranch.RepoURL)
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(url2Path, ":branch", repoBranch.Branch)
	absDir := x.space.AbsPath(dir)

	// check if clone already exists
	repo, err := git.PlainOpen(absDir)
	if err == nil {
		bc := &BranchClone{repo: repo}
		x.clones[repoBranch] = bc
		return bc, nil
	}

	// clone the repository
	repo, err = git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{
		URL:           repoBranch.RepoURL,
		ReferenceName: plumbing.ReferenceName(repoBranch.Branch),
		SingleBranch:  true,
		NoCheckout:    true,
		Progress:      os.Stderr,
	})
	if err != nil {
		return nil, err
	}
	bc := &BranchClone{repo: repo}
	x.clones[repoBranch] = bc
	return bc, nil
}

type BranchClone struct {
	sync.Mutex
	repo *git.Repository
}
