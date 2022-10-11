package arb

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type PollIn struct {
	Path            string   `json:"path"` // path where poll will be persisted
	Choices         []string `json:"choices"`
	Group           string   `json:"group"`
	Strategy        string   `json:"strategy"`
	GoverningBranch string   `json:"governing_branch"`
}

func (x *PollIn) Sanitize() error {
	// sanitize path
	x.Path = git.MakeNonAbs(x.Path)
	if x.Path == "" {
		return fmt.Errorf("missing poll path")
	}
	return nil
}

type PollOut struct {
	CommunityURL      string `json:"community_url"`
	GoverningBranch   string `json:"governing_branch"`
	Path              string `json:"path"`
	PollBranch        string `json:"poll_branch"`
	PollGenesisCommit string `json:"poll_genesis_commit"`
}

func (x GovArbService) Poll(ctx context.Context, in *PollIn) (*PollOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.GoverningBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	out, err := x.PollLocal(ctx, community, in)
	if err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return out, nil
}

func (x GovArbService) PollLocal(ctx context.Context, community git.Local, in *PollIn) (*PollOut, error) {
	// XXX: verify path is not in use
	// XXX: verify poll branch is not in use
	// XXX: verify group exists
	if err := in.Sanitize(); err != nil {
		return nil, err
	}

	// get hash of current commit on branch
	head, err := community.HeadCommitHash(ctx)
	if err != nil {
		return nil, err
	}

	// checkout a new poll branch
	pollBranch := proto.PollBranch(in.Path)
	if err := community.CheckoutNewBranch(ctx, pollBranch); err != nil {
		return nil, err
	}

	// create and stage poll advertisement
	pollAdPath := proto.PollAdPath(in.Path)
	var pollAd proto.GovPollAd
	switch in.Strategy {
	case "prioritize":
		pollAd = proto.GovPollAd{
			Path:         in.Path,
			Choices:      in.Choices,
			Group:        in.Group,
			Strategy:     proto.GovPollStrategy{Prioritize: &proto.GovPollStrategyPrioritize{}},
			Branch:       in.GoverningBranch,
			ParentCommit: head,
		}
	default:
		return nil, fmt.Errorf("unknown poll strategy %v", in.Strategy)
	}
	stage := files.FormFiles{
		files.FormFile{
			Path: pollAdPath,
			Form: pollAd,
		},
	}
	if err := community.Dir().WriteFormFiles(ctx, stage); err != nil {
		return nil, err
	}
	if err := community.Add(ctx, stage.Paths()); err != nil {
		return nil, err
	}

	// commit changes and include poll ad in commit message
	out := &PollOut{
		CommunityURL:      x.GovConfig.CommunityURL,
		GoverningBranch:   in.GoverningBranch,
		Path:              pollAd.Path,
		PollBranch:        pollBranch,
		PollGenesisCommit: "", // populate after commit
	}
	hum := fmt.Sprintf(proto.PollGenesisCommitHeader(pollBranch)+`

Vote using:

   gov4git vote --community=%v --branch=%v

   `, out.CommunityURL, out.PollBranch)
	msg, err := git.PrepareCommitMsg(ctx, hum, pollAd)
	if err != nil {
		return nil, err
	}
	if err := community.Commit(ctx, msg); err != nil {
		return nil, err
	}

	// get hash of poll genesis commit
	if out.PollGenesisCommit, err = community.HeadCommitHash(ctx); err != nil {
		return nil, err
	}

	return out, nil
}
