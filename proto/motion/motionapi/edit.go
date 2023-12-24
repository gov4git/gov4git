package motionapi

import (
	"context"
	"slices"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func EditMotion(
	ctx context.Context,
	addr gov.OwnerAddress,
	id motionproto.MotionID,
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
	id motionproto.MotionID,
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

	motion := motionproto.MotionKV.Get(ctx, motionproto.MotionNS, cloned.PublicClone().Tree(), id)
	must.Assertf(ctx, !motion.Closed, "cannot edit a closed motion %v", id)

	motion.Author = author
	motion.TrackerURL = trackerURL
	motion.Title = title
	motion.Body = body
	motion.Labels = labels
	return motionproto.MotionKV.Set(ctx, motionproto.MotionNS, cloned.PublicClone().Tree(), id, motion)
}
