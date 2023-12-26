package motionapi

import (
	"context"
	"slices"
	"time"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicy"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/lib4git/must"
)

func OpenMotion(
	ctx context.Context,
	addr gov.OwnerAddress,
	id motionproto.MotionID,
	typ motionproto.MotionType,
	policy motionproto.PolicyName,
	author member.User,
	title string,
	desc string,
	trackerURL string,
	labels []string,

) (motionpolicy.Report, notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	report, notices := OpenMotion_StageOnly(ctx, cloned, id, typ, policy, author, title, desc, trackerURL, labels)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_open", "Open motion %v", id)
	return report, notices
}

func OpenMotion_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id motionproto.MotionID,
	typ motionproto.MotionType,
	policyName motionproto.PolicyName,
	author member.User,
	title string,
	desc string,
	trackerURL string,
	labels []string,
	args ...any,

) (motionpolicy.Report, notice.Notices) {

	t := cloned.Public.Tree()
	labels = slices.Clone(labels)
	slices.Sort(labels)

	// verify author is a user, or empty string
	must.Assertf(ctx, author == "" || member.IsUser_Local(ctx, cloned.PublicClone(), author), "motion author %v is not in the community", author)

	must.Assert(ctx, !IsMotion_Local(ctx, t, id), motionproto.ErrMotionAlreadyExists)
	motion := motionproto.Motion{
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
	motionproto.MotionKV.Set(ctx, motionproto.MotionNS, t, id, motion)

	// apply policy
	pcy := motionpolicy.Get(ctx, policyName)
	report, notices := pcy.Open(
		ctx,
		cloned,
		motion,
		motionpolicy.MotionPolicyNS(id),
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