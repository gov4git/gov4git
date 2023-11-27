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
