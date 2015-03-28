# Experimental datastructure

At this point, this is a simple 2-3-Tree implementing an ordered sequence, the goal is to expand this into a ordered map-like structure.

# Functional

The datastructure is fully functional, always returning a new tree on insert/delete operations.

# Examlpe

This is an example from the tests:

```go

func TestRemove (t *testing.T) {
  count := 100
  tree := TreeRoot{}

  for i := 0 ; i < count ; i++ {
    tree = tree.InsertAt(i, two3.NewLeaf(s.Itoa(i)))
  }

  for i := 0 ; i < count ; i++ {

    actual := tree.Count()
    if actual != count-i {
      t.Errorf("Should have count() equal to %v, is %v", count-i, actual)
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

```