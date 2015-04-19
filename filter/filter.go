package filter

import (
	bf "github.com/ipfs/go-ipfs/blocks/bloom"
)

type Filter map[string]bf.Filter

func NewFilter(size int) bf.Filter {
	return bf.NewFilter(size)
}

func (fs1 Filter) Merge(fs2 Filter) Filter {

	newfilt := Filter{}

	for k, v := range fs1 {
		newfilt[k] = v
	}

	for k, v := range fs2 {
		if newfilt[k] != nil {
			merged, err := newfilt[k].Merge(fs2[k])
			if err != nil {
				panic(err)
			}
			newfilt[k] = merged
		} else {
			newfilt[k] = v
		}
	}
	return newfilt
}

func (f1 Filter) HammingDistance(f2 Filter) int {
	acc := 0

	for k := range f1 {
		if f2[k] != nil {
			dist, _ := f1[k].HammingDistance(f2[k])
			acc += dist
		}
	}
	return acc
}

func (bigger Filter) MayContain(smaller Filter) bool {
	for k, _ := range smaller {
		may, _ := bigger[k].MayContain(smaller[k])
		if !may {
			return false
		}
	}
	return true
}
