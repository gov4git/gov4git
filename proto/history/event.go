package history

import "github.com/gov4git/gov4git/v2/proto"

var (
	historyNS = proto.RootNS.Append("history")
	history   = History{Root: historyNS}
)

type Event struct {
	Op      *Op           `json:"op"`
	Join    *JoinEvent    `json:"join"`
	Motion  *MotionEvent  `json:"motion"`
	Account *AccountEvent `json:"account"`
	Vote    *VoteEvent    `json:"vote"`
}

type Op struct {
	Op     string `json:"op"`
	Note   string `json:"note"`
	Args   M      `json:"args,omitempty"`
	Result M      `json:"result,omitempty"`
}

type M = map[string]any
