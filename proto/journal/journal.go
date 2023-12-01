package journal

import (
	"context"
	"fmt"
	"time"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

type Journal[X form.Form] struct {
	Root ns.NS // root in repo
}

type Entry[X form.Form] struct {
	ID      id.ID     `json:"id"`
	Stamp   time.Time `json:"stamp"`
	Payload X         `json:"payload"`
}

func (j Journal[X]) Log_StageOnly(
	ctx context.Context,
	t *git.Tree,
	x X,
) {

	now := time.Now()
	id := id.GenerateRandomID()
	dirNS := j.Root.Append(
		fmt.Sprintf("%4d", now.Year()),
		fmt.Sprintf("%2d", now.Month()),
		fmt.Sprintf("%2d", now.Day()),
	)
	git.TreeMkdirAll(ctx, t, dirNS)
	entry := Entry[X]{
		ID:      id,
		Stamp:   now,
		Payload: x,
	}
	filebase := fmt.Sprintf(
		"%2d:%2d:%2d-%v.json",
		now.Hour(),
		now.Minute(),
		now.Second(),
		id,
	)
	form.ToFile[Entry[X]](ctx, t.Filesystem, dirNS.Append(filebase), entry)
}
