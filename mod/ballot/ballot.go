package ballot

import (
	"github.com/gov4git/gov4git/lib/ns"
	"github.com/gov4git/gov4git/mod/gov"
)

type Strategy interface {
	StrategyName() string
}

type BallotAddress[S Strategy] struct {
	Gov  gov.CommunityAddress
	Name ns.NS
}
