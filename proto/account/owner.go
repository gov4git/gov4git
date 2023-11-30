package account

import (
	"github.com/gov4git/lib4git/ns"
)

type OwnerID string

func OwnerIDFromNS(p ns.NS) OwnerID {
	return OwnerID(p.GitPath())
}

func OwnerIDFromLine(line Line) OwnerID {
	return OwnerID(line)
}

func (x OwnerID) String() string {
	return string(x)
}
