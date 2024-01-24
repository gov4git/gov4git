package concern

import "github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"

type CloseReport struct {
}

type CancelReport struct {
	PriorityPollOutcome ballotproto.Outcome `json:"priority_poll_outcome"`
}
