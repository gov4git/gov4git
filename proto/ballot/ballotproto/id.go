package ballotproto

import (
	"strings"

	"github.com/gov4git/lib4git/ns"
)

type BallotID string

func (x BallotID) ToNS() ns.NS {
	return ns.ParseFromGitPath(x.String())
}

func (x BallotID) String() string {
	return string(x)
}

func (x BallotID) GitPath() string {
	return string(x)
}

func (x BallotID) TallyNS() ns.NS {
	return x.GitNS().Append(TallyFilebase)
}

func (x BallotID) AdNS() ns.NS {
	return x.GitNS().Append(AdFilebase)
}

func (x BallotID) OutcomeNS() ns.NS {
	return x.GitNS().Append(OutcomeFilebase)
}

func (x BallotID) StrategyNS() ns.NS {
	return x.GitNS().Append(StrategyFilebase)
}

func (x BallotID) GitNS() ns.NS {
	return BallotKV.KeyNS(BallotNS, x)
}

func ParseBallotID(p string) BallotID {
	return BallotID(strings.TrimSpace(p))
}
