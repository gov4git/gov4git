package concern

import "github.com/gov4git/gov4git/v2/proto/ballot/common"

type CloseReport struct {
}

type CancelReport struct {
	PriorityPollOutcome common.Outcome `json:"priority_poll_outcome"`
}
