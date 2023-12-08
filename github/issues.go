package github

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func ImportIssuesForPrioritization(
	ctx context.Context,
	repo Repo,
	githubClient *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OwnerAddress,
) git.Change[form.Map, ImportedIssues] {

	base.Infof("importing issues for prioritization ...")
	govCloned := gov.CloneOwner(ctx, govAddr)
	issuesCausingChange := ImportIssuesForPrioritization_StageOnly(ctx, repo, githubClient, govAddr, govCloned)
	chg := git.NewChange[form.Map, ImportedIssues](
		fmt.Sprintf("Import %d GitHub issues", len(issuesCausingChange)),
		"github_import",
		form.Map{},
		issuesCausingChange,
		nil,
	)
	return proto.CommitIfChanged(ctx, govCloned.Public, chg)
}

func ImportIssuesForPrioritization_StageOnly(
	ctx context.Context,
	repo Repo,
	githubClient *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OwnerAddress,
	govCloned gov.OwnerCloned,
) ImportedIssues {

	// load github issues and governance ballots, and
	// index them under a common key space
	loadPR := func(ctx context.Context,
		repo Repo,
		issue *github.Issue,
	) bool {

		return false
	}
	_, issues := LoadIssues(ctx, githubClient, repo, loadPR)
	ballots := filterIssuesForPrioritization(ballot.List_Local(ctx, govCloned.PublicClone()))

	// ensure every issue has a corresponding up-to-date ballot
	causedChange := ImportedIssues{}
	for k, ghIssue := range issues {
		if ghIssue.ForPrioritization {
			if govBallot, ok := ballots[k]; ok { // ballot for issue already exists, update it

				must.Assertf(ctx, ns.Equal(ghIssue.BallotName().NS(), govBallot.Name.NS()),
					"issue ballot name %v and actual ballot name %v mismatch", ghIssue.BallotName(), govBallot.Name)

				switch {
				case ghIssue.Closed && govBallot.Closed:
					// nothing to do
				case ghIssue.Closed && !govBallot.Closed:
					UpdateMeta_StageOnly(ctx, repo, govAddr, govCloned, ghIssue, govBallot)
					UpdateFrozen_StageOnly(ctx, repo, govAddr, govCloned, ghIssue, govBallot)
					ballot.Close_StageOnly(ctx, govCloned, ghIssue.BallotName(), account.BurnAccountID)
					causedChange = append(causedChange, ghIssue)
				case !ghIssue.Closed && govBallot.Closed:
					UpdateMeta_StageOnly(ctx, repo, govAddr, govCloned, ghIssue, govBallot)
					ballot.Reopen_StageOnly(ctx, govCloned, ghIssue.BallotName())
					UpdateFrozen_StageOnly(ctx, repo, govAddr, govCloned, ghIssue, govBallot)
					causedChange = append(causedChange, ghIssue)
				case !ghIssue.Closed && !govBallot.Closed:
					c1 := UpdateMeta_StageOnly(ctx, repo, govAddr, govCloned, ghIssue, govBallot)
					c2 := UpdateFrozen_StageOnly(ctx, repo, govAddr, govCloned, ghIssue, govBallot)
					if c1 || c2 {
						causedChange = append(causedChange, ghIssue)
					}
				}

			} else { // no ballot for this issue, create it
				ballot.Open_StageOnly(
					ctx,
					load.QVStrategyName,
					govCloned,
					ghIssue.BallotName(),
					ghIssue.Title,
					ghIssue.Body,
					[]string{PrioritizeBallotChoice},
					member.Everybody,
				)
				if ghIssue.Locked {
					ballot.Freeze_StageOnly(ctx, govCloned, ghIssue.BallotName())
				}
				if ghIssue.Closed {
					ballot.Close_StageOnly(ctx, govCloned, ghIssue.BallotName(), account.BurnAccountID)
				}
				causedChange = append(causedChange, ghIssue)
			}
		} else { // issue is not for prioritization, freeze ballot if it exists and is open
			if govBallot, ok := ballots[k]; ok { // ballot for issue already exists, update it
				// if ballot closed, do nothing
				// if ballot frozen, do nothing
				// otherwise, freeze ballot
				if !govBallot.Closed && !govBallot.Frozen {
					ballot.Freeze_StageOnly(ctx, govCloned, ghIssue.BallotName())
					causedChange = append(causedChange, ghIssue)
				}
			}
		}
	}
	causedChange.Sort()

	// don't touch ballots that have no corresponding issue

	return causedChange
}

func filterIssuesForPrioritization(ads []common.Advertisement) map[string]common.Advertisement {
	filtered := map[string]common.Advertisement{}
	for _, ad := range ads {
		if len(ad.Name) == 3 && ad.Name[0] == ImportedGithubPrefix && (ad.Name[1] == ImportedIssuePrefix || ad.Name[1] == ImportedPullPrefix) {
			key := ad.Name[2]
			if _, err := strconv.Atoi(key); err == nil {
				filtered[key] = ad
			}
		}
	}
	return filtered
}

func UpdateMeta_StageOnly(
	ctx context.Context,
	repo Repo,
	govAddr gov.OwnerAddress,
	govCloned gov.OwnerCloned,
	ghIssue ImportedIssue,
	govBallot common.Advertisement,
) (changed bool) {
	if ghIssue.Title == govBallot.Title && ghIssue.Body == govBallot.Description {
		return false
	}
	ballot.Change_StageOnly(ctx, govCloned, ghIssue.BallotName(), ghIssue.Title, ghIssue.Body)
	return true
}

func UpdateFrozen_StageOnly(
	ctx context.Context,
	repo Repo,
	govAddr gov.OwnerAddress,
	govCloned gov.OwnerCloned,
	ghIssue ImportedIssue,
	govBallot common.Advertisement,
) (changed bool) {
	switch {
	case ghIssue.Locked && govBallot.Frozen:
		return false
	case ghIssue.Locked && !govBallot.Frozen:
		ballot.Freeze_StageOnly(ctx, govCloned, ghIssue.BallotName())
		return true
	case !ghIssue.Locked && govBallot.Frozen:
		ballot.Unfreeze_StageOnly(ctx, govCloned, ghIssue.BallotName())
		return true
	case !ghIssue.Locked && !govBallot.Frozen:
		return false
	}
	panic("unreachable")
}
