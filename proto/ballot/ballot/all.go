package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func TallyAll(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
) git.Change[form.Map, []common.Tally] {

	govOwner := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := TallyAllStageOnly(ctx, govAddr, govOwner)
	if len(chg.Result) == 0 {
		return chg
	}
	proto.Commit(ctx, govOwner.Public.Tree(), chg)
	govOwner.Public.Push(ctx)
	return chg
}

func TallyAllStageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govOwner id.OwnerCloned,
) git.Change[form.Map, []common.Tally] {

	// list all open ballots
	communityTree := govOwner.Public.Tree()
	ads := ListLocal(ctx, communityTree, false)

	// compute union of all voter accounts from all open ballots
	adVoters := make([]adVoters, len(ads))
	allVoters := map[member.User]member.Account{}
	for i, ad := range ads {
		adVoters[i].ad = &ad
		adVoters[i].voterAccounts = map[member.User]member.Account{}
		adVoters[i].voterClones = map[member.User]git.Cloned{}
		adVoters[i].voters = member.ListGroupUsersLocal(ctx, communityTree, ad.Participants)
		for _, user := range adVoters[i].voters {
			if _, ok := allVoters[user]; !ok {
				account := member.GetUserLocal(ctx, communityTree, user)
				adVoters[i].voterAccounts[user] = account
				allVoters[user] = account
			}
		}
	}

	// fetch repos of all participating users
	allVoterClones := map[member.User]git.Cloned{}
	for u, a := range allVoters {
		allVoterClones[u] = git.CloneOne(ctx, git.Address(a.PublicAddress))
	}

	// populate ad voter structures
	for i, ad := range adVoters {
		for u := range ad.voterAccounts {
			adVoters[i].voterClones[u] = allVoterClones[u]
		}
	}

	// perform tallies for all open ballots
	tallyChanges := []git.Change[map[string]form.Form, common.Tally]{}
	tallies := []common.Tally{}
	for _, adv := range adVoters {
		if tallyChg, changed := tallyVotersClonedStageOnly(ctx, govAddr, govOwner, adv.ad.Name, adv.voterAccounts, adv.voterClones); changed {
			tallyChanges = append(tallyChanges, tallyChg)
			tallies = append(tallies, tallyChg.Result)
		}
	}

	return git.NewChange(
		fmt.Sprintf("Tallied votes on all ballots"),
		"ballot_tally_all",
		form.Map{},
		tallies,
		form.ToForms(tallyChanges),
	)
}

type adVoters struct {
	ad            *common.Advertisement
	voters        []member.User
	voterAccounts map[member.User]member.Account
	voterClones   map[member.User]git.Cloned
}
