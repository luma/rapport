package causality

import "sync/atomic"

type LamportTime uint64

type LamportClock struct {
	value uint64
}

func CreateLamportClock(value LamportTime) LamportClock {
	return LamportClock{value: uint64(value)}
}

// Value retrieves the current value of the clock
func (l *LamportClock) Value() LamportTime {
	return LamportTime(atomic.LoadUint64(&l.value))
}

// Incr increments the value by 1
func (l *LamportClock) Incr() LamportTime {
	return LamportTime(atomic.AddUint64(&l.value, 1))
}

// Dominates indicates whether this clock dominates, i.e. is causually
// greater, than the other.
//
func (l *LamportClock) Dominates(other LamportClock) bool {
	return l.Value() > other.Value()
}

// Merge moves the clock value to newer than the current value and the value
// in +time+.
func (l *LamportClock) Merge(time LamportTime) {
ATTEMPT:
	ours := atomic.LoadUint64(&l.value)
	theres := uint64(time)

	if theres <= ours {
		// If it's in the past or is identical then we don't need
		// to do anything
		return
	}

	if !atomic.CompareAndSwapUint64(&l.value, ours, theres+1) {
		// This will always eventually suceed because either the CAS suceeds or
		// our new time will move into the past and will not be applied at all.
		goto ATTEMPT
	}
}
