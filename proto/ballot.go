package proto

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"path/filepath"
)

type GovBallotAd struct {
	Path            string            `json:"path"`          // path within repo where ballot will be persisted, also unique ballot name
	Choices         []string          `json:"choices"`       // ballot choices
	Group           string            `json:"group"`         // community group eligible to participate
	Strategy        GovBallotStrategy `json:"strategy"`      // balloting strategy
	GoverningBranch string            `json:"branch"`        // branch governing the ballot
	ParentCommit    string            `json:"parent_commit"` // commit before ballot
}

var (
	GovBallotAdFilebase    = "ballot_advertisement"
	GovBallotTallyFilebase = "ballot_tally"

	GovBallotRoot = filepath.Join(GovRoot, "ballots")

	GovBallotVoteFilepath          = "vote"
	GovBallotVoteSignatureFilepath = "vote.signature.ed25519"
)

func BallotAdPath(ballotPath string) string {
	return filepath.Join(GovBallotRoot, ballotPath, GovBallotAdFilebase)
}

func BallotTallyPath(ballotPath string) string {
	return filepath.Join(GovBallotRoot, ballotPath, GovBallotTallyFilebase)
}

func BallotGenesisCommitHeader(ballotBranch string) string {
	return fmt.Sprintf("Create ballot on branch %v", ballotBranch)
}

func BallotVoteCommitHeader(communityURL string, ballotBranch string, ballotPath string) string {
	return fmt.Sprintf("Submitted vote in community %v at branch %v on ballot %v", communityURL, ballotBranch, ballotPath)
}

// GovBallotVote describes the contents of a vote on a ballot.
type GovBallotVote struct {
	BallotAd GovBallotAd `json:"ballot_advertisement"`
	Choice   string      `json:"vote_choice"`
	Strength float64     `json:"vote_strength"`
}

var (
	GovBallotVoteBranchPrefix = "vote#"
)

func BallotVoteBranch(ctx context.Context, ballotAdBytes []byte) (string, error) {
	h := sha256.New()
	if _, err := h.Write(ballotAdBytes); err != nil {
		return "", err
	}
	return GovBallotVoteBranchPrefix + base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

// tally results

type GovBallotTally struct {
	Ad           GovBallotAd     `json:"ballot_ad"`
	TallyUsers   GovTallyUsers   `json:"tally_users"`
	TallyChoices GovTallyChoices `json:"tally_choices"`
}

type GovTallyUser struct {
	UserName       string         `json:"user_name"`
	UserPublicURL  string         `json:"user_public_url"`
	UserVote       *GovBallotVote `json:"user_vote"` // nil indicates vote was not accessible
	UserFetchError error          `json:"user_fetch_error"`
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
