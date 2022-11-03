package arb

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/proto/idproto"
	"github.com/gov4git/gov4git/services/gov"
	"github.com/gov4git/gov4git/services/gov/arb/strategy"
	"github.com/gov4git/gov4git/services/gov/member"
	"github.com/gov4git/gov4git/services/gov/user"
)

type TallyIn struct {
	BallotBranch string `json:"ballot_branch"`
	BallotPath   string `json:"ballot_path"`
}

type TallyOut struct {
	BallotRepo   string                  `json:"ballot_repo"`
	BallotBranch string                  `json:"ballot_branch"`
	BallotTally  govproto.GovBallotTally `json:"ballot_tally"`
}

func (x GovArbService) Tally(ctx context.Context, in *TallyIn) (*TallyOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.BallotBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	out, err := x.TallyLocal(ctx, community, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (x GovArbService) TallyLocal(ctx context.Context, community git.Local, in *TallyIn) (*TallyOut, error) {
	// find ballot ad
	findAd, err := x.FindBallotAdLocal(ctx, community,
		&FindBallotAdIn{BallotBranch: in.BallotBranch, BallotPath: in.BallotPath})
	if err != nil {
		return nil, err
	}

	// list users participating in ballot
	memberService := member.GovMemberService{GovConfig: x.GovConfig}
	participants, err := memberService.ListLocal(ctx, community, "", findAd.BallotAd.Group)
	if err != nil {
		return nil, err
	}

	// get participants' user infos (public_url, etc.)
	userInfo, err := user.GetInfos(ctx, community, member.ExtractUsersFromMembership(participants))
	if err != nil {
		return nil, err
	}

	// checkout referendum branch latest
	if err := community.CheckoutBranch(ctx, in.BallotBranch); err != nil {
		return nil, err
	}

	// fetch votes and compute tally
	out, err := x.FetchVotesAndTallyLocal(ctx, community, in, findAd, userInfo)
	if err != nil {
		return nil, err
	}

	// push to community origin repo
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}

	return out, nil
}

func (x GovArbService) FetchVotesAndTallyLocal(
	ctx context.Context,
	community git.Local,
	in *TallyIn,
	findBallot *FindBallotAdOut,
	userInfo user.UserInfos,
) (*TallyOut, error) {

	out := &TallyOut{
		BallotRepo:   x.GovConfig.CommunityURL,
		BallotBranch: in.BallotBranch,
		BallotTally: govproto.GovBallotTally{
			Ad:         findBallot.BallotAd,
			TallyUsers: make(govproto.GovTallyUsers, len(userInfo)),
		},
	}

	// snapshot votes from user repos
	// TODO: parallelize snapshots
	for i, info := range userInfo {
		userVote, err := x.snapshotParseVerifyUserVote(ctx, community, findBallot, info)
		errstr := ""
		if err != nil {
			errstr = err.Error()
		}
		out.BallotTally.TallyUsers[i] = govproto.GovTallyUser{
			UserName:       info.UserName,
			UserPublicURL:  info.UserInfo.PublicURL,
			UserVote:       userVote,
			UserFetchError: errstr,
		}
	}

	// aggregate votes to choices
	out.BallotTally.TallyChoices = govproto.AggregateVotes(out.BallotTally.TallyUsers)

	// invoke ballot strategy
	strat, err := strategy.ParseStrategy(findBallot.BallotAd.Strategy)
	if err != nil {
		return nil, err
	}
	if err := strat.Tally(ctx, community, out.BallotTally); err != nil {
		return nil, err
	}

	// write/stage snapshots and tally to community repo
	tallyPath := govproto.OpenBallotTallyFilepath(findBallot.BallotAd.Path)
	stage := files.FormFiles{
		files.FormFile{Path: tallyPath, Form: out.BallotTally},
	}
	if err := community.Dir().WriteFormFiles(ctx, stage); err != nil {
		return nil, err
	}
	if err := community.Add(ctx, stage.Paths()); err != nil {
		return nil, err
	}

	// commit snapshots and tally to community repo
	if err := community.Commitf(ctx, "Tally votes on referendum branch %v", in.BallotBranch); err != nil {
		return nil, err
	}

	return out, nil
}

func (x GovArbService) snapshotParseVerifyUserVote(
	ctx context.Context,
	community git.Local,
	findBallot *FindBallotAdOut,
	userInfo user.UserInfo,
) (*govproto.BallotVote, error) {

	// compute the name of the vote branch in the user's repo
	voteBranch, err := govproto.BallotVoteBranch(ctx, findBallot.BallotAdBytes)
	if err != nil {
		return nil, err
	}

	// snapshot user's vote branch into the community repo
	snap, err := x.GovService().SnapshotBranchLatest(ctx,
		&gov.SnapshotBranchLatestIn{
			SourceRepo:   userInfo.UserInfo.PublicURL,
			SourceBranch: voteBranch,
			Community:    community,
		})
	if err != nil {
		return nil, err
	}

	// parse and verify user's vote
	snapDir := gov.GetSnapshotDirLocal(community, snap.In.SourceRepo, snap.SourceCommit)
	var signature idproto.SignedPlaintext

	if _, err := snapDir.ReadFormFile(ctx, govproto.BallotVoteSignatureFilepath, &signature); err != nil {
		return nil, err
	}
	if !signature.Verify() {
		return nil, fmt.Errorf("signature is not valid")
	}
	var vote govproto.BallotVote
	if err := form.DecodeForm(ctx, signature.Plaintext, &vote); err != nil {
		return nil, err
	}

	// TODO: verify the user's vote is for this ballot

	return &vote, nil
}
