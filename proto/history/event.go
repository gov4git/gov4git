package history

import "github.com/gov4git/gov4git/v2/proto"

var (
	historyNS = proto.RootNS.Append("history")
	history   = History{Root: historyNS}
)

type Event struct {
	*Op `json:"op"`
}

type Op struct {
	Op     string `json:"op"`
	Note   string `json:"note"`
	Args   M      `json:"args,omitempty"`
	Result M      `json:"result,omitempty"`
}

type M = map[string]any
