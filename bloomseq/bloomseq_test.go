package bloomseq

import (
	ds "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-datastore"
	"github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-datastore/sync"
	"github.com/ipfs/go-ipfs/blocks/blockstore"
	bs "github.com/ipfs/go-ipfs/blockservice"
	"github.com/ipfs/go-ipfs/exchange/offline"
	mdag "github.com/ipfs/go-ipfs/merkledag"

	"encoding/binary"
	"math/rand"
	"testing"
)

func getMockDagServ(t testing.TB) mdag.DAGService {
	dstore := ds.NewMapDatastore()
	tsds := sync.MutexWrap(dstore)
	bstore := blockstore.NewBlockstore(tsds)
	bserv, err := bs.New(bstore, offline.Exchange(bstore))
	if err != nil {
		t.Fatal(err)
	}
	dserv := mdag.NewDAGService(bserv)
	return dserv
}

func BytesFromInt(i uint64) []byte {
	b := make([]byte, 8)
	binary.PutUvarint(b, i)
	return b
}

func IntFromBytes(b []byte) uint64 {
	res, _ := binary.Uvarint(b)
	return res
}

func TestAdd(t *testing.T) {
	var count uint64 = 100
	tree := BloomSeq{}

	var i uint64
	for i = 0; i < count; i++ {

		if tree.Count() != i {
			t.Fatalf("Should have count() equal to %v", count)
		}
		tree = tree.InsertAt(0, BytesFromInt(i))
		if !tree.InvariantAllLeavesAtSameDepth() {
			t.Fatal("invariant all leaves at same depth broken")
		}
	}

	if tree.Count() != count {
		t.Fatalf("Should have count() equal to %v", count)
	}
}

func TestAddReverse(t *testing.T) {
	var count uint64 = 100
	tree := BloomSeq{}

	var i uint64
	for i = 0; i < count; i++ {

		if tree.Count() != i {
			t.Fatalf("Should have count() equal to %v", count)
		}
		tree = tree.InsertAt(i, BytesFromInt(i))
		if !tree.InvariantAllLeavesAtSameDepth() {
			t.Fatal("invariant all leaves at same depth broken")
		}
	}

	if tree.Count() != count {
		t.Fatalf("Should have count() equal to %v", count)
	}
}

func TestAddRandom(t *testing.T) {
	var count uint64 = 100
	rounds := 100

	for r := 0; r < rounds; r++ {
		tree := BloomSeq{}

		var i uint64
		for i = 0; i < count; i++ {

			if tree.Count() != i {
				t.Fatalf("Should have count() equal to %v", count)
			}

			var pos uint64
			if i == 0 {
				pos = 0
			} else {
				pos = uint64(rand.Int63()) % i
			}
			tree = tree.InsertAt(pos, BytesFromInt(i))
			if !tree.InvariantAllLeavesAtSameDepth() {
				t.Fatal("invariant all leaves at same depth broken")
			}
		}

		actual := tree.Count()
		if actual != count {
			t.Fatalf("Should have count equal to %v, is %v", count, actual)
		}
	}
}

func TestGet(t *testing.T) {
	var count uint64 = 100
	tree := BloomSeq{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(0, BytesFromInt(i))
	}

	for i = 0; i < count; i++ {
		res, _ := tree.GetAt(i)

		if IntFromBytes(res) != -i-1+count {
			t.Fatalf("Got %v from index %v, expected %v", res, i, -i-1+count)
		}
	}
}

func TestGetAddingFromEnd(t *testing.T) {
	var count uint64 = 100
	tree := BloomSeq{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(i, BytesFromInt(i))
	}

	for i = 0; i < count; i++ {
		res, _ := tree.GetAt(i)

		if IntFromBytes(res) != i {
			t.Fatalf("Got %v from index %v, expected %v", res, i, i)
		}
	}
}

func TestRemoveSpecific(t *testing.T) {

	var count uint64 = 4
	tree := BloomSeq{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(i, BytesFromInt(i))
	}

	tree = tree.RemoveAt(1)

	got, _ := tree.GetAt(0)

	if IntFromBytes(got) != 0 {
		t.Fatal("did not get 0 from tree with removed 1:st element")
	}
}

