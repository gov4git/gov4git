package notice

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func SaveNoticeQueue(
	ctx context.Context,
	addr gov.Address,
	filepath ns.NS,
	queue *NoticeQueue,
) git.Change[*NoticeQueue, form.None] {

	cloned := gov.Clone(ctx, addr)
	chg := SaveNoticeQueue_StageOnly(ctx, cloned, filepath, queue)
	return proto.CommitIfChanged(ctx, cloned, chg)
}

func SaveNoticeQueue_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	filepath ns.NS,
	queue *NoticeQueue,
) git.Change[*NoticeQueue, form.None] {

	git.ToFileStage[*NoticeQueue](ctx, cloned.Tree(), filepath, queue)
	return git.NewChange[*NoticeQueue, form.None](
		"Save notice queue",
		"notice_save_queue",
		queue,
		form.None{},
		nil,
	)
}

func LoadNoticeQueue(
	ctx context.Context,
	addr gov.Address,
	filepath ns.NS,
) *NoticeQueue {

	return LoadNoticeQueue_Local(ctx, gov.Clone(ctx, addr), filepath)
}

func LoadNoticeQueue_Local(
	ctx context.Context,
	cloned gov.Cloned,
	filepath ns.NS,
) *NoticeQueue {

	queue, err := git.TryFromFile[*NoticeQueue](ctx, cloned.Tree(), filepath)
	if git.IsNotExist(err) {
		return NewNoticeQueue()
	}
	must.NoError(ctx, err)
	return queue
}
