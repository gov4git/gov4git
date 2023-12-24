package github

import (
	"context"
	"strings"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

func syncRefs(
	ctx context.Context,
	cloned gov.OwnerCloned,
	chg *SyncManagedChanges,
	issues map[string]ImportedIssue,
	motions map[motionproto.MotionID]motionproto.Motion,

) {

	motionRefs := motionproto.RefSet{} // index of current refs between motions
	issueRefs := motionproto.RefSet{}  // index of current refs between issues, corresponding to existing motions
	ids := motionproto.MotionIDSet{}   // index of existing motions

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
				ref := motionproto.Ref{
					From: from,
					To:   to,
					Type: motionproto.RefType(strings.ToLower(importedRef.Type)),
				}
				issueRefs[ref] = true
			}
		}
	}

	// update edge differences

	// add refs in issues, not in motions
	for issueRef := range issueRefs {
		if !motionRefs[issueRef] {
			motionapi.LinkMotions_StageOnly(ctx, cloned, issueRef.From, issueRef.To, issueRef.Type)
			chg.AddedRefs.Add(issueRef)
		}
	}

	// remove refs in motions, not in issues
	for motionRef := range motionRefs {
		if !issueRefs[motionRef] {
			motionapi.UnlinkMotions_StageOnly(ctx, cloned, motionRef.From, motionRef.To, motionRef.Type)
			chg.RemovedRefs.Add(motionRef)
		}
	}

}
