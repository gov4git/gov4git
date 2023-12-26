package motionapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicy"
	"github.com/gov4git/gov4git/v2/proto/notice"
)

func UpdateMotions(
	ctx context.Context,
	addr gov.OwnerAddress,
	args ...any,

) ([]motionpolicy.Report, []notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	report, notices := UpdateMotions_StageOnly(ctx, cloned, args...)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_update", "Update motions")
	return report, notices
}

func UpdateMotions_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	args ...any,

) ([]motionpolicy.Report, []notice.Notices) {

	t := cloned.Public.Tree()
	motions := ListMotions_Local(ctx, t)
	reportList := []motionpolicy.Report{}
	noticesList := []notice.Notices{}
	for i, motion := range motions {
		// only update open motions
		if motion.Closed {
			continue
		}
		p := motionpolicy.GetMotionPolicy(ctx, motion)
		report, notices := p.Update(
			ctx,
			cloned,
			motion,
			motionpolicy.MotionPolicyNS(motions[i].ID),
			args...,
		)
		reportList = append(reportList, report)
		noticesList = append(noticesList, notices)
		AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), motions[i].ID, notices)
	}

	motions.Sort()

	return reportList, noticesList
}