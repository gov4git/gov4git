package journal

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
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

type Entries[X form.Form] []Entry[X]

func (x Entries[X]) Len() int {
	return len(x)
}

func (x Entries[X]) Less(i, j int) bool {
	return x[i].Stamp.Before(x[j].Stamp)
}

func (x Entries[X]) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x Entries[X]) Sort() {
	sort.Sort(x)
}

func (j Journal[X]) List_Local(
	ctx context.Context,
	t *git.Tree,
) Entries[X] {

	es := Entries[X]{}

	root := j.Root.GitPath()
	yinfos, err := t.Filesystem.ReadDir(root)
	must.Assert(ctx, err == nil || git.IsNotExist(err), err)
	// years
	for _, yinfo := range yinfos {
		if !yinfo.IsDir() {
			continue
		}
		if _, err := strconv.Atoi(yinfo.Name()); err != nil {
			continue
		}
		yearNS := j.Root.Append(yinfo.Name())
		minfos, err := t.Filesystem.ReadDir(yearNS.GitPath())
		must.Assert(ctx, err == nil || git.IsNotExist(err), err)
		// months
		for _, minfo := range minfos {
			if !minfo.IsDir() {
				continue
			}
			if _, err := strconv.Atoi(minfo.Name()); err != nil {
				continue
			}
			monthNS := yearNS.Append(minfo.Name())
			dinfos, err := t.Filesystem.ReadDir(monthNS.GitPath())
			must.Assert(ctx, err == nil || git.IsNotExist(err), err)
			// days
			for _, dinfo := range dinfos {
				if !dinfo.IsDir() {
					continue
				}
				if _, err := strconv.Atoi(dinfo.Name()); err != nil {
					continue
				}
				dayNS := monthNS.Append(dinfo.Name())
				loginfos, err := t.Filesystem.ReadDir(dayNS.GitPath())
				must.Assert(ctx, err == nil || git.IsNotExist(err), err)
				// entries
				for _, loginfo := range loginfos {
					if loginfo.IsDir() || filepath.Ext(loginfo.Name()) != ".json" {
						continue
					}
					e := form.FromFile[Entry[X]](ctx, t.Filesystem, dayNS.Append(loginfo.Name()))
					es = append(es, e)
				}
			}
		}
	}
	es.Sort()
	return es
}

func (j Journal[X]) Log_StageOnly(
	ctx context.Context,
	t *git.Tree,
	x X,
) {

	now := time.Now()
	id := id.GenerateRandomID()
	dirNS := j.Root.Append(
		fmt.Sprintf("%04d", now.Year()),
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()),
	)
	git.TreeMkdirAll(ctx, t, dirNS)
	entry := Entry[X]{
		ID:      id,
		Stamp:   now,
		Payload: x,
	}
	filebase := fmt.Sprintf(
		"%04d-%02d-%02d_%02d:%02d:%02d_%v.json",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		id,
	)
	git.ToFileStage[Entry[X]](ctx, t, dirNS.Append(filebase), entry)
}
