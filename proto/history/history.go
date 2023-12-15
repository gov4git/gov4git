package history

import (
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/journal"
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

func List(
	ctx context.Context,
	addr gov.Address,
) journal.Entries[*Event] {

	cloned := gov.Clone(ctx, addr)
	return List_Local(ctx, cloned)
}

func List_Local(
	ctx context.Context,
	cloned gov.Cloned,
) journal.Entries[*Event] {

	return history.Journal().List_Local(ctx, cloned.Tree())
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
