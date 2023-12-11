package id

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/materials"
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
	privChg := Init_Local(ctx, ownerCloned)

	ownerCloned.Public.Push(ctx)
	ownerCloned.Private.Push(ctx)
	return privChg
}

func Init_Local(
	ctx context.Context,
	ownerCloned OwnerCloned,
) git.Change[form.None, PrivateCredentials] {

	privChg := initPrivate_StageOnly(ctx, ownerCloned.Private.Tree(), ownerCloned.Address())
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
	git.StringToFileStage(ctx, priv, proto.ReadmeNS, PrivateReadmeMD)
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
	git.StringToFileStage(ctx, pub, proto.ReadmeNS, PublicReadmeMD)
	return git.NewChangeNoResult("Initialized public credentials.", "id_init_public")
}

var (
	PublicReadmeMD  = readmeMDHeader("This is a Gov4Git public identity repository.") + readmeBody
	PrivateReadmeMD = readmeMDHeader("This is a Gov4Git private identity repository.") + readmeBody

	readmeBody = fmt.Sprintf(`
		[Gov4Git](%s) is a decentralized governance and management system for git projects.

		Learn about Gov4Git:
		[Gov4Git on GitHub](%s).
		[Gov4Git on Twitter/X](%s).
		`,
		materials.Gov4GitWebsiteURL, materials.Gov4GitGithubURL, materials.Gov4GitXURL)
)

func readmeMDHeader(title string) string {
	return fmt.Sprintf(
		"## <a href=%q><img src=%q alt=%q width=\"65\" /></a> %s\n\n",
		materials.Gov4GitWebsiteURL, materials.Gov4GitAvatarURL, title, title)
}
