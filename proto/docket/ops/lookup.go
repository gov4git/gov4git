package ops

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/docket/schema"
	"github.com/gov4git/gov4git/v2/proto/gov"
)

func LookupMotion(
	ctx context.Context,
	addr gov.Address,
	id schema.MotionID,
	args ...any,

) schema.Motion {

	return LookupMotion_Local(ctx, gov.Clone(ctx, addr), id, args...)
}

func LookupMotion_Local(
	ctx context.Context,
	cloned gov.Cloned,
	id schema.MotionID,
	args ...any,

) schema.Motion {

	return schema.MotionKV.Get(ctx, schema.MotionNS, cloned.Tree(), id)
}
