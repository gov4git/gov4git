package arb

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/services/gov"
	"github.com/gov4git/gov4git/services/gov/member"
	"github.com/gov4git/gov4git/services/gov/user"
)

type TallyIn struct {
	ReferendumBranch string `json:"referendum_branch"`
}

type TallyOut struct {
	ReferendumRepo   string             `json:"referendum_repo"`
	ReferendumBranch string             `json:"referendum_branch"`
	ReferendumTally  proto.GovPollTally `json:"referendum_tally"`
}

func (x TallyOut) Human(ctx context.Context) string {
	data, err := form.EncodeForm(ctx, x)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (x GovArbService) Tally(ctx context.Context, in *TallyIn) (*TallyOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.ReferendumBranch); err != nil {
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
	// find poll ad and leave local repo checked out at the genesis commit
	findAd, err := x.FindPollAdLocal(ctx, community, &FindPollAdIn{PollBranch: in.ReferendumBranch})
	if err != nil {
		return nil, err
	}

	// list users participating in poll
	memberService := member.GovMemberService{GovConfig: x.GovConfig}
	participants, err := memberService.ListLocal(ctx, community, "", findAd.PollAd.Group)
	if err != nil {
		return nil, err
	}

	// get participants' user infos (public_url, etc.)
	userInfo, err := user.GetInfos(ctx, community, member.ExtractUsersFromMembership(participants))
	if err != nil {
		return nil, err
	}

	// checkout referendum branch latest
	if err := community.CheckoutBranch(ctx, in.ReferendumBranch); err != nil {
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
	findPoll *FindPollAdOut,
	userInfo user.UserInfos,
) (*TallyOut, error) {

	out := &TallyOut{
		ReferendumRepo:   x.GovConfig.CommunityURL,
		ReferendumBranch: in.ReferendumBranch,
		ReferendumTally: proto.GovPollTally{
			Ad:         findPoll.PollAd,
			TallyUsers: make(proto.GovTallyUsers, len(userInfo)),
		},
	}

	// snapshot votes from user repos
	// TODO: parallelize snapshots
	for i, info := range userInfo {
		userVote, err := x.snapshotParseVerifyUserVote(ctx, community, findPoll, info)
		out.ReferendumTally.TallyUsers[i] = proto.GovTallyUser{
			UserName:       info.UserName,
			UserPublicURL:  info.UserInfo.URL,
			UserVote:       userVote,
			UserFetchError: err,
		}
	}

	// aggregate votes to choices
	out.ReferendumTally.TallyChoices = proto.AggregateVotes(out.ReferendumTally.TallyUsers)

	// write/stage snapshots and tally to community repo
	tallyPath := proto.PollTallyPath(findPoll.PollAd.Path)
	stage := files.FormFiles{
		files.FormFile{Path: tallyPath, Form: out.ReferendumTally},
	}
	if err := community.Dir().WriteFormFiles(ctx, stage); err != nil {
		return nil, err
	}
	if err := community.Add(ctx, stage.Paths()); err != nil {
		return nil, err
	}

	// commit snapshots and tally to community repo
	if err := community.Commitf(ctx, "tally votes on referendum branch %v", in.ReferendumBranch); err != nil {
		return nil, err
	}

	return out, nil
}

func (x GovArbService) snapshotParseVerifyUserVote(
	ctx context.Context,
	community git.Local,
	findPoll *FindPollAdOut,
	userInfo user.UserInfo,
) (*proto.GovPollVote, error) {

	// compute the name of the vote branch in the user's repo
	voteBranch, err := proto.PollVoteBranch(ctx, findPoll.PollAdBytes)
	if err != nil {
		return nil, err
	}

	// snapshot user's vote branch into the community repo
	snap, err := x.GovService().SnapshotBranchLatest(ctx,
		&gov.SnapshotBranchLatestIn{
			SourceRepo:   userInfo.UserInfo.URL,
			SourceBranch: voteBranch,
			Community:    community,
		})
	if err != nil {
		return nil, err
	}

	// parse and verify user's vote
	snapDir := gov.GetSnapshotDirLocal(community, snap.In.SourceRepo, snap.SourceCommit)
	var signature proto.SignedPlaintext

	if _, err := snapDir.ReadFormFile(ctx, proto.GovPollVoteSignatureFilepath, &signature); err != nil {
		return nil, err
	}
	if !signature.Verify() {
		return nil, fmt.Errorf("signature is not valid")
	}
	var vote proto.GovPollVote
	if err := form.DecodeForm(ctx, signature.Plaintext, &vote); err != nil {
		return nil, err
	}

	// TODO: verify the user's vote is for this poll

	return &vote, nil
}
