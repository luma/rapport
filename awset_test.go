package rapport_test

import (
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/luma/pith/rapport"
	"github.com/luma/pith/rapport/causality"
)

var _ = Describe("AWSet", func() {
	var set *AWSet

	JustBeforeEach(func() {
		set = CreateAWSet()
		set.AddOne("foo", "replica1")
	})

	Describe("AddOne()", func() {
		It("increments the version vector", func() {
			t, exists := set.Version.Get("replica1")
			Expect(t).To(Equal(causality.LamportTime(1)))
			Expect(exists).To(BeTrue())
		})

		It("records the entry", func() {
			Expect(set.GetEntry("foo")).ToNot(BeNil())
		})

		It("witnesses the new entry", func() {
			version := set.GetEntry("foo")
			t, exists := version.Get("replica1")
			Expect(t).To(Equal(causality.LamportTime(1)))
			Expect(exists).To(BeTrue())
		})

		It("increases the version and rewitnesses the entry when adding the same value", func() {
			set.AddOne("foo", "replica1")

			// Set version records the new add
			t, _ := set.Version.Get("replica1")
			Expect(t).To(Equal(causality.LamportTime(2)))

			// entry version records the change
			version := set.GetEntry("foo")
			t, _ = version.Get("replica1")
			Expect(t).To(Equal(causality.LamportTime(2)))
		})

		It("records values from different replicas separately", func() {
			set.AddOne("foo", "replica2")

			// Set version records the new add
			t1, _ := set.Version.Get("replica1")
			Expect(t1).To(Equal(causality.LamportTime(1)))
			t2, _ := set.Version.Get("replica2")
			Expect(t2).To(Equal(causality.LamportTime(1)))

			// entry version records the change
			version := set.GetEntry("foo")
			t2, _ = version.Get("replica2")
			Expect(t2).To(Equal(causality.LamportTime(1)))
		})

		It("adds the right number of values", func() {
			Expect(set.Cardinality()).To(Equal(1))
		})
	})

	Describe("Add()", func() {
		JustBeforeEach(func() {
			set.Add([]string{"foo", "bar", "baz"}, "replica1")
		})

		It("adds several values for the same replica", func() {
			values := set.Values()
			sort.Strings(values)
			Expect(values).To(Equal([]string{"bar", "baz", "foo"}))
		})

		It("increments the version once for each add", func() {
			t1, _ := set.Version.Get("replica1")
			Expect(t1).To(Equal(causality.LamportTime(4)))
		})

		It("adds the right number of values", func() {
			Expect(set.Cardinality()).To(Equal(3))
		})
	})

	Describe("RemoveOne()", func() {
		var removedVersion *causality.VersionVector

		JustBeforeEach(func() {
			removedVersion = set.RemoveOne("foo")
		})

		// Removes do not increment the vv, only adds
		It("does not increment the version vector", func() {
			t, exists := set.Version.Get("replica1")
			Expect(t).To(Equal(causality.LamportTime(1)))
			Expect(exists).To(BeTrue())
		})

		It("removes the entry", func() {
			Expect(set.GetEntry("foo")).To(BeNil())
		})
	})

	Describe("RemoveOneWithContext()", func() {
		var context *causality.VersionVector

		JustBeforeEach(func() {
			context = causality.CreateVersionVector()
			context.Witness("replica1", causality.LamportTime(2))

			set.RemoveOneWithContext("foo", context)
		})

		Context("Version dominates supplied context", func() {
			BeforeEach(func() {
				set.Add([]string{"foo", "bar", "baz"}, "replica1")
			})

			// Removes do not increment the vv, only adds
			It("does not increment the version vector", func() {
				t, exists := set.Version.Get("replica1")
				Expect(t).To(Equal(causality.LamportTime(1)))
				Expect(exists).To(BeTrue())
			})

			It("removes the entry", func() {
				Expect(set.GetEntry("foo")).To(BeNil())
			})
		})

		Context("Supplied context dominates the Version", func() {
			// Removes do not increment the vv, only adds
			It("does not increment the version vector", func() {
				t, exists := set.Version.Get("replica1")
				Expect(t).To(Equal(causality.LamportTime(1)))
				Expect(exists).To(BeTrue())
			})

			It("removes the entry", func() {
				Expect(set.GetEntry("foo")).To(BeNil())
			})
		})

		Context("Supplied context dominates the Version and replaces another context", func() {
			JustBeforeEach(func() {
				set.RemoveOneWithContext("foo", context)
			})

			// Removes do not increment the vv, only adds
			It("does not increment the version vector", func() {
				t, exists := set.Version.Get("replica1")
				Expect(t).To(Equal(causality.LamportTime(1)))
				Expect(exists).To(BeTrue())
			})

			It("removes the entry", func() {
				Expect(set.GetEntry("foo")).To(BeNil())
			})
		})
	})

	Describe("Remove()", func() {
		var preVersion causality.LamportTime
		var removedCount int

		JustBeforeEach(func() {
			set.Add([]string{"foo", "bar", "baz"}, "replica1")
			preVersion, _ = set.Version.Get("replica1")
			removedCount = set.Remove([]string{"foo", "bar"})
		})

		// Removes do not increment the vv, only adds
		It("does not increment the version vector", func() {
			t, _ := set.Version.Get("replica1")
			Expect(t).To(Equal(preVersion))
		})

		It("removes the right values", func() {
			Expect(set.Values()).To(Equal([]string{"baz"}))
		})

		It("returns how many elements were removed", func() {
			Expect(removedCount).To(Equal(2))
		})
	})

	Describe("RemoveWithContext()", func() {
		var preVersion causality.LamportTime
		var removedCount int
		var context *causality.VersionVector
		var expectValues []string

		JustBeforeEach(func() {
			set.Add([]string{"foo", "bar", "baz"}, "replica1")

			context = causality.CreateVersionVector()
			context.Witness("replica1", causality.LamportTime(5))

			preVersion, _ = set.Version.Get("replica1")
			removedCount = set.RemoveWithContext([]string{"foo", "bar"}, context)

			if expectValues == nil {
				expectValues = []string{"baz"}
			}
		})

		shouldVerifyRemoved := func() {
			// Removes do not increment the vv, only adds
			It("does not increment the version vector", func() {
				t, _ := set.Version.Get("replica1")
				Expect(t).To(Equal(preVersion))
			})

			It("removes the right values", func() {
				Expect(set.Values()).To(Equal([]string{"baz"}))
			})

			It("returns how many elements were removed", func() {
				Expect(removedCount).To(Equal(2))
			})
		}

		Context("Supplied context dominates the Version", shouldVerifyRemoved)

		Context("Version dominates supplied context", func() {
			BeforeEach(func() {
				set.Add([]string{"foo2", "bar2", "baz2"}, "replica1")
				expectValues = []string{"bar2", "baz", "baz2", "foo2"}
			})

			shouldVerifyRemoved()
		})
	})

	Describe("Merge()", func() {
		It("merges an empty into a non-empty one", func() {
			beforeVer, _ := set.Version.Get("replica1")
			beforeValueVer, _ := set.GetEntry("foo").Get("replica1")

			set2 := CreateAWSet()
			set.Merge(set2)

			afterVer, _ := set.Version.Get("replica1")
			afterValueVer, _ := set.GetEntry("foo").Get("replica1")

			Expect(set.Values()).To(Equal([]string{"foo"}))
			Expect(afterVer).To(Equal(beforeVer))
			Expect(afterValueVer).To(Equal(beforeValueVer))
		})

		It("merges a non-empty set into an empty one", func() {
			set2 := CreateAWSet()
			beforeVer, _ := set2.Version.Get("replica1")

			set2.Merge(set)

			afterVer, _ := set2.Version.Get("replica1")
			afterValueVer, _ := set2.GetEntry("foo").Get("replica1")

			Expect(set2.Values()).To(Equal([]string{"foo"}))
			Expect(afterVer).To(Equal(beforeVer + 1))
			Expect(afterValueVer).To(Equal(causality.LamportTime(1)))
		})

		It("merges values that have been added in both replicas, then removed in one", func() {
			set2 := CreateAWSet()
			set.AddOne("bar", "replica1")

			// Sync set2 with set
			set2.Merge(set)

			set2.RemoveOne("foo")

			Expect(set.Contains("foo")).To(BeTrue())
			Expect(set2.Contains("foo")).To(BeFalse())

			// Sync set with set2
			set.Merge(set2)

			ver1, _ := set.Version.Get("replica1")
			ver2, _ := set2.Version.Get("replica1")

			Expect(ver1).To(Equal(ver2))
			Expect(set.Contains("foo")).To(BeFalse())
			Expect(set2.Contains("foo")).To(BeFalse())
		})

		It("correctly merges adds and removes when they happen in different replicas", func() {
			set2 := CreateAWSet()
			set3 := CreateAWSet()
			set.AddOne("bar", "replica1")

			// Sync set2 with set
			set2.Merge(set)

			set2.RemoveOne("bar")

			set3.Merge(set2)
			set3.Add([]string{"bar", "baz"}, "replica3")

			set.Merge(set3)
			set.Merge(set2)

			values := set.Values()
			sort.Strings(values)

			Expect(values).To(Equal([]string{"bar", "baz", "foo"}))

			ver1, _ := set.Version.Get("replica1")
			ver2, _ := set2.Version.Get("replica2")
			ver3, _ := set3.Version.Get("replica3")

			Expect(ver1).To(Equal(causality.LamportTime(2)))
			Expect(ver2).To(Equal(causality.LamportTime(0)))
			Expect(ver3).To(Equal(causality.LamportTime(2)))
		})

		XIt("applies any deferred removals", func() {
			// TODO
		})
	})

	Describe("Contains()", func() {
		It("returns true when the set contains the value", func() {
			Expect(set.Contains("foo")).To(BeTrue())
		})

		It("returns false when the set has never contained the value", func() {
			Expect(set.Contains("wut")).To(BeFalse())
		})

		It("returns false when the value was removed from the set", func() {
			Expect(set.RemoveOne("foo")).ToNot(BeNil())
			Expect(set.Contains("foo")).To(BeFalse())
		})
	})

	Describe("Each()", func() {
		It("calls the function for each value", func() {
			visitedValues := make([]string, 0)
			set.Add([]string{"bar", "baz"}, "replica1")

			set.Each(func(value string) {
				visitedValues = append(visitedValues, value)
			})

			sort.Strings(visitedValues)
			Expect(visitedValues).To(Equal([]string{"bar", "baz", "foo"}))
		})

		It("doesn't explode when called on an empty set", func() {
			visitedValues := make([]string, 0)
			set2 := CreateAWSet()

			set2.Each(func(value string) {
				visitedValues = append(visitedValues, value)
			})

			Expect(visitedValues).To(HaveLen(0))
		})
	})

	Describe("Union()", func() {
		var set2 *AWSet

		JustBeforeEach(func() {
			set2 = CreateAWSet()
		})

		It("returns a other set if the receiver set is empty", func() {
			set3 := set2.Union(set, "replica1")
			Expect(set3.Values()).To(Equal(set.Values()))
		})

		It("returns a receiver set if the other set is empty", func() {
			set3 := set.Union(set2, "replica1")
			Expect(set3.Values()).To(Equal(set.Values()))
		})

		It("returns the union of both sets", func() {
			set2.Add([]string{"foo", "bar", "baz"}, "replica1")
			set.Add([]string{"wut"}, "replica1")

			values := set.Union(set2, "replica1").Values()
			sort.Strings(values)

			Expect(values).To(Equal([]string{
				"bar", "baz", "foo", "wut",
			}))
		})
	})

	Describe("Intersect()", func() {
		var set2 *AWSet

		JustBeforeEach(func() {
			set2 = CreateAWSet()
		})

		It("returns an empty set if the receiver is empty", func() {
			set3 := set2.Intersect(set, "replica1")
			Expect(set3.IsEmpty()).To(BeTrue())
		})

		It("returns an empty set if the other set is empty", func() {
			set3 := set.Intersect(set2, "replica1")
			Expect(set3.IsEmpty()).To(BeTrue())
		})

		It("returns the intersection of both sets", func() {
			set2.Add([]string{"foo", "bar", "baz"}, "replica1")
			set.Add([]string{"wut", "baz"}, "replica1")
			values := set.Intersect(set2, "replica1").Values()
			sort.Strings(values)

			Expect(values).To(Equal([]string{
				"baz", "foo",
			}))
		})
	})

	Describe("IsSubsetOf()", func() {
		var set2 *AWSet

		JustBeforeEach(func() {
			set2 = CreateAWSet()
			set2.AddOne("foo", "replica1")
		})

		It("returns true when the sets are equal", func() {
			Expect(set.IsSubsetOf(set2)).To(BeTrue())
		})

		It("returns true when the tested set is a subset", func() {
			set2.Add([]string{"foo", "bar", "baz"}, "replica1")
			Expect(set.IsSubsetOf(set2)).To(BeTrue())
		})

		It("returns false when the tested set is not a subset", func() {
			set.Add([]string{"foo", "bar", "baz"}, "replica1")
			Expect(set.IsSubsetOf(set2)).To(BeFalse())
		})
	})

	Describe("Difference()", func() {
		var set2 *AWSet

		JustBeforeEach(func() {
			set2 = CreateAWSet()
		})

		It("returns an empty set when the sets are equal", func() {
			set2.AddOne("foo", "replica1")
			values := set.Difference(set2)
			Expect(values).To(HaveLen(0))
		})

		It("returns an empty set when the receiver set is a subset", func() {
			set2.Add([]string{"foo", "bar", "baz"}, "replica1")
			values := set.Difference(set2)
			Expect(values).To(HaveLen(0))
		})

		It("returns the receiver set if the sets are disjoint", func() {
			set2.AddOne("wut", "replica1")
			values := set.Difference(set2)
			Expect(values).To(Equal(set.Values()))
		})

		It("returns the differences between the sets", func() {
			set.Add([]string{"bar", "baz"}, "replica1")
			set2.Add([]string{"bar"}, "replica1")

			values := set.Difference(set2)
			sort.Strings(values)

			Expect(values).To(Equal([]string{"baz", "foo"}))
		})
	})
})
