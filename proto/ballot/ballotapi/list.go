package ballotapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func List(
	ctx context.Context,
	govAddr gov.Address,

) ballotproto.Advertisements {

	return List_Local(ctx, gov.Clone(ctx, govAddr))
}

func List_Local(
	ctx context.Context,
	cloned gov.Cloned,

) ballotproto.Advertisements {

	ballotsNS := ballotproto.BallotPath(ballotproto.BallotName{})

	files, err := git.ListFilesRecursively(cloned.Tree(), ballotsNS)
	must.NoError(ctx, err)

	ads := ballotproto.Advertisements{}
	for _, file := range files {
		if file.Base() != ballotproto.AdFilebase {
			continue
		}
		var ad ballotproto.Advertisement
		err := must.Try(
			func() {
				ad = git.FromFile[ballotproto.Advertisement](ctx, cloned.Tree(), file)
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

) ballotproto.Advertisements {

	return ListFilter_Local(ctx, gov.Clone(ctx, govAddr), onlyOpen, onlyClosed, onlyFrozen, withParticipant)
}

func ListFilter_Local(
	ctx context.Context,
	cloned gov.Cloned,
	onlyOpen bool,
	onlyClosed bool,
	onlyFrozen bool,
	withParticipant member.User,

) ballotproto.Advertisements {

	ads := List_Local(ctx, cloned)
	if onlyOpen {
		ads = ballotproto.FilterOpenClosedAds(false, ads)
	}
	if onlyClosed {
		ads = ballotproto.FilterOpenClosedAds(true, ads)
	}
	if onlyFrozen {
		ads = ballotproto.FilterFrozenAds(true, ads)
	}
	if withParticipant != "" {
		userGroups := member.ListUserGroups_Local(ctx, cloned, withParticipant)
		ads = ballotproto.FilterWithParticipants(userGroups, ads)
	}
	ads.Sort()
	return ads
}
