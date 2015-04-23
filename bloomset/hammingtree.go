package bloomset

import (
	"bytes"
	proto "code.google.com/p/goprotobuf/proto"
	context "github.com/ipfs/go-ipfs/Godeps/_workspace/src/golang.org/x/net/context"
	mdag "github.com/ipfs/go-ipfs/merkledag"
	pb "github.com/krl/bloomtree/bloomset/pb"
	"github.com/krl/bloomtree/filter"
)

type tree interface {
	insert(leaf) tree
	remove(leaf) (tree, bool)
	getFilter() filter.Filter
	find(filter.Filter, chan []byte) tree
	persist(dserv mdag.DAGService) treeRef

	// for tests only
	getLeavesDepth(int) []int
	countUnreferencedNodes() int
}

type node struct {
	children [2]tree
	filter   filter.Filter
}

type leaf struct {
	bytes  []byte
	filter filter.Filter
}

func newNode(c1 tree, c2 tree) tree {
	return node{
		children: [2]tree{c1, c2},
		filter:   c1.getFilter().Merge(c2.getFilter()),
	}
}

func (l1 leaf) insert(l2 leaf) tree {
	if bytes.Equal(l1.bytes, l2.bytes) {
		return l1 // store no duplicates
	}
	return newNode(l1, l2)
}

func (l1 leaf) remove(l2 leaf) (tree, bool) {
	// false positives are still possible
	if bytes.Equal(l1.bytes, l2.bytes) {
		return nil, true
	}
	return l1, false
}

func (l leaf) getFilter() filter.Filter {
	return l.filter
}

func (l leaf) find(fs filter.Filter, c chan []byte) tree {
	// TODO check for false positives
	c <- l.bytes
	return l
}

// leaf test functions

func (l leaf) getLeavesDepth(depth int) []int {
	return []int{depth}
}

func (l leaf) countUnreferencedNodes() int {
	return 1
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

func (n node) remove(l leaf) (tree, bool) {
	lfilt := l.getFilter()

	for i := 0; i < 2; i++ {
		var res tree
		var success bool

		if n.children[i].getFilter().MayContain(lfilt) {
			res, success = n.children[i].remove(l)
		}

		if success && res == nil {
			// we deleted a leaf. Return the other one
			return n.children[(i+1)%2], true
		} else if success {
			// leaf was deleted further down
			// NB child order does not matter
			return newNode(n.children[(i+1)%2], res), true
		}
	}

	// if unsuccessful, return self
	return n, false
}

func (n node) getFilter() filter.Filter {
	return n.filter
}

func (n node) find(fs filter.Filter, c chan []byte) tree {
	for i := 0; i < 2; i++ {
		if n.children[i].getFilter().MayContain(fs) {
			n.children[i] = n.children[i].find(fs, c)
		}
	}
	return n
}

// test functions

func (t node) getLeavesDepth(depth int) []int {
	depths := make([]int, 0, 2)

	for i := 0; i < 2; i++ {
		depths = append(depths, t.children[i].getLeavesDepth(depth+1)...)
	}
	return depths
}

func (n node) countUnreferencedNodes() int {
	count := 0
	for i := 0; i < 2; i++ {
		count += n.children[i].countUnreferencedNodes()
	}

	return count
}

// persistance

func refFromTree(t tree, dserv mdag.DAGService) treeRef {

	var datatype pb.Tree_DataType

	mdagnode := new(mdag.Node)
	message := new(pb.Tree)

	switch s := t.(type) {
	case node:
		datatype = pb.Tree_Node
		mdagnode.AddRawLink("0", s.children[0].persist(dserv).link)
		mdagnode.AddRawLink("1", s.children[1].persist(dserv).link)
	case leaf:
		datatype = pb.Tree_Leaf
		message.Data = s.bytes
	}

	filtermap := t.getFilter()

	filter := make([]*pb.FilterElement, 0, len(filtermap))

	for k, v := range filtermap {
		name := k // need to provide unchanging pointer
		f := &pb.FilterElement{}
		f.Name = &name
		f.BloomFilter = v.GetBytes()
		filter = append(filter, f)
	}

	message.Filter = filter
	message.Type = &datatype

	marshalled, _ := proto.Marshal(message)
	mdagnode.Data = marshalled

	_, err := dserv.Add(mdagnode)
	if err != nil {
		panic(err)
	}

	link, err := mdag.MakeLink(mdagnode)

	return treeRef{
		link:  link,
		dserv: dserv,
	}
}

func (n node) persist(dserv mdag.DAGService) treeRef {
	return refFromTree(n, dserv)
}

func (l leaf) persist(dserv mdag.DAGService) treeRef {
	return refFromTree(l, dserv)
}

// Operations on tree references

type treeRef struct {
	link  *mdag.Link
	dserv mdag.DAGService
}

func (r treeRef) read() tree {
	mdagnode, err := r.link.GetNode(context.Background(), r.dserv)
	if err != nil {
		panic(err)
	}

	unmarshalled := new(pb.Tree)

	err = proto.Unmarshal(mdagnode.Data, unmarshalled)
	if err != nil {
		panic(err)
	}

	// both types have filters
	filter := filter.EmptyFilter()

	for _, v := range unmarshalled.Filter {
		filter = filter.AddBloom(*v.Name, v.BloomFilter)
	}

	// switch on the rest

	switch *unmarshalled.Type {
	case pb.Tree_Leaf:
		return leaf{
			bytes:  unmarshalled.Data,
			filter: filter,
		}

	case pb.Tree_Node:
		child0, err := mdagnode.GetNodeLink("0")
		if err != nil {
			panic(err)
		}

		child1, err := mdagnode.GetNodeLink("1")
		if err != nil {
			panic(err)
		}

		return node{
			children: [2]tree{
				treeRef{link: child0, dserv: r.dserv},
				treeRef{link: child1, dserv: r.dserv},
			},
			filter: filter,
		}
	}
	panic("unhandled case")
	return nil
}

// all indirected methods

func (r treeRef) persist(_ mdag.DAGService) treeRef {
	// already persisted!
	return r
}

func (r treeRef) insert(l leaf) tree {
	return r.read().insert(l)
}

func (r treeRef) remove(l leaf) (tree, bool) {
	return r.read().remove(l)
}

func (r treeRef) find(f filter.Filter, c chan []byte) tree {
	return r.read().find(f, c)
}

func (r treeRef) getFilter() filter.Filter {
	return r.read().getFilter()
}

func (r treeRef) getLeavesDepth(i int) []int {
	return r.read().getLeavesDepth(i)
}

func (r treeRef) countUnreferencedNodes() int {
	// already persisted!
	return 1
}
