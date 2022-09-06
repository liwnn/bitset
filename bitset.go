package bitset

import "math/bits"

const (
	unitByteSize        = 6
	unitBitsNum         = 1 << unitByteSize
	unitBitsMask uint64 = unitBitsNum - 1
	unitMask     uint64 = 1<<64 - 1
)

// BitSet manages a compact array of bit values, which are represented as bool,
// where true indicates that the bit is on (1) and false indicates the bit is off (0).
type BitSet struct {
	values    []uint64
	onesCount uint64
}

func New() *BitSet {
	return NewSize(1)
}

func NewSize(length uint64) *BitSet {
	n := length >> unitByteSize
	if length&unitBitsMask != 0 {
		n++
	}
	return &BitSet{
		values: make([]uint64, n),
	}
}

// Set index to 1.
func (b *BitSet) Set(index uint64) {
	unitIndex := index >> unitByteSize
	if unitIndex >= uint64(len(b.values)) {
		b.grow(unitIndex + 1)
	}
	x := uint64(1 << (index & unitBitsMask))
	if b.values[unitIndex]&x == 0 {
		b.values[unitIndex] |= x
		b.onesCount++
	}
}

func (b *BitSet) grow(size uint64) {
	if int(size) <= cap(b.values) {
		b.values = b.values[:size]
	} else {
		v := make([]uint64, size, size)
		copy(v, b.values)
		b.values = v
	}
}

// Get true if index is set 1, or return false.
func (b *BitSet) Get(index uint64) bool {
	unitIndex := index >> unitByteSize
	return unitIndex < uint64(len(b.values)) && (b.values[unitIndex]&(1<<(index&unitBitsMask))) != 0
}

// Clear sets the bit specified by the index to 0.
func (b *BitSet) Clear(index uint64) {
	unitIndex := index >> unitByteSize
	if unitIndex >= uint64(len(b.values)) {
		return
	}
	x := b.values[unitIndex] & ^uint64(1<<(index&unitBitsMask))
	if x == b.values[unitIndex] {
		return
	}
	b.values[unitIndex] = x
	b.onesCount--

	i := len(b.values) - 1
	for ; i >= 0 && b.values[i] == 0; i-- {
	}
	b.values = b.values[:i+1]
}

// Reset all bits to 0.
func (b *BitSet) Reset() {
	var r = make([]uint64, 16)
	var a = b.values
	for len(a) > 0 {
		a = a[copy(a, r):]
	}
	b.values = b.values[:0]
	b.onesCount = 0
}

// Size return the number of bits of space actually in use by this BitSet.
func (b BitSet) Size() uint64 {
	return uint64(len(b.values)) << unitByteSize
}

// Length return the "logical size": the index of the highest set bit plus one.
func (b BitSet) Length() int {
	if len(b.values) == 0 {
		return 0
	}
	return len(b.values)*unitBitsNum - bits.LeadingZeros64(b.values[len(b.values)-1])
}

// NextClearBit return the index of the first bit that is set to false that occurs on or after
// the specified starting index.
func (b BitSet) NextClearBit(fromIndex uint64) uint64 {
	index := fromIndex >> unitByteSize
	if index >= uint64(len(b.values)) {
		return fromIndex
	}
	v := b.values[index] | ((1 << (fromIndex & unitBitsMask)) - 1)
	for {
		if v != unitMask {
			return index<<unitByteSize +
				uint64(bits.TrailingZeros64(^v)) // find the first bit that is set to 0
		}
		index++
		if index >= uint64(len(b.values)) {
			return uint64(len(b.values)) << unitByteSize
		}
		v = b.values[index]
	}
}

// Cardinality returns the number of bits set to true.
func (b BitSet) Cardinality() uint64 {
	return b.onesCount
}
