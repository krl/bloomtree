package root

import (
	//  "fmt"
	two3 "github.com/krl/bloomtree/two3"
	"math/rand"
	s "strconv"
	"testing"
)

func TestAdd(t *testing.T) {
	count := 100
	tree := TreeRoot{}

	for i := 0; i < count; i++ {

		if tree.Count() != i {
			t.Errorf("Should have count() equal to %v", count)
		}
		tree = tree.InsertAt(0, two3.NewLeaf(s.Itoa(i)))
		if !tree.InvariantAllLeavesAtSameDepth() {
			t.Errorf("invariant all leaves at same depth broken")
		}
	}

	if tree.Count() != count {
		t.Errorf("Should have count() equal to %v", count)
	}
}

func TestAddReverse(t *testing.T) {
	count := 100
	tree := TreeRoot{}

	for i := 0; i < count; i++ {

		if tree.Count() != i {
			t.Errorf("Should have count() equal to %v", count)
		}
		tree = tree.InsertAt(i, two3.NewLeaf(s.Itoa(i)))
		if !tree.InvariantAllLeavesAtSameDepth() {
			t.Errorf("invariant all leaves at same depth broken")
		}
	}

	if tree.Count() != count {
		t.Errorf("Should have count() equal to %v", count)
	}
}

func TestAddRandom(t *testing.T) {
	count := 100
	rounds := 100

	for r := 0; r < rounds; r++ {
		tree := TreeRoot{}

		for i := 0; i < count; i++ {

			if tree.Count() != i {
				t.Errorf("Should have count() equal to %v", count)
				t.FailNow()
			}

			var pos int
			if i == 0 {
				pos = 0
			} else {
				pos = rand.Int() % i
			}
			tree = tree.InsertAt(pos, two3.NewLeaf(s.Itoa(i)))
			if !tree.InvariantAllLeavesAtSameDepth() {
				t.Errorf("invariant all leaves at same depth broken")
			}
		}

		actual := tree.Count()
		if actual != count {
			t.Errorf("Should have count equal to x is y", count, actual)
			t.FailNow()
		}
	}
}

func TestGet(t *testing.T) {
	count := 100
	tree := TreeRoot{}

	for i := 0; i < count; i++ {
		tree = tree.InsertAt(0, two3.NewLeaf(s.Itoa(i)))
	}

	for i := 0; i < count; i++ {
		res, _ := tree.GetAt(i)

		if res != s.Itoa(-i-1+count) {
			t.Errorf("Got %v from index %v, expected %v", res, i, -i-1+count)
		}
	}
}

func TestGetAddingFromEnd(t *testing.T) {
	count := 100
	tree := TreeRoot{}

	for i := 0; i < count; i++ {
		tree = tree.InsertAt(i, two3.NewLeaf(s.Itoa(i)))
	}

	for i := 0; i < count; i++ {
		res, _ := tree.GetAt(i)

		if res != s.Itoa(i) {
			t.Errorf("Got %v from index %v, expected %v", res, i, i)
		}
	}
}

func TestRemoveSpecific(t *testing.T) {

	count := 4
	tree := TreeRoot{}

	for i := 0; i < count; i++ {
		tree = tree.InsertAt(i, two3.NewLeaf(s.Itoa(i)))
	}

	tree = tree.RemoveAt(1)

	got, _ := tree.GetAt(0)

	if got != s.Itoa(0) {
		t.Errorf("did not get 0 from tree with removed 1:st element")
	}
}

func TestRemove(t *testing.T) {
	count := 100
	tree := TreeRoot{}

	for i := 0; i < count; i++ {
		tree = tree.InsertAt(i, two3.NewLeaf(s.Itoa(i)))
	}

	for i := 0; i < count; i++ {

		actual := tree.Count()
		if actual != count-i {
			t.Errorf("Should have count() equal to", count-i, "is", actual)
			t.FailNow()
		}

		tree = tree.RemoveAt(0)
		if !tree.InvariantAllLeavesAtSameDepth() {
			t.Errorf("invariant all leaves at same depth broken")
		}
	}

	if tree.Count() != 0 {
		t.Errorf("Should have count() equal to 0")
	}
}

func TestRemoveReverse(t *testing.T) {
	count := 100
	tree := TreeRoot{}

	for i := 0; i < count; i++ {
		tree = tree.InsertAt(i, two3.NewLeaf(s.Itoa(i)))
	}

	for i := 0; i < count; i++ {
		actual := tree.Count()
		if actual != count-i {
			t.Errorf("Should have count() equal to %v, is %v", count-i, actual)
			t.FailNow()
		}

		tree = tree.RemoveAt(count - i - 1)
		if !tree.InvariantAllLeavesAtSameDepth() {
			t.Errorf("invariant all leaves at same depth broken")
		}
	}

	if tree.Count() != 0 {
		t.Errorf("Should have count() equal to 0")
	}
}

func TestRemoveEveryOtherPreservingOrder(t *testing.T) {
	count := 100
	tree := TreeRoot{}

	for i := 0; i < count*2; i++ {
		tree = tree.InsertAt(i, two3.NewLeaf(s.Itoa(i)))
	}

	for i := 0; i < count; i++ {
		tree = tree.RemoveAt(i)
	}

	for i := 0; i < count; i++ {
		res, _ := tree.GetAt(i)
		if res != s.Itoa(i*2+1) {
			t.Errorf("Got %v from index %v, expected %v", res, i, count-i-1)
		}
	}
}

func TestRemoveRandom(t *testing.T) {
	count := 100
	rounds := 100

	for r := 0; r < rounds; r++ {
		tree := TreeRoot{}

		for i := 0; i < count; i++ {
			tree = tree.InsertAt(i, two3.NewLeaf(s.Itoa(i)))
		}

		for i := 0; i < count; i++ {
			if tree.Count() != count-i {
				t.Errorf("Should have count() equal %v", count-i)
				t.FailNow()
			}

			var pos int
			if count-i == 0 {
				pos = 0
			} else {
				pos = rand.Int()%count - i
			}
			tree = tree.RemoveAt(pos)
			if !tree.InvariantAllLeavesAtSameDepth() {
				t.Errorf("invariant all leaves at same depth broken")
			}
		}

		actual := tree.Count()
		if actual != 0 {
			t.Errorf("Should have count equal to %v, is %v", 0, actual)
			t.FailNow()
		}
	}
}

func BenchmarkAdd(*testing.B) {
	tree := TreeRoot{}

	for i := 0; i < 100000; i++ {
		tree = tree.InsertAt(0, two3.NewLeaf("fonk"))
	}
}
