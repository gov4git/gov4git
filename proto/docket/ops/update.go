package ops

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func UpdateMotions(
	ctx context.Context,
	addr gov.OwnerAddress,

) git.Change[form.Map, form.None] {

	cloned := gov.CloneOwner(ctx, addr)
	chg := UpdateMotions_StageOnly(ctx, cloned)
	return proto.CommitIfChanged(ctx, cloned.Public, chg)
}

func UpdateMotions_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,

) git.Change[form.Map, form.None] {

	t := cloned.Public.Tree()
	motions := ListMotions_Local(ctx, t)
	for i, motion := range motions {
		// only update open motions
		if motion.Closed {
			continue
		}
		p := policy.GetMotionPolicy(ctx, motion)
		p.Update(
			ctx,
			cloned,
			motion,
			policy.MotionPolicyNS(motions[i].ID),
		)
	}

	motions.Sort()

	return git.NewChange(
		fmt.Sprintf("Update all %d motions", len(motions)),
		"docket_update_motions",
		form.Map{},
		form.None{},
		form.Forms{},
	)
}
