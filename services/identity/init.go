package identity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type IdentityService struct {
	IdentityConfig proto.IdentityConfig
}

type InitIn struct{}

type InitOut struct {
	PrivateCredentials proto.PrivateCredentials `json:"private_credentials"`
}

func (x InitOut) Human(context.Context) string {
	data, _ := json.MarshalIndent(x.PrivateCredentials, "", "   ")
	return string(data)
}

func (x IdentityService) Init(ctx context.Context, in *InitIn) (*InitOut, error) {

	// generate private credentials

	localPrivate, err := git.MakeLocalCtx(ctx, "private")
	if err != nil {
		return nil, err
	}
	// clone or init repo
	if err := localPrivate.CloneOrInitBranch(ctx, x.IdentityConfig.PrivateURL, proto.IdentityBranch); err != nil {
		return nil, err
	}
	// check if key files already exist
	if _, err := localPrivate.Dir().Stat(proto.PrivateCredentialsPath); err == nil {
		return nil, fmt.Errorf("private credentials file already exists")
	}
	// generate credentials
	privateCredentials, err := proto.GenerateCredentials(x.IdentityConfig.PublicURL, x.IdentityConfig.PrivateURL)
	if err != nil {
		return nil, err
	}
	// write changes
	stagePrivate := files.FormFiles{
		files.FormFile{Path: proto.PrivateCredentialsPath, Form: privateCredentials},
	}
	if err = localPrivate.Dir().WriteFormFiles(ctx, stagePrivate); err != nil {
		return nil, err
	}
	// stage changes
	if err = localPrivate.Add(ctx, stagePrivate.Paths()); err != nil {
		return nil, err
	}
	// commit changes
	if err = localPrivate.Commit(ctx, "initializing private credentials"); err != nil {
		return nil, err
	}
	// push repo
	if err = localPrivate.PushUpstream(ctx); err != nil {
		return nil, err
	}

	// generate public credentials

	localPublic, err := git.MakeLocalCtx(ctx, "public")
	if err != nil {
		return nil, err
	}
	// clone or init repo
	if err := localPublic.CloneOrInitBranch(ctx, x.IdentityConfig.PublicURL, proto.IdentityBranch); err != nil {
		return nil, err
	}
	// write changes
	stagePublic := files.FormFiles{
		files.FormFile{Path: proto.PublicCredentialsPath, Form: privateCredentials.PublicCredentials},
	}
	if err = localPublic.Dir().WriteFormFiles(ctx, stagePublic); err != nil {
		return nil, err
	}
	// stage changes
	if err = localPublic.Add(ctx, stagePublic.Paths()); err != nil {
		return nil, err
	}
	// commit changes
	if err = localPublic.Commit(ctx, "initializing public credentials"); err != nil {
		return nil, err
	}
	// push repo
	if err = localPublic.PushUpstream(ctx); err != nil {
		return nil, err
	}

	return &InitOut{PrivateCredentials: *privateCredentials}, nil
}
