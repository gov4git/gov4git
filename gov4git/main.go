package main

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/v2/gov4git/cmd"
	_ "github.com/gov4git/gov4git/v2/runtime"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/must"
)

func main() {
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
