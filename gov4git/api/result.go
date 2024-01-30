package api

import (
	"fmt"
	"os"

	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/must"
)

type Result struct {
	Status   Status `json:"status"`
	Returned any    `json:"returned,omitempty"`
	Msg      string `json:"msg,omitempty"`
}

func Invoke(f func()) Result {
	err := must.TryThru(f)
	r := NewResult(nil, err.Wrapped())
	if err != nil && base.IsVerbose() {
		fmt.Fprint(os.Stderr, string(err.Stack))
	}
	fmt.Fprint(os.Stdout, form.SprintJSON(r))
	return r
}

func Invoke1[R1 any](f func() R1) Result {
	r1, err := must.Try1Thru[R1](f)
	r := NewResult(r1, err.Wrapped())
	if err != nil && base.IsVerbose() {
		fmt.Fprint(os.Stderr, string(err.Stack))
	}
	fmt.Fprint(os.Stdout, form.SprintJSON(r))
	return r
}

func NewResult(r any, err error) Result {
	var result Result
	if err == nil {
		result.Status = StatusSuccess
	} else {
		result.Status = StatusError
		result.Msg = err.Error()
	}
	result.Returned = r
	return result
}

type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)
