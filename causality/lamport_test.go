package causality_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/luma/pith/rapport/causality"
)

var _ = Describe("Lamport", func() {
	It("CreateLamportClock() creates a clock with a specific LamportTime", func() {
		time := LamportTime(10)
		l := CreateLamportClock(time)
		Expect(l.Value()).To(Equal(time))
	})

	It("Incr() increments the clock by 1", func() {
		l := CreateLamportClock(LamportTime(1))
		time := l.Incr()
		Expect(time).To(Equal(LamportTime(2)))
	})

	Describe("Merge()", func() {
		It("advances the clock if the time is newer than the current one", func() {
			l := CreateLamportClock(LamportTime(6))
			l.Merge(LamportTime(10))
			Expect(l.Value()).To(Equal(LamportTime(11)))
		})

		It("does nothing if the time to merge is in the past", func() {
			l := CreateLamportClock(LamportTime(6))
			l.Merge(LamportTime(4))
			Expect(l.Value()).To(Equal(LamportTime(6)))
		})
	})

	Describe("Dominates()", func() {
		It("returns true when the clock dominates the other clock", func() {
			l1 := CreateLamportClock(LamportTime(6))
			l2 := CreateLamportClock(LamportTime(4))
			Expect(l1.Dominates(l2)).To(BeTrue())
		})

		It("returns false when the clock equals the other clock", func() {
			l1 := CreateLamportClock(LamportTime(4))
			l2 := CreateLamportClock(LamportTime(4))
			Expect(l1.Dominates(l2)).To(BeFalse())
		})

		It("returns false when the clock is causually lesser than the other clock", func() {
			l1 := CreateLamportClock(LamportTime(4))
			l2 := CreateLamportClock(LamportTime(6))
			Expect(l1.Dominates(l2)).To(BeFalse())
		})
	})
})
