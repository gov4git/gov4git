package etc

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func SetSettings(
	ctx context.Context,
	addr gov.Address,
	config Settings,
) git.Change[Settings, form.None] {

	cloned := gov.Clone(ctx, addr)
	chg := SetSettings_StageOnly(ctx, cloned, config)
	return proto.CommitIfChanged(ctx, cloned, chg)
}

func SetSettings_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
	config Settings,
) git.Change[Settings, form.None] {

	git.ToFileStage[Settings](ctx, cloned.Tree(), SettingsNS, config)
	return git.NewChange[Settings, form.None](
		"Change settings",
		"etc_set_settings",
		config,
		form.None{},
		nil,
	)
}

func GetSettings(
	ctx context.Context,
	addr gov.Address,
) Settings {

	cloned := gov.Clone(ctx, addr)
	return GetSettings_StageOnly(ctx, cloned)
}

func GetSettings_StageOnly(
	ctx context.Context,
	cloned gov.Cloned,
) Settings {

	config, err := git.TryFromFile[Settings](ctx, cloned.Tree(), SettingsNS)
	if git.IsNotExist(err) {
		return DefaultSettings
	}
	must.NoError(ctx, err)
	return config
}
