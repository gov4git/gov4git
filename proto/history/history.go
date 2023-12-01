package history

import (
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/journal"
	"golang.org/x/net/context"
)

type History journal.Journal[*Event]

func (h History) Journal() journal.Journal[*Event] {
	return journal.Journal[*Event](h)
}

func Log_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	event *Event,
) {
	if IsMuted(ctx) {
		return
	}
	history.Journal().Log_StageOnly(ctx, cloned.Tree(), event)
}

type muteCtxKey struct{}

func Mute(ctx context.Context) context.Context {
	return context.WithValue(ctx, muteCtxKey{}, true)
}

func Unmute(ctx context.Context) context.Context {
	return context.WithValue(ctx, muteCtxKey{}, false)
}

func IsMuted(ctx context.Context) bool {
	v, ok := ctx.Value(muteCtxKey{}).(bool)
	return ok && v
}
