package root

import (
	"fmt"
	//"github.com/ipfs/go-ipfs/core"
	mdag "github.com/ipfs/go-ipfs/merkledag"
	two3 "github.com/krl/bloomtree/two3"
)

type TreeRoot struct {
	value two3.Tree
}

func (r TreeRoot) GetLeavesDepth() []int {
	if r.value == nil {
		return []int{}
	}
	return r.value.GetLeavesDepth(0)
}

func (r TreeRoot) CountUnreferencedNodes() int {
	if r.value != nil {
		return r.value.CountUnreferencedNodes()
	}
	fmt.Printf("empty tree case: 0")
	return 0
}

func (r TreeRoot) InvariantAllLeavesAtSameDepth() bool {
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

func (r TreeRoot) Count() uint64 {
	if r.value == nil {
		return 0
	} else {
		return r.value.Count()
	}
}

func (r TreeRoot) GetAt(i uint64) (TreeRoot, string, error) {
	if i > r.Count() {
		return r, "", fmt.Errorf("Index out of bounds")
	}

	tree, leaf := r.value.GetAt(i)

	return TreeRoot{value: tree}, leaf.Value, nil

}

func (r TreeRoot) RemoveAt(i uint64) TreeRoot {
	if i > r.Count()-1 {
		panic("Index out of bounds")
	}

	new, _ := r.value.RemoveAt(i)

	return TreeRoot{value: new}
}

func (r TreeRoot) InsertAt(i uint64, s string) TreeRoot {

	leaf := two3.NewLeaf(s)

	if r.value == nil {
		// if tree root is empty, just insert the leaf
		return TreeRoot{value: leaf}
	} else {
		// else insert in element
		ref1, ref2 := r.value.InsertAt(i, leaf)

		// do we have a split?
		if ref2 != nil {
			return TreeRoot{value: two3.NewNode2([]two3.Tree{ref1, ref2})}
		} else {
			return TreeRoot{value: ref1}
		}
	}
	return TreeRoot{}
}

// Persistance

func (r TreeRoot) Persist(dserv mdag.DAGService) TreeRoot {
	if r.value != nil {
		return TreeRoot{value: r.value.Persist(dserv)}
	} else {
		return TreeRoot{}
	}
}
