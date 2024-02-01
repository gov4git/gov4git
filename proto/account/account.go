package account

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/metric"
	"github.com/gov4git/gov4git/v2/proto/history/trace"
	"github.com/gov4git/gov4git/v2/proto/kv"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

type AccountID string

func AccountIDFromNS(p ns.NS) AccountID {
	return AccountID(p.GitPath())
}

func AccountIDFromLine(line Line) AccountID {
	return AccountID(line)
}

func (x AccountID) String() string {
	return string(x)
}

func (x AccountID) MetricAccountID() metric.AccountID {
	return metric.AccountID(x)
}

type Account struct {
	ID     AccountID     `json:"id"`
	Owner  AccountID     `json:"owner"`
	Assets AssetHoldings `json:"assets"`
}

func (a *Account) DepositOverDraft(ctx context.Context, h Holding) {
	a.Assets.DepositOverDraft(ctx, h)
}

func (a *Account) Deposit(ctx context.Context, h Holding) {
	a.Assets.Deposit(ctx, h)
}

func (a *Account) WithdrawOverDraft(ctx context.Context, h Holding) {
	a.Assets.WithdrawOverDraft(ctx, h)
}

func (a *Account) Withdraw(ctx context.Context, h Holding) {
	a.Assets.Withdraw(ctx, h)
}

func (a *Account) Balance(asset Asset) Holding {
	return a.Assets.Balance(asset)
}

type AssetHoldings map[Asset]Holding

func (x AssetHoldings) Balance(asset Asset) Holding {
	if h, ok := x[asset]; ok {
		return h
	} else {
		return ZeroHolding(asset)
	}
}

func (x AssetHoldings) DepositOverDraft(ctx context.Context, h Holding) {
	if g, ok := x[h.Asset]; ok {
		x[h.Asset] = SumHolding(ctx, g, h)
	} else {
		x[h.Asset] = h
	}
}

func (x AssetHoldings) Deposit(ctx context.Context, h Holding) {
	if g, ok := x[h.Asset]; ok {
		d := SumHolding(ctx, g, h)
		must.Assertf(ctx, d.Quantity >= 0, "insufficient funds")
		x[h.Asset] = d
	} else {
		d := h
		must.Assertf(ctx, d.Quantity >= 0, "no funds")
		x[h.Asset] = d
	}
}

func (x AssetHoldings) WithdrawOverDraft(ctx context.Context, h Holding) {
	x.DepositOverDraft(ctx, NegHolding(h))
}

func (x AssetHoldings) Withdraw(ctx context.Context, h Holding) {
	x.Deposit(ctx, NegHolding(h))
}

func NewAccount(id AccountID, owner AccountID) *Account {
	return &Account{
		ID:     id,
		Owner:  owner,
		Assets: AssetHoldings{},
	}
}

var (
	accountKV = kv.KV[AccountID, *Account]{}
	accountNS = proto.RootNS.Append("account")
)

func Create(
	ctx context.Context,
	addr gov.Address,
	id AccountID,
	owner AccountID,
	note string,

) {
	cloned := gov.Clone(ctx, addr)
	Create_StageOnly(ctx, cloned, id, owner, note)
	proto.Commitf(ctx, cloned, "account_create", "create account %v (%v)", id, note)
	cloned.Push(ctx)
}

func Create_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	id AccountID,
	owner AccountID,
	note string,

) {
	must.Assertf(ctx, !Exists_Local(ctx, cloned, id), "account %v already exists", id)
	set_StageOnly(ctx, cloned, id, NewAccount(id, owner))
	trace.Log_StageOnly(ctx, cloned, &trace.Event{
		Op:     "account_create",
		Note:   note,
		Args:   trace.M{"id": id, "owner": owner},
		Result: nil,
	})
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

) *Account {
	cloned := gov.Clone(ctx, addr)
	return Get_Local(ctx, cloned, id)
}

func Get_Local(
	ctx context.Context,
	cloned gov.Cloned,
	id AccountID,

) *Account {
	return accountKV.Get(ctx, accountNS, cloned.Tree(), id)
}

func Transfer(
	ctx context.Context,
	addr gov.Address,
	from AccountID,
	to AccountID,
	amount Holding,
	note string,

) {
	cloned := gov.Clone(ctx, addr)
	Transfer_StageOnly(ctx, cloned, from, to, amount, note)
	proto.Commitf(ctx, cloned, "account_transfer", "transfer %v from %v to %v (%v)", amount, from, to, note)
	cloned.Push(ctx)
}

