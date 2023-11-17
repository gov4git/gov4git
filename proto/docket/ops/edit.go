package ops

import (
	"context"
	"slices"

	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/lib4git/git"
)

func EditMotionMeta_StageOnly(
	ctx context.Context,
	t *git.Tree,
	id schema.MotionID,
	trackerURL string,
	title string,
	desc string,
	labels []string,
) git.ChangeNoResult {

	labels = slices.Clone(labels)
	slices.Sort(labels)

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, t, id)
	motion.TrackerURL = trackerURL
	motion.Title = title
	motion.Body = desc
	motion.Labels = labels
	return schema.MotionKV.Set(ctx, schema.MotionNS, t, id, motion)
}
