package policy

import "github.com/gov4git/gov4git/proto/mod"

type Policy interface{
	
}

var policyRegistry = mod.NewModuleRegistry[Policy]()

func Set(key string, policy Policy) {
	policyRegistry.Set(key, policy)
}

func Get(key string) Policy {
	return policyRegistry.Get(key)
}
