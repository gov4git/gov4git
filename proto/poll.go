package proto

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"path/filepath"
)

type GovPollAd struct {
	Path            string          `json:"path"`          // path within repo where poll will be persisted, also unique poll name
	Choices         []string        `json:"choices"`       // ballot choices
	Group           string          `json:"group"`         // community group eligible to participate
	Strategy        GovPollStrategy `json:"strategy"`      // polling strategy
	GoverningBranch string          `json:"branch"`        // branch governing the poll
	ParentCommit    string          `json:"parent_commit"` // commit before poll
}

type GovPollStrategy struct {
	Prioritize *GovPollStrategyPrioritize `json:"prioritize"`
}

type GovPollStrategyPrioritize struct{}

var (
	GovPollAdFilebase    = "poll_advertisement"
	GovPollTallyFilebase = "poll_tally"

	GovPollRoot = filepath.Join(GovRoot, "polls")

	GovPollVoteFilepath          = "vote"
	GovPollVoteSignatureFilepath = "vote.signature.ed25519"
)

func PollAdPath(pollPath string) string {
	return filepath.Join(GovPollRoot, pollPath, GovPollAdFilebase)
}

func PollTallyPath(pollPath string) string {
	return filepath.Join(GovPollRoot, pollPath, GovPollTallyFilebase)
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

func PollVoteBranch(ctx context.Context, pollAdBytes []byte) (string, error) {
	h := sha256.New()
	if _, err := h.Write(pollAdBytes); err != nil {
		return "", err
	}
	return GovPollVoteBranchPrefix + base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

// tally results

type GovPollTally struct {
	Ad           GovPollAd       `json:"ad"`
	TallyUsers   GovTallyUsers   `json:"tally_users"`
	TallyChoices GovTallyChoices `json:"tally_choices"`
}

type GovTallyUser struct {
	UserName       string       `json:"user_name"`
	UserPublicURL  string       `json:"user_public_url"`
	UserVote       *GovPollVote `json:"user_vote"` // nil indicates vote was not accessible
	UserFetchError error        `json:"user_fetch_error"`
}

type GovTallyUsers []GovTallyUser

// vote aggregation to choices

type GovTallyChoice struct {
	Choice        string  `json:"choice"`
	TallyStrength float64 `json:"tally_strength"`
}

type GovTallyChoices []GovTallyChoice

func AggregateVotes(tallyVotes GovTallyUsers) GovTallyChoices {
	s := map[string]float64{} // choice -> strength
	for _, tv := range tallyVotes {
		if tv.UserVote == nil {
			continue
		}
		choice, strength := tv.UserVote.Choice, tv.UserVote.Strength
		t, ok := s[choice]
		if !ok {
			t = 0.0
		}
		t += strength
		s[choice] = t
	}
	tallies := make(GovTallyChoices, 0, len(s))
	for choice, strength := range s {
		tallies = append(tallies, GovTallyChoice{Choice: choice, TallyStrength: strength})
	}
	// sort.Sort(tallies)
	return tallies
}
