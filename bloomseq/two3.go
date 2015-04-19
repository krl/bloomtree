package bloomseq

import (
	proto "code.google.com/p/goprotobuf/proto"
	mdag "github.com/ipfs/go-ipfs/merkledag"
	pb "github.com/krl/bloomtree/bloomseq/pb"
)

// tree stuff

type tree interface {
	insertAt(uint64, leaf) (tree, tree)
	getAt(uint64) (tree, leaf)
	removeAt(uint64) (tree, bool)
	count() uint64
	persist(dserv mdag.DAGService) treeRef

	// for tests only
	getLeavesDepth(int) []int
	countUnreferencedNodes() int
}

// newLeaf

func newLeaf(value []byte) leaf {
	leaf := leaf{Value: value}
	return leaf
}

// node2

type node2 struct {
	children []tree
	m_count  uint64
}

func newNode2(children []tree) node2 {
	node := node2{children: children}
	for i := 0; i < 2; i++ {
		node.m_count += children[i].count()
	}
	return node
}

func (t node2) getLeavesDepth(depth int) []int {
	depths := make([]int, 0, 2)

	for i := 0; i < len(t.children); i++ {
		depths = append(depths, t.children[i].getLeavesDepth(depth+1)...)
	}
	return depths
}

func (n node2) getAt(i uint64) (tree, leaf) {
	var index int = 0

	// find the right child to get from, and decrement i
	for i >= n.children[index].count() {
		i -= n.children[index].count()
		index++
	}

	tr, lf := n.children[index].getAt(i)

	new_children := make([]tree, 2)

	for i := 0; i < 2; i++ {
		new_children[i] = n.children[i]
	}
	new_children[index] = tr

	return newNode2(new_children), lf
}

func (t node2) removeAt(i uint64) (tree, bool) {

	var index int = 0

	// find the right child to remove from, and decrement i
	for i >= t.children[index].count() {
		i -= t.children[index].count()
		index++
	}

	var sibling_index int
	if index == 0 {
		sibling_index = 1
	} else {
		sibling_index = 0
	}

	// we're going deeper down the stack
	result, underflow := t.children[index].removeAt(i)

	if result == nil {
		// we deleted a leaf, and are short one child
		return t.children[sibling_index], true
	} else if underflow {

		// find the appropriate sibling to resolve imbalance
		switch s := t.children[sibling_index].(type) {
		case node3:

			new_sibling_children := make([]tree, 2)
			new_children := make([]tree, 2)

			if sibling_index > index {
				// take the first child
				for i := 1; i < 3; i++ {
					new_sibling_children[i-1] = s.children[i]
				}

				new_children[index] = newNode2([]tree{result, s.children[0]})
			} else {
				// take the last child
				for i := 0; i < 2; i++ {
					new_sibling_children[i] = s.children[i]
				}

				new_children[index] = newNode2([]tree{s.children[2], result})
			}
			new_children[sibling_index] = newNode2(new_sibling_children)
			return newNode2(new_children), false
		case node2:

			new_sibling_children := make([]tree, 3)

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

			return newNode3(new_sibling_children), true
		}
	} else {
		// no underflow, just update index
		new_children := make([]tree, 2)
		copy(new_children, t.children)
		new_children[index] = result
		return newNode2(new_children), false
	}

	return t, false
}

func (t node2) insertAt(i uint64, leaf leaf) (tree, tree) {
	var index int = 0

	// find the right child to insert into, and decrement i
	for i > t.children[index].count() {
		i -= t.children[index].count()
		index++
	}

	// recursive recursion
	result, extra := t.children[index].insertAt(i, leaf)

	// new child array

	if extra == nil {
		new_children2 := make([]tree, 2)

		for i := 0; i < 2; i++ {
			new_children2[i] = t.children[i]
		}

		new_children2[index] = result

		return newNode2(new_children2), nil
	} else {
		new_children3 := make([]tree, 3)

		if index == 0 {
			// insert extra before old value
			new_children3[0] = result
			new_children3[1] = extra
			new_children3[2] = t.children[1]

			return newNode3(new_children3), nil
		} else {
			// insert extra after old value
			new_children3[0] = t.children[0]
			new_children3[1] = result
			new_children3[2] = extra
			return newNode3(new_children3), nil
		}
	}
}

