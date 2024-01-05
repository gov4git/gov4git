package motionproto

import (
	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/kv"
	"github.com/gov4git/lib4git/ns"
)

var (
	MotionNS = proto.RootNS.Append("motion")
	MotionKV = kv.KV[MotionID, Motion]{}
)

func MotionNoticesNS(id MotionID) ns.NS {
	return MotionKV.KeyNS(MotionNS, id).Append("notices.json")
}

func MotionAccountID(motionID MotionID) account.AccountID {
	return account.AccountIDFromLine(account.Pair("motion", motionID.String()))
}
