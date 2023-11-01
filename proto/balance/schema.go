package balance

import (
	"github.com/gov4git/lib4git/ns"
)

type Balance ns.NS

func (x Balance) NS() ns.NS {
	return ns.NS(x)
}

// p must be a valid git path.
func ParseBalance(p string) Balance {
	return Balance(ns.ParseFromGitPath(p))
}

// balanceKey must be a valid git path.
func userPropKey(balanceKey Balance) string {

	balanceNS := ns.NS{"balance"}
	return balanceNS.Join(balanceKey.NS()).GitPath()

	// XXX: old version below; test for backwards compatibility
	// key := strings.TrimLeft(filepath.ToSlash(string(balanceKey)), "/")
	// return strings.Join([]string{"balance", key}, "/")
}
