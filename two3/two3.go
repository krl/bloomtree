package two3

import (
	proto "code.google.com/p/goprotobuf/proto"
	mdag "github.com/ipfs/go-ipfs/merkledag"
	pb "github.com/krl/bloomtree/persist/pb"
)

// tree stuff

type Tree interface {
	InsertAt(uint64, Leaf) (Tree, Tree)
	GetAt(uint64) (Tree, Leaf)
	RemoveAt(uint64) (Tree, bool)
	Count() uint64
	Persist(dserv mdag.DAGService) TreeRef

	// for tests only
	GetLeavesDepth(int) []int
	CountUnreferencedNodes() int
}

// NewLeaf

func NewLeaf(value []byte) Leaf {
	leaf := Leaf{Value: value}
	return leaf
}

// Node2

type Node2 struct {
	children []Tree
	m_count  uint64
}

func NewNode2(children []Tree) Node2 {
	node := Node2{children: children}
	for i := 0; i < 2; i++ {
		node.m_count += children[i].Count()
	}
	return node
}

func (t Node2) GetLeavesDepth(depth int) []int {
	depths := make([]int, 0, 2)

	for i := 0; i < len(t.children); i++ {
		depths = append(depths, t.children[i].GetLeavesDepth(depth+1)...)
	}
	return depths
}

func (t Node2) GetAt(i uint64) (Tree, Leaf) {
	var index int = 0

	// find the right child to get from, and decrement i
	for i >= t.children[index].Count() {
		i -= t.children[index].Count()
		index++
	}

	tree, leaf := t.children[index].GetAt(i)

	new_children := make([]Tree, 2)

	for i := 0; i < 2; i++ {
		new_children[i] = t.children[i]
	}
	new_children[index] = tree

	return NewNode2(new_children), leaf
}

func (t Node2) RemoveAt(i uint64) (Tree, bool) {

	var index int = 0

	// find the right child to remove from, and decrement i
	for i >= t.children[index].Count() {
		i -= t.children[index].Count()
		index++
	}

	var sibling_index int
	if index == 0 {
		sibling_index = 1
	} else {
		sibling_index = 0
	}

	// we're going deeper down the stack
	result, underflow := t.children[index].RemoveAt(i)

	if result == nil {
		// we deleted a leaf, and are short one child
		return t.children[sibling_index], true
	} else if underflow {

		// find the appropriate sibling to resolve imbalance
		switch s := t.children[sibling_index].(type) {
		case Node3:

			new_sibling_children := make([]Tree, 2)
			new_children := make([]Tree, 2)

			if sibling_index > index {
				// take the first child
				for i := 1; i < 3; i++ {
					new_sibling_children[i-1] = s.children[i]
				}

				new_children[index] = NewNode2([]Tree{result, s.children[0]})
			} else {
				// take the last child
				for i := 0; i < 2; i++ {
					new_sibling_children[i] = s.children[i]
				}

				new_children[index] = NewNode2([]Tree{s.children[2], result})
			}
			new_children[sibling_index] = NewNode2(new_sibling_children)
			return NewNode2(new_children), false
		case Node2:

			new_sibling_children := make([]Tree, 3)

			if sibling_index > index {
				// put at beginning of siblings children
				for i := 0; i < 2; i++ {
					new_sibling_children[i+1] = s.children[i]
				}
				new_sibling_children[0] = result
			} else {
				// put at end of siblings children
				for i := 0; i < 2; i++ {
					new_sibling_children[i] = s.children[i]
				}
				new_sibling_children[2] = result
			}

			return NewNode3(new_sibling_children), true
		}
	} else {
		// no underflow, just update index
		new_children := make([]Tree, 2)
		copy(new_children, t.children)
		new_children[index] = result
		return NewNode2(new_children), false
	}

	return t, false
}

func (t Node2) InsertAt(i uint64, leaf Leaf) (Tree, Tree) {
	var index int = 0

	// find the right child to insert into, and decrement i
	for i > t.children[index].Count() {
		i -= t.children[index].Count()
		index++
	}

	// recursive recursion
	result, extra := t.children[index].InsertAt(i, leaf)

	// new child array

	if extra == nil {
		new_children2 := make([]Tree, 2)

		for i := 0; i < 2; i++ {
			new_children2[i] = t.children[i]
		}

		new_children2[index] = result

		return NewNode2(new_children2), nil
	} else {
		new_children3 := make([]Tree, 3)

		if index == 0 {
			// insert extra before old value
			new_children3[0] = result
			new_children3[1] = extra
			new_children3[2] = t.children[1]

			return NewNode3(new_children3), nil
		} else {
			// insert extra after old value
			new_children3[0] = t.children[0]
			new_children3[1] = result
			new_children3[2] = extra
			return NewNode3(new_children3), nil
		}
	}
}

func (t Node2) Count() uint64 {
	return t.m_count
}

func (t Node2) CountUnreferencedNodes() int {
	count := 0
	for i := 0; i < 2; i++ {
		count += t.children[i].CountUnreferencedNodes()
	}

	return count
}

