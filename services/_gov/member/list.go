package member

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type ListIn struct {
	User            string `json:"user"`
	Group           string `json:"group"`
	CommunityBranch string `json:"community_branch"`
}

type ListOut struct {
	Memberships []ListMembership `json:"membership"`
}

type ListMembership struct {
	User  string `json:"user"`
	Group string `json:"group"`
}

func ExtractUsersFromMembership(ms []ListMembership) []string {
	u := make([]string, len(ms))
	for i, m := range ms {
		u[i] = m.User
	}
	return u
}

func (x GovMemberService) List(ctx context.Context, in *ListIn) (*ListOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	memberships, err := x.ListLocal(ctx, community, in.User, in.Group)
	if err != nil {
		return nil, err
	}
	return &ListOut{Memberships: memberships}, nil
}

func (x GovMemberService) ListLocal(ctx context.Context, community git.Local, user string, group string) (memberships []ListMembership, err error) {
	userGlob, groupGlob := user, group
	if user == "" {
		userGlob = "*"
	}
	if group == "" {
		groupGlob = "*"
	}
	userFileGlob := govproto.GroupMemberFilepath(groupGlob, userGlob)
	// glob for group files
	m, err := community.Dir().Glob(userFileGlob)
	if err != nil {
		return nil, err
	}
	// extract user names
	memberships = make([]ListMembership, len(m))
	for i := range m {
		x1, _ := filepath.Split(m[i])
		x2, _ := filepath.Split(x1)
		memberships[i].User = filepath.Base(m[i])
		memberships[i].Group = filepath.Base(x2)
	}
	return memberships, nil
}
