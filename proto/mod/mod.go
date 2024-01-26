package mod

import (
	"context"
	"slices"
	"sync"

	"github.com/gov4git/lib4git/must"
)

type ModuleRegistry[N ~string, M any] struct {
	lk   sync.Mutex
	mods map[N]M
}

func NewModuleRegistry[N ~string, M any]() *ModuleRegistry[N, M] {
	return &ModuleRegistry[N, M]{
		mods: map[N]M{},
	}
}

func (r *ModuleRegistry[N, M]) ListKeys() []N {
	r.lk.Lock()
	defer r.lk.Unlock()
	ns := []N{}
	for k := range r.mods {
		ns = append(ns, k)
	}
	slices.Sort[[]N](ns)
	return ns
}

func (r *ModuleRegistry[N, M]) List() ([]N, []M) {
	r.lk.Lock()
	defer r.lk.Unlock()
	ns := []N{}
	for k := range r.mods {
		ns = append(ns, k)
	}
	slices.Sort[[]N](ns)
	ms := []M{}
	for _, k := range ns {
		ms = append(ms, r.mods[k])
	}
	return ns, ms
}

func (r *ModuleRegistry[N, M]) Set(ctx context.Context, key N, module M) {
	r.lk.Lock()
	defer r.lk.Unlock()
	_, ok := r.mods[key]
	must.Assertf(ctx, !ok, "module %v already registered", key)
	r.mods[key] = module
}

func (r *ModuleRegistry[N, M]) Get(ctx context.Context, key N) M {
	r.lk.Lock()
	defer r.lk.Unlock()
	v, ok := r.mods[key]
	must.Assertf(ctx, ok, "module %v not found", key)
	return v
}
