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
