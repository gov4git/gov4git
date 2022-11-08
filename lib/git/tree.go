package git

import (
	"context"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/must"
)

func TreeMkdirAll(ctx context.Context, t *Tree, path string) {
	err := t.Filesystem.MkdirAll(path, 0755)
	must.NoError(ctx, err)
}

func ToFile[V form.Form](ctx context.Context, t *Tree, filepath string, value V) {
	form.ToFile(ctx, t.Filesystem, filepath, value)
}

func ToFileStage[V form.Form](ctx context.Context, t *Tree, filepath string, value V) {
	ToFile(ctx, t, filepath, value)
	Add(ctx, t, filepath)
}

func FromFile[V form.Form](ctx context.Context, t *Tree, filepath string) V {
	return form.FromFile[V](ctx, t.Filesystem, filepath)
}

func TryFromFile[V form.Form](ctx context.Context, t *Tree, filepath string) (v V, err error) {
	err = must.Try(
		func() {
			v = FromFile[V](ctx, t, filepath)
		},
	)
	return
}
