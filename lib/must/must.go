package must

import "context"

type Error struct {
	Ctx context.Context
	error
}

func mkErr(ctx context.Context, err error) Error {
	return Error{Ctx: ctx, error: err}
}

func Panic(ctx context.Context, err error) {
	panic(mkErr(ctx, err))
}

func NoError(ctx context.Context, err error) {
	if err != nil {
		Panic(ctx, err)
	}
}

func Try0(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(Error).error
		}
	}()
	f()
	return nil
}

func Try1[R1 any](f func() R1) (_ R1, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(Error).error
		}
	}()
	return f(), nil
}
