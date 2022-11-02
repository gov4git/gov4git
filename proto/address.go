package proto

import "github.com/gov4git/gov4git/lib/git"

type Address = git.Origin

type PairAddress struct {
	PublicAddress  Address
	PrivateAddress Address
}
