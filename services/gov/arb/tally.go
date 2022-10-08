package arb

import (
	"context"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/services/gov/member"
)

type TallyIn struct {
	ReferendumBranch string `json:"referendum_branch"`
}

type TallyOut struct {
	ReferendumRepo   string `json:"referendum_repo"`
	ReferendumBranch string `json:"referendum_branch"`
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
	community := git.LocalInDir(files.WorkDir(ctx).Subdir("community"))
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
	_ = participants

	//XXX: get user info (public_url, etc.)

	//XXX

	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}

	panic("not implemented")
}