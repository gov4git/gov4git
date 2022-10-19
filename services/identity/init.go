package identity

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/identityproto"
)

type IdentityService struct {
	IdentityConfig identityproto.IdentityConfig
}

type InitIn struct{}

type InitOut struct {
	PrivateCredentials identityproto.PrivateCredentials `json:"private_credentials"`
}

func (x IdentityService) Init(ctx context.Context, in *InitIn) (*InitOut, error) {

	// generate private credentials

	localPrivate, err := git.MakeLocalInCtx(ctx, "private")
	if err != nil {
		return nil, err
	}
	// clone or init repo
	if err := localPrivate.CloneOrInitBranch(ctx, x.IdentityConfig.PrivateURL, proto.IdentityBranch); err != nil {
		return nil, err
	}
	// check if key files already exist
	if _, err := localPrivate.Dir().Stat(identityproto.PrivateCredentialsPath); err == nil {
		return nil, fmt.Errorf("private credentials file already exists")
	}
	// generate credentials
	privateCredentials, err := identityproto.GenerateCredentials(x.IdentityConfig.PublicURL, x.IdentityConfig.PrivateURL)
	if err != nil {
		return nil, err
	}
	// write changes
	stagePrivate := files.FormFiles{
		files.FormFile{Path: identityproto.PrivateCredentialsPath, Form: privateCredentials},
	}
	if err = localPrivate.Dir().WriteFormFiles(ctx, stagePrivate); err != nil {
		return nil, err
	}
	// stage changes
	if err = localPrivate.Add(ctx, stagePrivate.Paths()); err != nil {
		return nil, err
	}
	// commit changes
	if err = localPrivate.Commit(ctx, "Initializing private credentials."); err != nil {
		return nil, err
	}
	// push repo
	if err = localPrivate.PushUpstream(ctx); err != nil {
		return nil, err
	}

	// generate public credentials

	localPublic, err := git.MakeLocalInCtx(ctx, "public")
	if err != nil {
		return nil, err
	}
	// clone or init repo
	if err := localPublic.CloneOrInitBranch(ctx, x.IdentityConfig.PublicURL, proto.IdentityBranch); err != nil {
		return nil, err
	}
	// write changes
	stagePublic := files.FormFiles{
		files.FormFile{Path: identityproto.PublicCredentialsPath, Form: privateCredentials.PublicCredentials},
	}
	if err = localPublic.Dir().WriteFormFiles(ctx, stagePublic); err != nil {
		return nil, err
	}
	// stage changes
	if err = localPublic.Add(ctx, stagePublic.Paths()); err != nil {
		return nil, err
	}
	// commit changes
	if err = localPublic.Commit(ctx, "Initializing public credentials."); err != nil {
		return nil, err
	}
	// push repo
	if err = localPublic.PushUpstream(ctx); err != nil {
		return nil, err
	}

	return &InitOut{PrivateCredentials: *privateCredentials}, nil
}
