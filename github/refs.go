package github

import (
	"context"

	"github.com/gov4git/gov4git/proto/collab"
	"github.com/gov4git/lib4git/git"
)

func syncRefs(
	ctx context.Context,
	t *git.Tree,
	issues map[string]ImportedIssue,
	motions map[collab.MotionID]collab.Motion,
) {

	motionRefs := map[collab.Ref]bool{} // index of current refs between motions
	issueRefs := map[collab.Ref]bool{}  // index of current refs between issues, corresponding to existing motions
	ids := map[collab.MotionID]bool{}   // index of existing motions

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
				ref := collab.Ref{
					From: from,
					To:   to,
					Type: collab.RefType(importedRef.Type),
				}
				issueRefs[ref] = true
			}
		}
	}

	// update edge differences; only update open motions

	// add refs in issues, not in motions
	for issueRef := range issueRefs {
		if !motionRefs[issueRef] {
			collab.LinkMotions_StageOnly(ctx, t, issueRef.From, issueRef.To, issueRef.Type)
		}
	}

	// remove refs in motions, not in issues
	for motionRef := range motionRefs {
		if !issueRefs[motionRef] {
			collab.UnlinkMotions_StageOnly(ctx, t, motionRef.From, motionRef.To, motionRef.Type)
		}
	}

}
