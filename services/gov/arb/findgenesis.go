package arb

import (
	"context"
	"fmt"
	"strings"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type FindBallotGenesisIn struct {
	BallotBranch string `json:"ballot_branch"`
}

type FindBallotGenesisOut struct {
	BallotGenesisCommit string `json:"ballot_genesis_commit"`
}

func (x GovArbService) FindBallotGenesisLocal(ctx context.Context, repo git.Local, in *FindBallotGenesisIn) (*FindBallotGenesisOut, error) {
	if err := repo.CheckoutBranch(ctx, in.BallotBranch); err != nil {
		return nil, err
	}
	commitLog, err := repo.LogOneline(ctx)
	if err != nil {
		return nil, err
	}
	// find the ballot genesis commit going in reverse chrono order
	genesisCommit := ""
	header := strings.TrimSpace(govproto.BallotGenesisCommitHeader(in.BallotBranch))
	for _, line := range strings.Split(commitLog, "\n") {
		line = strings.TrimSpace(line)
		commitHash, commitHeader, found := strings.Cut(line, " ")
		commitHeader = strings.TrimSpace(commitHeader)
		if !found || commitHash == "" || commitHeader == "" {
			continue
		}
		if !strings.HasPrefix(commitHeader, header) {
			continue
		}
		genesisCommit = commitHash
		break
	}
	if genesisCommit == "" {
		return nil, fmt.Errorf("cannot find ballot genesis commit")
	}

	return &FindBallotGenesisOut{BallotGenesisCommit: genesisCommit}, nil
}