func (t node2) count() uint64 {
	return t.m_count
}

func (t node2) countUnreferencedNodes() int {
	count := 0
	for i := 0; i < 2; i++ {
		count += t.children[i].countUnreferencedNodes()
	}

	return count
}

// node3

type node3 struct {
	children []tree
	m_count  uint64
}

func newNode3(children []tree) node3 {
	node := node3{children: children}
	for i := 0; i < 3; i++ {
		node.m_count = node.m_count + children[i].count()
	}
	return node
}

func (t node3) getLeavesDepth(depth int) []int {
	depths := make([]int, 0)

	for i := 0; i < len(t.children); i++ {
		depths = append(depths, t.children[i].getLeavesDepth(depth+1)...)
	}
	return depths
}

func (n node3) getAt(i uint64) (tree, leaf) {
	var index int = 0

	// find the right child to get from, and decrement i
	for i >= n.children[index].count() {
		i -= n.children[index].count()
		index++
	}

	tr, lf := n.children[index].getAt(i)

	new_children := make([]tree, 3)

	for i := 0; i < 3; i++ {
		new_children[i] = n.children[i]
	}
	new_children[index] = tr

	return newNode3(new_children), lf
}

func (t node3) removeAt(i uint64) (tree, bool) {

	var index int = 0

	// find the right child to remove from, and decrement i
	for i >= t.children[index].count() {
		i -= t.children[index].count()
		index++
	}

	result, underflow := t.children[index].removeAt(i)

	if result == nil {
		// we deleted a leaf, return a two-node

		new_children := make([]tree, 2)

		for i := 0; i < index; i++ {
			new_children[i] = t.children[i]
		}

		// append post
		for i := index; i < 2; i++ {
			new_children[i] = t.children[i+1]
		}

		return newNode2(new_children), false
	} else if underflow {

		// find the appropriate sibling to resolve imbalance
		var sibling_index int
		if index == 0 {
			sibling_index = 1
		} else {
			sibling_index = index - 1
		}

		switch s := t.children[sibling_index].(type) {
		case node3:

			// let's steal a child!

			new_children := make([]tree, 3)

			// copy old values
			for i := 0; i < 3; i++ {
				new_children[i] = t.children[i]
			}

			new_sibling_children := make([]tree, 2)

			if sibling_index > index {

				// take the first child
				for i := 1; i < 3; i++ {
					new_sibling_children[i-1] = s.children[i]
				}
				new_children[index] = newNode2([]tree{result, s.children[0]})
			} else {

				// take the last child
				for i := 0; i < 2; i++ {
					new_sibling_children[i] = s.children[i]
				}
				new_children[index] = newNode2([]tree{s.children[2], result})

			}
			new_children[sibling_index] = newNode2(new_sibling_children)

			return newNode3(new_children), false

		case node2:
			// if neighbour is node2, we shrink this node to a node2

			new_sibling_children := make([]tree, 3)
			new_children := make([]tree, 2)

			if sibling_index > index {
				// put at beginning of siblings children
				for i := 0; i < 2; i++ {
					new_sibling_children[i+1] = s.children[i]
				}
				new_sibling_children[0] = result

				new_children[0] = newNode3(new_sibling_children)
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
					new_children[1] = newNode3(new_sibling_children)
				} else if index == 1 {
					new_children[0] = newNode3(new_sibling_children)
					new_children[1] = t.children[2]
				} else {
					new_children[0] = t.children[0]
					new_children[1] = newNode3(new_sibling_children)
				}
			}
			return newNode2(new_children), false
		}
	} else {
		// no underflow, just update index
		new_children := make([]tree, 3)
		copy(new_children, t.children)
		new_children[index] = result
		return newNode3(new_children), false
	}

	// should not be reached
	panic("unhandled case")
	return t, false
}

func (t node3) insertAt(i uint64, leaf leaf) (tree, tree) {

	var index int = 0

	// find the right child to insert into, and decrement i
	for i > t.children[index].count() {
		i -= t.children[index].count()
		index++
	}

	// deeper
	result, extra := t.children[index].insertAt(i, leaf)

	if extra == nil {
		// without extra

		new_children := make([]tree, 3)
		for i := 0; i < 3; i++ {
			new_children[i] = t.children[i]
		}
		new_children[index] = result

		return newNode3(new_children), nil
	} else {
		// with extra

		new_children := make([]tree, 4)

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
		new_children_a := make([]tree, 2)
		new_children_b := make([]tree, 2)

		for i := 0; i < 2; i++ {
			new_children_a[i] = new_children[i]
		}

		for i := 2; i < 4; i++ {
			new_children_b[i-2] = new_children[i]
		}

		return newNode2(new_children_a), newNode2(new_children_b)
	}
}

