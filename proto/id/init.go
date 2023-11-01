package id

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Init(
	ctx context.Context,
	ownerAddr OwnerAddress,
) git.Change[form.None, PrivateCredentials] {
	ownerCloned := CloneOwner(ctx, ownerAddr)
	privChg := Init_Local(ctx, ownerAddr, ownerCloned)

	ownerCloned.Public.Push(ctx)
	ownerCloned.Private.Push(ctx)
	return privChg
}

func Init_Local(
	ctx context.Context,
	ownerAddr OwnerAddress,
	ownerCloned OwnerCloned,
) git.Change[form.None, PrivateCredentials] {

	privChg := initPrivate_StageOnly(ctx, ownerCloned.Private.Tree(), ownerAddr)
	pubChg := initPublic_StageOnly(ctx, ownerCloned.Public.Tree(), privChg.Result.PublicCredentials)
	proto.Commit(ctx, ownerCloned.Private.Tree(), privChg)
	proto.Commit(ctx, ownerCloned.Public.Tree(), pubChg)
	return privChg
}

func initPrivate_StageOnly(ctx context.Context, priv *git.Tree, ownerAddr OwnerAddress) git.Change[form.None, PrivateCredentials] {
	if _, err := git.TreeStat(ctx, priv, PrivateCredentialsNS); err == nil {
		must.Errorf(ctx, "private credentials file already exists")
	}
	cred, err := GenerateCredentials()
	must.NoError(ctx, err)
	git.ToFileStage(ctx, priv, PrivateCredentialsNS, cred)
	return git.NewChange(
		"Initialized private credentials.",
		"id_init_private",
		form.None{},
		cred,
		nil,
	)
}

func initPublic_StageOnly(ctx context.Context, pub *git.Tree, cred PublicCredentials) git.ChangeNoResult {
	if _, err := git.TreeStat(ctx, pub, PublicCredentialsNS); err == nil {
		must.Errorf(ctx, "public credentials file already exists")
	}
	git.ToFileStage(ctx, pub, PublicCredentialsNS, cred)
	return git.NewChangeNoResult("Initialized public credentials.", "id_init_public")
}
