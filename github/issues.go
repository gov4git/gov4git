package github

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
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
	cloned gov.OwnerCloned,
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
	ballots := filterIssuesForPrioritization(ballotapi.List_Local(ctx, cloned.PublicClone()))

	// ensure every issue has a corresponding up-to-date ballot
	causedChange := ImportedIssues{}
	for k, ghIssue := range issues {
		if ghIssue.ForPrioritization {
			if govBallot, ok := ballots[k]; ok { // ballot for issue already exists, update it

				must.Assertf(ctx, ghIssue.BallotName() == govBallot.ID,
					"issue ballot name %v and actual ballot name %v mismatch", ghIssue.BallotName(), govBallot.ID)

				switch {
				case ghIssue.Closed && govBallot.Closed:
					// nothing to do
				case ghIssue.Closed && !govBallot.Closed:
					UpdateMeta_StageOnly(ctx, repo, govAddr, cloned, ghIssue, govBallot)
					UpdateFrozen_StageOnly(ctx, repo, govAddr, cloned, ghIssue, govBallot)
					ballotapi.Close_StageOnly(ctx, cloned, ghIssue.BallotName(), account.BurnAccountID)
					causedChange = append(causedChange, ghIssue)
				case !ghIssue.Closed && govBallot.Closed:
					UpdateMeta_StageOnly(ctx, repo, govAddr, cloned, ghIssue, govBallot)
					ballotapi.Reopen_StageOnly(ctx, cloned, ghIssue.BallotName())
					UpdateFrozen_StageOnly(ctx, repo, govAddr, cloned, ghIssue, govBallot)
					causedChange = append(causedChange, ghIssue)
				case !ghIssue.Closed && !govBallot.Closed:
					c1 := UpdateMeta_StageOnly(ctx, repo, govAddr, cloned, ghIssue, govBallot)
					c2 := UpdateFrozen_StageOnly(ctx, repo, govAddr, cloned, ghIssue, govBallot)
					if c1 || c2 {
						causedChange = append(causedChange, ghIssue)
					}
				}

			} else { // no ballot for this issue, create it
				ballotapi.Open_StageOnly(
					ctx,
					ballotio.QVStrategyName,
					cloned,
					ghIssue.BallotName(),
					account.NobodyAccountID,
					ghIssue.Purpose(),
					"",
					ghIssue.Title,
					ghIssue.Body,
					[]string{PrioritizeBallotChoice},
					member.Everybody,
				)
				if ghIssue.Locked {
					ballotapi.Freeze_StageOnly(ctx, cloned, ghIssue.BallotName())
				}
				if ghIssue.Closed {
					ballotapi.Close_StageOnly(ctx, cloned, ghIssue.BallotName(), account.BurnAccountID)
				}
				causedChange = append(causedChange, ghIssue)
			}
		} else { // issue is not for prioritization, freeze ballot if it exists and is open
			if govBallot, ok := ballots[k]; ok { // ballot for issue already exists, update it
				// if ballot closed, do nothing
				// if ballot frozen, do nothing
				// otherwise, freeze ballot
				if !govBallot.Closed && !govBallot.Frozen {
					ballotapi.Freeze_StageOnly(ctx, cloned, ghIssue.BallotName())
					causedChange = append(causedChange, ghIssue)
				}
			}
		}
	}
	causedChange.Sort()

	// don't touch ballots that have no corresponding issue

	return causedChange
}

func filterIssuesForPrioritization(ads []ballotproto.Advertisement) map[string]ballotproto.Advertisement {
	filtered := map[string]ballotproto.Advertisement{}
	for _, ad := range ads {
		idNS := ad.ID.ToNS()
		if len(idNS) == 3 &&
			idNS[0] == ImportedGithubPrefix &&
			(idNS[1] == ImportedIssuePrefix || idNS[1] == ImportedPullPrefix) {
			key := idNS[2]
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
	govBallot ballotproto.Advertisement,
) (changed bool) {
	if ghIssue.Title == govBallot.Title && ghIssue.Body == govBallot.Description {
		return false
	}
	ballotapi.Change_StageOnly(ctx, govCloned, ghIssue.BallotName(), ghIssue.Title, ghIssue.Body)
	return true
}

func UpdateFrozen_StageOnly(
	ctx context.Context,
	repo Repo,
	govAddr gov.OwnerAddress,
	govCloned gov.OwnerCloned,
	ghIssue ImportedIssue,
	govBallot ballotproto.Advertisement,
) (changed bool) {
	switch {
	case ghIssue.Locked && govBallot.Frozen:
		return false
	case ghIssue.Locked && !govBallot.Frozen:
		ballotapi.Freeze_StageOnly(ctx, govCloned, ghIssue.BallotName())
		return true
	case !ghIssue.Locked && govBallot.Frozen:
		ballotapi.Unfreeze_StageOnly(ctx, govCloned, ghIssue.BallotName())
		return true
	case !ghIssue.Locked && !govBallot.Frozen:
		return false
	}
	panic("unreachable")
}
