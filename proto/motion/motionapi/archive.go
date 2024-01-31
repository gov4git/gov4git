package motionapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

func ArchiveMotions(
	ctx context.Context,
	addr gov.OwnerAddress,
	args ...any,

) {

	cloned := gov.CloneOwner(ctx, addr)
	ArchiveMotions_StageOnly(ctx, cloned, args...)
	proto.Commitf(ctx, cloned.PublicClone(), "motion_archive", "Archive motions")
}

func ArchiveMotions_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	args ...any,

) {

	t := cloned.Public.Tree()
	motions := ListMotions_Local(ctx, t)
	for _, motion := range motions {
		if !motion.Archived && motion.Closed {
			continue
		}
		motion.Archived = true
		motionproto.MotionKV.Set(ctx, motionproto.MotionNS, t, motion.ID, motion)
	}
}
