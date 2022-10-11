package arb

import (
	"context"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type FindPollAdIn struct {
	PollBranch string `json:"poll_branch"`
	PollPath   string `json:"poll_path"`
}

type FindPollAdOut struct {
	PollAd      proto.GovPollAd `json:"poll_ad"`
	PollAdBytes form.Bytes      `json:"poll_ad_bytes"`
}

// FindPollAdLocal finds the advertisement of a poll in a local clone of community repo (at the poll branch) and
// leaves the local repo checked out at the genesis commit.
func (x GovArbService) FindPollAdLocal(ctx context.Context, repo git.Local, in *FindPollAdIn) (*FindPollAdOut, error) {
	// read the poll advertisement
	pollAdPath := proto.PollAdPath(in.PollPath)
	pollAdFile, err := repo.Dir().ReadByteFile(pollAdPath)
	if err != nil {
		return nil, err
	}

	var pollAd proto.GovPollAd
	if err := form.DecodeForm(ctx, pollAdFile.Bytes, &pollAd); err != nil {
		return nil, err
	}

	return &FindPollAdOut{
		PollAd:      pollAd,
		PollAdBytes: pollAdFile.Bytes,
	}, nil
}
