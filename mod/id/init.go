package id

import (
	"context"

	"github.com/gov4git/gov4git/mod"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Init(
	ctx context.Context,
	ownerAddr OwnerAddress,
) git.Change[PrivateCredentials] {
	ownerRepo, ownerTree := CloneOwner(ctx, ownerAddr)
	privChg := InitLocal(ctx, ownerAddr, ownerTree)

	git.Push(ctx, ownerRepo.Private)
	git.Push(ctx, ownerRepo.Public)
	return privChg
}

func InitLocal(
	ctx context.Context,
	ownerAddr OwnerAddress,
	ownerTree OwnerTree,
) git.Change[PrivateCredentials] {
	privChg := initPrivateStageOnly(ctx, ownerTree.Private, ownerAddr)
	pubChg := initPublicStageOnly(ctx, ownerTree.Public, privChg.Result.PublicCredentials)
	mod.Commit(ctx, ownerTree.Private, privChg.Msg)
	mod.Commit(ctx, ownerTree.Public, pubChg.Msg)
	return privChg
}

func initPrivateStageOnly(ctx context.Context, priv *git.Tree, ownerAddr OwnerAddress) git.Change[PrivateCredentials] {
	if _, err := priv.Filesystem.Stat(PrivateCredentialsNS.Path()); err == nil {
		must.Errorf(ctx, "private credentials file already exists")
	}
	cred, err := GenerateCredentials(git.Address(ownerAddr.Public), git.Address(ownerAddr.Private))
	if err != nil {
		must.Panic(ctx, err)
	}
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
