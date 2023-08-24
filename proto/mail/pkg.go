package mail

import (
	"context"
	"encoding/json"

	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

type Pkg interface {
	stagePkg(ctx context.Context, tree *git.Tree, pathInTree ns.NS)
}

// Dir

type PkgDir map[string]Pkg

func (x PkgDir) stagePkg(ctx context.Context, tree *git.Tree, pathInTree ns.NS) {
	for k, v := range x {
		v.stagePkg(ctx, tree, pathInTree.Sub(k))
	}
}

// Blob

type PkgBlob []byte

func (x PkgBlob) MarshalJSON() ([]byte, error) {
	return json.Marshal(form.Bytes(x))
}

func (x *PkgBlob) UnmarshalJSON(d []byte) error {
	return (*form.Bytes)(x).UnmarshalJSON(d)
}

func (x PkgBlob) stagePkg(ctx context.Context, tree *git.Tree, pathInTree ns.NS) {
	git.BytesToFileStage(ctx, tree, pathInTree, x)
}

// File

func PkgFile[V form.Form](v V) pkgFile[V] {
	return pkgFile[V]{Content: v}
}

type pkgFile[V form.Form] struct {
	Content V `json:"file_content"`
}

func (x pkgFile[V]) stagePkg(ctx context.Context, tree *git.Tree, pathInTree ns.NS) {
	git.ToFileStage(ctx, tree, pathInTree.Path(), x.Content)
}
