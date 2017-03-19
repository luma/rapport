package rapport

import (
	"fmt"
	"sync"

	"github.com/luma/pith/keys"
	"github.com/luma/pith/rapport/causality"
)

var (
	// EntriesKey is the sigil used to deliminate a key that is for a set's entries
	EntriesKey = []byte("E")

	// DeferredKey is the sigil used to deliminate a key that is for a set's deferreds
	DeferredKey = []byte("D")
)

// AWSet is a Add-Wins Set. AKA An addition-based, OR-set without tombstones,
// ORSWOT is harder to type though.
//
// This was ported from riak_dt
//
type AWSet struct {
	Version  *causality.VersionVector
	entries  map[string]*causality.VersionVector
	deferred DeferredMap
	l        sync.RWMutex
}

// CreateAWSet returns a new, empty AWSet.
//
func CreateAWSet() *AWSet {
	return &AWSet{
		Version:  causality.CreateVersionVector(),
		entries:  make(map[string]*causality.VersionVector),
		deferred: make(DeferredMap),
	}
}

// AddOne adds a single element to the set for a specific replica. It returns
// true if the element was added, otherwise it returns false.
//
func (a *AWSet) AddOne(value string, replica string) bool {
	newTime := a.Version.Incr(replica)
	entry := causality.CreateVersionVector()
	entry.Witness(replica, newTime)

	a.l.Lock()
	_, alreadyExists := a.entries[value]
	a.entries[value] = entry
	a.l.Unlock()

	return !alreadyExists
}

// Add adds multiple elements to the set for a specific replica. It returns the
// number of elements that were added.
//
func (a *AWSet) Add(values []string, replica string) int {
	added := 0

	for _, value := range values {
		if a.AddOne(value, replica) {
			added++
		}
	}

	return added
}

// RemoveOne removes a single element from the set by value. It returns the
// VersionVector of the element that was removed.
//
func (a *AWSet) RemoveOne(value string) *causality.VersionVector {
	a.l.Lock()
	version := a.entries[value]
	delete(a.entries, value)
	a.l.Unlock()
	return version
}

// Remove removes a number of values from the set and returns the number
// of elements that were removed.
//
func (a *AWSet) Remove(values []string) int {
	removed := 0

	for _, value := range values {
		if a.RemoveOne(value) != nil {
			removed++
		}
	}

	return removed
}

// RemoveOneWithContext removes a values using a witnessing context. It returns
// the VersionVector of the element that was removed.
//
func (a *AWSet) RemoveOneWithContext(value string, context *causality.VersionVector) *causality.VersionVector {
	a.l.Lock()
	defer a.l.Unlock()

	if !context.Subtract(a.Version).IsEmpty() {
		// Context dominates at least some items in our version so we
		// should track it
		deferred := a.deferred[context]
		if deferred == nil {
			deferred = MakeDeferredSet()
		}

		deferred.Members[value] = true
		a.deferred[context] = deferred
	}

	existingContext, exists := a.entries[value]

	delete(a.entries, value)
	if !exists {
		return nil
	}

	domVersions := existingContext.Subtract(context)
	if !domVersions.IsEmpty() {
		// Re-add any of the versions for which we still dominate context
		a.entries[value] = domVersions
	}

	return domVersions
}

// RemoveWithContext removes a number of values using a witnessing context.
// It returns the number of elements that were removed.
//
func (a *AWSet) RemoveWithContext(values []string, context *causality.VersionVector) int {
	removed := 0

	for _, value := range values {
		if a.RemoveOneWithContext(value, context) != nil {
			removed++
		}
	}

	return removed
}

// Values returns the set elements
func (a *AWSet) Values() []string {
	a.l.RLock()
	values := make([]string, 0, len(a.entries))
	for value := range a.entries {
		values = append(values, value)
	}
	a.l.RUnlock()

	return values
}

// Cardinality returns the number of elements in the set
func (a *AWSet) Cardinality() int {
	return len(a.entries)
}

// IsEmpty returns true if the set contains no elements
func (a *AWSet) IsEmpty() bool {
	return len(a.entries) == 0
}

// Contains returns true if the value is in the set
func (a *AWSet) Contains(value string) bool {
	a.l.RLock()
	_, exists := a.entries[value]
	a.l.RUnlock()

	return exists
}

// Each iterates over the set calling the provided function at each iteraction
func (a *AWSet) Each(fn func(string)) {
	a.l.RLock()
	defer a.l.RUnlock()

	for value := range a.entries {
		fn(value)
	}
}

// Union returns a new set that is the union between this
// set and the other
func (a *AWSet) Union(other Set, replica string) Set {
	union := CreateAWSet()
	union.Add(a.Values(), replica)
	union.Add(other.Values(), replica)
	return union
}

// Intersect returns a new set that is the intersection between this
// set and the other
func (a *AWSet) Intersect(other Set, replica string) Set {
	intersection := CreateAWSet()

	a.Each(func(value string) {
		if other.Contains(value) {
			intersection.AddOne(value, replica)
		}
	})

	return intersection
}

// IsSubsetOf indicates whether this set is a subset of the other
func (a *AWSet) IsSubsetOf(other Set) bool {
	for _, value := range a.Values() {
		if !other.Contains(value) {
			return false
		}
	}

	return true
}

