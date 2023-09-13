package collab

import (
	"context"
	"time"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func IsMotion(ctx context.Context, addr gov.GovAddress, id MotionID) bool {
	return IsMotion_Local(ctx, gov.Clone(ctx, addr).Tree(), id)
}

func IsMotion_Local(ctx context.Context, t *git.Tree, id MotionID) bool {
	err := must.Try(func() { motionKV.Get(ctx, motionNS, t, id) })
	return err == nil
}

func OpenMotion(
	ctx context.Context,
	addr gov.GovAddress,
	id MotionID,
	title string,
	desc string,
	typ MotionType,
	trackerURL string,
	labels []string,
) git.ChangeNoResult {

	cloned := gov.Clone(ctx, addr)
	chg := OpenMotion_StageOnly(ctx, cloned.Tree(), id, title, desc, typ, trackerURL, labels)
	return proto.CommitIfChanged(ctx, cloned, chg)
}

func OpenMotion_StageOnly(
	ctx context.Context,
	t *git.Tree,
	id MotionID,
	title string,
	desc string,
	typ MotionType,
	trackerURL string,
	labels []string,
) git.ChangeNoResult {

	must.Assert(ctx, !IsMotion_Local(ctx, t, id), ErrMotionAlreadyExists)
	state := Motion{
		OpenedAt:   time.Now(),
		ID:         id,
		Title:      title,
		Desc:       desc,
		Type:       typ,
		TrackerURL: trackerURL,
		Labels:     labels,
		Closed:     false,
	}
	return motionKV.Set(ctx, motionNS, t, id, state)
}

func CloseMotion(ctx context.Context, addr gov.GovAddress, id MotionID) git.ChangeNoResult {

	cloned := gov.Clone(ctx, addr)
	chg := CloseMotion_StageOnly(ctx, cloned.Tree(), id)
	return proto.CommitIfChanged(ctx, cloned, chg)
}

func CloseMotion_StageOnly(ctx context.Context, t *git.Tree, id MotionID) git.ChangeNoResult {

	motion := motionKV.Get(ctx, motionNS, t, id)
	must.Assert(ctx, !motion.Closed, ErrMotionAlreadyClosed)
	motion.Closed = true
	motion.ClosedAt = time.Now()
	return motionKV.Set(ctx, motionNS, t, id, motion)
}

func ReopenMotion_StageOnly(ctx context.Context, t *git.Tree, id MotionID) git.ChangeNoResult {

	motion := motionKV.Get(ctx, motionNS, t, id)
	must.Assert(ctx, motion.Closed, ErrMotionNotClosed)
	motion.Closed = false
	return motionKV.Set(ctx, motionNS, t, id, motion)
}

func FreezeMotion_StageOnly(ctx context.Context, t *git.Tree, id MotionID) git.ChangeNoResult {

	motion := motionKV.Get(ctx, motionNS, t, id)
	must.Assert(ctx, !motion.Frozen, ErrMotionAlreadyFrozen)
	motion.Frozen = true
	return motionKV.Set(ctx, motionNS, t, id, motion)
}

func UnfreezeMotion_StageOnly(ctx context.Context, t *git.Tree, id MotionID) git.ChangeNoResult {

	motion := motionKV.Get(ctx, motionNS, t, id)
	must.Assert(ctx, motion.Frozen, ErrMotionNotFrozen)
	motion.Frozen = false
	return motionKV.Set(ctx, motionNS, t, id, motion)
}

func UpdateMotionMeta_StageOnly(
	ctx context.Context,
	t *git.Tree,
	id MotionID,
	trackerURL string,
	title string,
	desc string,
	labels []string,
) git.ChangeNoResult {

	motion := motionKV.Get(ctx, motionNS, t, id)
	motion.TrackerURL = trackerURL
	motion.Title = title
	motion.Desc = desc
	motion.Labels = labels
	return motionKV.Set(ctx, motionNS, t, id, motion)
}

func ListMotions(ctx context.Context, addr gov.GovAddress) Motions {
	return ListMotions_Local(ctx, gov.Clone(ctx, addr).Tree())
}

func ListMotions_Local(ctx context.Context, t *git.Tree) Motions {
	_, motions := motionKV.ListKeyValues(ctx, motionNS, t)
	return motions
}