func (t node3) count() uint64 {
	return t.m_count
}

func (t node3) countUnreferencedNodes() int {
	count := 0
	for i := 0; i < 3; i++ {
		count += t.children[i].countUnreferencedNodes()
	}
	return count
}

// leaf

type leaf struct {
	Value []byte
}

func (t leaf) getLeavesDepth(depth int) []int {
	return []int{depth}
}

func (l leaf) insertAt(i uint64, leaf leaf) (tree, tree) {
	if i == 0 {
		return leaf, l
	} else {
		return l, leaf
	}
}

func (t leaf) getAt(i uint64) (tree, leaf) {
	return t, t
}

func (t leaf) removeAt(i uint64) (tree, bool) {
	return nil, true // or false, ignored
}

func (t leaf) count() uint64 {
	return 1
}

func (t leaf) countUnreferencedNodes() int {
	return 1
}

// persistance

func refFromTree(t tree, dserv mdag.DAGService) treeRef {

	var datatype pb.Tree_DataType

	node := new(mdag.Node)
	message := new(pb.Tree)

	switch s := t.(type) {
	case node2:
		datatype = pb.Tree_Node2
		node.AddRawLink("0", s.children[0].persist(dserv).link)
		node.AddRawLink("1", s.children[1].persist(dserv).link)
	case node3:
		datatype = pb.Tree_Node3
		node.AddRawLink("0", s.children[0].persist(dserv).link)
		node.AddRawLink("1", s.children[1].persist(dserv).link)
		node.AddRawLink("2", s.children[2].persist(dserv).link)
	case leaf:
		datatype = pb.Tree_Leaf
		message.Data = s.Value
	}

	message.Type = &datatype
	message.Count = proto.Uint64(t.count())

	marshalled, _ := proto.Marshal(message)
	node.Data = marshalled

	_, err := dserv.Add(node)
	if err != nil {
		panic(err)
	}

	link, err := mdag.MakeLink(node)

	return treeRef{
		link:  link,
		dserv: dserv,
	}
}

func (t node2) persist(dserv mdag.DAGService) treeRef {
	return refFromTree(t, dserv)
}

func (t node3) persist(dserv mdag.DAGService) treeRef {
	return refFromTree(t, dserv)
}

func (t leaf) persist(dserv mdag.DAGService) treeRef {
	return refFromTree(t, dserv)
}

// Operations on tree references

type treeRef struct {
	link  *mdag.Link
	dserv mdag.DAGService
}

func (r treeRef) read() tree {
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
		return newLeaf(unmarshalled.Data)

	case pb.Tree_Node2:
		child0, err := node.GetNodeLink("0")
		if err != nil {
			panic(err)
		}

		child1, err := node.GetNodeLink("1")
		if err != nil {
			panic(err)
		}

		return newNode2([]tree{
			treeRef{link: child0, dserv: r.dserv},
			treeRef{link: child1, dserv: r.dserv},
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

		return newNode3([]tree{
			treeRef{link: child0, dserv: r.dserv},
			treeRef{link: child1, dserv: r.dserv},
			treeRef{link: child2, dserv: r.dserv},
		})
	}
	panic("unhandled case")
	return nil
}

// all indirected methods

func (r treeRef) persist(_ mdag.DAGService) treeRef {
	// already persisted!
	return r
}

func (r treeRef) count() uint64 {
	return r.read().count()
}

func (r treeRef) getAt(i uint64) (tree, leaf) {
	return r.read().getAt(i)
}

func (r treeRef) insertAt(i uint64, l leaf) (tree, tree) {
	return r.read().insertAt(i, l)
}

func (r treeRef) removeAt(i uint64) (tree, bool) {
	return r.read().removeAt(i)
}

func (r treeRef) getLeavesDepth(i int) []int {
	return r.read().getLeavesDepth(i)
}

func (r treeRef) countUnreferencedNodes() int {
	// already persisted!
	return 1
}
