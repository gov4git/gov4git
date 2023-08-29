package github

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/go-github/v54/github"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/util"
)

func FetchIssues(ctx context.Context, repo GithubRepo) []*github.Issue {

	c := GetGithubClient(ctx, repo)

	opt := &github.IssueListByRepoOptions{State: "all"}
	var allIssues []*github.Issue
	for {
		issues, resp, err := c.Issues.ListByRepo(ctx, repo.Owner, repo.Name, opt)
		must.NoError(ctx, err)
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allIssues
}

func labelsToStrings(labels []*github.Label) []string {
	var labelStrings []string
	for _, label := range labels {
		labelStrings = append(labelStrings, label.GetName())
	}
	return labelStrings
}

func IsIssueForPrioritization(issue *github.Issue) bool {
	return util.IsIn(PrioritizeIssueByGovernanceLabel, labelsToStrings(issue.Labels)...)
}

func TransformIssue(ctx context.Context, issue *github.Issue) GithubBallotIssue {
	return GithubBallotIssue{
		URL:       issue.GetURL(),
		Number:    int64(issue.GetNumber()),
		Title:     issue.GetTitle(),
		Body:      issue.GetBody(),
		Labels:    labelsToStrings(issue.Labels),
		ClosedAt:  unwrapTimestamp(issue.ClosedAt),
		CreatedAt: unwrapTimestamp(issue.CreatedAt),
		UpdatedAt: unwrapTimestamp(issue.UpdatedAt),
		Locked:    issue.GetLocked(),
		Closed:    issue.GetState() == "closed",
	}
}

func unwrapTimestamp(ts *github.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	return &ts.Time
}

func LoadIssuesForPrioritization(ctx context.Context, repo GithubRepo) (GithubBallotIssues, map[string]GithubBallotIssue) {

	issues := FetchIssues(ctx, repo)
	key := map[string]GithubBallotIssue{}
	order := GithubBallotIssues{}
	for _, issue := range issues {
		if IsIssueForPrioritization(issue) {
			ghIssue := TransformIssue(ctx, issue)
			key[ghIssue.Key()] = ghIssue
			order = append(order, ghIssue)
		}
	}
	order.Sort()
	return order, key
}

func ImportIssuesForPrioritization(
	ctx context.Context,
	repo GithubRepo,
	govAddr gov.OrganizerAddress,
) git.Change[form.Map, GithubBallotIssues] {

	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	ghIssues := ImportIssuesForPrioritization_Local(ctx, repo, govAddr, govCloned)
	chg := git.NewChange[form.Map, GithubBallotIssues](
		fmt.Sprintf("Import %d GitHub issues", len(ghIssues)),
		"github_import",
		form.Map{},
		ghIssues,
		nil,
	)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func ImportIssuesForPrioritization_Local(
	ctx context.Context,
	repo GithubRepo,
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
) GithubBallotIssues {

	// load github issues and governance ballots, and
	// index them under a common key space
	ghOrderedIssues, ghIssues := LoadIssuesForPrioritization(ctx, repo)
	govBallots := filterIssuesForPrioritization(ballot.ListLocal(ctx, govCloned.Public.Tree()))

	// ensure every issue has a corresponding up-to-date ballot
	for k, ghIssue := range ghIssues {
		if govBallot, ok := govBallots[k]; ok { // ballot for issue already exists, update it

			switch {
			case ghIssue.Closed && govBallot.Closed:
				// nothing to do
			case ghIssue.Closed && !govBallot.Closed:
				UpdateMeta_StageOnly(ctx, repo, govAddr, govCloned, ghIssue, govBallot)
				ballot.CloseStageOnly(ctx, govAddr, govCloned, ghIssue.BallotName(), false)
			case !ghIssue.Closed && govBallot.Closed:
				UpdateMeta_StageOnly(ctx, repo, govAddr, govCloned, ghIssue, govBallot)
				UpdateFrozen_StageOnly(ctx, repo, govAddr, govCloned, ghIssue, govBallot)
				XXX // reopen ballot
			case !ghIssue.Closed && !govBallot.Closed:
				UpdateMeta_StageOnly(ctx, repo, govAddr, govCloned, ghIssue, govBallot)
				UpdateFrozen_StageOnly(ctx, repo, govAddr, govCloned, ghIssue, govBallot)
			}

		} else { // no ballot for this issue, create it
			ballot.OpenStageOnly(
				ctx,
				qv.QV{},
				gov.GovAddress(govAddr.Public),
				govCloned.Public,
				ghIssue.BallotName(),
				ghIssue.Title,
				ghIssue.Body,
				[]string{PrioritizeBallotChoice},
				member.Everybody,
			)
		}
	}

	// don't touch ballots that have no corresponding issue

	return ghOrderedIssues
}

func filterIssuesForPrioritization(ads []common.Advertisement) map[string]common.Advertisement {
	filtered := map[string]common.Advertisement{}
	for _, ad := range ads {
		if len(ad.Name) == 2 && ad.Name[0] == "issue" {
			key := ad.Name[1]
			if _, err := strconv.Atoi(key); err == nil {
				filtered[key] = ad
			}
		}
	}
	return filtered
}

func UpdateMeta_StageOnly(
	ctx context.Context,
	repo GithubRepo,
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	ghIssue GithubBallotIssue,
	govBallot common.Advertisement,
) {
	if ghIssue.Title == govBallot.Title && ghIssue.Body == govBallot.Description {
		return
	}
	ballot.ChangeStageOnly(ctx, govAddr, govCloned, ghIssue.BallotName(), ghIssue.Title, ghIssue.Body)
}

func UpdateFrozen_StageOnly(
	ctx context.Context,
	repo GithubRepo,
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	ghIssue GithubBallotIssue,
	govBallot common.Advertisement,
) {
	switch {
	case ghIssue.Locked && govBallot.Frozen:
	case ghIssue.Locked && !govBallot.Frozen:
		ballot.FreezeStageOnly(ctx, govAddr, govCloned, ghIssue.BallotName())
	case !ghIssue.Locked && govBallot.Frozen:
		ballot.UnfreezeStageOnly(ctx, govAddr, govCloned, ghIssue.BallotName())
	case !ghIssue.Locked && !govBallot.Frozen:
	}
}
