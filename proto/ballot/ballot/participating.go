package ballot

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/common"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

type participatingVoters struct {
	Ad            common.Advertisement
	Voters        []member.User
	VoterAccounts map[member.User]member.UserProfile
	VoterClones   map[member.User]git.Cloned
}

func loadParticipatingVoters(ctx context.Context, cloned gov.Cloned, ad common.Advertisement) *participatingVoters {
	pv := &participatingVoters{}
	pv.Ad = ad
	pv.VoterAccounts = map[member.User]member.UserProfile{}
	pv.VoterClones = map[member.User]git.Cloned{}
	pv.Voters = member.ListGroupUsers_Local(ctx, cloned, ad.Participants)
	for _, user := range pv.Voters {
		account := member.GetUser_Local(ctx, cloned, user)
		pv.VoterAccounts[user] = account
	}
	return pv
}

// attachVoterClones adds clones for all participating voter accounts.
// It also removes voter accounts that were not cloned from the structure.
func (pv *participatingVoters) attachVoterClones(ctx context.Context, votersCloned map[member.User]git.Cloned) {
	for u := range pv.VoterAccounts {
		if cloned, clonedOK := votersCloned[u]; clonedOK {
			pv.VoterClones[u] = cloned
		} else {
			delete(pv.VoterAccounts, u)
		}
	}
	must.Assertf(ctx, len(pv.VoterAccounts) == len(pv.VoterClones), "voter accounts and clones must have the same keys")
}
