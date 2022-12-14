package test

import (
	"context"
	"strconv"
	"testing"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/git"
)

type TestCommunity struct {
	gov       gov.GovAddress
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
		gov:       gov.GovAddress(organizerID.PublicAddress()),
		organizer: gov.OrganizerAddress(organizerID.OwnerAddress()),
		members:   members,
	}

	comty.addEverybody(t, ctx)

	return comty
}

func (x *TestCommunity) addEverybody(t *testing.T, ctx context.Context) {

	govCloned := git.Clone(ctx, git.Address(x.gov))

	for i, m := range x.members {
		member.AddUserByPublicAddressStageOnly(ctx, govCloned.Tree(), x.MemberUser(i), m.Public)
	}

	proto.Commit(ctx, govCloned.Tree(), "add everybody")
	govCloned.Push(ctx)
}

func (x *TestCommunity) Gov() gov.GovAddress {
	return x.gov
}

func (x *TestCommunity) Organizer() gov.OrganizerAddress {
	return x.organizer
}

func (x *TestCommunity) MemberUser(i int) member.User {
	return member.User("m" + strconv.Itoa(i))
}

func (x *TestCommunity) MemberOwner(i int) id.OwnerAddress {
	return x.members[i]
}
