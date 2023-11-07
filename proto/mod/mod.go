package mod

import "sync"

type ModuleRegistry[M any] struct {
	lk   sync.Mutex
	mods map[string]M
}

func (r *ModuleRegistry[M]) Add(key string, module M) {
	r.lk.Lock()
	defer r.lk.Unlock()
	r.mods[key] = module
}

func (r *ModuleRegistry[M]) Get(key string) M {
	r.lk.Lock()
	defer r.lk.Unlock()
	return r.mods[key]
}
