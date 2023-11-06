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
	motionNS = collabNS.Append("motion")
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

type MotionIDs []MotionID

func (x MotionIDs) Len() int           { return len(x) }
func (x MotionIDs) Less(i, j int) bool { return x[i] < x[j] }
func (x MotionIDs) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x MotionIDs) Sort()              { sort.Sort(x) }

type MotionIDSet map[MotionID]bool

func (x MotionIDSet) Add(id MotionID) {
	x[id] = true
}

func (x MotionIDSet) MotionIDs() MotionIDs {
	s := make(MotionIDs, 0, len(x))
	for id := range x {
		s = append(s, id)
	}
	s.Sort()
	return s
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
	Body       string   `json:"description"`
	Labels     []string `json:"labels"`
	// state
	Frozen    bool `json:"frozen"`
	Closed    bool `json:"closed"`
	Cancelled bool `json:"cancelled"`
	// attention ranking
	Scoring Scoring `json:"scoring"`
	Score   float64 `json:"score"` // priority score for this motion, computed during sync after tallying
	// network
	RefBy Refs `json:"ref_by"`
	RefTo Refs `json:"ref_to"`
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

func (m *Motion) AddRefTo(ref Ref) {
	if !m.RefersTo(ref.To, ref.Type) {
		m.RefTo = append(m.RefTo, ref)
	}
	m.RefTo.Sort()
}

func (m *Motion) AddRefBy(ref Ref) {
	if !m.ReferredBy(ref.From, ref.Type) {
		m.RefBy = append(m.RefBy, ref)
	}
	m.RefBy.Sort()
}

func (m *Motion) RemoveRef(unref Ref) {
	m.RefTo = m.RefTo.Remove(unref)
	m.RefBy = m.RefBy.Remove(unref)
}

// Scoring describes how a concern or a proposal is assigned a priority score.
type Scoring struct {
	Fixed *float64           `json:"fixed"`
	Poll  *common.BallotName `json:"poll"`
}

type RefType string

type Ref struct {
	Type RefType  `json:"type"`
	From MotionID `json:"from"`
	To   MotionID `json:"to"`
}

func RefEqual(x, y Ref) bool {
	return x.Type == y.Type && x.From == y.From && x.To == y.To
}

func RefLess(p, q Ref) bool {
	if p.Type < q.Type {
		return true
	}
	if p.From < q.From {
		return true
	}
	if p.To < q.To {
		return true
	}
	return false
}

type Refs []Ref

func (x Refs) Len() int {
	return len(x)
}

func (x Refs) Less(i, j int) bool {
	return RefLess(x[i], x[j])
}

func (x Refs) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x Refs) Sort() { sort.Sort(x) }

func (x Refs) Remove(unref Ref) Refs {
	w := Refs{}
	for _, ref := range x {
		if !RefEqual(ref, unref) {
			w = append(w, ref)
		}
	}
	w.Sort()
	return w
}

type RefSet map[Ref]bool

func (x RefSet) Add(r Ref) {
	x[r] = true
}

func (x RefSet) Remove(r Ref) {
	delete(x, r)
}

func (x RefSet) Refs() Refs {
	s := make(Refs, 0, len(x))
	for r := range x {
		s = append(s, r)
	}
	s.Sort()
	return s
}

type Motions []Motion

func (x Motions) Sort() { sort.Sort(x) }

func (x Motions) Len() int { return len(x) }

func (x Motions) Less(i, j int) bool { return x[i].Score < x[j].Score }

func (x Motions) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
