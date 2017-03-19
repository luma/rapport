package causality

import "sync"

// CreateVersionVectorValue returns a new, empty VersionVector
func CreateVersionVectorValue() *VersionVectorValue {
	return &VersionVectorValue{Dots: make(map[string]*Dot)}
}

// Dots is a lamport time for a set of actors
type Dots map[string]LamportTime

// VersionVector is a set of replica ids (string) to newest time.
// VersionVectors can be used to combine events from multiple replicas
//
// Think ( {"ReplicaA", 2}, {"ReplicaA", 3} )
//
type VersionVector struct {
	dots Dots
	l    sync.RWMutex
}

// CreateVersionVector returns a new, empty VersionVector
func CreateVersionVector() *VersionVector {
	return &VersionVector{dots: make(Dots)}
}

// UnmarshalVersionVector is a helper method to create a new version vector
// and umarshal binary data into it
func UnmarshalVersionVector(data []byte) (*VersionVector, error) {
	version := CreateVersionVector()
	if err := version.Unmarshal(data); err != nil {
		return nil, err
	}

	return version, nil
}

// Merge another clock into this one, without regard to dominance.
func (v *VersionVector) Merge(other *VersionVector) {
	other.REach(func(actor string, t LamportTime) {
		v.Witness(actor, t)
	})
}

// Witness stores a LamportTime for a particular actor if it
// dominates any existing times for that actor.
//
func (v *VersionVector) Witness(actor string, otherTime LamportTime) bool {
	v.l.Lock()
	shouldWitness := !v.descendentOf(actor, otherTime)

	if shouldWitness {
		v.dots[actor] = otherTime
	}
	v.l.Unlock()

	return shouldWitness
}

// Incr increments the LamportTime for a particular actor and
// returns the new time.
//
func (v *VersionVector) Incr(actor string) LamportTime {
	v.l.Lock()

	t, exists := v.dots[actor]
	if !exists {
		t = LamportTime(1)
		v.dots[actor] = t
	} else {
		t++
		v.dots[actor] = t
	}

	v.l.Unlock()
	return t
}

// IsConcurrentWith indicates whether this vector clock
// is totally divergent to the other.
//
func (v *VersionVector) IsConcurrentWith(other *VersionVector) bool {
	return v.Compare(other) == OrderNone
}

// Get retrieves the LamportTime for a particular actor.
func (v *VersionVector) Get(actor string) (LamportTime, bool) {
	v.l.RLock()
	time, exists := v.dots[actor]
	v.l.RUnlock()
	return time, exists
}

// IsEmpty returns true if the version vector is empty
func (v *VersionVector) IsEmpty() bool {
	return len(v.dots) == 0
}

// Compare returns the causal order between this and another
// vector clock
//
func (v *VersionVector) Compare(other *VersionVector) CausalOrder {
	v.l.RLock()
	defer v.l.RUnlock()

	if dotsEql(v.dots, other.dots) {
		return OrderEqual
	}

	if v.DescendsFrom(other) {
		return OrderGreater
	}

	if other.DescendsFrom(v) {
		return OrderLess
	}

	return OrderNone
}

// Subtract returns a new VersionVector with only the entries that
// dominate the other VersionVector.
//
func (v *VersionVector) Subtract(other *VersionVector) *VersionVector {
	dominatingDots := make(Dots)

	v.l.RLock()
	for actor, time := range v.dots {
		otherTime, exists := other.Get(actor)
		if !exists || time > otherTime {
			dominatingDots[actor] = time
		}
	}
	v.l.RUnlock()

	return &VersionVector{
		dots: dominatingDots,
	}
}

// Intersection returns a new VersionVector that contains the
// common (same actor and time) elements for both VersionVectors.
//
func (v *VersionVector) Intersection(other *VersionVector) *VersionVector {
	dots := make(Dots)

	v.l.RLock()
	for actor, time := range v.dots {
		otherTime, exists := other.Get(actor)
		if exists && otherTime == time {
			dots[actor] = time
		}
	}
	v.l.RUnlock()

	return &VersionVector{
		dots: dots,
	}
}

// DescendsFrom indicates whether this clock is causually greater
// than the other clock
//
// This method is not thread safe
//
func (v *VersionVector) DescendsFrom(other *VersionVector) bool {
	for actor, otherTime := range other.dots {
		time, exists := v.dots[actor]
		if !exists || time <= otherTime {
			return false
		}
	}

	return true
}

// descendentOf indicates whether actor is present in this clock and whether
// it's a descendant of otherTime (i.e. it's causually greater)
//
// This method is not thread safe
//
func (v *VersionVector) descendentOf(actor string, otherTime LamportTime) bool {
	time, exists := v.dots[actor]
	return exists && time >= otherTime
}

// REach iterates over each dot, calling the supplied function for each actor.
// It holds a Read lock while doing so
//
func (v *VersionVector) REach(fn func(actor string, t LamportTime)) {
	v.l.RLock()
	defer v.l.RUnlock()

	for actor, t := range v.dots {
		fn(actor, t)
	}
}

// Clone is a deep-copy of the VersionVector
func (v *VersionVector) Clone() *VersionVector {
	dots := make(Dots)

	v.l.RLock()
	for actor, t := range v.dots {
		dots[actor] = t
	}
	v.l.RUnlock()

	return &VersionVector{
		dots: dots,
	}
}

// Marshal serialises this VersionVector to binary using protocol buffers
func (v *VersionVector) Marshal() (data []byte, err error) {
	value := CreateVersionVectorValue()

	v.l.RLock()
	for actor, t := range v.dots {
		value.Dots[actor] = &Dot{Time: t}
	}
	v.l.RUnlock()

	return value.Marshal()
}

// Unmarshal parses a protobuf encoded VersionVector and loads
// it's data into v.
func (v *VersionVector) Unmarshal(data []byte) error {
	value := CreateVersionVectorValue()

	if err := value.Unmarshal(data); err != nil {
		return err
	}

	v.l.Lock()
	for actor, dot := range value.Dots {
		v.dots[actor] = dot.Time
	}
	v.l.Unlock()

	return nil
}
