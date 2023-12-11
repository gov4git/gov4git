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

) common.Advertisements {

	return List_Local(ctx, gov.Clone(ctx, govAddr))
}

func List_Local(
	ctx context.Context,
	cloned gov.Cloned,

) common.Advertisements {

	ballotsNS := common.BallotPath(common.BallotName{})

	files, err := git.ListFilesRecursively(cloned.Tree(), ballotsNS)
	must.NoError(ctx, err)

	ads := common.Advertisements{}
	for _, file := range files {
		if file.Base() != common.AdFilebase {
			continue
		}
		var ad common.Advertisement
		err := must.Try(
			func() {
				ad = git.FromFile[common.Advertisement](ctx, cloned.Tree(), file)
			},
		)
		if err != nil {
			continue
		}
		ads = append(ads, ad)
	}

	ads.Sort()
	return ads
}

func ListFilter(
	ctx context.Context,
	govAddr gov.Address,
	onlyOpen bool,
	onlyClosed bool,
	onlyFrozen bool,
	withParticipant member.User,

) common.Advertisements {

	return ListFilter_Local(ctx, gov.Clone(ctx, govAddr), onlyOpen, onlyClosed, onlyFrozen, withParticipant)
}

func ListFilter_Local(
	ctx context.Context,
	cloned gov.Cloned,
	onlyOpen bool,
	onlyClosed bool,
	onlyFrozen bool,
	withParticipant member.User,

) common.Advertisements {

	ads := List_Local(ctx, cloned)
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
		userGroups := member.ListUserGroups_Local(ctx, cloned, withParticipant)
		ads = common.FilterWithParticipants(userGroups, ads)
	}
	ads.Sort()
	return ads
}
