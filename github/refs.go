package github

import (
	"context"

	"github.com/gov4git/gov4git/proto/docket/ops"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
)

func syncRefs(
	ctx context.Context,
	cloned gov.OwnerCloned,
	chg *SyncManagedChanges,
	issues map[string]ImportedIssue,
	motions map[schema.MotionID]schema.Motion,
) {

	motionRefs := schema.RefSet{} // index of current refs between motions
	issueRefs := schema.RefSet{}  // index of current refs between issues, corresponding to existing motions
	ids := schema.MotionIDSet{}   // index of existing motions

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
				ref := schema.Ref{
					From: from,
					To:   to,
					Type: schema.RefType(importedRef.Type),
				}
				issueRefs[ref] = true
			}
		}
	}

	// update edge differences; only update open motions

	// add refs in issues, not in motions
	for issueRef := range issueRefs {
		if !motionRefs[issueRef] {
			ops.LinkMotions_StageOnly(ctx, cloned, issueRef.From, issueRef.To, issueRef.Type)
			chg.AddedRefs.Add(issueRef)
		}
	}

	// remove refs in motions, not in issues
	for motionRef := range motionRefs {
		if !issueRefs[motionRef] {
			ops.UnlinkMotions_StageOnly(ctx, cloned, motionRef.From, motionRef.To, motionRef.Type)
			chg.RemovedRefs.Add(motionRef)
		}
	}

}
