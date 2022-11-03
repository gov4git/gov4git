package govproto

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gov4git/gov4git/lib/files"
)

type BallotAd struct {
	Path            string            `json:"path"`          // path within repo where ballot will be persisted, also unique ballot name
	Choices         []string          `json:"choices"`       // ballot choices
	Group           string            `json:"group"`         // community group eligible to participate
	Strategy        GovBallotStrategy `json:"strategy"`      // balloting strategy
	GoverningBranch string            `json:"branch"`        // branch governing the ballot
	ParentCommit    string            `json:"parent_commit"` // commit before ballot
}

var (
	BallotAdFilebase    = "ballot_advertisement"
	BallotTallyFilebase = "ballot_tally"

	OpenBallotRootpath   = filepath.Join(GovRoot, "open_ballots")
	ClosedBallotRootpath = filepath.Join(GovRoot, "closed_ballots")

	BallotVoteFilepath          = "vote"
	BallotVoteSignatureFilepath = "vote.signature.ed25519"
)

// open ballots

func OpenBallotDirpath(ballotPath string) string {
	return filepath.Join(OpenBallotRootpath, ballotPath)
}

func OpenBallotAdFilepath(ballotPath string) string {
	return filepath.Join(OpenBallotDirpath(ballotPath), BallotAdFilebase)
}

func OpenBallotTallyFilepath(ballotPath string) string {
	return filepath.Join(OpenBallotDirpath(ballotPath), BallotTallyFilebase)
}

func ExtractOpenBallotPathFromTally(tallyFilepath string) (string, error) {
	dir, _ := filepath.Split(tallyFilepath)
	if !strings.HasPrefix(dir, OpenBallotRootpath) {
		return "", fmt.Errorf("missing open ballot dir prefix: %v", dir)
	}
	return files.TrimSlashes(dir[len(OpenBallotRootpath):]), nil
}

// closed ballots

func ClosedBallotDirpath(ballotPath string) string {
	return filepath.Join(ClosedBallotRootpath, ballotPath)
}

func ClosedBallotAdFilepath(ballotPath string) string {
	return filepath.Join(ClosedBallotDirpath(ballotPath), BallotAdFilebase)
}

func ClosedBallotTallyFilepath(ballotPath string) string {
	return filepath.Join(ClosedBallotDirpath(ballotPath), BallotTallyFilebase)
}

func ExtractClosedBallotPathFromTally(tallyFilepath string) (string, error) {
	dir, _ := filepath.Split(tallyFilepath)
	if !strings.HasPrefix(dir, ClosedBallotRootpath) {
		return "", fmt.Errorf("missing closed ballot dir prefix")
	}
	return files.TrimSlashes(dir[len(ClosedBallotRootpath):]), nil
}

// other

func BallotGenesisCommitHeader(ballotBranch string) string {
	return fmt.Sprintf("Create ballot on branch %v", ballotBranch)
}

func BallotVoteCommitHeader(communityURL string, ballotBranch string, ballotPath string) string {
	return fmt.Sprintf("Submitted vote in community %v at branch %v on ballot %v", communityURL, ballotBranch, ballotPath)
}

// BallotVote describes the contents of a vote on a ballot.
type BallotVote struct {
	BallotAd  BallotAd   `json:"ballot_advertisement"`
	Elections []Election `json:"elections"`
}

type Election struct {
	Choice   string  `json:"vote_choice"`
	Strength float64 `json:"vote_strength"`
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
	Ad           BallotAd        `json:"ballot_ad"`
	TallyUsers   GovTallyUsers   `json:"tally_users"`
	TallyChoices GovTallyChoices `json:"tally_choices"`
}

type GovTallyUser struct {
	UserName       string      `json:"user_name"`
	UserPublicURL  string      `json:"user_public_url"`
	UserVote       *BallotVote `json:"user_vote"` // nil indicates vote was not accessible
	UserFetchError string      `json:"user_fetch_error"`
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
		for _, v := range tv.UserVote.Elections {
			t, ok := s[v.Choice]
			if !ok {
				t = 0.0
			}
			t += v.Strength
			s[v.Choice] = t
		}
	}
	tallies := make(GovTallyChoices, 0, len(s))
	for choice, strength := range s {
		tallies = append(tallies, GovTallyChoice{Choice: choice, TallyStrength: strength})
	}
	// sort.Sort(tallies)
	return tallies
}
