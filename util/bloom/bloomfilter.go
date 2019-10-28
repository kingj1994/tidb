package bloom

import (
	"fmt"
	"hash/fnv"
)

// Filter a simple abstraction of bloom filter
type Filter struct {
	bitSet   []uint64
	length   uint64
	unitSize uint64
}

// NewFilter return a bloom filter with a empty bitset of given length
func NewFilter(length int) (*Filter, error) {
	if length <= 0 {
		return nil, fmt.Errorf("length is not positive")
	}
	bitset := make([]uint64, length)
	bits := uint64(64)
	return &Filter{
		bitSet:   bitset,
		length:   bits * uint64(length),
		unitSize: bits,
	}, nil
}

// NewFilterBySlice create a bloom filter by the given slice
func NewFilterBySlice(bs []uint64) (*Filter, error) {
	if len(bs) == 0 {
		return nil, fmt.Errorf("len(bs) == 0")
	}

	bits := uint64(64)
	return &Filter{
		bitSet:   bs,
		length:   bits * uint64(len(bs)),
		unitSize: bits,
	}, nil
}

// Insert a key into the filter
func (bf *Filter) Insert(key []byte) {
	idx, shift := bf.hash(key)
	bf.bitSet[idx] |= 1 << shift
}

// Probe check whether the given key is in the filter
func (bf *Filter) Probe(key []byte) bool {
	idx, shift := bf.hash(key)

	return bf.bitSet[idx]&(1<<shift) != 0
}

func (bf *Filter) hash(key []byte) (uint64, uint64) {
	hash := ihash(key) % uint64(bf.length)
	idx := hash / bf.unitSize
	shift := hash % bf.unitSize

	return idx, shift
}

func ihash(key []byte) uint64 {
	h := fnv.New64a()
	h.Write(key)
	return h.Sum64()
}