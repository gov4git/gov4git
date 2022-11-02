package entity

import (
	"context"
	"strconv"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/entityproto"
)

func (x EntityService[V]) GetProp(ctx context.Context, name string, key string) ([]byte, error) {
	local, err := git.CloneOrigin(ctx, x.Address)
	if err != nil {
		return nil, err
	}
	value, err := x.GetPropLocal(ctx, local, name, key)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (x EntityService[V]) GetPropLocal(ctx context.Context, local git.Local, name string, key string) ([]byte, error) {
	propFilepath := entityproto.PropValueFilepath(x.Namespace, name, key)
	data, err := local.Dir().ReadByteFile(propFilepath)
	if err != nil {
		return nil, err
	}
	return data.Bytes, nil
}

func (x EntityService[V]) GetPropFloat64Local(ctx context.Context, local git.Local, name string, key string) (float64, error) {
	v, err := x.GetPropLocal(ctx, local, name, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(string(v), 64)
}
