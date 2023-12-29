package metric

type MotionPolicy string

type MotionDecision string

type MotionEvent struct {
	Open   *MotionOpen   `json:"open"`
	Close  *MotionClose  `json:"close"`
	Cancel *MotionCancel `json:"cancel"`
}

type MotionID string

type MotionOpen struct {
	ID     MotionID     `json:"id"`
	Type   string       `json:"type"`
	Policy MotionPolicy `json:"policy"`
}

type MotionClose struct {
	ID       MotionID       `json:"id"`
	Type     string         `json:"type"`
	Policy   MotionPolicy   `json:"policy"`
	Decision MotionDecision `json:"decision"`
	Receipts Receipts       `json:"receipts"`
}

type MotionCancel struct {
	ID       MotionID     `json:"id"`
	Type     string       `json:"type"`
	Policy   MotionPolicy `json:"policy"`
	Receipts Receipts     `json:"receipts"`
}
