package schema

import "github.com/gov4git/gov4git/proto/ballot/common"

// Scoring describes how a concern or a proposal is assigned a priority score.
type Scoring struct {
	Fixed *float64           `json:"fixed"`
	Poll  *common.BallotName `json:"poll"`
}
