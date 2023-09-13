package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/collab"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func SyncGovernedIssues(
	ctx context.Context,
	repo Repo,
	githubClient *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OrganizerAddress,
) git.Change[form.Map, ImportedIssues] {

	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	ghIssues := SyncGovernedIssues_Local(ctx, repo, githubClient, govAddr, govCloned)
	chg := git.NewChange[form.Map, ImportedIssues](
		fmt.Sprintf("Sync %d GitHub issues participating in governance", len(ghIssues)),
		"github_sync",
		form.Map{},
		ghIssues,
		nil,
	)
	status, err := govCloned.Public.Tree().Status()
	must.NoError(ctx, err)
	if !status.IsClean() {
		proto.Commit(ctx, govCloned.Public.Tree(), chg)
		govCloned.Public.Push(ctx)
	}
	return chg
}

func SyncGovernedIssues_Local(
	ctx context.Context,
	repo Repo,
	githubClient *github.Client, // if nil, a new client for repo will be created
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
) ImportedIssues {

	t := govCloned.Public.Tree()

	// load github issues and governance motions, and
	// index them under a common key space
	ghOrderedIssues, ghIssues := LoadIssues(ctx, repo, githubClient)
	govMotions := indexMotions(collab.ListMotions_Local(ctx, t))

	// ensure every issue has a corresponding up-to-date motion
	for key, issue := range ghIssues {
		id := collab.MotionID(key)
		if issue.IsGoverned {
			if govMotion, ok := govMotions[id]; ok { // if motion for issue already exists, update it
				syncUpdateMeta(ctx, t, issue, id)
				switch {
				case issue.Closed && govMotion.Closed:
					XXX
					// nothing to do
				case issue.Closed && !govMotion.Closed:
					UpdateFrozen_StageOnly(ctx, repo, govAddr, govCloned, issue, govMotion)
					ballot.Close_StageOnly(ctx, govAddr, govCloned, issue.BallotName(), false)
				case !issue.Closed && govMotion.Closed:
					ballot.Reopen_StageOnly(ctx, govAddr, govCloned, issue.BallotName())
					UpdateFrozen_StageOnly(ctx, repo, govAddr, govCloned, issue, govMotion)
				case !issue.Closed && !govMotion.Closed:
					UpdateFrozen_StageOnly(ctx, repo, govAddr, govCloned, issue, govMotion)
				}
			} else { // otherwise, no motion for this issue exists, so create one
				syncCreateMotionForIssue(ctx, t, issue, id)
			}
		} else { // issue is not for prioritization, freeze motion if it exists and is open
			if govMotion, ok := govMotions[id]; ok { // motion for issue already exists, update it
				// if motion closed, do nothing
				// if motion frozen, do nothing
				// otherwise, freeze motion
				if !govMotion.Closed && !govMotion.Frozen {
					collab.FreezeMotion_StageOnly(ctx, t, id)
				}
			}
		}
	}

	// don't touch motions that have no corresponding issue

	// XXX: sync references

	return ghOrderedIssues
}

func syncUpdateMeta(ctx context.Context, t *git.Tree, issue ImportedIssue, id collab.MotionID) {
	XXX
}

func syncCreateMotionForIssue(ctx context.Context, t *git.Tree, issue ImportedIssue, id collab.MotionID) {
	collab.OpenMotion_StageOnly(
		ctx,
		t,
		id,
		issue.Title,
		issue.Body,
		issue.MotionType(),
		issue.URL,
		issue.Labels,
	)
	if issue.Locked {
		collab.FreezeMotion_StageOnly(ctx, t, id)
	}
	if issue.Closed {
		collab.CloseMotion_StageOnly(ctx, t, id)
	}
	// XXX: initialize poll scoring, if not closed or frozen
}

func indexMotions(ms collab.Motions) map[collab.MotionID]collab.Motion {
	x := map[string]collab.Motion{}
	for _, m := range ms {
		x[m.ID.String()] = m
	}
	return x
}
