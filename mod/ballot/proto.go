package ballot

import (
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/ns"
	"github.com/gov4git/gov4git/mod"
	"github.com/gov4git/gov4git/mod/member"
)

var (
	ballotNS   = mod.RootNS.Sub("ballots")
	adFilebase = "ballot_ad.json"
)

func OpenBallotNS[S Strategy](name ns.NS) ns.NS {
	var s S
	return ballotNS.Sub("open").Sub(s.StrategyName()).Join(name)
}

func ClosedBallotNS[S Strategy](name ns.NS) ns.NS {
	var s S
	return ballotNS.Sub("closed").Sub(s.StrategyName()).Join(name)
}

type Ad struct {
	Name         ns.NS          `json:"path"`
	Choices      []string       `json:"choices"`
	Participants member.Group   `json:"group"`
	ParentCommit git.CommitHash `json:"parent_commit"`
}
