package notice

import "context"

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
