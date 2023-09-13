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
) {

	cloned := gov.Clone(ctx, addr)
	chg := OpenMotion_StageOnly(ctx, cloned.Tree(), id, title, desc, typ, trackerURL, labels)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
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

func CloseMotion(ctx context.Context, addr gov.GovAddress, id MotionID) {

	cloned := gov.Clone(ctx, addr)
	chg := CloseMotion_StageOnly(ctx, cloned.Tree(), id)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func CloseMotion_StageOnly(ctx context.Context, t *git.Tree, id MotionID) git.ChangeNoResult {

	motion := motionKV.Get(ctx, motionNS, t, id)
	must.Assert(ctx, !motion.Closed, ErrMotionAlreadyClosed)
	motion.Closed = true
	motion.ClosedAt = time.Now()
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

func ListMotions(ctx context.Context, addr gov.GovAddress) Motions {
	return ListMotions_Local(ctx, gov.Clone(ctx, addr).Tree())
}

func ListMotions_Local(ctx context.Context, t *git.Tree) Motions {
	_, motions := motionKV.ListKeyValues(ctx, motionNS, t)
	return motions
}
