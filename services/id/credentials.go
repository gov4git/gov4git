package id

import (
	"context"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/idproto"
)

// private

type GetPrivateCredentialsIn struct{}

type GetPrivateCredentialsOut struct {
	PrivateCredentials idproto.PrivateCredentials `json:"private_credentials"`
}

func (x IdentityService) GetPrivateCredentials(ctx context.Context, in *GetPrivateCredentialsIn) (*GetPrivateCredentialsOut, error) {
	// clone private identity repo locally
	private, err := git.MakeLocalInCtx(ctx, "private")
	if err != nil {
		return nil, err
	}
	if err := private.CloneBranch(ctx, x.IdentityConfig.PrivateURL, idproto.IdentityBranch); err != nil {
		return nil, err
	}

	// read from the local clone
	out, err := GetPrivateCredentialsLocal(ctx, private, in)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func GetPrivateCredentialsLocal(ctx context.Context, private git.Local, in *GetPrivateCredentialsIn) (*GetPrivateCredentialsOut, error) {
	var credentials idproto.PrivateCredentials
	if _, err := private.Dir().ReadFormFile(ctx, idproto.PrivateCredentialsPath, &credentials); err != nil {
		return nil, err
	}
	return &GetPrivateCredentialsOut{PrivateCredentials: credentials}, nil
}

// public

type GetPublicCredentialsIn struct{}

type GetPublicCredentialsOut struct {
	PublicCredentials idproto.PublicCredentials `json:"public_credentials"`
}

func (x IdentityService) GetPublicCredentials(ctx context.Context, in *GetPublicCredentialsIn) (*GetPublicCredentialsOut, error) {
	// clone public identity repo locally
	public, err := git.MakeLocalInCtx(ctx, "public")
	if err != nil {
		return nil, err
	}
	if err := public.CloneBranch(ctx, x.IdentityConfig.PublicURL, idproto.IdentityBranch); err != nil {
		return nil, err
	}

	// read from the local clone
	out, err := GetPublicCredentialsLocal(ctx, public)
	if err != nil {
		return nil, err
	}

	return &GetPublicCredentialsOut{PublicCredentials: *out}, nil
}

func GetPublicCredentials(ctx context.Context, publicRepoURL string) (*idproto.PublicCredentials, error) {
	repo, err := git.MakeLocalInCtx(ctx, "")
	if err != nil {
		return nil, err
	}
	if err := repo.CloneBranch(ctx, publicRepoURL, idproto.IdentityBranch); err != nil {
		return nil, err
	}
	return GetPublicCredentialsLocal(ctx, repo)
}

func GetPublicCredentialsLocal(ctx context.Context, public git.Local) (*idproto.PublicCredentials, error) {
	var credentials idproto.PublicCredentials
	if _, err := public.Dir().ReadFormFile(ctx, idproto.PublicCredentialsPath, &credentials); err != nil {
		return nil, err
	}
	return &credentials, nil
}
