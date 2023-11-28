package account

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/kv"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

type AccountID string

func (x AccountID) String() string {
	return string(x)
}

type Account struct {
	ID      AccountID `json:"id"`
	Owner   OwnerID   `json:"owner"`
	Balance Holding   `json:"balance"`
}

func NewAccount(id AccountID, owner OwnerID, balance Holding) Account {
	return Account{ID: id, Owner: owner, Balance: balance}
}

var (
	accountKV = kv.KV[AccountID, Account]{}
	accountNS = proto.RootNS.Append("account")
)

func Create(
	ctx context.Context,
	addr gov.Address,
	id AccountID,
	owner OwnerID,
	holding Holding,
) {
	cloned := gov.Clone(ctx, addr)
	Create_StageOnly(ctx, cloned, id, owner, holding)
	proto.Commitf(ctx, cloned, "account_create", "create account %v", id)
	cloned.Push(ctx)
}

func Create_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	id AccountID,
	owner OwnerID,
	balance Holding,
) {
	must.Assertf(ctx, !Exists_Local(ctx, cloned, id), "account %v already exists", id)
	set_StageOnly(ctx, cloned, id, NewAccount(id, owner, balance))
}

func Exists_Local(
	ctx context.Context,
	cloned gov.Cloned,
	id AccountID,
) bool {
	err := must.Try(func() { Get_Local(ctx, cloned, id) })
	switch {
	case err == nil:
		return true
	case git.IsNotExist(err):
		return false
	default:
		must.Panic(ctx, err)
		return false
	}
}

func Get(
	ctx context.Context,
	addr gov.Address,
	id AccountID,
) Account {
	cloned := gov.Clone(ctx, addr)
	return Get_Local(ctx, cloned, id)
}

func Get_Local(
	ctx context.Context,
	cloned gov.Cloned,
	id AccountID,
) Account {
	return accountKV.Get(ctx, accountNS, cloned.Tree(), id)
}

func Transfer(
	ctx context.Context,
	addr gov.Address,
	from AccountID,
	to AccountID,
	amount Holding,
) {
	cloned := gov.Clone(ctx, addr)
	Transfer_StageOnly(ctx, cloned, from, to, amount)
	proto.Commitf(ctx, cloned, "account_transfer", "transfer %v from %v to %v", amount, from, to)
	cloned.Push(ctx)
}

func Transfer_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	from AccountID,
	to AccountID,
	amount Holding,
) {
	Withdraw_StageOnly(ctx, cloned, from, amount)
	Deposit_StageOnly(ctx, cloned, to, amount)
}

func Deposit(
	ctx context.Context,
	addr gov.Address,
	to AccountID,
	amount Holding,
) {
	cloned := gov.Clone(ctx, addr)
	Deposit_StageOnly(ctx, cloned, to, amount)
	proto.Commitf(ctx, cloned, "account_deposit", "deposit %v to %v", amount, to)
	cloned.Push(ctx)
}

func Deposit_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	to AccountID,
	amount Holding,
) {
	a := Get_Local(ctx, cloned, to)
	a.Balance = SumHolding(ctx, a.Balance, amount)
	set_StageOnly(ctx, cloned, to, a)
}

func Withdraw(
	ctx context.Context,
	addr gov.Address,
	from AccountID,
	amount Holding,
) {
	cloned := gov.Clone(ctx, addr)
	Withdraw_StageOnly(ctx, cloned, from, amount)
	proto.Commitf(ctx, cloned, "account_withdraw", "withdraw %v from %v", amount, from)
	cloned.Push(ctx)
}

func Withdraw_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	from AccountID,
	amount Holding,
) {
	a := Get_Local(ctx, cloned, from)
	d := SumHolding(ctx, a.Balance, NegHolding(amount))
	must.Assertf(ctx, d.Quantity >= 0, "insufficient funds")
	a.Balance = d
	set_StageOnly(ctx, cloned, from, a)
}

func set_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	id AccountID,
	account Account,
) {
	accountKV.Set(ctx, accountNS, cloned.Tree(), id, account)
}