// Node3

type Node3 struct {
	children []Tree
	m_count  uint64
}

func NewNode3(children []Tree) Node3 {
	node := Node3{children: children}
	for i := 0; i < 3; i++ {
		node.m_count = node.m_count + children[i].Count()
	}
	return node
}

func (t Node3) GetLeavesDepth(depth int) []int {
	depths := make([]int, 0)

	for i := 0; i < len(t.children); i++ {
		depths = append(depths, t.children[i].GetLeavesDepth(depth+1)...)
	}
	return depths
}

func (t Node3) GetAt(i uint64) (Tree, Leaf) {
	var index int = 0

	// find the right child to get from, and decrement i
	for i >= t.children[index].Count() {
		i -= t.children[index].Count()
		index++
	}

	tree, leaf := t.children[index].GetAt(i)

	new_children := make([]Tree, 3)

	for i := 0; i < 3; i++ {
		new_children[i] = t.children[i]
	}
	new_children[index] = tree

	return NewNode3(new_children), leaf
}

func (t Node3) RemoveAt(i uint64) (Tree, bool) {

	var index int = 0

	// find the right child to remove from, and decrement i
	for i >= t.children[index].Count() {
		i -= t.children[index].Count()
		index++
	}

	result, underflow := t.children[index].RemoveAt(i)

	if result == nil {
		// we deleted a leaf, return a two-node

		new_children := make([]Tree, 2)

		for i := 0; i < index; i++ {
			new_children[i] = t.children[i]
		}

		// append post
		for i := index; i < 2; i++ {
			new_children[i] = t.children[i+1]
		}

		return NewNode2(new_children), false
	} else if underflow {

		// find the appropriate sibling to resolve imbalance
		var sibling_index int
		if index == 0 {
			sibling_index = 1
		} else {
			sibling_index = index - 1
		}

		switch s := t.children[sibling_index].(type) {
		case Node3:

			// let's steal a child!

			new_children := make([]Tree, 3)

			// copy old values
			for i := 0; i < 3; i++ {
				new_children[i] = t.children[i]
			}

			new_sibling_children := make([]Tree, 2)

			if sibling_index > index {

				// take the first child
				for i := 1; i < 3; i++ {
					new_sibling_children[i-1] = s.children[i]
				}
				new_children[index] = NewNode2([]Tree{result, s.children[0]})
			} else {

				// take the last child
				for i := 0; i < 2; i++ {
					new_sibling_children[i] = s.children[i]
				}
				new_children[index] = NewNode2([]Tree{s.children[2], result})

			}
			new_children[sibling_index] = NewNode2(new_sibling_children)

			return NewNode3(new_children), false

		case Node2:
			// if neighbour is node2, we shrink this node to a node2

			new_sibling_children := make([]Tree, 3)
			new_children := make([]Tree, 2)

			if sibling_index > index {
				// put at beginning of siblings children
				for i := 0; i < 2; i++ {
					new_sibling_children[i+1] = s.children[i]
				}
				new_sibling_children[0] = result

				new_children[0] = NewNode3(new_sibling_children)
				new_children[1] = t.children[index+2]

			} else {
				// put at end of siblings children
				for i := 0; i < 2; i++ {
					new_sibling_children[i] = s.children[i]
				}
				new_sibling_children[2] = result

				// pick the right pieces to keep
				if index == 0 {
					new_children[0] = t.children[1]
					new_children[1] = NewNode3(new_sibling_children)
				} else if index == 1 {
					new_children[0] = NewNode3(new_sibling_children)
					new_children[1] = t.children[2]
				} else {
					new_children[0] = t.children[0]
					new_children[1] = NewNode3(new_sibling_children)
				}
			}
			return NewNode2(new_children), false
		}
	} else {
		// no underflow, just update index
		new_children := make([]Tree, 3)
		copy(new_children, t.children)
		new_children[index] = result
		return NewNode3(new_children), false
	}

	// should not be reached
	panic("unhandled case")
	return t, false
}

func (t Node3) InsertAt(i uint64, leaf Leaf) (Tree, Tree) {

	var index int = 0

	// find the right child to insert into, and decrement i
	for i > t.children[index].Count() {
		i -= t.children[index].Count()
		index++
	}

	// deeper
	result, extra := t.children[index].InsertAt(i, leaf)

	if extra == nil {
		// without extra

		new_children := make([]Tree, 3)
		for i := 0; i < 3; i++ {
			new_children[i] = t.children[i]
		}
		new_children[index] = result

		return NewNode3(new_children), nil
	} else {
		// with extra

		new_children := make([]Tree, 4)

		// append prior
		for i := 0; i < index; i++ {
			new_children[i] = t.children[i]
		}

		// append post
		for i := index + 2; i < 4; i++ {
			new_children[i] = t.children[i-1]
		}

		// append new and extra
		new_children[index] = result
		new_children[index+1] = extra

		// split the array in two
		new_children_a := make([]Tree, 2)
		new_children_b := make([]Tree, 2)

		for i := 0; i < 2; i++ {
			new_children_a[i] = new_children[i]
		}

		for i := 2; i < 4; i++ {
			new_children_b[i-2] = new_children[i]
		}

		return NewNode2(new_children_a), NewNode2(new_children_b)
	}
}

