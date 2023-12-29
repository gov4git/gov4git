package sv

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
)

func (qv SV) Margin(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Advertisement,
	current *ballotproto.Tally,

) *ballotproto.Margin {

	return &ballotproto.Margin{
		CalcJS: qv.Kernel.CalcJS(ctx, cloned, ad),
	}

}
