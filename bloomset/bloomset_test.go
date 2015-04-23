package bloomset

import (
	"fmt"
	"github.com/krl/bloomtree/filter"
	"testing"

	. "github.com/krl/bloomtree/common"
)

func TestSingletonTree(t *testing.T) {
	set := NewBloomSet(DeserializeTextValue)

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

	set := NewBloomSet(DeserializeTextValue)

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

	Content := (<-result1).(TextValue).Content

	if Content != "one two three four" {
		t.Fatalf("Should have found one two three four, found \"%v\"", Content)
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
	set := NewBloomSet(DeserializeTextValue)

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

func TestEmptyFilter(t *testing.T) {

	set := NewBloomSet(DeserializeTextValue)

	for i := 0; i < 10; i++ {
		set = set.Insert(NewTextValue(fmt.Sprintf("entry #%v", i)))
	}

	result := set.Find(filter.EmptyFilter())

	count := 0
	for v := range result {
		v = v // eeh, there should be a better way to do this
		count++
	}

	if count != 10 {
		t.Fatal("Emtpy filter should match everything")
	}
}

func TestHaystack(t *testing.T) {

	set := NewBloomSet(DeserializeTextValue)

	for i := 0; i < 1000; i++ {
		set = set.Insert(NewTextValue(fmt.Sprintf("haystrand #%v", i)))
	}

	set = set.Insert(NewTextValue("needle"))
	result := set.Find(TextFilter("needle"))

	Content := (<-result).(TextValue).Content

	if Content != "needle" {
		t.Fatal("did not find the needle!")
	}
}

func TestHaystackRemoving(t *testing.T) {

	set := NewBloomSet(DeserializeTextValue)

	count := 1000

	for i := 0; i < count; i++ {
		set = set.Insert(NewTextValue(fmt.Sprintf("haystrand #%v", i)))
	}

	set = set.Insert(NewTextValue("needle"))

	for i := 0; i < count; i++ {
		set = set.Remove(NewTextValue(fmt.Sprintf("haystrand #%v", i)))
	}

	result := set.Find(filter.EmptyFilter())

	Content := (<-result).(TextValue).Content

	if Content != "needle" {
		fmt.Println(Content)
		t.Fatal("did not find the needle!")
	}
}

// persistance test

func TestPersistSingletonRoot(t *testing.T) {
	dserv := GetMockDagServ(t)

	set := NewBloomSet(DeserializeTextValue)

	val := NewTextValue("wonk")
	set = set.Insert(val)

	persisted := set.Persist(dserv)

	result := persisted.Find(filter.EmptyFilter())

	if <-result != val {
		t.Fatal("Should have found the value in persisted tree")
	}
}

func TestPersistHaystack(t *testing.T) {

	dserv := GetMockDagServ(t)

	set := NewBloomSet(DeserializeTextValue)

	count := 1000

	for i := 0; i < count; i++ {
		set = set.Insert(NewTextValue(fmt.Sprintf("haystrand #%v", i)))
	}

	set = set.Insert(NewTextValue("needle"))

	persisted := set.Persist(dserv)

	if persisted.CountUnreferencedNodes() != 1 {
		t.Fatal("dereference fail")
	}

	result := persisted.Find(TextFilter("needle"))

	Content := (<-result).(TextValue).Content

	if Content != "needle" {
		t.Fatal("did not find the needle!")
	}

	unreffed := persisted.CountUnreferencedNodes()

	if unreffed != 12 {
		t.Fatalf("should have dereferenced log(n) nodes (got %v)", unreffed)
	}
}
