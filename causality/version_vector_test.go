package causality_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/luma/pith/rapport/causality"
)

var _ = Describe("VersionVector", func() {
	It("can create a new empty VersionVector with CreateVersionVector()", func() {
		clock := CreateVersionVector()
		Expect(clock.IsEmpty()).To(BeTrue())
	})

	verify := func(ver *VersionVector, values map[string]LamportTime) {
		for actor, expectedTime := range values {
			t, exists := ver.Get(actor)
			Expect(exists).To(BeTrue())
			Expect(t).To(Equal(expectedTime))
		}
	}

	Describe("Merge()", func() {
		var v1, v2 *VersionVector
		var a4, a5, a6, a7 string

		BeforeEach(func() {
			v1 = CreateVersionVector()
			v2 = CreateVersionVector()
			a4 = "Actor 4"
			a5 = "Actor 5"
			a6 = "Actor 6"
			a7 = "Actor 7"
		})

		It("merges", func() {
			v1.Witness(a4, LamportTime(4))
			v1.Witness(a5, LamportTime(5))
			v1.Witness(a7, LamportTime(7))

			v2.Witness(a6, LamportTime(6))
			v2.Witness(a7, LamportTime(7))

			v1.Merge(v2)

			verify(v1, map[string]LamportTime{
				a4: LamportTime(4),
				a5: LamportTime(5),
				a6: LamportTime(6),
				a7: LamportTime(7),
			})
		})

		It("merges less right", func() {
			v1.Witness(a6, LamportTime(6))
			v1.Witness(a7, LamportTime(7))

			v2.Witness(a5, LamportTime(5))

			v1.Merge(v2)

			verify(v1, map[string]LamportTime{
				a5: LamportTime(5),
				a6: LamportTime(6),
				a7: LamportTime(7),
			})
		})

		It("merges less left", func() {
			v1.Witness(a5, LamportTime(5))

			v2.Witness(a6, LamportTime(6))
			v2.Witness(a7, LamportTime(7))

			v1.Merge(v2)

			verify(v1, map[string]LamportTime{
				a5: LamportTime(5),
				a6: LamportTime(6),
				a7: LamportTime(7),
			})
		})

		It("merges with the same ids", func() {
			v1.Witness(a4, LamportTime(1))
			v1.Witness(a5, LamportTime(1))

			v2.Witness(a4, LamportTime(1))
			v2.Witness(a6, LamportTime(1))

			v1.Merge(v2)

			verify(v1, map[string]LamportTime{
				a4: LamportTime(1),
				a5: LamportTime(1),
				a6: LamportTime(1),
			})
		})
	})

	Describe("Subtract()", func() {
		var v1, v2 *VersionVector
		var a4, a5, a6 string

		BeforeEach(func() {
			v1 = CreateVersionVector()
			v2 = CreateVersionVector()
			a4 = "Actor 4"
			a5 = "Actor 5"
			a6 = "Actor 6"
		})

		It("returns any actors that exist in Version A, but not B", func() {
			v1.Witness(a4, LamportTime(1))
			v1.Witness(a5, LamportTime(1))
			v1.Witness(a6, LamportTime(1))
			v2.Witness(a4, LamportTime(1))

			v3 := v1.Subtract(v2)

			verify(v3, map[string]LamportTime{
				a5: LamportTime(1),
				a6: LamportTime(1),
			})
		})

		It("returns any actors that exist in Version A and dominate the same actor in B", func() {
			v1.Witness(a4, LamportTime(1))
			v1.Witness(a5, LamportTime(1))
			v1.Witness(a6, LamportTime(5))

			v2.Witness(a4, LamportTime(1))
			v2.Witness(a5, LamportTime(3))
			v2.Witness(a6, LamportTime(4))

			v3 := v1.Subtract(v2)

			verify(v3, map[string]LamportTime{
				a6: LamportTime(5),
			})
		})
	})

	Describe("Intersection()", func() {
		var v1, v2 *VersionVector
		var a4, a5, a6 string

		BeforeEach(func() {
			v1 = CreateVersionVector()
			v2 = CreateVersionVector()
			a4 = "Actor 4"
			a5 = "Actor 5"
			a6 = "Actor 6"
		})

		It("returns any actors that exist in both versions and are equal", func() {
			v1.Witness(a4, LamportTime(1))
			v1.Witness(a5, LamportTime(2))
			v1.Witness(a6, LamportTime(1))

			v2.Witness(a4, LamportTime(1))
			v2.Witness(a5, LamportTime(1))
			v2.Witness(a6, LamportTime(3))

			v3 := v1.Intersection(v2)

			verify(v3, map[string]LamportTime{
				a4: LamportTime(1),
			})
		})
	})

	Describe("Incr()", func() {
		var v1 *VersionVector
		var a4, a5, a6 string

		BeforeEach(func() {
			v1 = CreateVersionVector()
			a4 = "Actor 4"
			a5 = "Actor 5"
			a6 = "Actor 6"
		})

		It("increments the actor if it exists", func() {
			v1.Witness("foo", LamportTime(3))
			newTime := v1.Incr("foo")
			Expect(newTime).To(Equal(LamportTime(4)))
		})

		It("creates and increments the actor if it does not exist", func() {
			newTime := v1.Incr("foo")
			Expect(newTime).To(Equal(LamportTime(1)))
		})
	})

	Describe("Ordering", func() {
		var v1, v2 *VersionVector
		var a, b string

		BeforeEach(func() {
			v1 = CreateVersionVector()
			v2 = CreateVersionVector()
			a = "Actor A"
			b = "Actor B"
		})

		It("indicates when A dominates", func() {
			v1.Witness(a, LamportTime(1))
			v1.Witness(a, LamportTime(2))

			v2.Witness(a, LamportTime(1))

			Expect(v1.DescendsFrom(v2)).To(BeTrue())
			Expect(v2.DescendsFrom(v1)).To(BeFalse())
			Expect(v1.IsConcurrentWith(v2)).To(BeFalse())
			Expect(v1.Compare(v2)).To(Equal(OrderGreater))
		})

		It("indicates when B dominates", func() {
			v1.Witness(a, LamportTime(1))

			v2.Witness(a, LamportTime(1))
			v2.Witness(a, LamportTime(2))

			Expect(v1.DescendsFrom(v2)).To(BeFalse())
			Expect(v2.DescendsFrom(v1)).To(BeTrue())
			Expect(v1.IsConcurrentWith(v2)).To(BeFalse())
			Expect(v2.Compare(v1)).To(Equal(OrderGreater))
		})

		It("indicates when A and B are concurrent", func() {
			v1.Witness(a, LamportTime(2))
			v1.Witness(b, LamportTime(1))

			v2.Witness(a, LamportTime(3))

			Expect(v1.DescendsFrom(v2)).To(BeFalse())
			Expect(v2.DescendsFrom(v1)).To(BeFalse())
			Expect(v1.IsConcurrentWith(v2)).To(BeTrue())
			Expect(v2.Compare(v1)).To(Equal(OrderNone))
		})

		It("indicates when A and B are equal", func() {
			v1.Witness(a, LamportTime(1))

			v2.Witness(a, LamportTime(1))

			Expect(v1.DescendsFrom(v2)).To(BeFalse())
			Expect(v2.DescendsFrom(v1)).To(BeFalse())
			Expect(v1.IsConcurrentWith(v2)).To(BeFalse())
			Expect(v2.Compare(v1)).To(Equal(OrderEqual))
		})
	})
})
