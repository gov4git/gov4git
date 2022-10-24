package arb

import (
	"context"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type FindBallotAdIn struct {
	BallotBranch string `json:"ballot_branch"`
	BallotPath   string `json:"ballot_path"`
}

type FindBallotAdOut struct {
	BallotAd      govproto.BallotAd `json:"ballot_ad"`
	BallotAdBytes form.Bytes        `json:"ballot_ad_bytes"`
}

// FindBallotAdLocal finds the advertisement of a ballot in a local clone of community repo (at the ballot branch)
func (x GovArbService) FindBallotAdLocal(ctx context.Context, repo git.Local, in *FindBallotAdIn) (*FindBallotAdOut, error) {
	// read the ballot advertisement
	ballotAdPath := govproto.OpenBallotAdFilepath(in.BallotPath)
	ballotAdFile, err := repo.Dir().ReadByteFile(ballotAdPath)
	if err != nil {
		return nil, err
	}

	var ballotAd govproto.BallotAd
	if err := form.DecodeForm(ctx, ballotAdFile.Bytes, &ballotAd); err != nil {
		return nil, err
	}

	return &FindBallotAdOut{
		BallotAd:      ballotAd,
		BallotAdBytes: ballotAdFile.Bytes,
	}, nil
}
