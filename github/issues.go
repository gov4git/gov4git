package github

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
	"github.com/gov4git/lib4git/util"
)

func FetchIssues(ctx context.Context, repo GithubRepo, ghc *github.Client) []*github.Issue {
	opt := &github.IssueListByRepoOptions{State: "all"}
	var allIssues []*github.Issue
	for {
		issues, resp, err := ghc.Issues.ListByRepo(ctx, repo.Owner, repo.Name, opt)
		must.NoError(ctx, err)
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allIssues
}

func LabelsToStrings(labels []*github.Label) []string {
	var labelStrings []string
	for _, label := range labels {
		labelStrings = append(labelStrings, label.GetName())
	}
	return labelStrings
}

func IsIssueForPrioritization(issue *github.Issue) bool {
	return util.IsIn(PrioritizeIssueByGovernanceLabel, LabelsToStrings(issue.Labels)...)
}

func TransformIssue(ctx context.Context, issue *github.Issue) GithubIssueBallot {
	return GithubIssueBallot{
		ForPrioritization: IsIssueForPrioritization(issue),
		URL:               issue.GetURL(),
		Number:            int64(issue.GetNumber()),
		Title:             issue.GetTitle(),
		Body:              issue.GetBody(),
		Labels:            LabelsToStrings(issue.Labels),
		ClosedAt:          unwrapTimestamp(issue.ClosedAt),
		CreatedAt:         unwrapTimestamp(issue.CreatedAt),
		UpdatedAt:         unwrapTimestamp(issue.UpdatedAt),
		Locked:            issue.GetLocked(),
		Closed:            issue.GetState() == "closed",
		IsPullRequest:     issue.GetPullRequestLinks() != nil,
	}
}

func unwrapTimestamp(ts *github.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	return &ts.Time
}

func LoadIssuesForPrioritization(
	ctx context.Context,
	repo GithubRepo,
	githubClient *github.Client, // if nil, a new client for repo will be created
) (GithubIssueBallots, map[string]GithubIssueBallot) {

	issues := FetchIssues(ctx, repo, githubClient)
	key := map[string]GithubIssueBallot{}
	order := GithubIssueBallots{}
	for _, issue := range issues {
		ghIssue := TransformIssue(ctx, issue)
		key[ghIssue.Key()] = ghIssue
		order = append(order, ghIssue)
	}
	order.Sort()
	return order, key
}

func ImportIssuesForPrioritization(
	ctx context.Context,
	repo GithubRepo,
	githubClient *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OrganizerAddress,
) git.Change[form.Map, GithubIssueBallots] {

	base.Infof("importing issues for prioritization ...")
	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	ghIssues := ImportIssuesForPrioritization_StageOnly(ctx, repo, githubClient, govAddr, govCloned)
	chg := git.NewChange[form.Map, GithubIssueBallots](
		fmt.Sprintf("Import %d GitHub issues for prioritization", len(ghIssues)),
		"github_import_for_prioritization",
		form.Map{},
		ghIssues,
		nil,
	)
	status, err := govCloned.Public.Tree().Status()
	must.NoError(ctx, err)
	if !status.IsClean() {
		base.Infof("import from github caused changes")
		proto.Commit(ctx, govCloned.Public.Tree(), chg)
		govCloned.Public.Push(ctx)
	}
	return chg
}

func ImportIssuesForPrioritization_StageOnly(
	ctx context.Context,
	repo GithubRepo,
	githubClient *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
) GithubIssueBallots {

	// load github issues and governance ballots, and
	// index them under a common key space
	ghOrderedIssues, ghIssues := LoadIssuesForPrioritization(ctx, repo, githubClient)
	govBallots := filterIssuesForPrioritization(ballot.ListLocal(ctx, govCloned.Public.Tree()))

	// ensure every issue has a corresponding up-to-date ballot
	for k, ghIssue := range ghIssues {
		if ghIssue.ForPrioritization {
			if govBallot, ok := govBallots[k]; ok { // ballot for issue already exists, update it

				must.Assertf(ctx, ns.Equal(ghIssue.BallotName(), govBallot.Name),
					"issue ballot name %v and actual ballot name %v mismatch", ghIssue.BallotName(), govBallot.Name)

				switch {
				case ghIssue.Closed && govBallot.Closed:
					// nothing to do
				case ghIssue.Closed && !govBallot.Closed:
					UpdateMeta_StageOnly(ctx, repo, govAddr, govCloned, ghIssue, govBallot)
					UpdateFrozen_StageOnly(ctx, repo, govAddr, govCloned, ghIssue, govBallot)
					ballot.CloseStageOnly(ctx, govAddr, govCloned, ghIssue.BallotName(), false)
				case !ghIssue.Closed && govBallot.Closed:
					UpdateMeta_StageOnly(ctx, repo, govAddr, govCloned, ghIssue, govBallot)
					ballot.Reopen_StageOnly(ctx, govAddr, govCloned, ghIssue.BallotName())
					UpdateFrozen_StageOnly(ctx, repo, govAddr, govCloned, ghIssue, govBallot)
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
				if ghIssue.Locked {
					ballot.FreezeStageOnly(ctx, govAddr, govCloned, ghIssue.BallotName())
				}
				if ghIssue.Closed {
					ballot.CloseStageOnly(ctx, govAddr, govCloned, ghIssue.BallotName(), false)
				}
			}
		} else { // issue is not for prioritization, freeze ballot if it exists and is open
			if govBallot, ok := govBallots[k]; ok { // ballot for issue already exists, update it
				// if ballot closed, do nothing
				// if ballot frozen, do nothing
				// otherwise, freeze ballot
				if !govBallot.Closed && !govBallot.Frozen {
					ballot.FreezeStageOnly(ctx, govAddr, govCloned, ghIssue.BallotName())
				}
			}
		}
	}

	// don't touch ballots that have no corresponding issue

	return ghOrderedIssues
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
	repo GithubRepo,
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	ghIssue GithubIssueBallot,
	govBallot common.Advertisement,
) (changed bool) {
	if ghIssue.Title == govBallot.Title && ghIssue.Body == govBallot.Description {
		return false
	}
	ballot.Change_StageOnly(ctx, govAddr, govCloned, ghIssue.BallotName(), ghIssue.Title, ghIssue.Body)
	return true
}

func UpdateFrozen_StageOnly(
	ctx context.Context,
	repo GithubRepo,
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	ghIssue GithubIssueBallot,
	govBallot common.Advertisement,
) (changed bool) {
	switch {
	case ghIssue.Locked && govBallot.Frozen:
		return false
	case ghIssue.Locked && !govBallot.Frozen:
		ballot.FreezeStageOnly(ctx, govAddr, govCloned, ghIssue.BallotName())
		return true
	case !ghIssue.Locked && govBallot.Frozen:
		ballot.UnfreezeStageOnly(ctx, govAddr, govCloned, ghIssue.BallotName())
		return true
	case !ghIssue.Locked && !govBallot.Frozen:
		return false
	}
	panic("unreachable")
}
