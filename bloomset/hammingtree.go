package bloomset

import (
	"github.com/krl/bloomtree/filter"
	. "github.com/krl/bloomtree/value"
)

type tree interface {
	insert(leaf) tree
	getFilter() filter.Filter

	find(filter.Filter, chan Value)

	// Get() (tree, Value)
	// Count() uint64
	// Persist(dserv mdag.DAGService) treeRef

	// for tests only
	getLeavesDepth(int) []int
	// CountUnreferencednodes() int
}

type node struct {
	children [2]tree
	m_filter filter.Filter
}

type leaf struct {
	value Value
}

func newNode(c1 tree, c2 tree) tree {
	return node{
		children: [2]tree{c1, c2},
		m_filter: c1.getFilter().Merge(c2.getFilter()),
	}
}

func (l1 leaf) insert(l2 leaf) tree {
	// TODO, check for duplicates
	return newNode(l1, l2)
}

func (l leaf) getFilter() filter.Filter {
	return l.value.GetFilter()
}

func (l leaf) find(fs filter.Filter, c chan Value) {
	// TODO check for false positives
	c <- l.value
}

// nodes

func (n node) insert(l leaf) tree {

	// find the child node with the lowest hamming distance
	// and insert in that node
	// if the inserted leaf lacks a specific filter
	// this does not count towards the potential conflict

	leaffilter := l.getFilter()

	child0dist := n.children[0].getFilter().HammingDistance(leaffilter)
	child1dist := n.children[1].getFilter().HammingDistance(leaffilter)

	var insert int

	if child0dist == child1dist {
		if child0dist%2 == 0 {
			insert = 0
		} else {
			insert = 1
		}
	} else if child0dist < child1dist {
		insert = 0
	} else {
		insert = 1
	}

	if insert == 0 {
		return newNode(n.children[0].insert(l), n.children[1])
	} else {
		return newNode(n.children[0], n.children[1].insert(l))
	}
}

func (n node) getFilter() filter.Filter {
	return n.m_filter
}

func (n node) find(fs filter.Filter, c chan Value) {
	for i := 0; i < 2; i++ {
		if n.children[i].getFilter().MayContain(fs) {
			n.children[i].find(fs, c)
		}
	}
}

// test functions

func (t node) getLeavesDepth(depth int) []int {
	depths := make([]int, 0, 2)

	for i := 0; i < 2; i++ {
		depths = append(depths, t.children[i].getLeavesDepth(depth+1)...)
	}
	return depths
}

func (t leaf) getLeavesDepth(depth int) []int {
	return []int{depth}
}
