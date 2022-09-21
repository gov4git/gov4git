package base

/*

	ctx = base.Infof(ctx, "doing blah for %v", v)
	...
	return nil, base.ContextErr(ctx, err)

*/

import (
	"context"
)

type ErrorInContext struct {
	context.Context
	error
}

func (x ErrorInContext) Error() string {
	return x.Err.Error()
	// var w bytes.Buffer
	// XXX
	// return w.String()
}

func ErrInCtx(ctx context.Context, err error) error {
	return ErrorInContext{Context: ctx, error: err}
}

// stack

type Frame interface {
	String() string
}

type ctxStackFrame struct{}

func stackOf(ctx context.Context) *Stack {
	stk, _ := ctx.Value(ctxStackFrame{}).(*Stack)
	return stk
}

func pushFrame(ctx context.Context, frame Frame) context.Context {
	return context.WithValue(ctx, ctxStackFrame{}, stackOf(ctx).Push(frame))
}

type Stack struct {
	Frame  Frame
	Parent *Stack
}

func (x *Stack) Push(frame Frame) *Stack {
	return &Stack{Frame: frame, Parent: x}
}
