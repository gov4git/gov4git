package must

import "context"

type Error struct {
	Ctx context.Context
	Err error
}

func Err(ctx context.Context, err error) Error {
	return Error{Ctx: ctx, Err: err}
}

func Panic(ctx context.Context, err error) {
	panic(Err(ctx, err))
}

// func Must[R any](ctx context.Context, )
