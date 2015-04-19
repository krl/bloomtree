package bloomset

import (
	"github.com/krl/bloomtree/filter"
	. "github.com/krl/bloomtree/value"
)

type BloomSet struct {
	value tree
}

func NewBloomSet() BloomSet {
	return BloomSet{}
}

func (s BloomSet) Insert(v Value) BloomSet {
	leaf := leaf{value: v}

	if s.value == nil {
		return BloomSet{value: leaf}
	}
	return BloomSet{value: s.value.insert(leaf)}
}

func (s BloomSet) Remove(v Value) BloomSet {
	leaf := leaf{value: v}

	if s.value == nil {
		return BloomSet{}
	}
	value, _ := s.value.remove(leaf)
	return BloomSet{value: value}
}

func (s BloomSet) Find(f filter.Filter) <-chan Value {
	c := make(chan Value)

	go func() {
		if s.value.getFilter().MayContain(f) {
			s.value.find(f, c)
		}
		close(c)
	}()

	return c
}

func (s BloomSet) GetLeavesDepth() []int {
	if s.value == nil {
		return []int{}
	}
	return s.value.getLeavesDepth(0)
}
