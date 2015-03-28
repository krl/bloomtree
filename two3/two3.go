package two3

import (
// "fmt"
)

// tree stuff

type Tree interface {
	InsertAt(int, Leaf) (Tree, Tree)
	GetAt(int) Leaf
	RemoveAt(int) (Tree, bool)
	Count() int
	GetLeavesDepth(int) []int
}

// NewLeaf

func NewLeaf(pointer string) Leaf {
	leaf := Leaf{Pointer: pointer}
	return leaf
}

// Node2

type Node2 struct {
	children []Tree
	m_count  int
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

func (t Node2) GetAt(i int) Leaf {
	var index int = 0

	// find the right child to get from, and decrement i
	for i >= t.children[index].Count() {
		i -= t.children[index].Count()
		index++
	}

	return t.children[index].GetAt(i)
}

func (t Node2) RemoveAt(i int) (Tree, bool) {

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

func (t Node2) InsertAt(i int, leaf Leaf) (Tree, Tree) {
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

func (t Node2) Count() int {
	return t.m_count
}

// Node3

type Node3 struct {
	children []Tree
	m_count  int
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

func (t Node3) GetAt(i int) Leaf {
	var index int = 0

	// find the right child to get from, and decrement i
	for i >= t.children[index].Count() {
		i -= t.children[index].Count()
		index++
	}

	return t.children[index].GetAt(i)
}

func (t Node3) RemoveAt(i int) (Tree, bool) {

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

func (t Node3) InsertAt(i int, leaf Leaf) (Tree, Tree) {

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

func (t Node3) Count() int {
	return t.m_count
}

// Leaf

type Leaf struct {
	Pointer string
}

func (t Leaf) GetLeavesDepth(depth int) []int {
	return []int{depth}
}

func (l Leaf) InsertAt(i int, leaf Leaf) (Tree, Tree) {
	if i == 0 {
		return leaf, l
	} else {
		return l, leaf
	}
}

func (t Leaf) GetAt(i int) Leaf {
	return t
}

func (t Leaf) RemoveAt(i int) (Tree, bool) {
	return nil, true // or false, ignored
}

func (t Leaf) Count() int {
	return 1
}
