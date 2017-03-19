package rapport

import (
	"fmt"
	"time"
)

type LWWRegister struct {
	t     time.Time
	value string
}

func CreateLWWRegister(initialValue string) *LWWRegister {
	return &LWWRegister{
		value: initialValue,
		t:     time.Now().UTC(),
	}
}

func (l *LWWRegister) Set(value string, t time.Time) error {
	if t.Before(l.t) {
		return fmt.Errorf("Cannot set register to a value from the past: %v < %v", t, l.t)
	}

	l.t = t
	l.value = value
	return nil
}

func (l *LWWRegister) Get() string {
	return l.value
}

func (l *LWWRegister) Merge(crdt CRDT) {
	otherReg := crdt.(*LWWRegister)

	if l.t.Before(otherReg.t) {
		l.value = otherReg.value
		l.t = otherReg.t
	} else if l.t == otherReg.t && otherReg.value != l.value {
		// This is bad...
		panic("Merge found the same timestamp but different values, registers have diverged")
	}
}

// Marshal serialises the register data to bytes
func (l *LWWRegister) Marshal() ([]*Segment, error) {
	segment := &Segment{
		Value: []byte(l.value),
	}
	return []*Segment{segment}, nil
}

// Marshal deserialises the register data from bytes
func (l *LWWRegister) Unmarshal(data []*Segment) error {
	l.value = string(data[0].Value)
	return nil
}
