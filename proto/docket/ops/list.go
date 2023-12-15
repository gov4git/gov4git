package ops

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/docket/schema"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/lib4git/git"
)

func ListMotions(ctx context.Context, addr gov.Address) schema.Motions {
	return ListMotions_Local(ctx, gov.Clone(ctx, addr).Tree())
}

func ListMotions_Local(ctx context.Context, t *git.Tree) schema.Motions {
	_, motions := schema.MotionKV.ListKeyValues(ctx, schema.MotionNS, t)
	schema.MotionsByID(motions).Sort()
	return motions
}

func ListMotionViews(ctx context.Context, addr gov.Address) schema.MotionViews {
	return ListMotionViews_Local(ctx, gov.Clone(ctx, addr))
}

func ListMotionViews_Local(ctx context.Context, cloned gov.Cloned) schema.MotionViews {
	t := cloned.Tree()
	ids := schema.MotionKV.ListKeys(ctx, schema.MotionNS, t)
	mvs := make(schema.MotionViews, len(ids))
	for i, id := range ids {
		mvs[i] = ShowMotion_Local(ctx, cloned, id)
	}
	mvs.Sort()
	return mvs
}
