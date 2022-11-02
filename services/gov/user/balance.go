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

func (x GovUserService) BalanceSet(ctx context.Context, in *BalanceSetIn) (*BalanceSetOut, error) {
	community, err := git.CloneBranch(ctx, x.GovConfig.CommunityURL, in.Branch)
	if err != nil {
		return nil, err
	}
	if err := x.SetFloat64Local(ctx, community, in.User, govproto.BalanceKey(in.Balance), in.Value); err != nil {
		return nil, err
	}
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &BalanceSetOut{}, nil
}

type BalanceGetIn struct {
	User    string `json:"user"`
	Balance string `json:"balance"`
	Branch  string `json:"branch"`
}

type BalanceGetOut struct {
	Value float64 `json:"value"`
}

func (x GovUserService) BalanceGet(ctx context.Context, in *BalanceGetIn) (*BalanceGetOut, error) {
	community, err := git.CloneBranch(ctx, x.GovConfig.CommunityURL, in.Branch)
	if err != nil {
		return nil, err
	}
	value, err := x.GetFloat64Local(ctx, community, in.User, govproto.BalanceKey(in.Balance))
	if err != nil {
		return nil, err
	}
	return &BalanceGetOut{Value: value}, nil
}

type BalanceAddIn struct {
	User    string  `json:"user"`
	Balance string  `json:"balance"`
	Branch  string  `json:"branch"`
	Value   float64 `json:"value"`
}

type BalanceAddOut struct {
	ValueBefore float64 `json:"value_before"`
	ValueAfter  float64 `json:"value_after"`
}

func (x GovUserService) BalanceAdd(ctx context.Context, in *BalanceAddIn) (*BalanceAddOut, error) {
	community, err := git.CloneBranch(ctx, x.GovConfig.CommunityURL, in.Branch)
	if err != nil {
		return nil, err
	}
	valueBefore, err := x.GetFloat64Local(ctx, community, in.User, govproto.BalanceKey(in.Balance))
	if err != nil {
		return nil, err
	}
	valueAfter := valueBefore + in.Value
	if err := x.SetFloat64Local(ctx, community, in.User, govproto.BalanceKey(in.Balance), valueAfter); err != nil {
		return nil, err
	}
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &BalanceAddOut{ValueBefore: valueBefore, ValueAfter: valueAfter}, nil
}

type BalanceMulIn struct {
	User    string  `json:"user"`
	Balance string  `json:"balance"`
	Branch  string  `json:"branch"`
	Value   float64 `json:"value"`
}

type BalanceMulOut struct {
	ValueBefore float64 `json:"value_before"`
	ValueAfter  float64 `json:"value_after"`
}

func (x GovUserService) BalanceMul(ctx context.Context, in *BalanceMulIn) (*BalanceMulOut, error) {
	community, err := git.CloneBranch(ctx, x.GovConfig.CommunityURL, in.Branch)
	if err != nil {
		return nil, err
	}
	valueBefore, err := x.GetFloat64Local(ctx, community, in.User, govproto.BalanceKey(in.Balance))
	if err != nil {
		return nil, err
	}
	valueAfter := valueBefore * in.Value
	if err := x.SetFloat64Local(ctx, community, in.User, govproto.BalanceKey(in.Balance), valueAfter); err != nil {
		return nil, err
	}
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &BalanceMulOut{ValueBefore: valueBefore, ValueAfter: valueAfter}, nil
}
