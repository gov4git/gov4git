package ballotproto

import (
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/lib4git/util"
)

func AdsToBallotNames(ads []Advertisement) []string {
	names := make([]string, len(ads))
	for i := range ads {
		names[i] = ads[i].Name.GitPath()
	}
	return names
}

func FilterFrozenAds(frozen bool, ads []Advertisement) []Advertisement {
	r := []Advertisement{}
	for _, ad := range ads {
		if ad.Frozen == frozen {
			r = append(r, ad)
		}
	}
	return r
}

func FilterOpenClosedAds(closed bool, ads []Advertisement) []Advertisement {
	r := []Advertisement{}
	for _, ad := range ads {
		if ad.Closed == closed {
			r = append(r, ad)
		}
	}
	return r
}

func FilterWithParticipants(groups []member.Group, ads []Advertisement) []Advertisement {
	r := []Advertisement{}
	for _, ad := range ads {
		if util.IsIn(ad.Participants, groups...) {
			r = append(r, ad)
		}
	}
	return r
}
