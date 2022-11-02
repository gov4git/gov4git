package user

import (
	"context"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type BalanceSetIn struct {
	User    string  `json:"user"`
	Balance string  `json:"balance"`
	Value   float64 `json:"value"`
	Branch  string  `json:"branch"` // branch in community repo
}

type BalanceSetOut struct{}

func (x UserService) BalanceSet(
	ctx context.Context,
	user string,
	balance string,
	value float64,
) error {
	local, err := git.CloneOrigin(ctx, git.Origin(x))
	if err != nil {
		return err
	}
	if err := x.SetFloat64Local(ctx, local, user, govproto.BalanceKey(balance), value); err != nil {
		return err
	}
	if err := local.PushUpstream(ctx); err != nil {
		return err
	}
	return nil
}

func (x UserService) BalanceGet(
	ctx context.Context,
	user string,
	balance string,
) (float64, error) {
	local, err := git.CloneOrigin(ctx, git.Origin(x))
	if err != nil {
		return 0, err
	}
	value, err := x.GetFloat64Local(ctx, local, user, govproto.BalanceKey(balance))
	if err != nil {
		return 0, err
	}
	return value, nil
}

type BalanceAddResult struct {
	ValueBefore float64 `json:"value_before"`
	ValueAfter  float64 `json:"value_after"`
}

func (x UserService) BalanceAdd(
	ctx context.Context,
	user string,
	balance string,
	value float64,
) (*BalanceAddResult, error) {
	local, err := git.CloneOrigin(ctx, git.Origin(x))
	if err != nil {
		return nil, err
	}
	valueBefore, err := x.GetFloat64Local(ctx, local, user, govproto.BalanceKey(balance))
	if err != nil {
		return nil, err
	}
	valueAfter := valueBefore + value
	if err := x.SetFloat64Local(ctx, local, user, govproto.BalanceKey(balance), valueAfter); err != nil {
		return nil, err
	}
	if err := local.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &BalanceAddResult{ValueBefore: valueBefore, ValueAfter: valueAfter}, nil
}

type BalanceMulResult struct {
	ValueBefore float64 `json:"value_before"`
	ValueAfter  float64 `json:"value_after"`
}

func (x UserService) BalanceMul(
	ctx context.Context,
	user string,
	balance string,
	value float64,
) (*BalanceMulResult, error) {
	local, err := git.CloneOrigin(ctx, git.Origin(x))
	if err != nil {
		return nil, err
	}
	valueBefore, err := x.GetFloat64Local(ctx, local, user, govproto.BalanceKey(balance))
	if err != nil {
		return nil, err
	}
	valueAfter := valueBefore * value
	if err := x.SetFloat64Local(ctx, local, user, govproto.BalanceKey(balance), valueAfter); err != nil {
		return nil, err
	}
	if err := local.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &BalanceMulResult{ValueBefore: valueBefore, ValueAfter: valueAfter}, nil
}
