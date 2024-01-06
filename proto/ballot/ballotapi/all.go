package ballotapi

import (
	"context"
	"fmt"
	"sync"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func TallyAll(
	ctx context.Context,
	addr gov.OwnerAddress,
	maxPar int,
) git.Change[form.Map, []ballotproto.Tally] {

	base.Infof("fetching and tallying community votes ...")

	govOwner := gov.CloneOwner(ctx, addr)
	chg := TallyAll_StageOnly(ctx, govOwner, maxPar)
	if len(chg.Result) == 0 {
		return chg
	}
	proto.Commit(ctx, govOwner.Public.Tree(), chg)
	govOwner.Public.Push(ctx)
	return chg
}

func TallyAll_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	maxPar int,
) git.Change[form.Map, []ballotproto.Tally] {

	// list all open ballots
	ads := ballotproto.FilterOpenClosedAds(false, List_Local(ctx, cloned.PublicClone()))

	// compute union of all voter accounts from all open ballots
	participatingVoters := make([]participatingVoters, len(ads))
	allVoters := map[member.User]member.UserProfile{}
	for i, ad := range ads {
		participatingVoters[i] = *loadParticipatingVoters(ctx, cloned.PublicClone(), ad)
		for user, acct := range participatingVoters[i].VoterAccounts {
			allVoters[user] = acct
		}
	}

	// fetch repos of all participating users
	allVoterClones := clonePar(ctx, allVoters, maxPar)

	// populate participating voter clones
	for _, pv := range participatingVoters {
		pv.attachVoterClones(ctx, allVoterClones)
	}

	// perform tallies for all open ballots
	tallyChanges := []git.Change[map[string]form.Form, ballotproto.Tally]{}
	tallies := []ballotproto.Tally{}
	for _, pv := range participatingVoters {
		if tallyChg, changed := TallyVotersCloned_StageOnly(ctx, cloned, pv.Ad.ID, pv.VoterAccounts, pv.VoterClones); changed {
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

func clonePar(ctx context.Context, userAccounts map[member.User]member.UserProfile, maxPar int) map[member.User]git.Cloned {

	must.Assertf(ctx, maxPar > 0, "clone parallelism must be greater than zero")

	var wg sync.WaitGroup
	wg.Add(len(userAccounts))

	sem := make(chan bool, maxPar)

	var allLock sync.Mutex
	allClones := map[member.User]git.Cloned{}

	for u, a := range userAccounts {
		sem <- true
		go func(u member.User, a member.UserProfile) {

			base.Infof("cloning voter %v repository %v", u, a.PublicAddress)
			cloned, err := git.TryCloneOne(ctx, git.Address(a.PublicAddress))
			if err != nil {
				base.Infof("user %v repository %v unresponsive (%v)", u, a.PublicAddress, err)
			} else {
				base.Infof("user %v repository %v cloned successfully (%v)", u, a.PublicAddress, err)
				allLock.Lock()
				allClones[u] = cloned
				allLock.Unlock()
			}

			<-sem
			wg.Done()
		}(u, a)
	}

	wg.Wait()

	return allClones
}
