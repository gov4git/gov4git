package motionapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
)

func AppendMotionNotices_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	id motionproto.MotionID,
	notices notice.Notices,
) {

	noticesNS := motionproto.MotionNoticesNS(id)
	queue := notice.LoadNoticeQueue_Local(ctx, cloned, noticesNS)
	queue.Append(notices...)
	notice.SaveNoticeQueue_StageOnly(ctx, cloned, noticesNS, queue)
}

func LoadMotionNotices(
	ctx context.Context,
	addr gov.Address,
	id motionproto.MotionID,
) *notice.NoticeQueue {

	cloned := gov.Clone(ctx, addr)
	return LoadMotionNotices_Local(ctx, cloned, id)
}

func LoadMotionNotices_Local(
	ctx context.Context,
	cloned gov.Cloned,
	id motionproto.MotionID,
) *notice.NoticeQueue {

	noticesNS := motionproto.MotionNoticesNS(id)
	return notice.LoadNoticeQueue_Local(ctx, cloned, noticesNS)
}

func SaveMotionNotices_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	id motionproto.MotionID,
	queue *notice.NoticeQueue,
) {

	noticesNS := motionproto.MotionNoticesNS(id)
	notice.SaveNoticeQueue_StageOnly(ctx, cloned, noticesNS, queue)
}
