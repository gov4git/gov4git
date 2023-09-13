package proto

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

var RootNS = ns.NS{}

func Commit(ctx context.Context, t *git.Tree, chg git.Commitable) {
	var w bytes.Buffer
	fmt.Fprintln(&w, chg.Message())
	fmt.Fprintln(&w)
	fmt.Fprintln(&w, form.SprintJSON(chg))
	git.Commit(ctx, t, w.String())
}

func CommitIfChanged[C git.Commitable](ctx context.Context, cloned git.Cloned, commitable C) C {
	status, err := cloned.Tree().Status()
	must.NoError(ctx, err)
	if !status.IsClean() {
		Commit(ctx, cloned.Tree(), commitable)
		cloned.Push(ctx)
	}
	return commitable
}
