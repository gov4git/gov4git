package trace

import (
	"github.com/gov4git/gov4git/v2/proto/history"
)

var (
	traceHistoryNS = history.HistoryNS.Append("trace")
	traceHistory   = History{Root: traceHistoryNS}
)

type Event struct {
	Op     string `json:"op"`
	Note   string `json:"note"`
	Args   M      `json:"args,omitempty"`
	Result M      `json:"result,omitempty"`
}

type M = map[string]any
