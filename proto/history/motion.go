package history

type MotionEvent struct {
	Open   *MotionOpen   `json:"open"`
	Close  *MotionClose  `json:"close"`
	Cancel *MotionCancel `json:"cancel"`
}

type MotionID string

type MotionOpen struct {
	ID   MotionID `json:"id"`
	Type string   `json:"type"`
}

type MotionClose struct {
	ID       MotionID `json:"id"`
	Type     string   `json:"type"`
	Receipts Receipts `json:"receipts"`
}

type MotionCancel struct {
	ID       MotionID `json:"id"`
	Type     string   `json:"type"`
	Receipts Receipts `json:"receipts"`
}
