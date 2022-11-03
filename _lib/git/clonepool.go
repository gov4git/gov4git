package git

import (
	"context"
	"sync"
)

type ClonePool struct {
	sync.Mutex
	clones []*clone
}

type Clone interface {
	Local() Local
	Release()
}

type clone struct {
	repo   string
	branch string
	local  Local
	lk     sync.Mutex
	busy   bool
}

func (x *clone) Local() Local {
	return x.local
}

func (x *clone) TryLease() bool {
	x.lk.Lock()
	defer x.lk.Unlock()
	if x.busy {
		return false
	}
	x.busy = true
	return true
}

func (x *clone) Busy() bool {
	x.lk.Lock()
	defer x.lk.Unlock()
	return x.busy
}

func (x *clone) Release() {
	x.lk.Lock()
	defer x.lk.Unlock()
	if !x.busy {
		panic("releasing an idle clone")
	}
	x.busy = false
}

func NewClonePool() *ClonePool {
	return &ClonePool{}
}

func (x *ClonePool) Clone(ctx context.Context, repo string, branch string) (Clone, error) {
	for _, cand := range x.findIdleClones(repo, branch) {
		if !cand.TryLease() {
			continue
		}
		// switch to branch
		if err := cand.Local().CheckoutBranchForce(ctx, branch); err != nil {
			return nil, err
		}
		if err := cand.Local().PullUpstream(ctx); err != nil {
			return nil, err
		}
		return cand, nil
	}

	local, err := CloneOrigin(ctx, Origin{Repo: URL(repo), Branch: Branch(branch)})
	if err != nil {
		return nil, err
	}
	c := &clone{repo: repo, branch: branch, local: local, busy: true}
	x.Lock()
	defer x.Unlock()
	x.clones = append(x.clones, c)
	return c, nil
}

func (x *ClonePool) findIdleClones(repo string, branch string) []*clone {
	x.Lock()
	defer x.Unlock()
	r := []*clone{}
	for _, c := range x.clones {
		if c.repo == repo && c.branch == branch && !c.Busy() {
			r = append(r, c)
		}
	}
	return r
}
