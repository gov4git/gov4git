package balance

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/gov4git/mod/member"
	"github.com/gov4git/lib4git/git"
)

func Set(ctx context.Context, addr gov.CommunityAddress, user member.User, key Balance, value float64) {
	member.SetUserProp(ctx, addr, user, userPropKey(key), value)
}

func SetStageOnly(ctx context.Context, t *git.Tree, user member.User, key Balance, value float64) {
	member.SetUserPropStageOnly(ctx, t, user, userPropKey(key), value)
}

func Get(ctx context.Context, addr gov.CommunityAddress, user member.User, key Balance) float64 {
	return member.GetUserPropOrDefault(ctx, addr, user, userPropKey(key), 0.0)
}

func GetLocal(ctx context.Context, t *git.Tree, user member.User, key Balance) float64 {
	return member.GetUserPropLocalOrDefault(ctx, t, user, userPropKey(key), 0.0)
}

func Add(ctx context.Context, addr gov.CommunityAddress, user member.User, key Balance, value float64) float64 {
	r, t := gov.CloneCommunity(ctx, addr)
	prior := AddStageOnly(ctx, t, user, key, value)
	git.Commit(ctx, t, fmt.Sprintf("Add %v to balance %v of user %v", value, key, user))
	git.Push(ctx, r)
	return prior
}

func AddStageOnly(ctx context.Context, t *git.Tree, user member.User, key Balance, value float64) float64 {
	prior := GetLocal(ctx, t, user, key)
	SetStageOnly(ctx, t, user, key, prior+value)
	return prior
}

func Mul(ctx context.Context, addr gov.CommunityAddress, user member.User, key Balance, value float64) float64 {
	r, t := gov.CloneCommunity(ctx, addr)
	prior := MulStageOnly(ctx, t, user, key, value)
	git.Commit(ctx, t, fmt.Sprintf("Multiply %v into balance %v of user %v", value, key, user))
	git.Push(ctx, r)
	return prior
}

func MulStageOnly(ctx context.Context, t *git.Tree, user member.User, key Balance, value float64) float64 {
	prior := GetLocal(ctx, t, user, key)
	SetStageOnly(ctx, t, user, key, prior*value)
	return prior
}
