package marray

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMArrayBasic(t *testing.T) {
	m := &MArray[int]{}

	// Empty array
	require.Equal(t, 0, m.Len())
	_, ok := m.Last()
	require.False(t, ok)
	_, ok = m.First()
	require.False(t, ok)

	// Push single element
	m.Push(1)
	require.Equal(t, 1, m.Len())
	require.Equal(t, 1, m.Get(0))

	last, ok := m.Last()
	require.True(t, ok)
	require.Equal(t, 1, last)

	first, ok := m.First()
	require.True(t, ok)
	require.Equal(t, 1, first)

	// Push more elements
	m.Push(2)
	m.Push(3)
	m.Push(4)

	require.Equal(t, 4, m.Len())
	require.Equal(t, 1, m.Get(0))
	require.Equal(t, 2, m.Get(1))
	require.Equal(t, 3, m.Get(2))
	require.Equal(t, 4, m.Get(3))

	// Last is newest, First is oldest
	last, ok = m.Last()
	require.True(t, ok)
	require.Equal(t, 4, last)

	first, ok = m.First()
	require.True(t, ok)
	require.Equal(t, 1, first)
}

func TestMArrayPop(t *testing.T) {
	m := &MArray[int]{}

	// Pop from empty
	_, ok := m.Pop()
	require.False(t, ok)
	require.Equal(t, 0, m.Len())

	// Build array
	for i := 1; i <= 4; i++ {
		m.Push(i)
	}

	// Pop elements (LIFO - last in, first out)
	val, ok := m.Pop()
	require.True(t, ok)
	require.Equal(t, 4, val)
	require.Equal(t, 3, m.Len())

	val, ok = m.Pop()
	require.True(t, ok)
	require.Equal(t, 3, val)
	require.Equal(t, 2, m.Len())

	val, ok = m.Pop()
	require.True(t, ok)
	require.Equal(t, 2, val)
	require.Equal(t, 1, m.Len())

	val, ok = m.Pop()
	require.True(t, ok)
	require.Equal(t, 1, val)
	require.Equal(t, 0, m.Len())

	// Pop from empty again
	_, ok = m.Pop()
	require.False(t, ok)
}

func TestMArraySet(t *testing.T) {
	m := &MArray[int]{}

	// Build array
	for i := 1; i <= 10; i++ {
		m.Push(i)
	}

	require.Equal(t, 10, m.Len())

	// Verify initial values
	for i := 0; i < 10; i++ {
		require.Equal(t, i+1, m.Get(i))
	}

	// Set values
	m.Set(0, 100)
	require.Equal(t, 100, m.Get(0))

	m.Set(5, 500)
	require.Equal(t, 500, m.Get(5))

	m.Set(9, 900)
	require.Equal(t, 900, m.Get(9))

	// Verify other values unchanged
	require.Equal(t, 100, m.Get(0))
	require.Equal(t, 2, m.Get(1))
	require.Equal(t, 3, m.Get(2))
	require.Equal(t, 4, m.Get(3))
	require.Equal(t, 5, m.Get(4))
	require.Equal(t, 500, m.Get(5))
	require.Equal(t, 7, m.Get(6))
	require.Equal(t, 8, m.Get(7))
	require.Equal(t, 9, m.Get(8))
	require.Equal(t, 900, m.Get(9))

	// Length should remain unchanged
	require.Equal(t, 10, m.Len())
}

func TestMArraySetAllPositions(t *testing.T) {
	m := &MArray[int]{}

	// Build array
	for i := 0; i < 20; i++ {
		m.Push(i)
	}

	// Set all positions
	for i := 0; i < 20; i++ {
		m.Set(i, i*10)
	}

	// Verify all positions
	for i := 0; i < 20; i++ {
		require.Equal(t, i*10, m.Get(i))
	}
}

func TestMArrayLargeOperations(t *testing.T) {
	m := &MArray[int]{}

	// Push many elements
	n := 1000
	for i := 0; i < n; i++ {
		m.Push(i)
	}

	require.Equal(t, n, m.Len())

	// Verify all elements
	for i := 0; i < n; i++ {
		require.Equal(t, i, m.Get(i))
	}

	// Modify some elements
	for i := 0; i < n; i += 10 {
		m.Set(i, i*100)
	}

	// Verify modified and unmodified elements
	for i := 0; i < n; i++ {
		if i%10 == 0 {
			require.Equal(t, i*100, m.Get(i))
		} else {
			require.Equal(t, i, m.Get(i))
		}
	}

	// Pop some elements
	for i := 0; i < 100; i++ {
		val, ok := m.Pop()
		require.True(t, ok)
		expected := n - 1 - i
		if expected%10 == 0 {
			require.Equal(t, expected*100, val)
		} else {
			require.Equal(t, expected, val)
		}
	}

	require.Equal(t, n-100, m.Len())
}

