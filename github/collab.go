package github

import (
	"context"
	"fmt"
	"slices"

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
	githubClient *github.Client,
	govAddr gov.OrganizerAddress,
) git.Change[form.Map, ImportedIssues] {

	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	issuesCausingChange := SyncGovernedIssues_StageOnly(ctx, repo, githubClient, govAddr, govCloned)
	chg := git.NewChange[form.Map, ImportedIssues](
		fmt.Sprintf("Sync %d GitHub issues participating in governance", len(issuesCausingChange)),
		"github_sync",
		form.Map{},
		issuesCausingChange,
		nil,
	)
	return proto.CommitIfChanged(ctx, govCloned.Public, chg)
}

func SyncGovernedIssues_StageOnly(
	ctx context.Context,
	repo Repo,
	githubClient *github.Client,
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
) ImportedIssues {

	t := govCloned.Public.Tree()

	// load github issues and governance motions, and
	// index them under a common key space
	_, issues := LoadIssues(ctx, repo, githubClient)
	motions := indexMotions(collab.ListMotions_Local(ctx, t))

	// ensure every issue has a corresponding up-to-date motion
	causedChange := ImportedIssues{}
	for key, issue := range issues {
		id := collab.MotionID(key)
		if issue.IsGoverned {
			if motion, ok := motions[id]; ok { // if motion for issue already exists, update it
				changed := syncMeta(ctx, t, issue, motion)
				switch {
				case issue.Closed && motion.Closed:
					// nothing to do
				case issue.Closed && !motion.Closed:
					syncFrozen(ctx, t, issue, motion)
					collab.CloseMotion_StageOnly(ctx, t, id)
					changed = true
				case !issue.Closed && motion.Closed:
					collab.ReopenMotion_StageOnly(ctx, t, id)
					syncFrozen(ctx, t, issue, motion)
					changed = true
				case !issue.Closed && !motion.Closed:
					changed = changed || syncFrozen(ctx, t, issue, motion)
				}
				if changed {
					causedChange = append(causedChange, issue)
				}
			} else { // otherwise, no motion for this issue exists, so create one
				syncCreateMotionForIssue(ctx, t, issue, id)
				causedChange = append(causedChange, issue)
			}
		} else { // issue is not governed, freeze motion if it exists and is open
			if motion, ok := motions[id]; ok { // motion for issue already exists, update it
				// if motion closed, do nothing
				// if motion frozen, do nothing
				// otherwise, freeze motion
				if !motion.Closed && !motion.Frozen {
					collab.FreezeMotion_StageOnly(ctx, t, id)
					causedChange = append(causedChange, issue)
				}
			}
		}
	}

	// don't touch motions that have no corresponding issue

	matchingMotions := indexMotions(collab.ListMotions_Local(ctx, t))
	syncRefs(ctx, t, issues, matchingMotions)

	causedChange.Sort()
	return causedChange
}

func syncMeta(ctx context.Context, t *git.Tree, issue ImportedIssue, motion collab.Motion) bool {
	if motion.TrackerURL == issue.URL &&
		motion.Title == issue.Title &&
		motion.Body == issue.Body &&
		slices.Equal(motion.Labels, issue.Labels) {
		return false
	}
	collab.UpdateMotionMeta_StageOnly(
		ctx,
		t,
		motion.ID,
		issue.URL,
		issue.Title,
		issue.Body,
		issue.Labels,
	)
	return true
}

func syncFrozen(
	ctx context.Context,
	t *git.Tree,
	ghIssue ImportedIssue,
	govMotion collab.Motion,
) bool {
	switch {
	case ghIssue.Locked && govMotion.Frozen:
		return false
	case ghIssue.Locked && !govMotion.Frozen:
		collab.FreezeMotion_StageOnly(ctx, t, govMotion.ID)
		return true
	case !ghIssue.Locked && govMotion.Frozen:
		collab.UnfreezeMotion_StageOnly(ctx, t, govMotion.ID)
		return true
	case !ghIssue.Locked && !govMotion.Frozen:
		return false
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
