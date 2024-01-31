package motionapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/lib4git/base"
)

func Pipeline(
	ctx context.Context,
	cloned gov.OwnerCloned,

) {

	// update and aggregate motion policies
	for i := 0; i < 2; i++ {
		base.Infof("PIPELINE: updating motions")
		UpdateMotions_StageOnly(ctx, cloned)
		base.Infof("PIPELINE: aggregating motions")
		AggregateMotions_StageOnly(ctx, cloned)
	}

	// rescore motions to capture updated tallies
	base.Infof("PIPELINE: scoring motions")
	ScoreMotions_StageOnly(ctx, cloned)

	// clearance
	base.Infof("PIPELINE: clear motions")
	ClearMotions_StageOnly(ctx, cloned)

	// archive closed motions
	base.Infof("PIPELINE: archive motions")
	ArchiveMotions_StageOnly(ctx, cloned)

}
