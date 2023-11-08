package mod

import (
	"context"
	"slices"
	"sync"

	"github.com/gov4git/lib4git/must"
)

type ModuleRegistry[M any] struct {
	lk   sync.Mutex
	mods map[string]M
}

func NewModuleRegistry[M any]() *ModuleRegistry[M] {
	return &ModuleRegistry[M]{
		mods: map[string]M{},
	}
}

func (r *ModuleRegistry[M]) Keys() []string {
	r.lk.Lock()
	defer r.lk.Unlock()
	ks := []string{}
	for k := range r.mods {
		ks = append(ks, k)
	}
	slices.Sort(ks)
	return ks
}

func (r *ModuleRegistry[M]) Set(ctx context.Context, key string, module M) {
	r.lk.Lock()
	defer r.lk.Unlock()
	_, ok := r.mods[key]
	must.Assertf(ctx, !ok, "module %v already registered", key)
	r.mods[key] = module
}

func (r *ModuleRegistry[M]) Get(ctx context.Context, key string) M {
	r.lk.Lock()
	defer r.lk.Unlock()
	v, ok := r.mods[key]
	must.Assertf(ctx, ok, "module %v not found", key)
	return v
}