func TestMArrayPushPopCycle(t *testing.T) {
	m := &MArray[int]{}

	// Multiple cycles of push and pop
	for cycle := 0; cycle < 5; cycle++ {
		// Push elements
		for i := 0; i < 10; i++ {
			m.Push(cycle*10 + i)
		}

		// Pop half
		for i := 0; i < 5; i++ {
			_, ok := m.Pop()
			require.True(t, ok)
		}
	}

	// Should have 25 elements remaining
	require.Equal(t, 25, m.Len())

	// Verify they are in correct order (oldest first)
	expected := []int{
		0, 1, 2, 3, 4,
		10, 11, 12, 13, 14,
		20, 21, 22, 23, 24,
		30, 31, 32, 33, 34,
		40, 41, 42, 43, 44,
	}

	for i, exp := range expected {
		require.Equal(t, exp, m.Get(i), "index %d", i)
	}
}

func TestMArrayStrings(t *testing.T) {
	m := &MArray[string]{}

	// Test with string type
	m.Push("hello")
	m.Push("world")
	m.Push("foo")
	m.Push("bar")

	require.Equal(t, 4, m.Len())
	require.Equal(t, "hello", m.Get(0))
	require.Equal(t, "world", m.Get(1))
	require.Equal(t, "foo", m.Get(2))
	require.Equal(t, "bar", m.Get(3))

	m.Set(1, "WORLD")
	require.Equal(t, "WORLD", m.Get(1))

	val, ok := m.Pop()
	require.True(t, ok)
	require.Equal(t, "bar", val)

	first, ok := m.First()
	require.True(t, ok)
	require.Equal(t, "hello", first)

	last, ok := m.Last()
	require.True(t, ok)
	require.Equal(t, "foo", last)
}

func TestMArrayPointers(t *testing.T) {
	m := &MArray[*int]{}

	// Test with pointer type
	one := 1
	two := 2
	three := 3

	m.Push(&one)
	m.Push(&two)
	m.Push(&three)

	require.Equal(t, 3, m.Len())
	require.Equal(t, &one, m.Get(0))
	require.Equal(t, &two, m.Get(1))
	require.Equal(t, &three, m.Get(2))

	four := 4
	m.Set(1, &four)
	require.Equal(t, &four, m.Get(1))
}

func TestMArrayPowerOfTwo(t *testing.T) {
	m := &MArray[int]{}

	// Test power-of-two sizes (important for segment merging)
	sizes := []int{1, 2, 4, 8, 16, 32, 64, 128}

	for _, size := range sizes {
		m = &MArray[int]{}

		// Push size elements
		for i := 0; i < size; i++ {
			m.Push(i)
		}

		require.Equal(t, size, m.Len())

		// Verify all elements
		for i := 0; i < size; i++ {
			require.Equal(t, i, m.Get(i))
		}

		// Set all elements
		for i := 0; i < size; i++ {
			m.Set(i, i*2)
		}

		// Verify modified elements
		for i := 0; i < size; i++ {
			require.Equal(t, i*2, m.Get(i))
		}
	}
}

func TestMArrayFirstLastAfterOperations(t *testing.T) {
	m := &MArray[int]{}

	// Push elements
	m.Push(10)
	m.Push(20)
	m.Push(30)
	m.Push(40)
	m.Push(50)

	// First is oldest (10), Last is newest (50)
	first, ok := m.First()
	require.True(t, ok)
	require.Equal(t, 10, first)

	last, ok := m.Last()
	require.True(t, ok)
	require.Equal(t, 50, last)

	// Pop newest
	val, ok := m.Pop()
	require.True(t, ok)
	require.Equal(t, 50, val)

	// First unchanged, Last is now 40
	first, ok = m.First()
	require.True(t, ok)
	require.Equal(t, 10, first)

	last, ok = m.Last()
	require.True(t, ok)
	require.Equal(t, 40, last)

	// Set doesn't affect First/Last
	m.Set(1, 999)

	first, ok = m.First()
	require.True(t, ok)
	require.Equal(t, 10, first)

	last, ok = m.Last()
	require.True(t, ok)
	require.Equal(t, 40, last)
}

func TestMArrayEmpty(t *testing.T) {
	m := &MArray[int]{}

	// All operations on empty array
	require.Equal(t, 0, m.Len())

	_, ok := m.Last()
	require.False(t, ok)

	_, ok = m.First()
	require.False(t, ok)

	_, ok = m.Pop()
	require.False(t, ok)
}

// Benchmarks

func BenchmarkMArrayPush(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m := &MArray[int]{}
				for j := 0; j < size; j++ {
					m.Push(j)
				}
			}
		})
	}
}

func BenchmarkMArrayPop(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				m := &MArray[int]{}
				for j := 0; j < size; j++ {
					m.Push(j)
				}
				b.StartTimer()
				for j := 0; j < size; j++ {
					m.Pop()
				}
				b.StopTimer()
			}
		})
	}
}

func BenchmarkMArrayGet(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		m := &MArray[int]{}
		for j := 0; j < size; j++ {
			m.Push(j)
		}

		b.Run(fmt.Sprintf("Front_n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = m.Get(0)
			}
		})

		b.Run(fmt.Sprintf("Middle_n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = m.Get(size / 2)
			}
		})

		b.Run(fmt.Sprintf("Back_n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = m.Get(size - 1)
			}
		})
	}
}

func BenchmarkMArrayIteration(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		m := &MArray[int]{}
		for j := 0; j < size; j++ {
			m.Push(j)
		}

		b.Run(fmt.Sprintf("n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				sum := 0
				for j := 0; j < m.Len(); j++ {
					sum += m.Get(j)
				}
				_ = sum
			}
		})
	}
}
