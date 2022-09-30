package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/petar/gitty/lib/files"
	"github.com/petar/gitty/lib/git"
	"github.com/petar/gitty/proto"
)

type SoulService struct {
	SoulConfig proto.SoulConfig
}

type SoulInitIn struct{}

type SoulInitOut struct {
	PrivateCredentials proto.PrivateCredentials
}

func (x SoulInitOut) Human(context.Context) string {
	data, _ := json.MarshalIndent(x.PrivateCredentials, "", "   ")
	return string(data)
}

func (x SoulService) Init(ctx context.Context, in *SoulInitIn) (*SoulInitOut, error) {

	// generate private credentials

	localPrivate := git.LocalFromDir(files.WorkDir(ctx).Subdir("private"))
	// clone or init repo
	if err := localPrivate.CloneOrInitBranch(ctx, x.SoulConfig.PrivateURL, proto.IdentityBranch); err != nil {
		return nil, err
	}
	// check if key files already exist
	if _, err := localPrivate.Dir().Stat(proto.PrivateCredentialsPath); err == nil {
		return nil, fmt.Errorf("private credentials file already exists")
	}
	// generate credentials
	privateCredentials, err := proto.GenerateCredentials(x.SoulConfig.PublicURL, x.SoulConfig.PrivateURL)
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

	localPublic := git.LocalFromDir(files.WorkDir(ctx).Subdir("public"))
	// clone or init repo
	if err := localPublic.CloneOrInitBranch(ctx, x.SoulConfig.PublicURL, proto.IdentityBranch); err != nil {
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

	return &SoulInitOut{PrivateCredentials: *privateCredentials}, nil
}
