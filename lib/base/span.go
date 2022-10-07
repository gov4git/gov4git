package base

/*

	ctx = base.Enterf(ctx, "doing blah for %v", v)
	...
	return nil, base.Abort(ctx, err)

*/

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
)

type AbortError struct {
	Context context.Context
	Cause   error
}

func (x AbortError) RootCause() error {
	if x, ok := x.Cause.(AbortError); ok {
		return x.RootCause()
	}
	return x.Cause
}

func (x AbortError) Error() string {
	var w bytes.Buffer
	for i, f := range framesOf(x.Context) {
		fmt.Fprintf(&w, "<%d> ", i)
		w.WriteString(f.String())
		w.WriteString("\n")
	}
	w.WriteString(x.Cause.Error())
	return w.String()
}

func Abort(ctx context.Context, cause error) error {
	return AbortError{Context: ctx, Cause: cause}
}

// stack

type Frame interface {
	String() string
}

type ctxStackFrame struct{}

func stackOf(ctx context.Context) *stack {
	stk, _ := ctx.Value(ctxStackFrame{}).(*stack)
	return stk
}

func framesOf(ctx context.Context) []Frame {
	stk := stackOf(ctx)
	if stk == nil {
		return nil
	}
	return stk.Frames()
}

func pushFrame(ctx context.Context, frame Frame) context.Context {
	return context.WithValue(ctx, ctxStackFrame{}, stackOf(ctx).Push(frame))
}

type stack struct {
	Frame  Frame
	Parent *stack
}

func (x *stack) Push(frame Frame) *stack {
	return &stack{Frame: frame, Parent: x}
}

func (x *stack) Frames() []Frame {
	if x.Parent == nil {
		return []Frame{x.Frame}
	}
	return append(x.Parent.Frames(), x.Frame)
}

type frame struct {
	File string
	Line int
	Msg  string
}

func (f frame) String() string {
	return fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Msg)
}

func Spanf(ctx context.Context, f string, args ...any) context.Context {
	_, file, line, _ := runtime.Caller(1)
	return pushFrame(ctx, frame{File: file, Line: line, Msg: fmt.Sprintf(f, args...)})
}
