package errors

/*

	ctx = errors.Infof(ctx, "doing blah for %v", v)
	...
	return nil, errors.ContextErr(ctx, err)

*/

import (
	"context"
)

type Frame interface {
	String() string
}

type contextErr struct {
	Stack *Stack
	Err   error
}

func (x contextErr) Error() string {
	return x.Err.Error()
	// var w bytes.Buffer
	// XXX
	// return w.String()
}

func Error(ctx context.Context, err error) error {
	return contextErr{Stack: stackOf(ctx), Err: err}
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
