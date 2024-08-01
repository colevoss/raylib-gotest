package engine

import (
	"testing"
)

func TestSparseSetAdd(t *testing.T) {
	ss := NewSparseSet(10)

	ss.Add(1)

	if !ss.Has(1) {
		t.Errorf("Expected sparse set to have numer that was just added")
	}

	if ss.Has(0) {
		t.Errorf("Should not have 0")
	}
}

func TestAddGrow(t *testing.T) {
	ss := NewSparseSet(10)

	ss.AddGrow(20)

	if !ss.Has(20) {
		t.Errorf("Expected sparse set to have 20 after it grew")
	}

	if ss.Max() != 21 {
		t.Errorf("Expected sparse set to grow to 21")
	}
}

func TestSparseSetAddMany(t *testing.T) {
	ss := NewSparseSet(100)

	for i := range 10 {
		ss.Add(uint(i))
	}

	for i := range 10 {
		if !ss.Has(uint(i)) {
			t.Errorf("Expected sparse set to have %d that was just added", i)
		}
	}
}

func TestSparseSetHas(t *testing.T) {
	ss := NewSparseSet(11)

	ss.Add(10)

	if ss.Has(9) {
		t.Errorf("Expected SparseSet.Has to return false")
	}

	if !ss.Has(10) {
		t.Errorf("Expected SparseSet.has to return true")
	}
}

func TestSparseSetN(t *testing.T) {
	ss := NewSparseSet(11)

	ss.Add(10)
	ss.Add(1)
	ss.Add(2)

	if ss.N() != 3 {
		t.Errorf("expected N() to return 3, got %d", ss.N())
	}
}

func TestSparseRemove(t *testing.T) {
	ss := NewSparseSet(10)

	ss.Add(0)
	ss.Add(1)
	ss.Add(2)

	ss.Remove(0)

	if ss.N() != 2 {
		t.Errorf("Expected N to be 2 after removing")
	}

	if ss.Has(0) {
		t.Errorf("Expected sparse set to not have removed element")
	}
}

func TestSparseSetGrowByFactor(t *testing.T) {
	ss := NewSparseSet(2)

	ss.Add(0)
	ss.Add(1)

	ss.GrowByFactor(2)

	if ss.Max() != 4 {
		t.Errorf("Expexted new max to be 4. Got %d", ss.Max())
	}
}

func TestSparseSetGrow(t *testing.T) {
	ss := NewSparseSet(2)

	ss.Add(0)
	ss.Add(1)

	ss.Grow(20)

	if ss.Max() != 20 {
		t.Errorf("Expexted new max to be 4. Got %d", ss.Max())
	}
}
