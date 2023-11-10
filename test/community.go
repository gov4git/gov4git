package test

import (
	"context"
	"strconv"
	"testing"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/boot"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	_ "github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/git"
)

type TestCommunity struct {
	gov           gov.Address
	organizer     gov.OwnerAddress
	members       []id.OwnerAddress
	silentMembers []id.OwnerAddress
}

func NewTestCommunity(t *testing.T, ctx context.Context, numMembers int) *TestCommunity {

	// initialize organizer and community
	organizerID := id.NewTestID(ctx, t, git.MainBranch, true)
	boot.Boot(ctx, organizerID.OwnerAddress())
	base.Infof("gov_public=%v gov_private=%v", organizerID.PublicAddress(), organizerID.PrivateAddress())

	// initialize members
	members := make([]id.OwnerAddress, numMembers)
	for i := 0; i < numMembers; i++ {
		memberID := id.NewTestID(ctx, t, git.MainBranch, true)
		base.Infof("member_%d_public=%v member_%d_private=%v",
			i, memberID.Public.Address(), i, memberID.Private.Address())
		id.Init(ctx, memberID.OwnerAddress())
		members[i] = memberID.OwnerAddress()
	}

	// initialize silent members (their repos become unreachable after joining the community)
	silentMembers := make([]id.OwnerAddress, numMembers)
	silentTestIDs := make([]id.TestID, numMembers)
	for i := 0; i < numMembers; i++ {
		memberID := id.NewTestID(ctx, t, git.MainBranch, true)
		base.Infof("silent_member_%d_public=%v silent_member_%d_private=%v",
			i, memberID.Public.Address(), i, memberID.Private.Address())
		id.Init(ctx, memberID.OwnerAddress())
		silentMembers[i] = memberID.OwnerAddress()
		silentTestIDs[i] = memberID
	}

	comty := &TestCommunity{
		gov:           gov.Address(organizerID.PublicAddress()),
		organizer:     gov.OwnerAddress(organizerID.OwnerAddress()),
		members:       members,
		silentMembers: silentMembers,
	}

	comty.addEverybody(t, ctx)

	// erase silent memers' repos
	for _, tid := range silentTestIDs {
		tid.Erase(ctx)
	}

	return comty
}

func (x *TestCommunity) addEverybody(t *testing.T, ctx context.Context) {

	govCloned := git.CloneOne(ctx, git.Address(x.gov))

	for i, m := range x.members {
		member.AddUserByPublicAddress_StageOnly(ctx, govCloned.Tree(), x.MemberUser(i), m.Public)
	}

	for i, m := range x.silentMembers {
		member.AddUserByPublicAddress_StageOnly(ctx, govCloned.Tree(), x.InvalidMemberUser(i), m.Public)
	}

	chg := git.NewChangeNoResult("add everybody", "test_add_everybody")
	proto.Commit(ctx, govCloned.Tree(), chg)
	govCloned.Push(ctx)
}

func (x *TestCommunity) Gov() gov.Address {
	return x.gov
}

func (x *TestCommunity) Organizer() gov.OwnerAddress {
	return x.organizer
}

func (x *TestCommunity) MemberUser(i int) member.User {
	return member.User("member_" + strconv.Itoa(i))
}

func (x *TestCommunity) MemberOwner(i int) id.OwnerAddress {
	return x.members[i]
}

func (x *TestCommunity) InvalidMemberUser(i int) member.User {
	return member.User("silent_member_" + strconv.Itoa(i))
}

func (x *TestCommunity) InvalidMemberOwner(i int) id.OwnerAddress {
	return x.silentMembers[i]
}

func (x *TestCommunity) NonExistentMemberUser() member.User {
	return member.User("non_existent_member")
}
