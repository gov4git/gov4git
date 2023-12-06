package ops

import (
	"context"
	"slices"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func EditMotion(
	ctx context.Context,
	addr gov.OwnerAddress,
	id schema.MotionID,
	title string,
	body string,
	trackerURL string,
	labels []string,

) git.ChangeNoResult {

	cloned := gov.CloneOwner(ctx, addr)
	chg := EditMotionMeta_StageOnly(ctx, cloned, id, title, body, trackerURL, labels)
	return proto.CommitIfChanged(ctx, cloned.PublicClone(), chg)
}

func EditMotionMeta_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id schema.MotionID,
	title string,
	body string,
	trackerURL string,
	labels []string,

) git.ChangeNoResult {

	labels = slices.Clone(labels)
	slices.Sort(labels)

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, cloned.PublicClone().Tree(), id)
	must.Assertf(ctx, !motion.Closed, "cannot edit a closed motion %v", id)

	motion.TrackerURL = trackerURL
	motion.Title = title
	motion.Body = body
	motion.Labels = labels
	return schema.MotionKV.Set(ctx, schema.MotionNS, cloned.PublicClone().Tree(), id, motion)
}
