package arb

import (
	"context"

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
	out, err := x.FetchVotesAndTallyLocal(ctx, community, in, findAd.PollAd, userInfo)
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
	ad proto.GovPollAd,
	userInfo user.UserInfos,
) (*TallyOut, error) {

	out := &TallyOut{
		ReferendumRepo:   x.GovConfig.CommunityURL,
		ReferendumBranch: in.ReferendumBranch,
		ReferendumTally: proto.GovPollTally{
			Ad:         ad,
			TallyVotes: make(proto.GovTallyVotes, len(userInfo)),
		},
	}

	// snapshot votes from user repos
	// TODO: parallelize snapshots
	govService := gov.GovService{GovConfig: x.GovConfig}
	for i, info := range userInfo {
		out.ReferendumTally.TallyVotes[i].UserName = info.UserName
		out.ReferendumTally.TallyVotes[i].UserPublicURL = info.UserInfo.URL

		// compute the name of the vote branch in the user's repo
		voteBranch, err := proto.PollVoteBranch(ctx, ad)
		if err != nil {
			return nil, err
		}

		snapOut, err := govService.SnapshotBranchLatest(ctx, &gov.SnapshotBranchLatestIn{
			SourceRepo:   info.UserInfo.URL,
			SourceBranch: voteBranch,
			Community:    community,
		})
		if err != nil {
			// skip snapshotting unresponsive user repos
			continue
		}

		// XXX: prepare the votes tally
		_ = snapOut
		// out.ReferendumTally.TallyVotes[i].UserVote = userVoteXXX
	}

	// XXX: aggregate votes to choices

	// write/stage snapshots and tally to community repo
	tallyPath := proto.PollTallyPath(ad.Path)
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
