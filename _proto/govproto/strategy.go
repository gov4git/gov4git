package govproto

import (
	"fmt"
)

type GovBallotStrategy struct {
	PriorityPoll *PriorityPollStrategy `json:"priority_poll"`
}

type PriorityPollStrategy struct{}

type StrategyName string

const (
	PriorityPollStrategyName = "priority-poll"
)

func ParseBallotStrategy(s string) (GovBallotStrategy, error) {
	switch s {
	case PriorityPollStrategyName:
		return GovBallotStrategy{PriorityPoll: &PriorityPollStrategy{}}, nil
	}
	return GovBallotStrategy{}, fmt.Errorf("unknown strategy")
}
