package engine

import (
	"fmt"
	"raylib-gotest/assert"
)

type SparseSet struct {
	max    uint
	n      uint
	dense  []uint
	sparse []uint
}

func NewSparseSet(size uint) *SparseSet {
	return &SparseSet{
		n:      0,
		max:    size,
		dense:  make([]uint, size),
		sparse: make([]uint, size),
	}
}

func (ss *SparseSet) N() uint {
	return ss.n
}

func (ss *SparseSet) Max() uint {
	return ss.max
}

func (ss *SparseSet) Add(id uint) {
	assert.Assert(
		id <= ss.max,
		fmt.Sprintf("Attempted to add id %d to sparse set with size %d", id, ss.max),
	)

	ss.dense[ss.n] = id
	ss.sparse[id] = ss.n

	ss.n += 1
}

func (ss *SparseSet) AddGrow(id uint) {
	if id > ss.max {
		ss.Grow(id + 1)
	}

	ss.Add(id)
}

func (ss *SparseSet) Get(denseIndex uint) uint {
	return ss.dense[denseIndex]
}

func (ss *SparseSet) Has(id uint) bool {
	denseIndex := ss.sparse[id]

	return denseIndex < ss.n && ss.dense[denseIndex] == id
}

func (ss *SparseSet) Remove(id uint) {
	if !ss.Has(id) {
		return
	}

	// shift the current n back on
	ss.n -= 1

	// get the dense index of the item to remove
	denseIndex := ss.sparse[id]

	// get the new item to put in its place
	item := ss.dense[ss.n]

	// replace id with last item
	ss.dense[denseIndex] = item

	// set the new index for the replacing item in sparse
	ss.sparse[item] = denseIndex
}

func (ss *SparseSet) GrowByFactor(factor uint) {
	newMax := ss.max * factor
	ss.Grow(newMax)
}

func (ss *SparseSet) Grow(newMax uint) {
	newSparse := make([]uint, newMax)
	newDense := make([]uint, newMax)

	for i := range ss.max {
		newSparse[i] = ss.sparse[i]
		newDense[i] = ss.dense[i]
	}

	ss.dense = newDense
	ss.sparse = newSparse
	ss.max = newMax
}

// func (ss *SparseSet) Iter() func(yield func(uint) bool) {
// 	fmt.Printf("ss: %v\n", ss)
// 	return func(yield func(uint) bool) {
// 		for i := 0; i < len(ss.dense); i++ {
// 			if !yield(ss.sparse[i]) {
// 				fmt.Printf("i: %v\n", i)
// 				return
// 			}
// 		}
// 	}
// }
