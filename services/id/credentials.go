package id

import (
	"context"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/idproto"
)

func (x IdentityPrivateService) GetPrivateCredentials(ctx context.Context) (*idproto.PrivateCredentials, error) {
	private, err := git.MakeLocal(ctx)
	if err != nil {
		return nil, err
	}
	if err := private.CloneOrigin(ctx, git.Origin(x.PrivateAddress)); err != nil {
		return nil, err
	}
	out, err := GetPrivateCredentialsLocal(ctx, private)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func GetPrivateCredentialsLocal(ctx context.Context, private git.Local) (*idproto.PrivateCredentials, error) {
	var credentials idproto.PrivateCredentials
	if _, err := private.Dir().ReadFormFile(ctx, idproto.PrivateCredentialsPath, &credentials); err != nil {
		return nil, err
	}
	return &credentials, nil
}

func (x IdentityPublicService) GetPublicCredentials(ctx context.Context) (*idproto.PublicCredentials, error) {
	public, err := git.MakeLocal(ctx)
	if err != nil {
		return nil, err
	}
	if err := public.CloneOrigin(ctx, git.Origin(x)); err != nil {
		return nil, err
	}
	return GetPublicCredentialsLocal(ctx, public)
}

func GetPublicCredentialsLocal(ctx context.Context, public git.Local) (*idproto.PublicCredentials, error) {
	var credentials idproto.PublicCredentials
	if _, err := public.Dir().ReadFormFile(ctx, idproto.PublicCredentialsPath, &credentials); err != nil {
		return nil, err
	}
	return &credentials, nil
}
