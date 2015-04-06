package root

import (
	// "fmt"
	"math/rand"
	s "strconv"
	"testing"
)

func TestAdd(t *testing.T) {
	var count uint64 = 100
	tree := TreeRoot{}

	var i uint64
	for i = 0; i < count; i++ {

		if tree.Count() != i {
			t.Fatal("Should have count() equal to %v", count)
		}
		tree = tree.InsertAt(0, s.FormatUint(i, 10))
		if !tree.InvariantAllLeavesAtSameDepth() {
			t.Fatal("invariant all leaves at same depth broken")
		}
	}

	if tree.Count() != count {
		t.Fatal("Should have count() equal to %v", count)
	}
}

func TestAddReverse(t *testing.T) {
	var count uint64 = 100
	tree := TreeRoot{}

	var i uint64
	for i = 0; i < count; i++ {

		if tree.Count() != i {
			t.Fatal("Should have count() equal to %v", count)
		}
		tree = tree.InsertAt(i, s.FormatUint(i, 10))
		if !tree.InvariantAllLeavesAtSameDepth() {
			t.Fatal("invariant all leaves at same depth broken")
		}
	}

	if tree.Count() != count {
		t.Fatal("Should have count() equal to %v", count)
	}
}

func TestAddRandom(t *testing.T) {
	var count uint64 = 100
	rounds := 100

	for r := 0; r < rounds; r++ {
		tree := TreeRoot{}

		var i uint64
		for i = 0; i < count; i++ {

			if tree.Count() != i {
				t.Fatal("Should have count() equal to %v", count)
			}

			var pos uint64
			if i == 0 {
				pos = 0
			} else {
				pos = uint64(rand.Int63()) % i
			}
			tree = tree.InsertAt(pos, s.FormatUint(i, 10))
			if !tree.InvariantAllLeavesAtSameDepth() {
				t.Fatal("invariant all leaves at same depth broken")
			}
		}

		actual := tree.Count()
		if actual != count {
			t.Fatal("Should have count equal to x is y", count, actual)
		}
	}
}

func TestGet(t *testing.T) {
	var count uint64 = 100
	tree := TreeRoot{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(0, s.FormatUint(i, 10))
	}

	for i = 0; i < count; i++ {
		res, _ := tree.GetAt(i)

		if res != s.FormatUint(-i-1+count, 10) {
			t.Fatal("Got %v from index %v, expected %v", res, i, -i-1+count)
		}
	}
}

func TestGetAddingFromEnd(t *testing.T) {
	var count uint64 = 100
	tree := TreeRoot{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(i, s.FormatUint(i, 10))
	}

	for i = 0; i < count; i++ {
		res, _ := tree.GetAt(i)

		if res != s.FormatUint(i, 10) {
			t.Fatal("Got %v from index %v, expected %v", res, i, i)
		}
	}
}

func TestRemoveSpecific(t *testing.T) {

	var count uint64 = 4
	tree := TreeRoot{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(i, s.FormatUint(i, 10))
	}

	tree = tree.RemoveAt(1)

	got, _ := tree.GetAt(0)

	if got != "0" {
		t.Fatal("did not get 0 from tree with removed 1:st element")
	}
}

func TestRemove(t *testing.T) {
	var count uint64 = 100
	tree := TreeRoot{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(i, s.FormatUint(i, 10))
	}

	for i = 0; i < count; i++ {

		actual := tree.Count()
		if actual != count-i {
			t.Fatal("Should have count() equal to", count-i, "is", actual)
		}

		tree = tree.RemoveAt(0)
		if !tree.InvariantAllLeavesAtSameDepth() {
			t.Fatal("invariant all leaves at same depth broken")
		}
	}

	if tree.Count() != 0 {
		t.Fatal("Should have count() equal to 0")
	}
}

func TestRemoveReverse(t *testing.T) {
	var count uint64 = 100
	tree := TreeRoot{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(i, s.FormatUint(i, 10))
	}

	for i = 0; i < count; i++ {
		actual := tree.Count()
		if actual != count-i {
			t.Fatal("Should have count() equal to %v, is %v", count-i, actual)
		}

		tree = tree.RemoveAt(count - i - 1)
		if !tree.InvariantAllLeavesAtSameDepth() {
			t.Fatal("invariant all leaves at same depth broken")
		}
	}

	if tree.Count() != 0 {
		t.Fatal("Should have count() equal to 0")
	}
}

func TestRemoveEveryOtherPreservingOrder(t *testing.T) {
	var count uint64 = 100
	tree := TreeRoot{}

	var i uint64
	for i = 0; i < count*2; i++ {
		tree = tree.InsertAt(i, s.FormatUint(i, 10))
	}

	for i = 0; i < count; i++ {
		tree = tree.RemoveAt(i)
	}

	for i = 0; i < count; i++ {
		res, _ := tree.GetAt(i)
		if res != s.FormatUint(i*2+1, 10) {
			t.Fatal("Got %v from index %v, expected %v", res, i, count-i-1)
		}
	}
}

func TestRemoveRandom(t *testing.T) {
	var count uint64 = 100
	rounds := 100

	for r := 0; r < rounds; r++ {
		tree := TreeRoot{}

		var i uint64
		for i = 0; i < count; i++ {
			tree = tree.InsertAt(i, s.FormatUint(i, 10))
		}

		for i = 0; i < count; i++ {
			if tree.Count() != count-i {
				t.Fatal("Should have count() equal %v", count-i)
			}

			var pos uint64
			if count-i == 0 {
				pos = 0
			} else {
				pos = uint64(rand.Int63()) % (count - i)

			}
			tree = tree.RemoveAt(pos)
			if !tree.InvariantAllLeavesAtSameDepth() {
				t.Fatal("invariant all leaves at same depth broken")
			}
		}

		actual := tree.Count()
		if actual != 0 {
			t.Fatal("Should have count equal to %v, is %v", 0, actual)
		}
	}
}

func BenchmarkAdd(*testing.B) {
	tree := TreeRoot{}

	var i uint64
	for i = 0; i < 100000; i++ {
		tree = tree.InsertAt(0, "fonk")
	}
}
