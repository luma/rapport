package rapport

import (
	"time"

	"github.com/luma/pith/rapport/causality"
)

// ProtoVersion is of the form major.minor.patch-label
const ProtoVersion = "1.0.0-alpha.0"

// Marshaler exposes method to Marshal and Unmarshal
type Marshaler interface {
	Marshal() (data []*Segment, err error)
	Unmarshal(data []*Segment) error
}

// CRDT exposes a CRDT compliant Merge method
type CRDT interface {
	Merge(other CRDT)
}

// Value encapsulates Rapport Value
type Value interface {
	CRDT
	Marshaler
}

// SetOperations encapsulates the common set operations
type SetOperations interface {
	Contains(value string) bool
	Cardinality() int
	Difference(other Set) []string
	Union(other Set, replica string) Set
	Intersect(other Set, replica string) Set
	IsSubsetOf(other Set) bool
	IsEmpty() bool
}

type Register interface {
	CRDT
	Marshaler

	Set(value string, t time.Time) error
	Get() string
}

// Set is the contract that all Pith sets must abide by
type Set interface {
	CRDT
	Marshaler
	SetOperations

	Add(values []string, replica string) int
	AddOne(value string, replica string) bool
	Remove(values []string) int
	RemoveOne(value string) *causality.VersionVector

	Values() []string
	Each(fn func(string))
}

// Counter is the contract that all Pith counters must abide by
type Counter interface {
	CRDT
	Marshaler

	Incr() int64
	IncrBy(amount int64) int64
	Decr() (int64, error)
	DecrBy(amount int64) (int64, error)

	Value() int64
}
