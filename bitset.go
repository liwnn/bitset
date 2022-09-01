package bitset

import "math/bits"

const (
	unitByteSize        = 6
	unitMax      uint64 = 1<<unitByteSize - 1
	uint64Mask   uint64 = 1<<64 - 1
)

// BitSet manages a compact array of bit values, which are represented as bool,
// where true indicates that the bit is on (1) and false indicates the bit is off (0).
type BitSet struct {
	values []uint64
	size   uint64
}

func New() *BitSet {
	return NewSize(1)
}

func NewSize(length uint64) *BitSet {
	b := &BitSet{}
	b.grow(length)
	return b
}

// Set index to 1.
func (b *BitSet) Set(index uint64) {
	if index >= b.size {
		b.grow(index + 1)
	}
	b.values[index>>unitByteSize] |= 1 << (index & unitMax)
}

// Get true if index is set 1, or return false.
func (b *BitSet) Get(index uint64) bool {
	if index >= b.size {
		return false
	}
	return (b.values[index>>unitByteSize] & (1 << (index & unitMax))) != 0
}

// Clear sets the bit specified by the index to 0.
func (b *BitSet) Clear(index uint64) {
	if index >= b.size {
		return
	}
	b.values[index>>unitByteSize] &^= 1 << (index & unitMax)
}

// Reset all bits to 0.
func (b *BitSet) Reset() {
	var r = make([]uint64, 16)
	var a = b.values
	for len(a) > 0 {
		a = a[copy(a, r):]
	}
}

// Size return the number of bits of space actually in use by this BitSet.
func (b *BitSet) Size() uint64 {
	return b.size
}

// NextClearBit return the index of the first bit that is set to false that occurs on or after
// the specified starting index.
func (b *BitSet) NextClearBit(fromIndex uint64) uint64 {
	begin := fromIndex >> unitByteSize
	for index := begin; index < uint64(len(b.values)); index++ {
		v := b.values[index]
		if index == begin { // first unit
			v |= (1 << (fromIndex & unitMax)) - 1
		}
		if v^uint64Mask != 0 {
			return index<<unitByteSize +
				uint64(bits.TrailingZeros64(^v)) // find the first bit that is set to 0
		}
	}
	return b.size
}

func (b *BitSet) grow(n uint64) {
	size := n >> unitByteSize
	if n&unitMax > 0 {
		size++
	}

	if b.values == nil {
		b.values = make([]uint64, size)
	} else if int(size) <= cap(b.values) {
		b.values = b.values[:size]
	} else {
		capacity := size
		if size >= 1024 {
			capacity += size >> 2
		} else {
			capacity <<= 1
		}
		v := make([]uint64, size, capacity)
		copy(v, b.values)
		b.values = v
	}
	b.size = size << unitByteSize
}
