package regime

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/history/metric"
	"github.com/gov4git/gov4git/v2/proto/history/trace"
	"github.com/gov4git/gov4git/v2/proto/notice"
	"github.com/gov4git/lib4git/git"
)

func Dry(ctx context.Context) context.Context {
	ctx = metric.Mute(ctx)
	ctx = trace.Mute(ctx)
	ctx = notice.Mute(ctx)
	ctx = git.MuteStaging(ctx)
	return ctx
}
