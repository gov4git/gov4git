package qv

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/gov4git/gov4git/mod/balance"
	"github.com/gov4git/gov4git/mod/ballot/proto"
	"github.com/gov4git/gov4git/mod/id"
	"github.com/gov4git/gov4git/mod/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/git"
)

type PriorityPoll struct {
	UseVotingCredits bool `json:"use_voting_credits"`
}

const PriorityPollName = "priority_poll"

func (x PriorityPoll) Name() string {
	return PriorityPollName
}

func (x PriorityPoll) Tally(
	ctx context.Context,
	govRepo id.OwnerRepo,
	govTree id.OwnerTree,
	ad *proto.Advertisement,
	prior *proto.Tally,
	fetched []proto.FetchedVote,
) git.Change[proto.Tally] {

	// TODO: key on member+address to account for changes in user → address mapping
	fetchedVotesMap := map[member.User]proto.FetchedVote{}

	// load prior participant votes
	if prior != nil {
		for _, fv := range prior.Votes {
			fetchedVotesMap[fv.Voter] = fv
		}
	}

	// pay for voter elections with voting credits
	paid := []proto.FetchedVote{}
	if x.UseVotingCredits {
		for _, fv := range fetched {
			paidElections := proto.Elections{}
			for _, el := range fv.Elections {
				err := balance.TryTransferStageOnly(
					ctx,
					govTree.Public,
					fv.Voter, VotingCredits,
					fv.Voter, VotingCreditsOnHold,
					el.VoteStrengthChange,
				)
				if err != nil {
					base.Infof("not enough voting credits for voter %v election", fv.Voter)
					continue
				}
				paidElections = append(paidElections, el)
			}
			paid = append(paid, proto.FetchedVote{Voter: fv.Voter, Address: fv.Address, Elections: paidElections})
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
	fetchedVotes := proto.FetchedVotes{}
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
	choiceScores := proto.ChoiceScores{}
	for choice, score := range choiceScoresMap {
		choiceScores = append(choiceScores, proto.ChoiceScore{Choice: choice, Score: score})
	}
	sort.Sort(choiceScores)

	return git.Change[proto.Tally]{
		Result: proto.Tally{
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
	ad *proto.Advertisement,
	tally *proto.Tally,
	summary proto.Summary,
) git.Change[proto.Outcome] {

	return git.Change[proto.Outcome]{
		Result: proto.Outcome{
			Summary: summary,
			Scores:  tally.Scores,
		},
		Msg: fmt.Sprintf("closed ballot %v with outcome %v", ad.Name, summary),
	}
}
