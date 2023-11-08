package ops

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func IsMotion(ctx context.Context, addr gov.GovAddress, id schema.MotionID) bool {
	return IsMotion_Local(ctx, gov.Clone(ctx, addr).Tree(), id)
}

func IsMotion_Local(ctx context.Context, t *git.Tree, id schema.MotionID) bool {
	err := must.Try(func() { schema.MotionKV.Get(ctx, schema.MotionNS, t, id) })
	return err == nil
}

func OpenMotion(
	ctx context.Context,
	addr gov.GovAddress,
	id schema.MotionID,
	policy schema.PolicyName,
	title string,
	desc string,
	typ schema.MotionType,
	trackerURL string,
	labels []string,
) git.ChangeNoResult {

	cloned := gov.Clone(ctx, addr)
	chg := OpenMotion_StageOnly(ctx, cloned.Tree(), id, policy, title, desc, typ, trackerURL, labels)
	return proto.CommitIfChanged(ctx, cloned, chg)
}

func OpenMotion_StageOnly(
	ctx context.Context,
	t *git.Tree,
	id schema.MotionID,
	policyName schema.PolicyName,
	title string,
	desc string,
	typ schema.MotionType,
	trackerURL string,
	labels []string,
) git.ChangeNoResult {

	labels = slices.Clone(labels)
	slices.Sort(labels)

	must.Assert(ctx, !IsMotion_Local(ctx, t, id), schema.ErrMotionAlreadyExists)
	state := schema.Motion{
		OpenedAt:   time.Now(),
		ID:         id,
		Type:       typ,
		Policy:     policyName,
		TrackerURL: trackerURL,
		Title:      title,
		Body:       desc,
		Labels:     labels,
		Closed:     false,
	}
	schema.MotionKV.Set(ctx, schema.MotionNS, t, id, state)

	// apply policy
	pcy := policy.Get(ctx, policyName.String())
	pcy.Open(ctx, schema.MotionKV.KeyNS(schema.MotionNS, id).Append("policy"), state)

	return git.NewChangeNoResult(fmt.Sprintf("Open motion %v", id), "docket_open_motion")
}

func CloseMotion(ctx context.Context, addr gov.GovAddress, id schema.MotionID) git.ChangeNoResult {

	cloned := gov.Clone(ctx, addr)
	return CloseMotion_StageOnly(ctx, cloned.Tree(), id)
}

func CloseMotion_StageOnly(ctx context.Context, t *git.Tree, id schema.MotionID) git.ChangeNoResult {

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, t, id)
	must.Assert(ctx, !motion.Closed, schema.ErrMotionAlreadyClosed)
	motion.Closed = true
	motion.ClosedAt = time.Now()
	schema.MotionKV.Set(ctx, schema.MotionNS, t, id, motion)

	// apply policy
	pcy := policy.Get(ctx, motion.Policy.String())
	pcy.Close(ctx, schema.MotionKV.KeyNS(schema.MotionNS, id).Append("policy"), motion)

	return git.NewChangeNoResult(fmt.Sprintf("Close motion %v", id), "docket_close_motion")
}

func ReopenMotion_StageOnly(ctx context.Context, t *git.Tree, id schema.MotionID) git.ChangeNoResult {

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, t, id)
	must.Assert(ctx, motion.Closed, schema.ErrMotionNotClosed)
	motion.Closed = false
	return schema.MotionKV.Set(ctx, schema.MotionNS, t, id, motion)
}

func FreezeMotion_StageOnly(ctx context.Context, t *git.Tree, id schema.MotionID) git.ChangeNoResult {

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, t, id)
	must.Assert(ctx, !motion.Frozen, schema.ErrMotionAlreadyFrozen)
	motion.Frozen = true
	return schema.MotionKV.Set(ctx, schema.MotionNS, t, id, motion)
}

func UnfreezeMotion_StageOnly(ctx context.Context, t *git.Tree, id schema.MotionID) git.ChangeNoResult {

	motion := schema.MotionKV.Get(ctx, schema.MotionNS, t, id)
	must.Assert(ctx, motion.Frozen, schema.ErrMotionNotFrozen)
	motion.Frozen = false
	return schema.MotionKV.Set(ctx, schema.MotionNS, t, id, motion)
}

func UpdateMotionMeta_StageOnly(
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

func ListMotions(ctx context.Context, addr gov.GovAddress) schema.Motions {
	return ListMotions_Local(ctx, gov.Clone(ctx, addr).Tree())
}

func ListMotions_Local(ctx context.Context, t *git.Tree) schema.Motions {
	_, motions := schema.MotionKV.ListKeyValues(ctx, schema.MotionNS, t)
	return motions
}
