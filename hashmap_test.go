package hashmap

import (
	"strconv"
	"testing"
	"unsafe"
)

type Animal struct {
	name string
}

func TestMapCreation(t *testing.T) {
	m := New()
	if m == nil {
		t.Error("map is null.")
	}

	if m.Len() != 0 {
		t.Error("new map should be empty.")
	}
}

func TestOverwrite(t *testing.T) {
	m := New()
	elephant := &Animal{"elephant"}
	monkey := &Animal{"monkey"}

	m.Set(1<<62, unsafe.Pointer(elephant))
	m.Set(1<<62, unsafe.Pointer(monkey))

	if m.Len() != 1 {
		t.Error("map should contain exactly one element.")
	}

	tmp, ok := m.Get(1 << 62) // Retrieve inserted element.
	if ok == false {
		t.Error("ok should be true for item stored within the map.")
	}

	item := (*Animal)(tmp) // Type assertion.
	if item != monkey {
		t.Error("wrong item returned.")
	}
}

func TestInsert(t *testing.T) {
	m := NewSize(4)
	elephant := &Animal{"elephant"}
	monkey := &Animal{"monkey"}

	m.Set(4, unsafe.Pointer(elephant))
	m.Set(3, unsafe.Pointer(elephant))
	m.Set(2, unsafe.Pointer(monkey))
	m.Set(1, unsafe.Pointer(monkey))

	if m.Len() != 4 {
		t.Error("map should contain exactly 4 elements.")
	}
}

func TestGet(t *testing.T) {
	m := New()

	val, ok := m.Get(0) // Get a missing element.
	if ok == true {
		t.Error("ok should be false when item is missing from map.")
	}
	if val != nil {
		t.Error("Missing values should return as nil.")
	}

	elephant := &Animal{"elephant"}
	m.Set(1, unsafe.Pointer(elephant))

	_, ok = m.Get(2)
	if ok == true {
		t.Error("ok should be false when item is missing from map.")
	}

	_, ok = m.Get(0) // Get a missing element.
	if ok == true {
		t.Error("ok should be false when item is missing from map.")
	}

	tmp, ok := m.Get(1) // Retrieve inserted element.
	if ok == false {
		t.Error("ok should be true for item stored within the map.")
	}

	elephant = (*Animal)(tmp) // Type assertion.
	if elephant == nil {
		t.Error("expecting an element, not null.")
	}
	if elephant.name != "elephant" {
		t.Error("item was modified.")
	}
}

func TestResize(t *testing.T) {
	m := NewSize(2)
	itemCount := 16
	log := log2(uint64(itemCount))

	for i := 0; i < itemCount; i++ {
		m.Set(uint64(i)<<(64-log), unsafe.Pointer(&Animal{strconv.Itoa(i)}))
	}

	if m.Len() != uint64(itemCount) {
		t.Error("Expected etelemnt count did not match.")
	}

	// Using keys, the fill rate is less than 50
	if m.Fillrate() > 50 {
		t.Errorf("Expecting 50 or lower percent fillrate. got: %d", m.Fillrate())
	}

	for i := 0; i < itemCount; i++ {
		_, ok := m.Get(uint64(i) << (64 - log))
		if !ok {
			t.Error("Getting inserted item failed.")
		}
	}
}

func TestStringer(t *testing.T) {
	m := New()
	elephant := &Animal{"elephant"}
	monkey := &Animal{"monkey"}

	s := m.String()
	if s != "[]" {
		t.Error("empty map as string does not match.")
	}

	m.Set(0<<62, unsafe.Pointer(elephant))
	s = m.String()
	if s != "[0]" {
		t.Error("1 item map as string does not match.")
	}

	m.Set(1<<62, unsafe.Pointer(monkey))
	s = m.String()
	if s != "[4611686018427387904,0]" {
		t.Error("2 item map as string does not match.")
	}
}

func TestDelete(t *testing.T) {
	m := New()
	m.Del(0)

	elephant := &Animal{"elephant"}
	monkey := &Animal{"monkey"}
	m.Set(1, unsafe.Pointer(elephant))
	m.Set(2, unsafe.Pointer(monkey))
	m.Del(0)
	m.Del(3)
	if m.Len() != 2 {
		t.Error("map should contain exactly two elements.")
	}

	m.Del(1)
	m.Del(1)
	m.Del(2)
	if m.Len() != 0 {
		t.Error("map should be empty.")
	}

	val, ok := m.Get(1) // Get a missing element.
	if ok == true {
		t.Error("ok should be false when item is missing from map.")
	}
	if val != nil {
		t.Error("Missing values should return as nil.")
	}

	m.Set(1, unsafe.Pointer(elephant))
}

func TestIterator(t *testing.T) {
	m := New()
	itemCount := 16
	log := log2(uint64(itemCount))

	for i := itemCount; i > 0; i-- {
		m.Set(uint64(i)<<(64-log), unsafe.Pointer(&Animal{strconv.Itoa(i)}))
	}

	counter := 0
	for item := range m.Iter() {
		val := item.Value
		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != itemCount {
		t.Error("Returned item count did not match.")
	}
}

func TestCompareAndSwap(t *testing.T) {
	m := New()
	elephant := &Animal{"elephant"}
	monkey := &Animal{"monkey"}

	m.Set(1<<62, unsafe.Pointer(elephant))
	if m.Len() != 1 {
		t.Error("map should contain exactly one element.")
	}
	if !m.Cas(1<<62, unsafe.Pointer(elephant), unsafe.Pointer(monkey)) {
		t.Error("Cas should success if expectation met")
	}
	if m.Cas(1<<62, unsafe.Pointer(elephant), unsafe.Pointer(monkey)) {
		t.Error("Cas should fail if expectation didn't meet")
	}
	tmp, ok := m.Get(1 << 62)
	if ok == false {
		t.Error("ok should be true for item stored within the map.")
	}
	item := (*Animal)(tmp)
	if item != monkey {
		t.Error("wrong item returned.")
	}
}
