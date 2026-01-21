package pflist

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

var benchSizesSlow = []int{
	1 << 10, // 1,024
	1 << 12, // 4,096
	1 << 14, // 16,384
}

func buildPFListInt(n int) PFList[int] {
	p := PFList[int]{}
	for i := 0; i < n; i++ {
		p = p.Push(i)
	}
	return p
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

func BenchmarkPushPop_PFList(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			p := buildPFListInt(n)
			b.ResetTimer()

			var sink int
			for i := 0; i < b.N; i++ {
				p = p.Push(i)
				var ok bool
				p, sink, ok = p.Pop()
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

func BenchmarkGetRandom_PFList(b *testing.B) {
	for _, n := range benchSizes {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			p := buildPFListInt(n)

			indices := make([]int, 1<<14)
			var x uint64 = 0xBADC0DE
			for i := range indices {
				// LCG (fast, deterministic)
				x = x*2862933555777941757 + 3037000493
				indices[i] = int(x % uint64(n))
			}

			b.ResetTimer()
			var sink int
			for i := 0; i < b.N; i++ {
				sink = p.Get(indices[i&(len(indices)-1)])
			}
			_ = sink
		})
	}
}

func BenchmarkBranching_PFList(b *testing.B) {
	const branchOps = 64
	for _, n := range benchSizesSlow {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			base := buildPFListInt(n)
			b.ResetTimer()

			var sink int
			for i := 0; i < b.N; i++ {
				b1 := base
				for k := 0; k < branchOps; k++ {
					b1 = b1.Push(k)
				}
				b2 := base
				for k := 0; k < branchOps; k++ {
					b2 = b2.Push(k + branchOps)
				}
				v, _ := b1.Last()
				sink ^= v
				v, _ = b2.Last()
				sink ^= v
			}
			_ = sink
		})
	}
}

func BenchmarkBranching_SliceCopy(b *testing.B) {
	const branchOps = 64
	for _, n := range benchSizesSlow {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			base := buildSliceInt(n)
			b.ResetTimer()

			var sink int
			for i := 0; i < b.N; i++ {
				b1 := append([]int(nil), base...)
				for k := 0; k < branchOps; k++ {
					b1 = append(b1, k)
				}
				b2 := append([]int(nil), base...)
				for k := 0; k < branchOps; k++ {
					b2 = append(b2, k+branchOps)
				}
				sink ^= b1[len(b1)-1]
				sink ^= b2[len(b2)-1]
			}
			_ = sink
		})
	}
}