func TestRemove(t *testing.T) {
	var count uint64 = 100
	tree := BloomSeq{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(i, BytesFromInt(i))
	}

	for i = 0; i < count; i++ {

		actual := tree.Count()
		if actual != count-i {
			t.Fatalf("Should have count() equal to %v, is %v", count-i, actual)
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
	tree := BloomSeq{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(i, BytesFromInt(i))
	}

	for i = 0; i < count; i++ {
		actual := tree.Count()
		if actual != count-i {
			t.Fatalf("Should have count() equal to %v, is %v", count-i, actual)
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
	tree := BloomSeq{}

	var i uint64
	for i = 0; i < count*2; i++ {
		tree = tree.InsertAt(i, BytesFromInt(i))
	}

	for i = 0; i < count; i++ {
		tree = tree.RemoveAt(i)
	}

	for i = 0; i < count; i++ {
		res, _ := tree.GetAt(i)
		if IntFromBytes(res) != i*2+1 {
			t.Fatalf("Got %v from index %v, expected %v", IntFromBytes(res), i, count-i-1)
		}
	}
}

func TestRemoveRandom(t *testing.T) {
	var count uint64 = 100
	rounds := 100

	for r := 0; r < rounds; r++ {
		tree := BloomSeq{}

		var i uint64
		for i = 0; i < count; i++ {
			tree = tree.InsertAt(i, BytesFromInt(i))
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
	tree := BloomSeq{}

	var i uint64
	for i = 0; i < 100000; i++ {
		tree = tree.InsertAt(0, []byte("fonk"))
	}
}

// persistance tests

func TestPersistEmtpyRoot(t *testing.T) {
	dserv := getMockDagServ(t)

	tree := BloomSeq{}

	persisted := tree.Persist(dserv)

	if tree != persisted {
		t.Fatal("Did not get same BloomSeq back")
	}

	if tree.Count() != persisted.Count() {
		t.Fatal("Count differs between trees")
	}

}

func TestPersistSingletonRoot(t *testing.T) {
	dserv := getMockDagServ(t)

	tree := BloomSeq{}
	tree = tree.InsertAt(0, []byte("leafy!"))

	persisted := tree.Persist(dserv)

	if tree.Count() != persisted.Count() {
		t.Fatal("Count differs between trees")
	}

	get0, _ := tree.GetAt(0)
	get1, _ := persisted.GetAt(0)

	if IntFromBytes(get0) != IntFromBytes(get1) {
		t.Fatal("Item retrieved does not match")
	}
}

func PersistAndGetAllValues(t *testing.T) {

	dserv := getMockDagServ(t)

	var count uint64 = 100
	tree := BloomSeq{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(0, BytesFromInt(i))
	}

	tree = tree.Persist(dserv)

	for i = 0; i < count; i++ {
		res, _ := tree.GetAt(i)

		if IntFromBytes(res) != -i-1+count {
			t.Fatalf("Got %v from index %v, expected %v", IntFromBytes(res), i, -i-1+count)
		}
	}
}

func TestPersistReplacingRoot(t *testing.T) {

	dserv := getMockDagServ(t)

	var count uint64 = 100
	tree := BloomSeq{}

	var i uint64
	for i = 0; i < count; i++ {

		// persist at each step
		tree = tree.Persist(dserv)
		res := tree.Count()

		if res != i {
			t.Fatalf("Should have count() equal to %v, is %v", i, res)
		}
		tree = tree.InsertAt(0, BytesFromInt(i))
	}
}

func TestPersistAndGetFirstValue(t *testing.T) {

	dserv := getMockDagServ(t)

	var count uint64 = 100
	tree := BloomSeq{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(0, BytesFromInt(i))
	}

	// persist and insert at beginning

	tree = tree.Persist(dserv)

	if tree.CountUnreferencedNodes() != 1 {
		t.Fatal("dereference fail")
	}

	tree = tree.InsertAt(0, []byte("beep boop"))

	if tree.CountUnreferencedNodes() != 10 {
		t.Fatal("Should have 10 unreferenced members")
	}

	if tree.Count() != count+1 {
		t.Fatal("Should have count + 1 members")
	}
}

func TestPersistAndGetAllValues(t *testing.T) {

	dserv := getMockDagServ(t)

	var count uint64 = 100
	tree := BloomSeq{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(0, BytesFromInt(i))
	}

	// persist and insert at beginning

	tree = tree.Persist(dserv)

	if tree.CountUnreferencedNodes() != 1 {
		t.Fatal("dereference fail")
	}

	// access all the values
	for i = 0; i < count; i++ {
		_, _ = tree.GetAt(i)
	}

	// now local tree should have all elements unreferenced
	if tree.Count() != uint64(tree.CountUnreferencedNodes()) {
		t.Fatal("Should have all nodes in memory")
	}
}
