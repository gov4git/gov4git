package arb

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/petar/gitty/lib/files"
	"github.com/petar/gitty/lib/git"
	"github.com/petar/gitty/proto"
)

type GovArbPollIn struct {
	Path            string   `json:"path"` // path where poll will be persisted
	Alternatives    []string `json:"alternatives"`
	Group           string   `json:"group"`
	Strategy        string   `json:"strategy"`
	GoverningBranch string   `json:"community_branch"`
}

type GovArbPollOut struct {
	Path   string `json:"path"`
	Branch string `json:"branch"`
	Commit string `json:"commit"`
}

func (x GovArbPollOut) Human(context.Context) string {
	return fmt.Sprintf("poll=%v branch=%v poll-commit=%v")
}

func (x GovArbService) ArbPoll(ctx context.Context, in *GovArbPollIn) (*GovArbPollOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.GoverningBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := GovArbPoll(ctx, community, in); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &GovArbPollOut{}, nil
}

func GovArbPoll(ctx context.Context, community git.Local, in *GovArbPollIn) error {
	// checkout a new poll branch
	// create and stage poll advertisement
	// commit poll advertisement, including poll ad in commit message
	XXX
	groupFile := filepath.Join(proto.GovGroupsDir, name, proto.GovGroupInfoFilebase)
	// write group file
	stage := files.FormFiles{
		files.FormFile{Path: groupFile, Form: proto.GovGroupInfo{}},
	}
	if err := community.Dir().WriteFormFiles(ctx, stage); err != nil {
		return err
	}
	// stage changes
	if err := community.Add(ctx, stage.Paths()); err != nil {
		return err
	}
	// commit changes
	if err := community.Commitf(ctx, "gov: add group %v", name); err != nil {
		return err
	}
	return nil
}
