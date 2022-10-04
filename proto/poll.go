package proto

type GovPollAd struct {
	Path         string   `json:"path"`          // path within repo where poll will be persisted, also unique poll name
	Choices      []string `json:"choices"`       // ballot choices
	Group        string   `json:"group"`         // community group eligible to participate
	Strategy     string   `json:"strategy"`      // polling strategy
	Branch       string   `json:"branch"`        // branch governing the poll
	ParentCommit string   `json:"parent_commit"` // commit before poll
}

// XXX: polling strategies?
