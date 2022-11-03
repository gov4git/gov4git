package govproto

type GovDirPolicy struct {
	Change GovArbitration `json:"change"`
}

type GovArbitration struct {
	Quorum *GovQuorum `json:"quorum"`
}

type GovQuorum struct {
	Group     string `json:"group"`     // quorum participants
	Threshold uint32 `json:"threshold"` // minimum number of approvals required for quorum
}
