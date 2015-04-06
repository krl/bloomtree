package persist_test

import (
	"fmt"
	ds "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-datastore"
	"github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-datastore/sync"
	"github.com/ipfs/go-ipfs/blocks/blockstore"
	bs "github.com/ipfs/go-ipfs/blockservice"
	"github.com/ipfs/go-ipfs/exchange/offline"
	mdag "github.com/ipfs/go-ipfs/merkledag"
	s "strconv"

	. "github.com/krl/bloomtree/root"
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

func TestPersistEmtpyRoot(t *testing.T) {
	dserv := getMockDagServ(t)

	tree := TreeRoot{}

	persisted := tree.Persist(dserv)

	if tree != persisted {
		t.Fatal("Did not get same TreeRoot back")
	}

	if tree.Count() != persisted.Count() {
		t.Fatal("Count differs between trees")
	}

}

func TestPersistSingletonRoot(t *testing.T) {
	dserv := getMockDagServ(t)

	tree := TreeRoot{}
	tree = tree.InsertAt(0, "leafy!")

	persisted := tree.Persist(dserv)

	if tree.Count() != persisted.Count() {
		t.Fatal("Count differs between trees")
	}

	get0, _ := tree.GetAt(0)
	get1, _ := persisted.GetAt(0)

	if get0 != get1 {
		t.Fatal("Item retrieved does not match")
	}
}

func PersistAndGetAllValues(t *testing.T) {

	dserv := getMockDagServ(t)

	var count uint64 = 100
	tree := TreeRoot{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(0, s.FormatUint(i, 10))
	}

	tree = tree.Persist(dserv)

	for i = 0; i < count; i++ {
		res, _ := tree.GetAt(i)

		if res != s.FormatUint(-i-1+count, 10) {
			t.Errorf("Got %v from index %v, expected %v", res, i, -i-1+count)
		}
	}
}

func TestPersistReplacingRoot(t *testing.T) {

	dserv := getMockDagServ(t)

	var count uint64 = 100
	tree := TreeRoot{}

	var i uint64
	for i = 0; i < count; i++ {

		// persist at each step
		tree = tree.Persist(dserv)
		res := tree.Count()

		if res != i {
			t.Errorf("Should have count() equal to %v, is %v", i, res)
		}
		tree = tree.InsertAt(0, s.FormatUint(i, 10))
	}
}

func TestPersistAndGetFirstValue(t *testing.T) {

	dserv := getMockDagServ(t)

	var count uint64 = 100
	tree := TreeRoot{}

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

func TestPersistAndGetAllValues(t *testing.T) {

	dserv := getMockDagServ(t)

	var count uint64 = 100
	tree := TreeRoot{}

	var i uint64
	for i = 0; i < count; i++ {
		tree = tree.InsertAt(0, s.FormatUint(i, 10))
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
