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
) git.Change[form.Map, SyncChanges] {

	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	syncChanges := SyncGovernedIssues_StageOnly(ctx, repo, githubClient, govAddr, govCloned)
	chg := git.NewChange[form.Map, SyncChanges](
		fmt.Sprintf("Sync %d GitHub issues participating in governance", len(syncChanges.IssuesCausingChange)),
		"github_sync",
		form.Map{},
		syncChanges,
		nil,
	)
	return proto.CommitIfChanged(ctx, govCloned.Public, chg)
}

type SyncChanges struct {
	IssuesCausingChange ImportedIssues     `json:"issues_causing_change"`
	Updated             collab.MotionIDSet `json:"updated_motions"`
	Opened              collab.MotionIDSet `json:"opened_motions"`
	Closed              collab.MotionIDSet `json:"closed_motions"`
	Froze               collab.MotionIDSet `json:"froze_motions"`
	Unfroze             collab.MotionIDSet `json:"unfroze_motions"`
	AddedRefs           collab.RefSet      `json:"added_refs"`
	RemovedRefs         collab.RefSet      `json:"removed_refs"`
}

func SyncGovernedIssues_StageOnly(
	ctx context.Context,
	repo Repo,
	githubClient *github.Client,
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
) (syncChanges SyncChanges) {

	syncChanges.AddedRefs = collab.RefSet{}
	syncChanges.RemovedRefs = collab.RefSet{}

	t := govCloned.Public.Tree()

	// load github issues and governance motions, and
	// index them under a common key space
	_, issues := LoadIssues(ctx, repo, githubClient)
	motions := indexMotions(collab.ListMotions_Local(ctx, t))

	// ensure every issue has a corresponding up-to-date motion
	for key, issue := range issues {
		id := collab.MotionID(key)
		if issue.IsGoverned {
			if motion, ok := motions[id]; ok { // if motion for issue already exists, update it
				changed := syncMeta(ctx, t, &syncChanges, issue, motion)
				switch {
				case issue.Closed && motion.Closed:
				case issue.Closed && !motion.Closed:
					syncFrozen(ctx, t, &syncChanges, issue, motion)
					collab.CloseMotion_StageOnly(ctx, t, id)
					syncChanges.Closed.Add(id)
					changed = true
				case !issue.Closed && motion.Closed:
					//XXX: reopening is prohibited
					collab.ReopenMotion_StageOnly(ctx, t, id)
					syncFrozen(ctx, t, &syncChanges, issue, motion)
					changed = true
				case !issue.Closed && !motion.Closed:
					changed = changed || syncFrozen(ctx, t, &syncChanges, issue, motion)
				}
				if changed {
					syncChanges.IssuesCausingChange = append(syncChanges.IssuesCausingChange, issue)
				}
			} else { // otherwise, no motion for this issue exists, so create one
				syncCreateMotionForIssue(ctx, t, &syncChanges, issue, id)
				syncChanges.IssuesCausingChange = append(syncChanges.IssuesCausingChange, issue)
			}
		} else { // issue is not governed, freeze motion if it exists and is open
			if motion, ok := motions[id]; ok { // motion for issue already exists, update it
				// if motion closed, do nothing
				// if motion frozen, do nothing
				// otherwise, freeze motion
				if !motion.Closed && !motion.Frozen {
					collab.FreezeMotion_StageOnly(ctx, t, id)
					syncChanges.Froze.Add(id)
					syncChanges.IssuesCausingChange = append(syncChanges.IssuesCausingChange, issue)
				}
			}
		}
	}

	// don't touch motions that have no corresponding issue

	matchingMotions := indexMotions(collab.ListMotions_Local(ctx, t))
	syncRefs(ctx, t, &syncChanges, issues, matchingMotions)

	syncChanges.IssuesCausingChange.Sort()
	return
}

func syncMeta(
	ctx context.Context,
	t *git.Tree,
	chg *SyncChanges,
	issue ImportedIssue,
	motion collab.Motion,
) bool {
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
	chg.Updated.Add(motion.ID)
	return true
}

func syncFrozen(
	ctx context.Context,
	t *git.Tree,
	chg *SyncChanges,
	ghIssue ImportedIssue,
	govMotion collab.Motion,
) bool {
	switch {
	case ghIssue.Locked && govMotion.Frozen:
		return false
	case ghIssue.Locked && !govMotion.Frozen:
		collab.FreezeMotion_StageOnly(ctx, t, govMotion.ID)
		chg.Froze.Add(govMotion.ID)
		return true
	case !ghIssue.Locked && govMotion.Frozen:
		collab.UnfreezeMotion_StageOnly(ctx, t, govMotion.ID)
		chg.Unfroze.Add(govMotion.ID)
		return true
	case !ghIssue.Locked && !govMotion.Frozen:
		return false
	}
	panic("unreachable")
}

func syncCreateMotionForIssue(
	ctx context.Context,
	t *git.Tree,
	chg *SyncChanges,
	issue ImportedIssue,
	id collab.MotionID,
) {
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
	chg.Opened.Add(id)
	if issue.Locked {
		collab.FreezeMotion_StageOnly(ctx, t, id)
		chg.Froze.Add(id)
	}
	if issue.Closed {
		collab.CloseMotion_StageOnly(ctx, t, id)
		chg.Closed.Add(id)
	}
}

func indexMotions(ms collab.Motions) map[collab.MotionID]collab.Motion {
	x := map[collab.MotionID]collab.Motion{}
	for _, m := range ms {
		x[m.ID] = m
	}
	return x
}
