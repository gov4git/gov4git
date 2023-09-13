package collab

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/kv"
	"github.com/gov4git/lib4git/must"
)

var (
	motionNS = collabNS.Sub("motion")
	motionKV = kv.KV[MotionID, Motion]{}
)

var (
	MotionPollBallotChoice = "prioritize"
)

func MotionPollBallotName(id MotionID) common.BallotName {
	return common.BallotName{"collab", "motion", "poll", id.String()}
}

type MotionID string

func (x MotionID) String() string {
	return string(x)
}

type MotionType string

const (
	MotionConcernType  MotionType = "concern"
	MotionProposalType MotionType = "proposal"
)

func ParseMotionType(ctx context.Context, s string) MotionType {
	switch s {
	case string(MotionConcernType):
		return MotionConcernType
	case string(MotionProposalType):
		return MotionProposalType
	}
	must.Panic(ctx, fmt.Errorf("unknown motion type"))
	return MotionType("")
}

type Motion struct {
	OpenedAt time.Time `json:"opened_at"`
	ClosedAt time.Time `json:"closed_at"`
	//
	ID   MotionID   `json:"id"`
	Type MotionType `json:"type"`
	// meta
	TrackerURL string   `json:"tracker_url"` // link to concern on an external concern tracker, such as a GitHub issue
	Title      string   `json:"title"`
	Desc       string   `json:"description"`
	Labels     []string `json:"labels"`
	// state
	Frozen    bool `json:"frozen"`
	Closed    bool `json:"closed"`
	Cancelled bool `json:"cancelled"`
	// attention ranking
	Scoring Scoring `json:"scoring"`
	Score   float64 `json:"score"` // priority score for this motion, computed during sync after tallying
	// network
	RefBy []*Ref `json:"ref_by"`
	RefTo []*Ref `json:"ref_to"`
}

func (m Motion) IsConcern() bool {
	return m.Type == MotionConcernType
}

func (m Motion) IsProposal() bool {
	return m.Type == MotionProposalType
}

func (m Motion) RefersTo(toID MotionID, typ RefType) bool {
	for _, ref := range m.RefTo {
		if ref.To == toID && ref.Type == typ {
			return true
		}
	}
	return false
}

func (m Motion) ReferredBy(fromID MotionID, typ RefType) bool {
	for _, ref := range m.RefBy {
		if ref.From == fromID && ref.Type == typ {
			return true
		}
	}
	return false
}

// Scoring describes how a concern or a proposal is assigned a priority score.
type Scoring struct {
	Fixed *float64           `json:"fixed"`
	Poll  *common.BallotName `json:"poll"`
}

type RefType string

type Ref struct {
	From MotionID `json:"from"`
	To   MotionID `json:"to"`
	Type RefType  `json:"type"`
}

type Motions []Motion

func (x Motions) Sort()              { sort.Sort(x) }
func (x Motions) Len() int           { return len(x) }
func (x Motions) Less(i, j int) bool { return x[i].Score < x[j].Score }
func (x Motions) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
