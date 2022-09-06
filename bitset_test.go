package bitset

import (
	"math/rand"
	"testing"
)

func perm(n int) []uint64 {
	var l []uint64
	for _, v := range rand.Perm(n) {
		l = append(l, uint64(v))
	}
	return l
}

func TestBitSet(t *testing.T) {
	const size = 10000
	b := NewSize(size)

	b.Set(size - 1)
	b.Clear(size - 1)
	if b.Size() != 0 {
		t.Error("Clear")
	}

	for _, v := range perm(size) {
		b.Set(v)
	}
	for _, v := range perm(size) {
		if !b.Get(v) {
			t.Error("Set bit error")
		}
	}

	for _, v := range perm(size / 2) {
		b.Clear(v)
	}

	for i := 0; i < size/2; i++ {
		if b.Get(uint64(i)) {
			t.Error("Clear bit error")
		}
	}

	for i := size / 2; i < size; i++ {
		if !b.Get(uint64(i)) {
			t.Error("Clear bit error")
		}
	}
}

func TestBitSet_NextClearBit(t *testing.T) {
	bs := New()
	bs.Set(1)
	bs.Set(2)

	if i := bs.NextClearBit(0); i != 0 {
		t.Errorf("NextClearBit(0) = %d, want 1", i)
	}

	if i := bs.NextClearBit(1); i != 3 {
		t.Errorf("NextClearBit(1) = %d, want 3", i)
	}

	for i := 0; i < 100; i++ {
		bs.Set(uint64(i))
	}
	if i := bs.NextClearBit(0); i != 100 {
		t.Errorf("NextClearBit(0) = %d, want 100", i)
	}
}

var N = 1000000

func newBitSet() *BitSet {
	s := NewSize(uint64(N))
	for i := 0; i < len(s.values); i++ {
		for j := 0; j < 8; j++ {
			s.Set(uint64(i>>unitByteSize + j))
		}
	}
	return s
}

func BenchmarkSet(b *testing.B) {
	s := newBitSet()
	n := perm(N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Set(n[i%N])
	}
}

func BenchmarkGet(b *testing.B) {
	s := newBitSet()
	n := perm(N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Get(n[i%N])
	}
}

func BenchmarkClear(b *testing.B) {
	s := newBitSet()
	n := perm(N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Clear(n[i%N])
	}
}

func BenchmarkNextClearBit(b *testing.B) {
	s := newBitSet()
	n := perm(N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.NextClearBit(n[i%N])
	}
}

func BenchmarkCardinality(b *testing.B) {
	s := newBitSet()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Cardinality()
	}
}