func Transfer_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	fromID AccountID,
	toID AccountID,
	amount Holding,
	note string,

) {

	from := Get_Local(ctx, cloned, fromID)
	from.Withdraw(ctx, amount)

	to := Get_Local(ctx, cloned, toID)
	to.Deposit(ctx, amount)

	set_StageOnly(ctx, cloned, fromID, from)
	set_StageOnly(ctx, cloned, toID, to)

	trace.Log_StageOnly(ctx, cloned, &trace.Event{
		Op:     "account_transfer",
		Note:   note,
		Args:   trace.M{"from": fromID, "to": toID, "amount": amount},
		Result: nil,
	})
	metric.Log_StageOnly(ctx, cloned, &metric.Event{
		Account: &metric.AccountEvent{
			Transfer: &metric.AccountTransferEvent{
				From:   fromID.MetricAccountID(),
				To:     toID.MetricAccountID(),
				Amount: amount.MetricHolding(),
			},
		},
	})
}

func TryTransfer_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	from AccountID,
	to AccountID,
	amount Holding,
	note string,

) error {
	return must.Try(func() { Transfer_StageOnly(ctx, cloned, from, to, amount, note) })
}

func Issue(
	ctx context.Context,
	addr gov.Address,
	to AccountID,
	amount Holding,
	note string,

) {

	cloned := gov.Clone(ctx, addr)
	Issue_StageOnly(ctx, cloned, to, amount, note)
	proto.Commitf(ctx, cloned, "account_issue", "issue %v to %v (%v)", amount, to, note)
	cloned.Push(ctx)
}

func Issue_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	toID AccountID,
	amount Holding,
	note string,

) {

	TransferOverDraft_StageOnly(
		metric.Mute(ctx),
		cloned,
		IssueAccountID,
		toID,
		amount,
		note,
	)
	trace.Log_StageOnly(ctx, cloned, &trace.Event{
		Op:     "account_issue",
		Note:   note,
		Args:   trace.M{"from": IssueAccountID, "to": toID, "amount": amount},
		Result: nil,
	})
	metric.Log_StageOnly(ctx, cloned, &metric.Event{
		Account: &metric.AccountEvent{
			Issue: &metric.AccountIssueEvent{
				To:     toID.MetricAccountID(),
				Amount: amount.MetricHolding(),
			},
		},
	})
}

func Burn(
	ctx context.Context,
	addr gov.Address,
	fromID AccountID,
	amount Holding,
	note string,

) {

	cloned := gov.Clone(ctx, addr)
	Burn_StageOnly(ctx, cloned, fromID, amount, note)
	proto.Commitf(ctx, cloned, "account_burn", "burn %v from %v (%v)", amount, fromID, note)
	cloned.Push(ctx)
}

func Burn_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	fromID AccountID,
	amount Holding,
	note string,

) {
	TransferOverDraft_StageOnly(
		metric.Mute(ctx),
		cloned,
		fromID,
		BurnAccountID,
		amount,
		note,
	)
	trace.Log_StageOnly(ctx, cloned, &trace.Event{
		Op:     "account_burn",
		Note:   note,
		Args:   trace.M{"from": fromID, "to": BurnAccountID, "amount": amount},
		Result: nil,
	})
	metric.Log_StageOnly(ctx, cloned, &metric.Event{
		Account: &metric.AccountEvent{
			Burn: &metric.AccountBurnEvent{
				From:   fromID.MetricAccountID(),
				Amount: amount.MetricHolding(),
			},
		},
	})
}

func set_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	id AccountID,
	account *Account,

) {
	accountKV.Set(ctx, accountNS, cloned.Tree(), id, account)
}

func List(
	ctx context.Context,
	addr gov.Address,

) []AccountID {
	cloned := gov.Clone(ctx, addr)
	return List_Local(ctx, cloned)
}

func List_Local(
	ctx context.Context,
	cloned gov.Cloned,

) []AccountID {
	return accountKV.ListKeys(ctx, accountNS, cloned.Tree())
}

func TransferOverDraft_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	fromID AccountID,
	toID AccountID,
	amount Holding,
	note string,

) {

	from := Get_Local(ctx, cloned, fromID)
	from.WithdrawOverDraft(ctx, amount)

	to := Get_Local(ctx, cloned, toID)
	to.DepositOverDraft(ctx, amount)

	set_StageOnly(ctx, cloned, fromID, from)
	set_StageOnly(ctx, cloned, toID, to)

	trace.Log_StageOnly(ctx, cloned, &trace.Event{
		Op:     "account_transfer_overdraft",
		Note:   note,
		Args:   trace.M{"from": fromID, "to": toID, "amount": amount},
		Result: nil,
	})
}

func TryTransferOverDraft_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	from AccountID,
	to AccountID,
	amount Holding,
	note string,

) error {
	return must.Try(func() { TransferOverDraft_StageOnly(ctx, cloned, from, to, amount, note) })
}
