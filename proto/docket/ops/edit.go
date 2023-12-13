package ops

import (
	"context"
	"slices"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func EditMotion(
	ctx context.Context,
	addr gov.OwnerAddress,
	id schema.MotionID,
	author member.User,
	title string,
	body string,
	trackerURL string,
	labels []string,

) git.ChangeNoResult {

	cloned := gov.CloneOwner(ctx, addr)
	chg := EditMotionMeta_StageOnly(ctx, cloned, id, author, title, body, trackerURL, labels)
	return proto.CommitIfChanged(ctx, cloned.PublicClone(), chg)
}

func EditMotionMeta_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id schema.MotionID,
	author member.User,
	title string,
	body string,
	trackerURL string,
	labels []string,

) git.ChangeNoResult {

	// verify author is a user, or empty string
	must.Assertf(ctx, author == "" || member.IsUser_Local(ctx, cloned.PublicClone(), author), "motion author %v is not in the community", author)

	labels = slices.Clone(labels)
	slices.Sort(labels)

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, cloned.PublicClone().Tree(), id)
	must.Assertf(ctx, !motion.Closed, "cannot edit a closed motion %v", id)

	motion.Author = author
	motion.TrackerURL = trackerURL
	motion.Title = title
	motion.Body = body
	motion.Labels = labels
	return schema.MotionKV.Set(ctx, schema.MotionNS, cloned.PublicClone().Tree(), id, motion)
}
