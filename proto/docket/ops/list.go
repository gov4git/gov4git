package ops

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
)

func ListMotions(ctx context.Context, addr gov.GovAddress) schema.Motions {
	return ListMotions_Local(ctx, gov.Clone(ctx, addr).Tree())
}

func ListMotions_Local(ctx context.Context, t *git.Tree) schema.Motions {
	_, motions := schema.MotionKV.ListKeyValues(ctx, schema.MotionNS, t)
	return motions
}
