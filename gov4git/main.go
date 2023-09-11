package main

import (
	"fmt"
	"os"

	gov4git_root "github.com/gov4git/gov4git"
	"github.com/gov4git/gov4git/gov4git/cmd"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/must"
)

func main() {
	base.Infof("gov4git version %v", gov4git_root.VersionDotTxt)
	err, stk := must.TryWithStack(
		func() { cmd.Execute() },
	)
	if err != nil {
		if base.IsVerbose() {
			fmt.Fprintln(os.Stderr, string(stk))
		}
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
