package balance

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Set(ctx context.Context, addr gov.GovAddress, user member.User, key Balance, value float64) {
	member.SetUserProp(ctx, addr, user, userPropKey(key), value)
}

func SetStageOnly(ctx context.Context, t *git.Tree, user member.User, key Balance, value float64) {
	member.SetUserPropStageOnly(ctx, t, user, userPropKey(key), value)
}

func Get(ctx context.Context, addr gov.GovAddress, user member.User, key Balance) float64 {
	return member.GetUserPropOrDefault(ctx, addr, user, userPropKey(key), 0.0)
}

func GetLocal(ctx context.Context, t *git.Tree, user member.User, key Balance) float64 {
	return member.GetUserPropLocalOrDefault(ctx, t, user, userPropKey(key), 0.0)
}

func TryTransferStageOnly(
	ctx context.Context,
	t *git.Tree,
	fromUser member.User,
	fromBal Balance,
	toUser member.User,
	toBal Balance,
	amount float64,
) error {
	return must.Try(func() { TransferStageOnly(ctx, t, fromUser, fromBal, toUser, toBal, amount) })
}

func TransferStageOnly(
	ctx context.Context,
	t *git.Tree,
	fromUser member.User,
	fromBal Balance,
	toUser member.User,
	toBal Balance,
	amount float64,
) {
	base.Infof("transfering %v units from %v:%v to %v:%v", amount, fromUser, fromBal, toUser, toBal)
	must.Assertf(ctx, amount >= 0, "negative transfer")
	prior := GetLocal(ctx, t, fromUser, fromBal)
	must.Assertf(ctx, prior >= amount, "insufficient balance")
	AddStageOnly(ctx, t, fromUser, fromBal, -amount)
	AddStageOnly(ctx, t, toUser, toBal, amount)
}

func Add(ctx context.Context, addr gov.GovAddress, user member.User, key Balance, value float64) float64 {
	r, t := gov.Clone(ctx, addr)
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

func Mul(ctx context.Context, addr gov.GovAddress, user member.User, key Balance, value float64) float64 {
	r, t := gov.Clone(ctx, addr)
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
