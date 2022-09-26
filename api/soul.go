package soul

import (
	"context"
	"fmt"

	. "github.com/petar/gitty/lib/base"
	. "github.com/petar/gitty/lib/files"
	"github.com/petar/gitty/lib/git"
	"github.com/petar/gitty/proto"
)

type SoulAPI struct {
	SoulConfig proto.SoulConfig
}

func (x SoulAPI) Init(ctx context.Context) ContextError {

	// generate private credentials

	localPrivate := git.LocalFromDir(WorkDir(ctx).Subdir("private"))
	// clone or init repo
	if err := localPrivate.CloneOrInitBranch(ctx, x.SoulConfig.PrivateURL, proto.IdentityBranch); err != nil {
		return DoneErr(ctx, err)
	}
	// check if key files already exist
	if _, err := localPrivate.Dir().Stat(proto.PrivateCredentialsPath); err == nil {
		return DoneErr(ctx, fmt.Errorf("private credentials file already exists"))
	}
	// generate credentials
	privateCredentials, err := proto.GenerateCredentials(x.SoulConfig.PublicURL, x.SoulConfig.PrivateURL)
	if err != nil {
		return DoneErr(ctx, err)
	}
	// write changes
	stagePrivate := FormFiles{
		FormFile{Path: proto.PrivateCredentialsPath, Form: privateCredentials},
	}
	if err = localPrivate.Dir().WriteFormFiles(ctx, stagePrivate); err != nil {
		return DoneErr(ctx, err)
	}
	// stage changes
	if err = localPrivate.Add(ctx, stagePrivate.Paths()); err != nil {
		return DoneErr(ctx, err)
	}
	// commit changes
	if err = localPrivate.Commit(ctx, "initializing private credentials"); err != nil {
		return DoneErr(ctx, err)
	}
	// push repo
	if err = localPrivate.PushUpstream(ctx); err != nil {
		return DoneErr(ctx, err)
	}

	// generate public credentials

	localPublic := git.LocalFromDir(WorkDir(ctx).Subdir("public"))
	// clone or init repo
	if err := localPublic.CloneOrInitBranch(ctx, x.SoulConfig.PublicURL, proto.IdentityBranch); err != nil {
		return DoneErr(ctx, err)
	}
	// write changes
	stagePublic := FormFiles{
		FormFile{Path: proto.PublicCredentialsPath, Form: privateCredentials.PublicCredentials},
	}
	if err = localPublic.Dir().WriteFormFiles(ctx, stagePublic); err != nil {
		return DoneErr(ctx, err)
	}
	// stage changes
	if err = localPublic.Add(ctx, stagePublic.Paths()); err != nil {
		return DoneErr(ctx, err)
	}
	// commit changes
	if err = localPublic.Commit(ctx, "initializing public credentials"); err != nil {
		return DoneErr(ctx, err)
	}
	// push repo
	if err = localPublic.PushUpstream(ctx); err != nil {
		return DoneErr(ctx, err)
	}

	return DoneOk(ctx)
}

// func (x SoulAPI) CheckoutSendChannel(ctx context.Context, to XXX, topic string) error {
// }

// func (x SoulAPI) SyncCheckoutReceiveChannel(ctx context.Context, from XXX, topic string) error {
// }
