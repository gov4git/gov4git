package entity

import (
	"context"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/entityproto"
)

func (x EntityService[V]) SetBalance(
	ctx context.Context,
	user string,
	balance string,
	value float64,
) error {
	local, err := git.CloneOrigin(ctx, x.Address)
	if err != nil {
		return err
	}
	if err := x.SetPropFloat64Local(ctx, local, user, entityproto.BalancePropKey(balance), value); err != nil {
		return err
	}
	if err := local.PushUpstream(ctx); err != nil {
		return err
	}
	return nil
}

func (x EntityService[V]) GetBalance(
	ctx context.Context,
	user string,
	balance string,
) (float64, error) {
	local, err := git.CloneOrigin(ctx, x.Address)
	if err != nil {
		return 0, err
	}
	value, err := x.GetPropFloat64Local(ctx, local, user, entityproto.BalancePropKey(balance))
	if err != nil {
		return 0, err
	}
	return value, nil
}

type AddBalanceResult struct {
	BalanceBefore float64 `json:"balance_before"`
	BalanceAfter  float64 `json:"balance_after"`
}

func (x EntityService[V]) AddBalance(
	ctx context.Context,
	user string,
	balance string,
	value float64,
) (*AddBalanceResult, error) {
	local, err := git.CloneOrigin(ctx, x.Address)
	if err != nil {
		return nil, err
	}
	valueBefore, err := x.GetPropFloat64Local(ctx, local, user, entityproto.BalancePropKey(balance))
	if err != nil {
		return nil, err
	}
	valueAfter := valueBefore + value
	if err := x.SetPropFloat64Local(ctx, local, user, entityproto.BalancePropKey(balance), valueAfter); err != nil {
		return nil, err
	}
	if err := local.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &AddBalanceResult{BalanceBefore: valueBefore, BalanceAfter: valueAfter}, nil
}

type MulBalanceResult struct {
	BalanceBefore float64 `json:"balance_before"`
	BalanceAfter  float64 `json:"balance_after"`
}

func (x EntityService[V]) MulBalance(
	ctx context.Context,
	user string,
	balance string,
	value float64,
) (*MulBalanceResult, error) {
	local, err := git.CloneOrigin(ctx, x.Address)
	if err != nil {
		return nil, err
	}
	valueBefore, err := x.GetPropFloat64Local(ctx, local, user, entityproto.BalancePropKey(balance))
	if err != nil {
		return nil, err
	}
	valueAfter := valueBefore * value
	if err := x.SetPropFloat64Local(ctx, local, user, entityproto.BalancePropKey(balance), valueAfter); err != nil {
		return nil, err
	}
	if err := local.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &MulBalanceResult{BalanceBefore: valueBefore, BalanceAfter: valueAfter}, nil
}
