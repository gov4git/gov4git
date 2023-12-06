package ops

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/notice"
)

func UpdateMotions(
	ctx context.Context,
	addr gov.OwnerAddress,
	args ...any,

) ([]policy.Report, []notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	report, notices := UpdateMotions_StageOnly(ctx, cloned, args...)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_update", "Update motions")
	return report, notices
}

func UpdateMotions_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	args ...any,

) ([]policy.Report, []notice.Notices) {

	t := cloned.Public.Tree()
	motions := ListMotions_Local(ctx, t)
	reportList := []policy.Report{}
	noticesList := []notice.Notices{}
	for i, motion := range motions {
		// only update open motions
		if motion.Closed {
			continue
		}
		p := policy.GetMotionPolicy(ctx, motion)
		report, notices := p.Update(
			ctx,
			cloned,
			motion,
			policy.MotionPolicyNS(motions[i].ID),
			args...,
		)
		reportList = append(reportList, report)
		noticesList = append(noticesList, notices)
		AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), motions[i].ID, notices)
	}

	motions.Sort()

	return reportList, noticesList
}
