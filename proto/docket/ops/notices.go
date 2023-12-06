package ops

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/notice"
)

func AppendMotionNotices_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	id schema.MotionID,
	notices notice.Notices,
) {

	noticesNS := schema.MotionNoticesNS(id)
	queue := notice.LoadNoticeQueue_Local(ctx, cloned, noticesNS)
	queue.Append(notices...)
	notice.SaveNoticeQueue_StageOnly(ctx, cloned, noticesNS, queue)
}

func LoadMotionNotices(
	ctx context.Context,
	addr gov.Address,
	id schema.MotionID,
) *notice.NoticeQueue {

	cloned := gov.Clone(ctx, addr)
	return LoadMotionNotices_Local(ctx, cloned, id)
}

func LoadMotionNotices_Local(
	ctx context.Context,
	cloned gov.Cloned,
	id schema.MotionID,
) *notice.NoticeQueue {

	noticesNS := schema.MotionNoticesNS(id)
	return notice.LoadNoticeQueue_Local(ctx, cloned, noticesNS)
}

func SaveMotionNotices_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	id schema.MotionID,
	queue *notice.NoticeQueue,
) {

	noticesNS := schema.MotionNoticesNS(id)
	notice.SaveNoticeQueue_StageOnly(ctx, cloned, noticesNS, queue)
}
