package api

import (
	"os"
	"testing"

	"github.com/gov4git/gov4git/gov4git/cmd"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	resRun := testscript.RunMain(m, map[string]func() int{
		"gov4git": cmd.Execute,
	})
	os.Exit(resRun)
}

func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/script",
	})
}
