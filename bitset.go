package bitset

const (
	unitByteSize = 6
	unitMax      = 1<<unitByteSize - 1
)

// BitSet manages a compact array of bit values, which are represented as bool,
// where true indicates that the bit is on (1) and false indicates the bit is off (0).
type BitSet struct {
	values []uint64
	maxBit uint64
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
	if index >= b.maxBit {
		n := index + 1
		b.grow(n)
	}
	b.values[index>>unitByteSize] |= 1 << (index & unitMax)
}

// Get true if index is set 1, or return false.
func (b *BitSet) Get(index uint64) bool {
	if index >= b.maxBit {
		return false
	}
	return (b.values[index>>unitByteSize] & (1 << (index & unitMax))) != 0
}

// Clear sets the bit specified by the index to 0.
func (b *BitSet) Clear(index uint64) {
	if index >= b.maxBit {
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

// NextClearBit return the index of the first bit that is set to false that occurs on or after
// the specified starting index.
func (b *BitSet) NextClearBit(fromIndex uint64) uint64 {
	for index := fromIndex; index < b.maxBit; index++ {
		if (b.values[index>>unitByteSize]>>(index&unitMax))&1 == 0 {
			return index
		}
	}
	return b.maxBit
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
			capacity += size / 4
		} else {
			capacity += size
		}
		v := make([]uint64, size, capacity)
		copy(v, b.values)
		b.values = v
	}
	b.maxBit = size << unitByteSize
}
