package qv

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/gov4git/gov4git/mod/ballot"
	"github.com/gov4git/gov4git/mod/id"
	"github.com/gov4git/gov4git/mod/member"
	"github.com/gov4git/lib4git/git"
)

type PriorityPoll struct{}

func (x PriorityPoll) StrategyName() string {
	return "priority_poll"
}

func (x PriorityPoll) Tally(
	ctx context.Context,
	govRepo id.OwnerRepo,
	govTree id.OwnerTree,
	ad *ballot.Advertisement,
	prior *ballot.TallyForm,
	fetched []ballot.FetchedVote,
) git.Change[ballot.TallyForm] {

	// TODO: key on member+address to account for changes in user → address changes
	fetchedVotesMap := map[member.User]ballot.FetchedVote{}

	// load prior participant votes
	if prior != nil {
		for _, fv := range prior.FetchedVotes {
			fetchedVotesMap[fv.Voter] = fv
		}
	}

	// update votes
	for _, fv := range fetched {
		prior := fetchedVotesMap[fv.Voter]
		fv.Elections = append(fv.Elections, prior.Elections...)
		fetchedVotesMap[fv.Voter] = fv
	}

	// sort fetched votes
	fetchedVotes := ballot.FetchedVotes{}
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
	choiceScores := ballot.ChoiceScores{}
	for choice, score := range choiceScoresMap {
		choiceScores = append(choiceScores, ballot.ChoiceScore{Choice: choice, Score: score})
	}
	sort.Sort(choiceScores)

	return git.Change[ballot.TallyForm]{
		Result: ballot.TallyForm{
			Ad:           *ad,
			FetchedVotes: fetchedVotes,
			ChoiceScores: choiceScores,
		},
		Msg: fmt.Sprintf("Tallied QV priority poll scores for %v", ad.Name),
	}
}
