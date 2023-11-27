package schema

import (
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/kv"
	"github.com/gov4git/lib4git/ns"
)

var (
	DocketNS = proto.RootNS.Append("docket")
	MotionNS = DocketNS.Append("motion")
	MotionKV = kv.KV[MotionID, Motion]{}
)

func MotionNoticesNS(id MotionID) ns.NS {
	return MotionKV.KeyNS(MotionNS, id).Append("notices.json")
}
