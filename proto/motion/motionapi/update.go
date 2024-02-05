package motionapi

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
)

func UpdateMotions(
	ctx context.Context,
	addr gov.OwnerAddress,
	args ...any,

) ([]motionproto.Report, []notice.Notices) {

	cloned := gov.CloneOwner(ctx, addr)
	report, notices := UpdateMotions_StageOnly(ctx, cloned, args...)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_update", "Update motions")
	return report, notices
}

func UpdateMotions_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	args ...any,

) ([]motionproto.Report, []notice.Notices) {

	t := cloned.Public.Tree()
	motions := ListMotions_Local(ctx, t)
	reportList := []motionproto.Report{}
	noticesList := []notice.Notices{}
	for i, motion := range motions {
		fmt.Printf("UPDATING motion %v\n", motion.ID)
		if motion.Archived || motion.Closed {
			continue
		}
		fmt.Printf("UPDATING motion %v continuing\n", motion.ID)
		p := motionproto.GetMotionPolicy(ctx, motion)
		report, notices := p.Update(
			ctx,
			cloned,
			motion,
			args...,
		)
		reportList = append(reportList, report)
		noticesList = append(noticesList, notices)
		AppendMotionNotices_StageOnly(ctx, cloned.PublicClone(), motions[i].ID, notices)
	}

	motions.Sort()

	return reportList, noticesList
}
