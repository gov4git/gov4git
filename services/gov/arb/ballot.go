package arb

import (
	"context"
	"fmt"
	"strings"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/services/gov/group"
)

type CreateBallotIn struct {
	Path            string   `json:"ballot_path"` // path where ballot will be persisted
	Choices         []string `json:"ballot_choices"`
	Group           string   `json:"ballot_group"`
	Strategy        string   `json:"ballot_strategy"`
	GoverningBranch string   `json:"ballot_governing_branch"`
	BallotBranch    string   `json:"ballot_branch"`
}

type CreateBallotOut struct {
	CreateBallotIn      *CreateBallotIn `json:"create_ballot_in"`
	BallotCommunityURL  string          `json:"ballot_community_url"`
	BallotBranch        string          `json:"ballot_branch"`
	BallotGenesisCommit string          `json:"ballot_genesis_commit"`
}

func (x GovArbService) CreateBallot(ctx context.Context, in *CreateBallotIn) (*CreateBallotOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.GoverningBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	out, err := x.CreateBallotLocal(ctx, community, in)
	if err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return out, nil
}

func (x GovArbService) CreateBallotLocal(ctx context.Context, community git.Local, in *CreateBallotIn) (*CreateBallotOut, error) {
	// verify path is not in use
	ballotPath := strings.TrimSpace(git.MakeNonAbs(in.Path))
	ballotAdPath := govproto.OpenBallotAdFilepath(ballotPath)
	if _, err := community.Dir().Stat(ballotAdPath); err == nil {
		return nil, fmt.Errorf("ballot already exists")
	}

	// XXX: verify no closed ballots with same name

	// verify group exists
	if _, err := group.GetInfo(ctx, community, in.Group); err != nil {
		return nil, fmt.Errorf("ballot group does not exist")
	}

	// get hash of current commit on branch
	head, err := community.HeadCommitHash(ctx)
	if err != nil {
		return nil, err
	}

	// checkout the ballot branch
	ballotBranch := strings.TrimSpace(in.BallotBranch)
	if ballotBranch == "" { // use governing branch
		ballotBranch = in.GoverningBranch
	} else {
		if err := community.CheckoutNewBranch(ctx, ballotBranch); err != nil {
			return nil, err
		}
	}

	// create and stage ballot advertisement
	strategy, err := govproto.ParseBallotStrategy(in.Strategy)
	if err != nil {
		return nil, err
	}
	ballotAd := govproto.BallotAd{
		Path:            ballotPath,
		Choices:         in.Choices,
		Group:           in.Group,
		Strategy:        strategy,
		GoverningBranch: in.GoverningBranch,
		ParentCommit:    head,
	}
	stage := files.FormFiles{
		files.FormFile{Path: ballotAdPath, Form: ballotAd},
	}
	if err := community.Dir().WriteFormFiles(ctx, stage); err != nil {
		return nil, err
	}
	if err := community.Add(ctx, stage.Paths()); err != nil {
		return nil, err
	}

	// commit changes and include ballot ad in commit message
	out := &CreateBallotOut{
		CreateBallotIn:      in,
		BallotCommunityURL:  x.GovConfig.CommunityURL,
		BallotBranch:        ballotBranch,
		BallotGenesisCommit: "", // populate after commit
	}
	hum := govproto.BallotGenesisCommitHeader(ballotBranch)
	msg, err := git.PrepareCommitMsg(ctx, hum, ballotAd)
	if err != nil {
		return nil, err
	}
	if err := community.Commit(ctx, msg); err != nil {
		return nil, err
	}

	// get hash of ballot genesis commit
	if out.BallotGenesisCommit, err = community.HeadCommitHash(ctx); err != nil {
		return nil, err
	}

	return out, nil
}