// Difference returns a new set that is the difference between this
// set and the other
func (a *AWSet) Difference(other Set) []string {
	diff := make([]string, 0)

	for _, value := range a.Values() {
		if !other.Contains(value) {
			diff = append(diff, value)
		}
	}

	return diff
}

// GetEntry returns the VersionVector associated with a specfic set value.
// If the set does not contain the value then it returns nil.
func (a *AWSet) GetEntry(value string) *causality.VersionVector {
	a.l.RLock()
	version := a.entries[value]
	a.l.RUnlock()

	if version == nil {
		return nil
	}

	return version.Clone()
}

// Merge another AWSet into this one
//
func (a *AWSet) Merge(crdt CRDT) {
	other := crdt.(*AWSet)

	a.l.Lock()
	defer func() {
		a.l.Unlock()
		a.applyDeferred()
	}()

	finalEntries := make(map[string]*causality.VersionVector)
	otherRemaining := make(map[string]*causality.VersionVector)

	other.l.Lock()
	for value, version := range other.entries {
		otherRemaining[value] = version
	}
	other.l.Unlock()

	for value, version := range a.entries {
		if otherVer := other.GetEntry(value); otherVer == nil {
			// Other doesn't know about this value because it:
			//  1. Has never added it
			//  2. Has added it, but it was then removed
			if version.Subtract(other.Version).IsEmpty() {
				// The other set knew about the value being added, we
				// know this because they have a "newer" version which
				// must have seen our older one. Consequently, this means
				// that they have removed the value
			} else {
				// The other set hasn't seen this value yet. So we should
				// keep it.
				finalEntries[value] = version
			}

		} else {
			// https://github.com/basho/riak_dt/blob/develop/src/riak_dt_orswot.erl#L310
			// The value is present in both but may still have been removed
			common := version.Intersection(otherVer)
			luniq := version.Subtract(common)
			runiq := otherVer.Subtract(common)
			lkeep := luniq.Subtract(other.Version)
			rkeep := runiq.Subtract(a.Version)

			common.Merge(lkeep)
			common.Merge(rkeep)
			if common.IsEmpty() {
				// nothing to drop
			} else {
				finalEntries[value] = common
			}

			// This is common so let's remove it so we avoid double work below
			delete(otherRemaining, value)
		}
	}

	for value, version := range otherRemaining {
		uniq := version.Subtract(a.Version)
		if !uniq.IsEmpty() {
			// Other has witnessed additions that we don't have. Add them.
			finalEntries[value] = uniq
		}
	}

	// merge deferred removals
	for version, otherDeferred := range other.deferred {
		deferred := a.deferred[version]
		if deferred == nil {
			deferred = MakeDeferredSet()
		}

		for removedValue := range otherDeferred.Members {
			deferred.Members[removedValue] = true
		}

		a.deferred[version] = deferred
	}

	a.entries = finalEntries
	a.Version.Merge(other.Version)
}

func (a *AWSet) applyDeferred() {
	deferredMap := a.deferred.Clone()
	a.deferred = make(DeferredMap)
	for version, entries := range *deferredMap {
		a.RemoveWithContext(entries.Values(), version)
	}
}

// Marshal serialises the set data to bytes
func (a *AWSet) Marshal() (data []*Segment, err error) {
	a.l.RLock()

	segments := make([]*Segment, 0, 1+len(a.entries)+len(a.deferred))

	v, err := a.Version.Marshal()
	if err != nil {
		return nil, err
	}

	segments = append(segments, &Segment{
		Value: v,
	})

	for value, version := range a.entries {
		b, err := version.Marshal()
		if err != nil {
			return nil, err
		}

		segments = append(segments, &Segment{
			KeySuffix: keys.Make(EntriesKey, []byte(value)),
			Value:     b,
		})
	}

	for version, deferredSet := range a.deferred {
		v, err := version.Marshal()
		if err != nil {
			return nil, err
		}

		b, err := deferredSet.Marshal()
		if err != nil {
			return nil, err
		}

		segments = append(segments, &Segment{
			KeySuffix: keys.Make(DeferredKey, v),
			Value:     b,
		})
	}

	a.l.RUnlock()

	return segments, nil
}

// Marshal deserialises the set data from bytes
func (a *AWSet) Unmarshal(data []*Segment) error {
	version := causality.CreateVersionVector()
	entries := make(map[string]*causality.VersionVector)
	deferred := make(DeferredMap)

	a.l.Lock()
	defer a.l.Unlock()

	err := version.Unmarshal(data[0].Value)
	if err != nil {
		return err
	}

	for _, s := range data[1:] {
		if s.KeySuffix[0] == EntriesKey[0] {
			entryVersion, err := causality.UnmarshalVersionVector(s.Value)
			if err != nil {
				return err
			}

			// Strip off the key sigil and add the entry
			entryKey := string(s.KeySuffix[2:])
			entries[entryKey] = entryVersion

		} else if s.KeySuffix[0] == DeferredKey[0] {
			deferredVersion, err := causality.UnmarshalVersionVector([]byte(s.KeySuffix[2:]))
			if err != nil {
				return err
			}

			deferredSet, err := UnmarshalDeferredSet(s.Value)
			if err != nil {
				return err
			}

			deferred[deferredVersion] = deferredSet

		} else {
			return fmt.Errorf("Unexpected key suffix for set: %s", s.KeySuffix)
		}
	}

	a.Version = version
	a.entries = entries
	a.deferred = deferred

	return nil
}
