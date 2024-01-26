package gov

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/mod"
)

type PostCloner interface {
	PostClone(
		ctx context.Context,
		cloned OwnerCloned,
	)
}

var postCloneRegistry = mod.NewModuleRegistry[string, PostCloner]()

func InstallPostClone(ctx context.Context, name string, pc PostCloner) {
	postCloneRegistry.Set(ctx, name, pc)
}

func invokePostCloners(ctx context.Context, cloned OwnerCloned) {
	_, pcs := postCloneRegistry.List()
	for _, pc := range pcs {
		pc.PostClone(ctx, cloned)
	}
}
