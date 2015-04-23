package bloomset

import (
	mdag "github.com/ipfs/go-ipfs/merkledag"
	"github.com/krl/bloomtree/filter"
	. "github.com/krl/bloomtree/value"
)

type BloomSet struct {
	value   tree
	valfunc func([]byte) Value
}

func NewBloomSet(valfunc func([]byte) Value) BloomSet {
	return BloomSet{valfunc: valfunc}
}

func (s BloomSet) Insert(v Value) BloomSet {
	lf := leaf{
		bytes:  v.Serialize(),
		filter: v.GetFilter(),
	}

	if s.value == nil {
		return BloomSet{
			value:   lf,
			valfunc: s.valfunc,
		}
	}

	return BloomSet{
		value:   s.value.insert(lf),
		valfunc: s.valfunc,
	}
}

func (s BloomSet) Remove(v Value) BloomSet {
	leaf := leaf{
		bytes:  v.Serialize(),
		filter: v.GetFilter(),
	}

	if s.value == nil {
		return BloomSet{}
	}
	value, _ := s.value.remove(leaf)
	return BloomSet{
		value:   value,
		valfunc: s.valfunc,
	}
}

// we pass a pointer here, so that the tree is mutated as
// new values are pulled into cache

func (s *BloomSet) Find(f filter.Filter) <-chan Value {
	bytechan := make(chan []byte)
	valuechan := make(chan Value)

	go func() {
		for bytes := range bytechan {
			valuechan <- s.valfunc(bytes)
		}
		close(valuechan)
	}()

	go func() {
		if s.value.getFilter().MayContain(f) {
			s.value = s.value.find(f, bytechan)
		}
		close(bytechan)
	}()

	return valuechan
}

func (s BloomSet) GetLeavesDepth() []int {
	if s.value == nil {
		return []int{}
	}
	return s.value.getLeavesDepth(0)
}

func (s BloomSet) Persist(dserv mdag.DAGService) BloomSet {
	if s.value != nil {
		return BloomSet{
			value:   s.value.persist(dserv),
			valfunc: s.valfunc,
		}
	} else {
		return BloomSet{}
	}
}

func (r BloomSet) CountUnreferencedNodes() int {
	if r.value != nil {
		return r.value.countUnreferencedNodes()
	}
	return 0
}