func (t Node3) Count() uint64 {
	return t.m_count
}

func (t Node3) CountUnreferencedNodes() int {
	count := 0
	for i := 0; i < 3; i++ {
		count += t.children[i].CountUnreferencedNodes()
	}
	return count
}

// Leaf

type Leaf struct {
	Value []byte
}

func (t Leaf) GetLeavesDepth(depth int) []int {
	return []int{depth}
}

func (l Leaf) InsertAt(i uint64, leaf Leaf) (Tree, Tree) {
	if i == 0 {
		return leaf, l
	} else {
		return l, leaf
	}
}

func (t Leaf) GetAt(i uint64) (Tree, Leaf) {
	return t, t
}

func (t Leaf) RemoveAt(i uint64) (Tree, bool) {
	return nil, true // or false, ignored
}

func (t Leaf) Count() uint64 {
	return 1
}

func (t Leaf) CountUnreferencedNodes() int {
	return 1
}

// Persistance

func RefFromTree(t Tree, dserv mdag.DAGService) TreeRef {

	var datatype pb.Tree_DataType

	node := new(mdag.Node)
	message := new(pb.Tree)

	switch s := t.(type) {
	case Node2:
		datatype = pb.Tree_Node2
		node.AddRawLink("0", s.children[0].Persist(dserv).link)
		node.AddRawLink("1", s.children[1].Persist(dserv).link)
	case Node3:
		datatype = pb.Tree_Node3
		node.AddRawLink("0", s.children[0].Persist(dserv).link)
		node.AddRawLink("1", s.children[1].Persist(dserv).link)
		node.AddRawLink("2", s.children[2].Persist(dserv).link)
	case Leaf:
		datatype = pb.Tree_Leaf
		message.Data = s.Value
	}

	message.Type = &datatype
	message.Count = proto.Uint64(t.Count())

	marshalled, _ := proto.Marshal(message)
	node.Data = marshalled

	_, err := dserv.Add(node)
	if err != nil {
		panic(err)
	}

	link, err := mdag.MakeLink(node)

	return TreeRef{
		link:  link,
		dserv: dserv,
	}
}

func (t Node2) Persist(dserv mdag.DAGService) TreeRef {
	return RefFromTree(t, dserv)
}

func (t Node3) Persist(dserv mdag.DAGService) TreeRef {
	return RefFromTree(t, dserv)
}

func (t Leaf) Persist(dserv mdag.DAGService) TreeRef {
	return RefFromTree(t, dserv)
}

// Operations on tree references

type TreeRef struct {
	link  *mdag.Link
	dserv mdag.DAGService
}

func (r TreeRef) Read() Tree {
	node, err := r.link.GetNode(r.dserv)
	if err != nil {
		panic(err)
	}

	unmarshalled := new(pb.Tree)

	err = proto.Unmarshal(node.Data, unmarshalled)
	if err != nil {
		panic(err)
	}

	switch *unmarshalled.Type {
	case pb.Tree_Leaf:
		return NewLeaf(unmarshalled.Data)

	case pb.Tree_Node2:
		child0, err := node.GetNodeLink("0")
		if err != nil {
			panic(err)
		}

		child1, err := node.GetNodeLink("1")
		if err != nil {
			panic(err)
		}

		return NewNode2([]Tree{
			TreeRef{link: child0, dserv: r.dserv},
			TreeRef{link: child1, dserv: r.dserv},
		})

	case pb.Tree_Node3:
		child0, err := node.GetNodeLink("0")
		if err != nil {
			panic(err)
		}

		child1, err := node.GetNodeLink("1")
		if err != nil {
			panic(err)
		}

		child2, err := node.GetNodeLink("2")
		if err != nil {
			panic(err)
		}

		return NewNode3([]Tree{
			TreeRef{link: child0, dserv: r.dserv},
			TreeRef{link: child1, dserv: r.dserv},
			TreeRef{link: child2, dserv: r.dserv},
		})
	}
	panic("unhandled case")
	return nil
}

// all indirected methods

func (r TreeRef) Persist(_ mdag.DAGService) TreeRef {
	// already persisted!
	return r
}

func (r TreeRef) Count() uint64 {
	return r.Read().Count()
}

func (r TreeRef) GetAt(i uint64) (Tree, Leaf) {
	return r.Read().GetAt(i)
}

func (r TreeRef) InsertAt(i uint64, l Leaf) (Tree, Tree) {
	return r.Read().InsertAt(i, l)
}

func (r TreeRef) RemoveAt(i uint64) (Tree, bool) {
	return r.Read().RemoveAt(i)
}

func (r TreeRef) GetLeavesDepth(i int) []int {
	return r.Read().GetLeavesDepth(i)
}

func (r TreeRef) CountUnreferencedNodes() int {
	// already persisted!
	return 1
}
