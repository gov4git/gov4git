package entity

import (
	"context"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/entityproto"
)

func (x EntityService[V]) Remove(ctx context.Context, name string) error {
	local, err := git.MakeLocal(ctx)
	if err != nil {
		return err
	}
	if err := local.CloneOrigin(ctx, x.Address); err != nil {
		return err
	}
	if err := x.RemoveLocal(ctx, local, name); err != nil {
		return err
	}
	if err := local.PushUpstream(ctx); err != nil {
		return err
	}
	return nil
}

func (x EntityService[V]) RemoveLocal(ctx context.Context, local git.Local, name string) error {
	valueFilepath := entityproto.ValueFilepath(x.Namespace, name)
	if err := local.Dir().Remove(valueFilepath); err != nil {
		return err
	}
	if err := local.Remove(ctx, []string{valueFilepath}); err != nil {
		return err
	}
	if err := local.Commitf(ctx, "Remove %v entity %v", x.Namespace, name); err != nil {
		return err
	}
	return nil
}
