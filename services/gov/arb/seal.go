package arb

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type SealIn struct {
	BallotBranch string `json:"ballot_branch"`
	BallotPath   string `json:"ballot_path"`
}

type SealOut struct {
	In          *SealIn                 `json:"in"`
	BallotRepo  string                  `json:"ballot_repo"`
	BallotTally govproto.GovBallotTally `json:"ballot_tally"`
}

func (x GovArbService) Seal(ctx context.Context, in *SealIn) (*SealOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.BallotBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	out, err := x.SealLocal(ctx, community, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (x GovArbService) SealLocal(ctx context.Context, community git.Local, in *SealIn) (*SealOut, error) {
	out, err := x.SealLocalStageOnly(ctx, community, in)
	if err != nil {
		return nil, err
	}

	// push to community origin repo
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}

	return out, nil
}

func (x GovArbService) SealLocalStageOnly(ctx context.Context, community git.Local, in *SealIn) (*SealOut, error) {
	// verify tally file is present
	ballotDirpath := govproto.OpenBallotDirpath(in.BallotPath)
	tallyFilepath := govproto.OpenBallotTallyFilepath(in.BallotPath)
	var tally govproto.GovBallotTally
	if _, err := community.Dir().ReadFormFile(ctx, tallyFilepath, &tally); err != nil {
		return nil, err
	}

	sealedDirpath := govproto.ClosedBallotDirpath(in.BallotPath)
	parent, _ := filepath.Split(sealedDirpath)
	if err := community.Dir().Mkdir(parent); err != nil {
		return nil, err
	}

	if err := files.Rename(community.Dir().Subdir(ballotDirpath), community.Dir().Subdir(sealedDirpath)); err != nil {
		return nil, err
	}

	if err := community.Add(ctx, []string{sealedDirpath}); err != nil {
		return nil, err
	}
	if err := community.Remove(ctx, []string{ballotDirpath}); err != nil {
		return nil, err
	}

	return &SealOut{
		In:          in,
		BallotRepo:  x.GovConfig.CommunityURL,
		BallotTally: tally,
	}, nil
}
