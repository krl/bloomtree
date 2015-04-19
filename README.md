# Experimental datastructure

At this point, this is a simple 2-3-Tree implementing an ordered sequence, the goal is to expand this into a ordered map-like structure.

# Functional

The datastructure is fully functional, always returning a new tree on insert/delete operations.

# Persist to IPFS

Here's an example on how to create collections that get persisted to disk using IPFS

```go
func TestPersistAndGetFirstValue(t *testing.T) {

	dserv := getMockDagServ(t)

	var count uint64 = 100
	tree := BloomSeq{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(0, s.FormatUint(i, 10))
	}

	// persist and insert at beginning

	tree = tree.Persist(dserv)

	if tree.CountUnreferencedNodes() != 1 {
		t.Fatal("dereference fail")
	}

	tree = tree.InsertAt(0, "beep boop")

	if tree.CountUnreferencedNodes() != 10 {
		t.Fatal("Should have 10 unreferenced members")
	}

	if tree.Count() != count+1 {
		t.Fatal("Should have count + 1 members")
	}
}
```

As you can see, adding 100 elements to the sequence, persisting it to disk, and then reading the first element back, only 10 trees of the node are re-created in memory. The rest of the tree just contains pointers to ipfs multi-hashes

# Examlpe

This is an example from the tests:

```go

func TestRemove(t *testing.T) {
	var count uint64 = 100
	tree := BloomSeq{}

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

```