package arb

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/petar/gov4git/lib/files"
	"github.com/petar/gov4git/lib/git"
	"github.com/petar/gov4git/proto"
)

type GovArbPollIn struct {
	Path            string   `json:"path"` // path where poll will be persisted
	Choices         []string `json:"choices"`
	Group           string   `json:"group"`
	Strategy        string   `json:"strategy"`
	GoverningBranch string   `json:"governing_branch"`
}

func (x *GovArbPollIn) Sanitize() error {
	// sanitize path
	x.Path = git.MakeNonAbs(x.Path)
	if x.Path == "" {
		return fmt.Errorf("missing poll path")
	}
	return nil
}

type GovArbPollOut struct {
	CommunityURL      string `json:"community_url"`
	GoverningBranch   string `json:"governing_branch"`
	Path              string `json:"path"`
	PollBranch        string `json:"poll_branch"`
	PollGenesisCommit string `json:"poll_genesis_commit"`
}

func (x GovArbPollOut) Human(context.Context) string {
	return fmt.Sprintf(`
community_url=%v
governing_branch=%v
poll_path=%v
poll_branch=%v
poll_genesis_commit=%v

Vote using:

   gov4git vote --community=%v --branch=%v
`,
		x.CommunityURL, x.GoverningBranch, x.Path, x.PollBranch, x.PollGenesisCommit,
		x.CommunityURL, x.PollBranch,
	)
}

func (x GovArbService) Poll(ctx context.Context, in *GovArbPollIn) (*GovArbPollOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
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

func (x GovArbService) PollLocal(ctx context.Context, community git.Local, in *GovArbPollIn) (*GovArbPollOut, error) {
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
	pollAdPath := filepath.Join(in.Path, proto.GovPollAdFilebase)
	pollAd := proto.GovPollAd{
		Path:         in.Path,
		Choices:      in.Choices,
		Group:        in.Group,
		Strategy:     in.Strategy,
		Branch:       in.GoverningBranch,
		ParentCommit: head,
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
	out := &GovArbPollOut{
		CommunityURL:      x.GovConfig.CommunityURL,
		GoverningBranch:   in.GoverningBranch,
		Path:              pollAd.Path,
		PollBranch:        pollBranch,
		PollGenesisCommit: "", // populate after commit
	}
	hum := fmt.Sprintf(proto.PollGenesisCommitHeader(pollBranch)+`

Vote using:

   gov4git vote --community=%v --branch=%v

   `, out.Path, out.PollBranch, out.CommunityURL, out.PollBranch)
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
