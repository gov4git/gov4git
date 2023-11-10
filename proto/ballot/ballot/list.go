package ballot

import (
	"context"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func List(
	ctx context.Context,
	govAddr gov.Address,
) []common.Advertisement {

	return List_Local(ctx, git.CloneOne(ctx, git.Address(govAddr)).Tree())
}

func List_Local(
	ctx context.Context,
	govTree *git.Tree,
) []common.Advertisement {

	ballotsNS := common.BallotPath(common.BallotName{})

	files, err := git.ListFilesRecursively(govTree, ballotsNS)
	must.NoError(ctx, err)

	ads := []common.Advertisement{}
	for _, file := range files {
		if file.Base() != common.AdFilebase {
			continue
		}
		var ad common.Advertisement
		err := must.Try(
			func() {
				ad = git.FromFile[common.Advertisement](ctx, govTree, file)
			},
		)
		if err != nil {
			continue
		}
		ads = append(ads, ad)
	}

	return ads
}

func ListFilter(
	ctx context.Context,
	govAddr gov.Address,
	onlyOpen bool,
	onlyClosed bool,
	onlyFrozen bool,
	withParticipant member.User,
) []common.Advertisement {

	return ListFilter_Local(ctx, gov.Clone(ctx, govAddr).Tree(), onlyOpen, onlyClosed, onlyFrozen, withParticipant)
}

func ListFilter_Local(
	ctx context.Context,
	govTree *git.Tree,
	onlyOpen bool,
	onlyClosed bool,
	onlyFrozen bool,
	withParticipant member.User,
) []common.Advertisement {

	ads := List_Local(ctx, govTree)
	if onlyOpen {
		ads = common.FilterOpenClosedAds(false, ads)
	}
	if onlyClosed {
		ads = common.FilterOpenClosedAds(true, ads)
	}
	if onlyFrozen {
		ads = common.FilterFrozenAds(true, ads)
	}
	if withParticipant != "" {
		userGroups := member.ListUserGroups_Local(ctx, govTree, withParticipant)
		ads = common.FilterWithParticipants(userGroups, ads)
	}
	return ads
}
