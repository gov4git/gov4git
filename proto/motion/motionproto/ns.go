package motionproto

import (
	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/kv"
	"github.com/gov4git/gov4git/v2/proto/motion"
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

var (
	// PoliciesNS is a namespace for holding individual policy class namespaces.
	PoliciesNS = proto.PolicyNS.Append("motion")

	PolicyStateFilebase = "state.json"
)

func PolicyNS(policyName motion.PolicyName) ns.NS {
	return PoliciesNS.Append(policyName.String())
}
