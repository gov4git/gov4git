package proto

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/gov4git/gov4git/lib/form"
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

	GovPollVoteFilepath          = "vote"
	GovPollVoteSignatureFilepath = "vote.signature.ed25519"
)

func PollBranch(path string) string {
	return GovPollBranchPrefix + path
}

func PollPathFromBranch(branch string) (string, error) {
	if len(branch) < len(GovPollBranchPrefix) {
		return "", fmt.Errorf("invalid poll branch")
	}
	return branch[len(GovPollBranchPrefix):], nil
}

func PollGenesisCommitHeader(pollBranch string) string {
	return fmt.Sprintf("Create poll on branch %v", pollBranch)
}

func PollVoteCommitHeader(communityURL string, referendumBranch string) string {
	return fmt.Sprintf("Submitted vote for community %v and branch %v", communityURL, referendumBranch)
}

// GovPollVote describes the contents of a vote on a poll.
type GovPollVote struct {
	PollAd   GovPollAd `json:"poll_advertisement"`
	Choice   string    `json:"vote_choice"`
	Strength float64   `json:"vote_strength"`
}

var (
	GovPollVoteBranchPrefix = "vote#"
)

// TODO: if PollVoteBranch depands on the pollAd byte representation,
// it will enable interoperation between different program versions at the voter and the community
func PollVoteBranch(ctx context.Context, pollAd GovPollAd) (string, error) {
	data, err := form.EncodeForm(ctx, pollAd)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := h.Write(data); err != nil {
		return "", err
	}
	return GovPollVoteBranchPrefix + base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
