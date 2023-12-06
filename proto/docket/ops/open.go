package ops

import (
	"context"
	"slices"
	"time"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/history"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/gov4git/proto/notice"
	"github.com/gov4git/lib4git/must"
)

func OpenMotion(
	ctx context.Context,
	addr gov.OwnerAddress,
	id schema.MotionID,
	typ schema.MotionType,
	policy schema.PolicyName,
	author member.User,
	title string,
	desc string,
	trackerURL string,
	labels []string,

) (policy.Report, notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	report, notices := OpenMotion_StageOnly(ctx, cloned, id, typ, policy, author, title, desc, trackerURL, labels)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_open", "Open motion %v", id)
	return report, notices
}

func OpenMotion_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id schema.MotionID,
	typ schema.MotionType,
	policyName schema.PolicyName,
	author member.User,
	title string,
	desc string,
	trackerURL string,
	labels []string,
	args ...any,

) (policy.Report, notice.Notices) {

	t := cloned.Public.Tree()
	labels = slices.Clone(labels)
	slices.Sort(labels)

	// verify author is a user, or empty string
	must.Assertf(ctx, author == "" || member.IsUser_Local(ctx, cloned.PublicClone(), author), "motion author %v is not in the community", author)

	must.Assert(ctx, !IsMotion_Local(ctx, t, id), schema.ErrMotionAlreadyExists)
	motion := schema.Motion{
		OpenedAt:   time.Now(),
		ID:         id,
		Type:       typ,
		Policy:     policyName,
		Author:     author,
		TrackerURL: trackerURL,
		Title:      title,
		Body:       desc,
		Labels:     labels,
		Closed:     false,
	}
	schema.MotionKV.Set(ctx, schema.MotionNS, t, id, motion)

	// apply policy
	pcy := policy.Get(ctx, policyName)
	report, notices := pcy.Open(
		ctx,
		cloned,
		motion,
		policy.MotionPolicyNS(id),
		args...,
	)
	AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), id, notices)

	// log
	history.Log_StageOnly(ctx, cloned.PublicClone(), &history.Event{
		Op: &history.Op{
			Op:     "motion_open",
			Args:   history.M{"id": id},
			Result: history.M{"motion": motion},
		},
	})

	return report, notices
}
