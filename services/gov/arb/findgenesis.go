package arb

import (
	"context"
	"fmt"
	"strings"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type FindPollGenesisIn struct {
	PollBranch string `json:"poll_branch"`
}

type FindPollGenesisOut struct {
	PollGenesisCommit string `json:"poll_genesis_commit"`
}

func (x GovArbService) FindPollGenesisLocal(ctx context.Context, repo git.Local, in *FindPollGenesisIn) (*FindPollGenesisOut, error) {
	if err := repo.CheckoutBranch(ctx, in.PollBranch); err != nil {
		return nil, err
	}
	commitLog, err := repo.LogOneline(ctx)
	if err != nil {
		return nil, err
	}
	// find the poll genesis commit going in reverse chrono order
	genesisCommit := ""
	header := strings.TrimSpace(proto.PollGenesisCommitHeader(in.PollBranch))
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
		return nil, fmt.Errorf("cannot find poll genesis commit")
	}

	return &FindPollGenesisOut{PollGenesisCommit: genesisCommit}, nil
}
