package sv

import (
	"context"
)

func (qv SV) CalcJS(
	ctx context.Context,
) string {

	return qv.Kernel.CalcJS(ctx)
}
