package schema

import (
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/account"
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

func MotionAccountID(motionID MotionID) account.AccountID {
	return account.AccountIDFromLine(account.Pair("motion", motionID.String()))
}

func MotionOwnerID(motionID MotionID) account.OwnerID {
	return account.OwnerIDFromLine(account.Pair("motion", motionID.String()))
}
