package id

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/idproto"
)

func (x IdentityPrivateService) Init(ctx context.Context) (*idproto.PrivateCredentials, error) {

	// generate private credentials

	localPrivate, err := git.MakeLocal(ctx)
	if err != nil {
		return nil, err
	}
	// clone or init repo
	if err := localPrivate.CloneOrInitBranch(ctx, string(x.PrivateAddress.Repo), string(x.PrivateAddress.Branch)); err != nil {
		// if err := localPrivate.CloneOrInitBranch(ctx, x.IdentityConfig.PrivateURL, idproto.IdentityBranch); err != nil {
		return nil, err
	}
	// check if key files already exist
	if _, err := localPrivate.Dir().Stat(idproto.PrivateCredentialsPath); err == nil {
		return nil, fmt.Errorf("private credentials file already exists")
	}
	// generate credentials
	privateCredentials, err := idproto.GenerateCredentials(x.PublicAddress, x.PrivateAddress)
	if err != nil {
		return nil, err
	}
	// write changes
	stagePrivate := files.FormFiles{
		files.FormFile{Path: idproto.PrivateCredentialsPath, Form: privateCredentials},
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

	localPublic, err := git.MakeLocal(ctx)
	if err != nil {
		return nil, err
	}
	// clone or init repo
	if err := localPublic.CloneOrInitBranch(ctx, string(x.PublicAddress.Repo), string(x.PublicAddress.Branch)); err != nil {
		return nil, err
	}
	// write changes
	stagePublic := files.FormFiles{
		files.FormFile{Path: idproto.PublicCredentialsPath, Form: privateCredentials.PublicCredentials},
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

	return privateCredentials, nil
}
