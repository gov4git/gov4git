package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v55/github"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/collab"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
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
		ghIssues, // XXX: may be too verbose
		nil,
	)
	return proto.CommitIfChanged(ctx, govCloned.Public, chg)
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
			if motion, ok := govMotions[id]; ok { // if motion for issue already exists, update it
				syncMeta(ctx, t, issue, id)
				switch {
				case issue.Closed && motion.Closed:
					// nothing to do
				case issue.Closed && !motion.Closed:
					syncFrozen(ctx, t, issue, motion)
					collab.CloseMotion_StageOnly(ctx, t, id)
				case !issue.Closed && motion.Closed:
					collab.ReopenMotion_StageOnly(ctx, t, id)
					syncFrozen(ctx, t, issue, motion)
				case !issue.Closed && !motion.Closed:
					syncFrozen(ctx, t, issue, motion)
				}
			} else { // otherwise, no motion for this issue exists, so create one
				syncCreateMotionForIssue(ctx, t, issue, id)
			}
		} else { // issue is not for prioritization, freeze motion if it exists and is open
			if motion, ok := govMotions[id]; ok { // motion for issue already exists, update it
				// if motion closed, do nothing
				// if motion frozen, do nothing
				// otherwise, freeze motion
				if !motion.Closed && !motion.Frozen {
					collab.FreezeMotion_StageOnly(ctx, t, id)
				}
			}
		}
	}

	// don't touch motions that have no corresponding issue

	matchingMotions := indexMotions(collab.ListMotions_Local(ctx, t))
	syncRefs(ctx, t, ghIssues, matchingMotions)

	return ghOrderedIssues
}

func syncRefs(ctx context.Context, t *git.Tree, issues map[string]ImportedIssue, motions map[collab.MotionID]collab.Motion) {

	panic("XXX")
	// XXX
}

func syncMeta(ctx context.Context, t *git.Tree, issue ImportedIssue, id collab.MotionID) {
	collab.UpdateMotionMeta_StageOnly(
		ctx,
		t,
		id,
		issue.URL,
		issue.Title,
		issue.Body,
		issue.Labels,
	)
}

func syncFrozen(
	ctx context.Context,
	t *git.Tree,
	ghIssue ImportedIssue,
	govMotion collab.Motion,
) {
	switch {
	case ghIssue.Locked && govMotion.Frozen:
		return
	case ghIssue.Locked && !govMotion.Frozen:
		collab.FreezeMotion_StageOnly(ctx, t, govMotion.ID)
		return
	case !ghIssue.Locked && govMotion.Frozen:
		collab.UnfreezeMotion_StageOnly(ctx, t, govMotion.ID)
		return
	case !ghIssue.Locked && !govMotion.Frozen:
		return
	}
	panic("unreachable")
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
	x := map[collab.MotionID]collab.Motion{}
	for _, m := range ms {
		x[m.ID] = m
	}
	return x
}
