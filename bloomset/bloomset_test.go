package bloomset

import (
	"encoding/binary"
	"fmt"
	"github.com/krl/bloomtree/filter"
	. "github.com/krl/bloomtree/value"
	"strings"
	"testing"
)

type TextValue struct {
	content string
}

func NewTextValue(s string) Value {
	return TextValue{content: s}
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

	for _, word := range strings.Split(t.content, " ") {
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

func TestSingletonTree(t *testing.T) {
	set := NewBloomSet()

	val1 := NewTextValue("wonk")
	val2 := NewTextValue("donk")

	set = set.Insert(val1)

	result1 := set.Find(val1.GetFilter())

	if <-result1 != val1 {
		t.Fatal("Should have found the value")
	}

	result2 := set.Find(val2.GetFilter())

	if <-result2 != nil {
		t.Fatal("Should not have found anything")
	}
}

func TestQueries(t *testing.T) {

	set := NewBloomSet()

	set = set.Insert(NewTextValue("one"))
	set = set.Insert(NewTextValue("one two"))
	set = set.Insert(NewTextValue("one two three"))
	set = set.Insert(NewTextValue("one two three four"))

	set = set.Insert(NewTextValue("ett"))
	set = set.Insert(NewTextValue("ett två"))
	set = set.Insert(NewTextValue("ett två tre"))
	set = set.Insert(NewTextValue("ett två tre fyra"))

	set = set.Insert(NewTextValue("eins"))
	set = set.Insert(NewTextValue("eins zwei"))
	set = set.Insert(NewTextValue("eins zwei drei"))
	set = set.Insert(NewTextValue("eins zwei drei vier"))

	// should get one result
	result1 := set.Find(TextFilter("four"))

	content := (<-result1).(TextValue).content

	if content != "one two three four" {
		t.Fatalf("Should have found one two three four, found \"%v\"", content)
	}

	// multiple results

	result2 := set.Find(TextFilter("två"))

	test := []TextValue{}

	for v := range result2 {
		test = append(test, v.(TextValue))
	}

	if len(test) != 3 {
		t.Fatalf("Should have returned 3, got %v", len(test))
	}

	fmt.Println("with två")
	fmt.Println(test)

	result3 := set.Find(TextFilter("eins"))

	test = []TextValue{}

	for v := range result3 {
		test = append(test, v.(TextValue))
	}

	if len(test) != 4 {
		t.Fatalf("Should have returned 4, got %v", len(test))
	}

	fmt.Println("with eins")
	fmt.Println(test)

	// word count

	result4 := set.Find(CountFilter(3))

	test = []TextValue{}

	for v := range result4 {
		test = append(test, v.(TextValue))
	}

	if len(test) != 3 {
		t.Fatalf("Should have returned 3, got %v", len(test))
	}

	fmt.Println("with wordcount 3")
	fmt.Println(test)

}

func TestReasonableBalance(t *testing.T) {
	set := NewBloomSet()

	for i := 0; i < 1000; i++ {
		set = set.Insert(NewTextValue(fmt.Sprintf("element #%v", i)))
	}

	depths := set.GetLeavesDepth()

	max := 0
	min := 10000

	for _, v := range depths {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}

	// semi-arbitrary definition of balance
	if (max - min) > min/2 {
		t.Fatalf("Tree is not very well balanced!")
	}
}

func TestHaystack(t *testing.T) {

	set := NewBloomSet()

	for i := 0; i < 10000; i++ {
		set = set.Insert(NewTextValue(fmt.Sprintf("haystrand #%v", i)))
	}

	set = set.Insert(NewTextValue("needle"))
	result := set.Find(TextFilter("needle"))

	content := (<-result).(TextValue).content

	if content != "needle" {
		t.Fatal("did not find the needle!")
	}
}
