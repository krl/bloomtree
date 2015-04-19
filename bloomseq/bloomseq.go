package bloomseq

import (
	"fmt"
	mdag "github.com/ipfs/go-ipfs/merkledag"
	// . "github.com/krl/bloomtree/value"
)

type BloomSeq struct {
	value tree
}

func (r BloomSeq) GetLeavesDepth() []int {
	if r.value == nil {
		return []int{}
	}
	return r.value.getLeavesDepth(0)
}

func (r BloomSeq) CountUnreferencedNodes() int {
	if r.value != nil {
		return r.value.countUnreferencedNodes()
	}
	fmt.Printf("empty tree case: 0")
	return 0
}

func (r BloomSeq) InvariantAllLeavesAtSameDepth() bool {
	depths := r.GetLeavesDepth()
	track := 0

	for i := 0; i < len(depths); i++ {
		if track == 0 {
			track = depths[i]
		}
		if depths[i] != track {
			return false
		}
	}
	return true
}

func (r BloomSeq) Count() uint64 {
	if r.value == nil {
		return 0
	} else {
		return r.value.count()
	}
}

// takes a pointer, to be mutable
// because we want to cache the nodes loaded from disk
// the logical value will still be immutable, and all non-test
// functions will report the same value
func (r *BloomSeq) GetAt(i uint64) ([]byte, error) {
	if i > r.Count() {
		return nil, fmt.Errorf("Index out of bounds")
	}

	tree, leaf := r.value.getAt(i)

	// scary mutation
	r.value = tree

	return leaf.Value, nil
}

func (r BloomSeq) RemoveAt(i uint64) BloomSeq {
	if i > r.Count()-1 {
		panic("Index out of bounds")
	}

	new, _ := r.value.removeAt(i)

	return BloomSeq{value: new}
}

func (r BloomSeq) InsertAt(i uint64, s []byte) BloomSeq {

	leaf := newLeaf(s)

	if r.value == nil {
		// if tree root is empty, just insert the leaf
		return BloomSeq{value: leaf}
	} else {
		// else insert in element
		ref1, ref2 := r.value.insertAt(i, leaf)

		// do we have a split?
		if ref2 != nil {
			return BloomSeq{value: newNode2([]tree{ref1, ref2})}
		} else {
			return BloomSeq{value: ref1}
		}
	}
	return BloomSeq{}
}

// Persistance

func (r BloomSeq) Persist(dserv mdag.DAGService) BloomSeq {
	if r.value != nil {
		return BloomSeq{value: r.value.persist(dserv)}
	} else {
		return BloomSeq{}
	}
}
