package test

import (
	"context"
	"strconv"
	"testing"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/gov4git/mod/id"
	"github.com/gov4git/gov4git/mod/member"
)

type TestCommunity struct {
	community gov.CommunityAddress
	organizer gov.OrganizerAddress
	members   []id.OwnerAddress
}

func NewTestCommunity(t *testing.T, ctx context.Context, numMembers int) *TestCommunity {

	// initialize organizer and community
	organizerID := id.NewTestID(ctx, t, git.MainBranch, true)
	id.Init(ctx, organizerID.OwnerAddress())
	base.Infof("gov_public=%v gov_private=%v", organizerID.PublicAddress(), organizerID.PrivateAddress())

	// initialize members
	members := make([]id.OwnerAddress, numMembers)
	for i := 0; i < numMembers; i++ {
		memberID := id.NewTestID(ctx, t, git.MainBranch, true)
		base.Infof("member_%d_public=%v member_%d_private=%v",
			i, organizerID.PublicAddress(), i, organizerID.PrivateAddress())
		id.Init(ctx, memberID.OwnerAddress())
		members[i] = memberID.OwnerAddress()
	}

	comty := &TestCommunity{
		community: gov.CommunityAddress(organizerID.PublicAddress()),
		organizer: gov.OrganizerAddress(organizerID.OwnerAddress()),
		members:   members,
	}

	comty.addEverybody(t, ctx)

	return comty
}

func (x *TestCommunity) addEverybody(t *testing.T, ctx context.Context) {

	govRepo, govTree := git.Clone(ctx, git.Address(x.community))

	for i, m := range x.members {
		member.AddUser(ctx, govTree, x.Member(i), member.Account{Home: m.Public})
	}

	git.Commit(ctx, govTree, "add everybody")
	git.Push(ctx, govRepo)
}

func (x *TestCommunity) Member(i int) member.User {
	return member.User("m" + strconv.Itoa(i))
}
