package common

import (
	"encoding/binary"
	ds "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-datastore"
	"github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-datastore/sync"
	"github.com/ipfs/go-ipfs/blocks/blockstore"
	bs "github.com/ipfs/go-ipfs/blockservice"
	"github.com/ipfs/go-ipfs/exchange/offline"
	mdag "github.com/ipfs/go-ipfs/merkledag"
	"github.com/krl/bloomtree/filter"
	. "github.com/krl/bloomtree/value"
	"strings"
	"testing"
)

func GetMockDagServ(t testing.TB) mdag.DAGService {
	dstore := ds.NewMapDatastore()
	tsds := sync.MutexWrap(dstore)
	bstore := blockstore.NewBlockstore(tsds)
	bserv, err := bs.New(bstore, offline.Exchange(bstore))
	if err != nil {
		t.Fatal(err)
	}
	dserv := mdag.NewDAGService(bserv)
	return dserv
}

func BytesFromInt(i uint64) []byte {
	b := make([]byte, 8)
	binary.PutUvarint(b, i)
	return b
}

func IntFromBytes(b []byte) uint64 {
	res, _ := binary.Uvarint(b)
	return res
}

type TextValue struct {
	Content string
}

func NewTextValue(s string) Value {
	return TextValue{Content: s}
}

func TextFilter(word string) filter.Filter {
	filt := filter.NewFilter(32)
	filt.Add([]byte(word))
	return filter.Filter{
		"words": filt,
	}
}

func CountFilter(i uint64) filter.Filter {
	b := make([]byte, 8)
	binary.PutUvarint(b, i)

	filt := filter.NewFilter(32)
	filt.Add(b)

	return filter.Filter{
		"count": filt,
	}
}

func (t TextValue) GetFilter() filter.Filter {

	wordfilter := filter.NewFilter(32)
	countfilter := filter.NewFilter(32)

	var count uint64 = 0

	for _, word := range strings.Split(t.Content, " ") {
		wordfilter.Add([]byte(word))
		count++
	}

	b := make([]byte, 8)
	binary.PutUvarint(b, count)

	countfilter.Add(b)

	return filter.Filter{
		"words": wordfilter,
		"count": countfilter,
	}
}

func (t TextValue) Serialize() []byte {
	return []byte(t.Content)
}

func DeserializeTextValue(bytes []byte) Value {
	return TextValue{Content: string(bytes)}
}
