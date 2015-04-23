// Code generated by protoc-gen-go.
// source: tree.proto
// DO NOT EDIT!

package bloomset_pb

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type Tree_DataType int32

const (
	Tree_Node Tree_DataType = 1
	Tree_Leaf Tree_DataType = 2
)

var Tree_DataType_name = map[int32]string{
	1: "Node",
	2: "Leaf",
}
var Tree_DataType_value = map[string]int32{
	"Node": 1,
	"Leaf": 2,
}

func (x Tree_DataType) Enum() *Tree_DataType {
	p := new(Tree_DataType)
	*p = x
	return p
}
func (x Tree_DataType) String() string {
	return proto.EnumName(Tree_DataType_name, int32(x))
}
func (x Tree_DataType) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.String())
}
func (x *Tree_DataType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Tree_DataType_value, data, "Tree_DataType")
	if err != nil {
		return err
	}
	*x = Tree_DataType(value)
	return nil
}

type FilterElement struct {
	Name             *string `protobuf:"bytes,1,req" json:"Name,omitempty"`
	BloomFilter      []byte  `protobuf:"bytes,2,req" json:"BloomFilter,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *FilterElement) Reset()         { *m = FilterElement{} }
func (m *FilterElement) String() string { return proto.CompactTextString(m) }
func (*FilterElement) ProtoMessage()    {}

func (m *FilterElement) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *FilterElement) GetBloomFilter() []byte {
	if m != nil {
		return m.BloomFilter
	}
	return nil
}

type Tree struct {
	Type             *Tree_DataType   `protobuf:"varint,1,req,enum=bloomset.pb.Tree_DataType" json:"Type,omitempty"`
	Filter           []*FilterElement `protobuf:"bytes,2,rep" json:"Filter,omitempty"`
	Data             []byte           `protobuf:"bytes,3,opt" json:"Data,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *Tree) Reset()         { *m = Tree{} }
func (m *Tree) String() string { return proto.CompactTextString(m) }
func (*Tree) ProtoMessage()    {}

func (m *Tree) GetType() Tree_DataType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}

func (m *Tree) GetFilter() []*FilterElement {
	if m != nil {
		return m.Filter
	}
	return nil
}

func (m *Tree) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterEnum("bloomset.pb.Tree_DataType", Tree_DataType_name, Tree_DataType_value)
}
