package github

import (
	"context"
	"fmt"
	"strings"

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
		fmt.Println("REFXXX ISSUE/PR", issue.Number)
		for _, importedRef := range issue.Refs {
			from := IssueNumberToMotionID(issue.Number)
			to := IssueNumberToMotionID(importedRef.To)
			fmt.Println("REF--- FROM", from, "TO", to)
			// only include issue refs between existing motions
			if ids[from] && ids[to] {
				fmt.Println("REF--- adding to issueRefs")
				ref := schema.Ref{
					From: from,
					To:   to,
					Type: schema.RefType(strings.ToLower(importedRef.Type)),
				}
				issueRefs[ref] = true
			}
		}
	}

	// update edge differences

	// add refs in issues, not in motions
	fmt.Println("REF ALL motionRefs", motionRefs)
	for issueRef := range issueRefs {
		fmt.Println("REF--- considering issueRef", issueRef)
		if !motionRefs[issueRef] {
			fmt.Println("REF--- merging issueRef", issueRef)
			ops.LinkMotions_StageOnly(ctx, cloned, issueRef.From, issueRef.To, issueRef.Type)
			chg.AddedRefs.Add(issueRef)
		}
	}

	// remove refs in motions, not in issues
	for motionRef := range motionRefs {
		fmt.Println("REF--- considering motionRef", motionRef)
		if !issueRefs[motionRef] {
			fmt.Println("REF--- removing motionRef", motionRef)
			ops.UnlinkMotions_StageOnly(ctx, cloned, motionRef.From, motionRef.To, motionRef.Type)
			chg.RemovedRefs.Add(motionRef)
		}
	}

}
