package ballot

import (
	"github.com/gov4git/gov4git/lib/ns"
	"github.com/gov4git/gov4git/mod"
)

type Strategy interface {
	StrategyName() string
}

type BallotAddress[S Strategy] struct {
	Gov  mod.GovAddress
	Name ns.NS
}
