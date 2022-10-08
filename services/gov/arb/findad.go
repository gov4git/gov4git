package arb

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type FindPollAdIn struct {
	PollBranch string `json:"poll_branch"`
}

type FindPollAdOut struct {
	PollGenesisCommit string          `json:"poll_genesis_commit"`
	PollAd            proto.GovPollAd `json:"poll_ad"`
}

// FindPollAdLocal finds the advertisement of a poll in a local clone of community repo (at the poll branch) and
// leaves the local repo checked out at the genesis commit.
func (x GovArbService) FindPollAdLocal(ctx context.Context, repo git.Local, in *FindPollAdIn) (*FindPollAdOut, error) {
	// find the genesis commit of the referendum
	findGenesis, err := x.FindPollGenesisLocal(ctx, repo, &FindPollGenesisIn{PollBranch: in.PollBranch})
	if err != nil {
		return nil, err
	}

	// read the poll advertisement
	pollPath, err := proto.PollPathFromBranch(in.PollBranch)
	if err != nil {
		return nil, err
	}
	pollAdPath := filepath.Join(pollPath, proto.GovPollAdFilebase)
	if err := repo.CheckoutBranch(ctx, findGenesis.PollGenesisCommit); err != nil {
		return nil, err
	}
	var pollAd proto.GovPollAd
	if _, err := repo.Dir().ReadFormFile(ctx, pollAdPath, &pollAd); err != nil {
		return nil, err
	}

	return &FindPollAdOut{
		PollGenesisCommit: findGenesis.PollGenesisCommit,
		PollAd:            pollAd,
	}, nil
}
