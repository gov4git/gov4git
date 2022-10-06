package member

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type GovMemberListIn struct {
	User            string `json:"user"`
	Group           string `json:"group"`
	CommunityBranch string `json:"community_branch"` // branch in community repo where group will be added
}

type GovMemberListOut struct {
	Memberships []GovMemberListMembership `json:"membership"`
}

type GovMemberListMembership struct {
	User  string `json:"user"`
	Group string `json:"group"`
}

func (x GovMemberListOut) Human(context.Context) string {
	var w bytes.Buffer
	for _, u := range x.Memberships {
		fmt.Fprintln(&w, u.User, u.Group)
	}
	return w.String()
}

func (x GovMemberService) MemberList(ctx context.Context, in *GovMemberListIn) (*GovMemberListOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	memberships, err := GovMemberList(ctx, community, in.User, in.Group)
	if err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &GovMemberListOut{Memberships: memberships}, nil
}

func GovMemberList(ctx context.Context, community git.Local, user string, group string) (memberships []GovMemberListMembership, err error) {
	userGlob, groupGlob := user, group
	if user == "" {
		userGlob = "*"
	}
	if group == "" {
		groupGlob = "*"
	}
	userFileGlob := filepath.Join(proto.GovGroupsDir, groupGlob, proto.GovMembersDirbase, userGlob)
	// glob for group files
	m, err := community.Dir().Glob(userFileGlob)
	if err != nil {
		return nil, err
	}
	// extract user names
	memberships = make([]GovMemberListMembership, len(m))
	for i := range m {
		x1, _ := filepath.Split(m[i])
		x2, _ := filepath.Split(x1)
		memberships[i].User = filepath.Base(m[i])
		memberships[i].Group = filepath.Base(x2)
	}
	return memberships, nil
}
