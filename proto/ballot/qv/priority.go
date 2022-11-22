package qv

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

type PriorityPoll struct {
	UseVotingCredits bool `json:"use_voting_credits"`
}

const PriorityPollName = "priority_poll"

func (x PriorityPoll) Name() string {
	return PriorityPollName
}

func (x PriorityPoll) VerifyElections(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	govAddr gov.GovAddress,
	voterTree id.OwnerTree,
	govTree *git.Tree,
	ad common.Advertisement,
	elections common.Elections,
) {
	if x.UseVotingCredits {
		spend := 0.0
		for _, el := range elections {
			spend += math.Abs(el.VoteStrengthChange)
		}
		user := member.LookupUserByAddressLocal(ctx, govTree, voterAddr.Public)
		if len(user) == 0 {
			must.Errorf(ctx, "cannot find user with address %v in the community", voterAddr.Public)
		}
		available := balance.GetLocal(ctx, govTree, user[0], VotingCredits)
		must.Assertf(ctx, available >= spend, "insufficient voting credits %v for elections costing %v", available, spend)
	}
}

func (x PriorityPoll) Tally(
	ctx context.Context,
	govRepo id.OwnerRepo,
	govTree id.OwnerTree,
	ad *common.Advertisement,
	prior *common.Tally,
	fetched []common.FetchedVote,
) git.Change[common.Tally] {

	// TODO: key on member+address to account for changes in user → address mapping
	fetchedVotesMap := map[member.User]common.FetchedVote{}

	// load prior participant votes
	if prior != nil {
		for _, fv := range prior.Votes {
			fetchedVotesMap[fv.Voter] = fv
		}
	}

	// pay for voter elections with voting credits
	paid := []common.FetchedVote{}
	if x.UseVotingCredits {
		for _, fv := range fetched {
			paidElections := common.Elections{}
			for _, el := range fv.Elections {
				err := balance.TryTransferStageOnly(
					ctx,
					govTree.Public,
					fv.Voter, VotingCredits,
					fv.Voter, VotingCreditsOnHold,
					math.Abs(el.VoteStrengthChange),
				)
				if err != nil {
					base.Infof("not enough voting credits for voter %v election", fv.Voter)
					continue
				}
				paidElections = append(paidElections, el)
			}
			paid = append(paid, common.FetchedVote{Voter: fv.Voter, Address: fv.Address, Elections: paidElections})
		}
	} else {
		paid = fetched
	}

	// update votes
	for _, fv := range paid {
		prior := fetchedVotesMap[fv.Voter]
		fv.Elections = append(fv.Elections, prior.Elections...)
		fetchedVotesMap[fv.Voter] = fv
	}

	// sort fetched votes
	fetchedVotes := common.FetchedVotes{}
	for _, fetchedVote := range fetchedVotesMap {
		fetchedVotes = append(fetchedVotes, fetchedVote)
	}
	sort.Sort(fetchedVotes)

	voterElectionsMap := map[member.User]map[string]float64{}
	for voter, fetchedVote := range fetchedVotesMap {

		voterChoiceScoresMap, ok := voterElectionsMap[voter]
		if !ok {
			voterChoiceScoresMap = map[string]float64{}
			voterElectionsMap[voter] = voterChoiceScoresMap
		}

		for _, el := range fetchedVote.Elections {
			voterChoiceScoresMap[el.VoteChoice] += el.VoteStrengthChange
		}
	}

	choiceScoresMap := map[string]float64{}
	for _, voterChoiceScoresMap := range voterElectionsMap {
		for choice, score := range voterChoiceScoresMap {
			sign := 1.0
			if score < 0 {
				sign = -1.0
			}
			choiceScoresMap[choice] += sign * math.Sqrt(math.Abs(score)) //qv
		}
	}

	// sort choice scores
	choiceScores := common.ChoiceScores{}
	for choice, score := range choiceScoresMap {
		choiceScores = append(choiceScores, common.ChoiceScore{Choice: choice, Score: score})
	}
	sort.Sort(choiceScores)

	return git.Change[common.Tally]{
		Result: common.Tally{
			Ad:     *ad,
			Votes:  fetchedVotes,
			Scores: choiceScores,
		},
		Msg: fmt.Sprintf("Tallied QV priority poll scores for %v", ad.Name),
	}
}

func (x PriorityPoll) Close(
	ctx context.Context,
	govRepo id.OwnerRepo,
	govTree id.OwnerTree,
	ad *common.Advertisement,
	tally *common.Tally,
	summary common.Summary,
) git.Change[common.Outcome] {

	if x.UseVotingCredits && summary == SummaryAdopted {

		// compute credits spent by each user
		us := map[member.User]float64{}
		for _, vote := range tally.Votes {
			u := us[vote.Voter]
			for _, el := range vote.Elections {
				u += el.VoteStrengthChange
			}
			us[vote.Voter] = u
		}

		// refund users
		for user, spent := range us {
			if spent < 0 {
				continue // don't refund voters against
			}
			onHold := balance.GetLocal(ctx, govTree.Public, user, VotingCreditsOnHold)
			refund := min(spent, onHold)
			balance.TransferStageOnly(ctx, govTree.Public, user, VotingCreditsOnHold, user, VotingCredits, refund)
		}
	}

	return git.Change[common.Outcome]{
		Result: common.Outcome{
			Summary: summary,
			Scores:  tally.Scores,
		},
		Msg: fmt.Sprintf("closed ballot %v with outcome %v", ad.Name, summary),
	}
}

func min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}
