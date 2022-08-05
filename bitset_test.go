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

var N = 1000000

func newBitSet() *BitSet {
	s := NewSize(uint64(N))
	for i := 0; i < len(s.values); i++ {
		s.values[i] = 1<<8 - 1
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