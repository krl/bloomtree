package root

import (
	"fmt"
	two3 "github.com/krl/bloomtree/two3"
)

type TreeRoot struct {
	ref two3.Tree
}

func (r TreeRoot) GetLeavesDepth() []int {
	if r.ref == nil {
		return []int{}
	}
	return r.ref.GetLeavesDepth(0)
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

func (r TreeRoot) Count() int {
	if r.ref == nil {
		return 0
	} else {
		return r.ref.Count()
	}
}

func (r TreeRoot) GetAt(i int) (string, error) {
	if i > r.Count() {
		return "", fmt.Errorf("Index out of bounds")
	}

	return r.ref.GetAt(i).Pointer, nil
}

func (r TreeRoot) RemoveAt(i int) TreeRoot {
	if i > r.Count()-1 {
		panic("Index out of bounds")
	}

	new, _ := r.ref.RemoveAt(i)

	return TreeRoot{ref: new}
}

func (r TreeRoot) InsertAt(i int, leaf two3.Leaf) TreeRoot {
	if r.ref == nil {
		// if tree root is empty, just insert the leaf
		return TreeRoot{ref: leaf}
	} else {
		// else insert in element
		ref1, ref2 := r.ref.InsertAt(i, leaf)

		// do we have a split?
		if ref2 != nil {
			return TreeRoot{ref: two3.NewNode2([]two3.Tree{ref1, ref2})}
		} else {
			return TreeRoot{ref: ref1}
		}
	}
	return TreeRoot{}
}
