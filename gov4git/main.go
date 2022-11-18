package main

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/gov4git/cmd"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/must"
)

func main() {
	if base.IsVerbose() {
		cmd.Execute()
	} else {
		err := must.Try(
			func() { cmd.Execute() },
		)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}
}
