package account

import (
	"context"
	"fmt"

	"github.com/gov4git/lib4git/must"
)

type Holding struct {
	Asset    Asset   `json:"asset"`
	Quantity float64 `json:"quantity"`
}

func (h Holding) String() string {
	return fmt.Sprintf("%v:%v", h.Asset, h.Quantity)
}

func ZeroHolding(asset Asset) Holding {
	return Holding{
		Asset:    asset,
		Quantity: 0,
	}
}

func NegHolding(p Holding) Holding {
	return Holding{
		Asset:    p.Asset,
		Quantity: -p.Quantity,
	}
}

func SumHolding(ctx context.Context, p, q Holding) Holding {
	must.Assertf(ctx, p.Asset == q.Asset, "cannot add different assets")
	return Holding{
		Asset:    p.Asset,
		Quantity: p.Quantity + q.Quantity,
	}
}
