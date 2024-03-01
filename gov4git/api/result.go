package api

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/must"
)

type Result struct {
	Status   Status `json:"status"`
	Returned any    `json:"returned,omitempty"`
	Msg      string `json:"msg,omitempty"`   // summary of error
	Error    error  `json:"error,omitempty"` // structure of error
}

func Invoke(f func()) Result {

	// mem profile
	defer func() {
		if memProfilePath != "" {
			f, err := os.Create(memProfilePath)
			if err != nil {
				base.Fatalf("could not create memory profile (%v)", err)
			}
			defer f.Close() // error handling omitted for example
			runtime.GC()    // get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				base.Fatalf("could not write memory profile (%v)", err)
			}
		}
	}()

	// cpu profile
	if cpuProfilePath != "" {
		f, err := os.Create(cpuProfilePath)
		if err != nil {
			base.Fatalf("could not create CPU profile (%v)", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			base.Fatalf("could not start CPU profile (%v)", err)
		}
		defer pprof.StopCPUProfile()
	}

	//
	err := must.TryThru(f)
	r := NewResult(nil, err)
	if err != nil && base.IsVerbose() {
		fmt.Fprint(os.Stderr, string(err.Stack))
	}
	fmt.Fprint(os.Stdout, form.SprintJSON(r))
	if err != nil {
		os.Exit(1)
	}
	return r
}

func Invoke1[R1 any](f func() R1) Result {

	// mem profile
	defer func() {
		if memProfilePath != "" {
			f, err := os.Create(memProfilePath)
			if err != nil {
				base.Fatalf("could not create memory profile (%v)", err)
			}
			defer f.Close() // error handling omitted for example
			runtime.GC()    // get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				base.Fatalf("could not write memory profile (%v)", err)
			}
		}
	}()

	// cpu profile
	if cpuProfilePath != "" {
		f, err := os.Create(cpuProfilePath)
		if err != nil {
			base.Fatalf("could not create CPU profile (%v)", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			base.Fatalf("could not start CPU profile (%v)", err)
		}
		defer pprof.StopCPUProfile()
	}

	//
	r1, err := must.Try1Thru[R1](f)
	r := NewResult(r1, err)
	if err != nil && base.IsVerbose() {
		fmt.Fprint(os.Stderr, string(err.Stack))
	}
	fmt.Fprint(os.Stdout, form.SprintJSON(r))
	if err != nil {
		os.Exit(1)
	}
	return r
}

func NewResult(r any, err *must.Error) Result {
	var result Result
	if err == nil {
		result.Status = StatusSuccess
	} else {
		result.Status = StatusError
		result.Msg = err.Error()
		result.Error = err.Wrapped()
	}
	result.Returned = r
	return result
}

type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

var cpuProfilePath string

func SetCPUProfilePath(filepath string) {
	cpuProfilePath = filepath
}

var memProfilePath string

func SetMemProfilePath(filepath string) {
	memProfilePath = filepath
}
