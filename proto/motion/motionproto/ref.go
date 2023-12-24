package motionproto

import (
	"encoding/json"
	"sort"
)

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

func (x Refs) Contains(y Ref) bool {
	for _, x := range x {
		if RefEqual(x, y) {
			return true
		}
	}
	return false
}

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

func (x Refs) RefSet() RefSet {
	rs := RefSet{}
	for _, ref := range x {
		rs[ref] = true
	}
	return rs
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

func (x RefSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Refs())
}

func (x *RefSet) UnmarshalJSON(b []byte) error {
	var refs Refs
	err := json.Unmarshal(b, &refs)
	if err != nil {
		return err
	}
	(*x) = RefSet{}
	for _, ref := range refs {
		(*x)[ref] = true
	}
	return nil
}
