package id

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Init(
	ctx context.Context,
	ownerAddr OwnerAddress,
) git.Change[PrivateCredentials] {
	ownerCloned := CloneOwner(ctx, ownerAddr)
	privChg := InitLocal(ctx, ownerAddr, ownerCloned)

	ownerCloned.Public.Push(ctx)
	ownerCloned.Private.Push(ctx)
	return privChg
}

func InitLocal(
	ctx context.Context,
	ownerAddr OwnerAddress,
	ownerCloned OwnerCloned,
) git.Change[PrivateCredentials] {

	privChg := initPrivateStageOnly(ctx, ownerCloned.Private.Tree(), ownerAddr)
	pubChg := initPublicStageOnly(ctx, ownerCloned.Public.Tree(), privChg.Result.PublicCredentials)
	proto.Commit(ctx, ownerCloned.Private.Tree(), privChg.Msg)
	proto.Commit(ctx, ownerCloned.Public.Tree(), pubChg.Msg)
	return privChg
}

func initPrivateStageOnly(ctx context.Context, priv *git.Tree, ownerAddr OwnerAddress) git.Change[PrivateCredentials] {
	if _, err := priv.Filesystem.Stat(PrivateCredentialsNS.Path()); err == nil {
		must.Errorf(ctx, "private credentials file already exists")
	}
	cred, err := GenerateCredentials()
	must.NoError(ctx, err)
	git.ToFileStage(ctx, priv, PrivateCredentialsNS.Path(), cred)
	return git.Change[PrivateCredentials]{
		Result: cred,
		Msg:    "Initialized private credentials.",
	}
}

func initPublicStageOnly(ctx context.Context, pub *git.Tree, cred PublicCredentials) git.ChangeNoResult {
	if _, err := pub.Filesystem.Stat(PublicCredentialsNS.Path()); err == nil {
		must.Errorf(ctx, "public credentials file already exists")
	}
	git.ToFileStage(ctx, pub, PublicCredentialsNS.Path(), cred)
	return git.ChangeNoResult{
		Msg: "Initialized public credentials.",
	}
}
