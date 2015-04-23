package value

import (
	"github.com/krl/bloomtree/filter"
)

// a value is something that
// * serializes to []byte
// * correspond to a set of bloomfilters

type Value interface {
	Serialize() []byte
	GetFilter() filter.Filter
}
