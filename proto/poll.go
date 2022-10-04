package proto

import (
	"fmt"
)

type GovPollAd struct {
	Path         string          `json:"path"`          // path within repo where poll will be persisted, also unique poll name
	Choices      []string        `json:"choices"`       // ballot choices
	Group        string          `json:"group"`         // community group eligible to participate
	Strategy     GovPollStrategy `json:"strategy"`      // polling strategy
	Branch       string          `json:"branch"`        // branch governing the poll
	ParentCommit string          `json:"parent_commit"` // commit before poll
}

type GovPollStrategy struct {
	Prioritize *GovPollStrategyPrioritize `json:"prioritize"`
}

type GovPollStrategyPrioritize struct{}

var (
	GovPollAdFilebase   = "poll_advertisement"
	GovPollBranchPrefix = "poll#"
)

func PollBranch(path string) string {
	return GovPollBranchPrefix + path
}

func PollGenesisCommitHeader(branch string) string {
	return fmt.Sprintf("Create poll on branch %v", branch)
}
