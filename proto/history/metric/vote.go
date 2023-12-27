package metric

type VoteEvent struct {
	By           User         `json:"by"`
	Purpose      VotePurpose  `json:"context"`
	MotionPolicy MotionPolicy `json:"motion_policy"`
	BallotPolicy BallotPolicy `json:"ballot_policy"`
	Receipts     Receipts     `json:"receipts"`
}

type BallotPolicy string

type VotePurpose string

const (
	VotePurposeUnspecified VotePurpose = "unspecified"
	VotePurposeConcern     VotePurpose = "concern"
	VotePurposeProposal    VotePurpose = "proposal"
)
