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
		adVoters[i].Ad = ad
		adVoters[i].VoterAccounts = map[member.User]member.Account{}
		adVoters[i].VoterClones = map[member.User]git.Cloned{}
		adVoters[i].Voters = member.ListGroupUsersLocal(ctx, communityTree, ad.Participants)
		for _, user := range adVoters[i].Voters {
			account := member.GetUserLocal(ctx, communityTree, user)
			adVoters[i].VoterAccounts[user] = account
			allVoters[user] = account
		}
	}

	// fetch repos of all participating users
	allVoterClones := map[member.User]git.Cloned{}
	for u, a := range allVoters {
		allVoterClones[u] = git.CloneOne(ctx, git.Address(a.PublicAddress))
	}

	// populate ad voter structures
	for i, ad := range adVoters {
		for u := range ad.VoterAccounts {
			adVoters[i].VoterClones[u] = allVoterClones[u]
		}
	}

	// perform tallies for all open ballots
	tallyChanges := []git.Change[map[string]form.Form, common.Tally]{}
	tallies := []common.Tally{}
	fmt.Println("adVoters: ", form.SprintJSON(adVoters))
	for _, adv := range adVoters {
		fmt.Println("XXX tallying", form.SprintJSON(adv))
		if tallyChg, changed := tallyVotersClonedStageOnly(ctx, govAddr, govOwner, adv.Ad.Name, adv.VoterAccounts, adv.VoterClones); changed {
			fmt.Println("	XXX changed", changed)
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
	Ad            common.Advertisement
	Voters        []member.User
	VoterAccounts map[member.User]member.Account
	VoterClones   map[member.User]git.Cloned
}
