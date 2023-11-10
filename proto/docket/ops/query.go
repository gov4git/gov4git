package ops

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func IsMotion(ctx context.Context, addr gov.Address, id schema.MotionID) bool {
	return IsMotion_Local(ctx, gov.Clone(ctx, addr).Tree(), id)
}

func IsMotion_Local(ctx context.Context, t *git.Tree, id schema.MotionID) bool {
	err := must.Try(func() { schema.MotionKV.Get(ctx, schema.MotionNS, t, id) })
	return err == nil
}
