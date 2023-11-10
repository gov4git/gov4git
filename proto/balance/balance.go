package balance

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Set(ctx context.Context, addr gov.GovAddress, user member.User, key Balance, value float64) {
	member.SetUserProp(ctx, addr, user, userPropKey(key), value)
}

func Set_StageOnly(ctx context.Context, t *git.Tree, user member.User, key Balance, value float64) {
	member.SetUserProp_StageOnly(ctx, t, user, userPropKey(key), value)
}

func Get(ctx context.Context, addr gov.GovAddress, user member.User, key Balance) float64 {
	return member.GetUserPropOrDefault[float64](ctx, addr, user, userPropKey(key), 0.0)
}

func Get_Local(ctx context.Context, t *git.Tree, user member.User, key Balance) float64 {
	return member.GetUserPropOrDefault_Local[float64](ctx, t, user, userPropKey(key), 0.0)
}

func TryTransfer_StageOnly(
	ctx context.Context,
	t *git.Tree,
	fromUser member.User,
	fromBal Balance,
	toUser member.User,
	toBal Balance,
	amount float64,
) error {
	return must.Try(func() { Transfer_StageOnly(ctx, t, fromUser, fromBal, toUser, toBal, amount) })
}

func Transfer_StageOnly(
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
	prior := Get_Local(ctx, t, fromUser, fromBal)
	must.Assertf(ctx, prior >= amount, "insufficient balance")
	Add_StageOnly(ctx, t, fromUser, fromBal, -amount)
	Add_StageOnly(ctx, t, toUser, toBal, amount)
}

func TryCharge_StageOnly(
	ctx context.Context,
	t *git.Tree,
	user member.User,
	bal Balance,
	amount float64,
) error {
	return must.Try(func() { Charge_StageOnly(ctx, t, user, bal, amount) })
}

func Charge_StageOnly(
	ctx context.Context,
	t *git.Tree,
	user member.User,
	bal Balance,
	amount float64,
) {
	base.Infof("charging %v units from %v:%v", amount, user, bal)
	prior := Get_Local(ctx, t, user, bal)
	must.Assertf(ctx, prior >= amount, "insufficient balance")
	Add_StageOnly(ctx, t, user, bal, -amount)
}

func Add(ctx context.Context, addr gov.GovAddress, user member.User, key Balance, value float64) float64 {
	cloned := gov.Clone(ctx, addr)
	prior := Add_StageOnly(ctx, cloned.Tree(), user, key, value)
	chg := git.NewChange[form.Map, float64](
		fmt.Sprintf("Add %v to balance %v of user %v", value, key, user),
		"balance_add",
		form.Map{"user": user, "balance": key, "value": value},
		prior,
		nil,
	)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
	return prior
}

func Add_StageOnly(ctx context.Context, t *git.Tree, user member.User, key Balance, value float64) float64 {
	prior := Get_Local(ctx, t, user, key)
	Set_StageOnly(ctx, t, user, key, prior+value)
	return prior
}

func Mul(ctx context.Context, addr gov.GovAddress, user member.User, key Balance, value float64) float64 {
	cloned := gov.Clone(ctx, addr)
	prior := Mul_StageOnly(ctx, cloned.Tree(), user, key, value)
	chg := git.NewChange[form.Map, float64](
		fmt.Sprintf("Multiply %v into balance %v of user %v", value, key, user),
		"balance_mul",
		form.Map{"user": user, "balance": key, "value": value},
		prior,
		nil,
	)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
	return prior
}

func Mul_StageOnly(ctx context.Context, t *git.Tree, user member.User, key Balance, value float64) float64 {
	prior := Get_Local(ctx, t, user, key)
	Set_StageOnly(ctx, t, user, key, prior*value)
	return prior
}
