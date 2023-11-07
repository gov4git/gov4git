package github

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket"
	"github.com/gov4git/lib4git/git"
)

func syncRefs(
	ctx context.Context,
	t *git.Tree,
	chg *SyncChanges,
	issues map[string]ImportedIssue,
	motions map[docket.MotionID]docket.Motion,
) {

	motionRefs := docket.RefSet{} // index of current refs between motions
	issueRefs := docket.RefSet{}  // index of current refs between issues, corresponding to existing motions
	ids := docket.MotionIDSet{}   // index of existing motions

	// index motion refs (directed edges)
	for id, motion := range motions {
		for _, ref := range motion.RefTo {
			motionRefs[ref] = true
		}
		ids[id] = true
	}

	// index issue refs (directed edges)
	for _, issue := range issues {
		for _, importedRef := range issue.Refs {
			from := IssueNumberToMotionID(issue.Number)
			to := IssueNumberToMotionID(importedRef.To)
			// only include issue refs between existing motions
			if ids[from] && ids[to] {
				ref := docket.Ref{
					From: from,
					To:   to,
					Type: docket.RefType(importedRef.Type),
				}
				issueRefs[ref] = true
			}
		}
	}

	// update edge differences; only update open motions

	// add refs in issues, not in motions
	for issueRef := range issueRefs {
		if !motionRefs[issueRef] {
			docket.LinkMotions_StageOnly(ctx, t, issueRef.From, issueRef.To, issueRef.Type)
			chg.AddedRefs.Add(issueRef)
		}
	}

	// remove refs in motions, not in issues
	for motionRef := range motionRefs {
		if !issueRefs[motionRef] {
			docket.UnlinkMotions_StageOnly(ctx, t, motionRef.From, motionRef.To, motionRef.Type)
			chg.RemovedRefs.Add(motionRef)
		}
	}

}
