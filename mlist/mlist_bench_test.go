package mlist

import (
	"container/list"
	"fmt"
	"testing"
)

var benchSizes = []int{
	1 << 10, // 1,024
	1 << 12, // 4,096
	1 << 14, // 16,384
	1 << 16, // 65,536
}

// For intentionally slow baselines (e.g. container/list random access), keep n smaller
// so "bench1" completes quickly.
var benchSizesSlow = []int{
	1 << 10, // 1,024
	1 << 12, // 4,096
	1 << 14, // 16,384
}

// simple xorshift64* for fast deterministic indices
type benchRNG uint64

func (r *benchRNG) next() uint64 {
	x := uint64(*r)
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	*r = benchRNG(x)
	return x * 2685821657736338717
}

func buildMListInt(n int) MList[int] {
	m := MList[int]{}
	for i := 0; i < n; i++ {
		m = m.Push(i)
	}
	return m
}

func buildSliceInt(n int) []int {
	s := make([]int, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, i)
	}
	return s
}

func buildContainerListInt(n int) *list.List {
	l := list.New()
	for i := 0; i < n; i++ {
		l.PushBack(i)
	}
	return l
}

func BenchmarkPushPop_MList(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			m := buildMListInt(n)
			b.ResetTimer()

			var sink int
			for i := 0; i < b.N; i++ {
				m = m.Push(i)
				var ok bool
				m, sink, ok = m.Pop()
				if !ok {
					b.Fatal("unexpected empty list")
				}
			}
			_ = sink
		})
	}
}

func BenchmarkPushPop_Slice(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			s := buildSliceInt(n)
			b.ResetTimer()

			var sink int
			for i := 0; i < b.N; i++ {
				s = append(s, i)
				sink = s[len(s)-1]
				s = s[:len(s)-1]
			}
			_ = sink
		})
	}
}

func BenchmarkPushPop_ContainerList(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			l := buildContainerListInt(n)
			b.ResetTimer()

			var sink int
			for i := 0; i < b.N; i++ {
				e := l.PushBack(i)
				sink = e.Value.(int)
				l.Remove(e)
			}
			_ = sink
		})
	}
}

func BenchmarkGetRandom_MList(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			m := buildMListInt(n)

			indices := make([]int, 1<<14)
			r := benchRNG(0xC0FFEE)
			for i := range indices {
				indices[i] = int(r.next() % uint64(n))
			}

			b.ResetTimer()
			var sink int
			for i := 0; i < b.N; i++ {
				sink = m.Get(indices[i&(len(indices)-1)])
			}
			_ = sink
		})
	}
}

func BenchmarkGetRandom_Slice(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			s := buildSliceInt(n)

			indices := make([]int, 1<<14)
			r := benchRNG(0xC0FFEE)
			for i := range indices {
				indices[i] = int(r.next() % uint64(n))
			}

			b.ResetTimer()
			var sink int
			for i := 0; i < b.N; i++ {
				sink = s[indices[i&(len(indices)-1)]]
			}
			_ = sink
		})
	}
}

func BenchmarkGetRandom_ContainerList(b *testing.B) {
	for _, n := range benchSizesSlow {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			l := buildContainerListInt(n)

			indices := make([]int, 1<<14)
			r := benchRNG(0xC0FFEE)
			for i := range indices {
				indices[i] = int(r.next() % uint64(n))
			}

			b.ResetTimer()
			var sink int
			for i := 0; i < b.N; i++ {
				target := indices[i&(len(indices)-1)]
				j := 0
				for e := l.Front(); e != nil; e = e.Next() {
					if j == target {
						sink = e.Value.(int)
						break
					}
					j++
				}
			}
			_ = sink
		})
	}
}
