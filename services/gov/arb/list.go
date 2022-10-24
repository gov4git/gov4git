package arb

import (
	"context"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type ListIn struct {
	BallotBranch string `json:"ballot_branch"`
}

type ListOut struct {
	In            *ListIn  `json:"in"`
	OpenBallots   []string `json:"open_ballots"`
	ClosedBallots []string `json:"closed_ballots"`
}

func (x GovArbService) List(ctx context.Context, in *ListIn) (*ListOut, error) {
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.BallotBranch); err != nil {
		return nil, err
	}
	out, err := x.ListLocal(ctx, community, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (x GovArbService) ListLocal(ctx context.Context, community git.Local, in *ListIn) (*ListOut, error) {
	openBallotGlob := govproto.OpenBallotAdFilepath("*")
	closedBallotGlob := govproto.ClosedBallotAdFilepath("*")

	base.Infof("globbing open ballots %v", openBallotGlob)
	base.Infof("globbing closed ballots %v", closedBallotGlob)

	// extract open ballot names
	openMatches, err := community.Dir().Glob(openBallotGlob)
	if err != nil {
		return nil, err
	}
	openBallots := make([]string, len(openMatches))
	for i := range openMatches {
		if openBallots[i], err = govproto.ExtractOpenBallotPathFromTally(openMatches[i]); err != nil {
			return nil, err
		}
	}

	// extract closed ballot names
	closedMatches, err := community.Dir().Glob(closedBallotGlob)
	if err != nil {
		return nil, err
	}
	closedBallots := make([]string, len(closedMatches))
	for i := range closedMatches {
		if closedBallots[i], err = govproto.ExtractClosedBallotPathFromTally(closedMatches[i]); err != nil {
			return nil, err
		}
	}

	return &ListOut{
		In:            in,
		OpenBallots:   openBallots,
		ClosedBallots: closedBallots,
	}, nil
}
