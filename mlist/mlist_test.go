package mlist

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMList(t *testing.T) {
	list := MList[int]{}

	list = list.Push(1)

	require.Equal(t, 1, list.Get(0))

	list = list.Push(2)

	require.Equal(t, 1, list.Get(0))
	require.Equal(t, 2, list.Get(1))

	list = list.Push(3)

	require.Equal(t, 1, list.Get(0))
	require.Equal(t, 2, list.Get(1))
	require.Equal(t, 3, list.Get(2))

	list = list.Push(4)

	require.Equal(t, 1, list.Get(0))
	require.Equal(t, 2, list.Get(1))
	require.Equal(t, 3, list.Get(2))
	require.Equal(t, 4, list.Get(3))

	list, value, ok := list.Pop()
	require.True(t, ok)
	require.Equal(t, 4, value)
	require.Equal(t, 1, list.Get(0))
	require.Equal(t, 2, list.Get(1))
	require.Equal(t, 3, list.Get(2))

	list, value, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, 3, value)
	require.Equal(t, 1, list.Get(0))
	require.Equal(t, 2, list.Get(1))

	list, value, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, 2, value)
	require.Equal(t, 1, list.Get(0))

	list, value, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, 1, value)
	require.Equal(t, 0, list.Len())

	list, value, ok = list.Pop()
	require.False(t, ok)
	require.Equal(t, 0, value)
	require.Equal(t, 0, list.Len())
}

func TestMListLastFirst(t *testing.T) {
	list := MList[int]{}

	// Empty list
	_, ok := list.Last()
	require.False(t, ok)
	_, ok = list.First()
	require.False(t, ok)

	// Single element
	list = list.Push(42)
	last, ok := list.Last()
	require.True(t, ok)
	require.Equal(t, 42, last)
	first, ok := list.First()
	require.True(t, ok)
	require.Equal(t, 42, first)

	// Multiple elements
	list = list.Push(100).Push(200)
	last, ok = list.Last()
	require.True(t, ok)
	require.Equal(t, 200, last) // newest
	first, ok = list.First()
	require.True(t, ok)
	require.Equal(t, 42, first) // oldest
}

func TestMListPartiallyPersistent(t *testing.T) {
	list := MList[int]{} // partially persistent mode

	// Build up a list
	for i := 1; i <= 10; i++ {
		list = list.Push(i)
	}

	require.Equal(t, 10, list.Len())

	// Verify order (oldest to newest via Get)
	for i := 0; i < 10; i++ {
		require.Equal(t, i+1, list.Get(i))
	}

	// Verify Last/First
	last, ok := list.Last()
	require.True(t, ok)
	require.Equal(t, 10, last)
	first, ok := list.First()
	require.True(t, ok)
	require.Equal(t, 1, first)
}

func TestMListPowerOfTwoMerging(t *testing.T) {
	list := MList[int]{}

	// Test power-of-two behavior: 1, 2, 4, 8, 16
	for i := 1; i <= 16; i++ {
		list = list.Push(i)
		require.Equal(t, i, list.Len())
	}

	// Verify all elements are accessible
	for i := 0; i < 16; i++ {
		require.Equal(t, i+1, list.Get(i))
	}

	// At 16 elements (power of 2), should be very efficient
	require.Equal(t, 16, list.Len())
	last, ok := list.Last()
	require.True(t, ok)
	require.Equal(t, 16, last)
	first, ok := list.First()
	require.True(t, ok)
	require.Equal(t, 1, first)
}

func TestMListLargeOperations(t *testing.T) {
	list := MList[int]{}
	const size = 1000

	// Build large list
	for i := 1; i <= size; i++ {
		list = list.Push(i)
	}

	require.Equal(t, size, list.Len())

	// Random access test
	require.Equal(t, 1, list.Get(0)) // oldest
	require.Equal(t, size/2, list.Get(size/2-1))
	require.Equal(t, size, list.Get(size-1)) // newest

	// Pop half the elements
	for i := 0; i < size/2; i++ {
		var val int
		var ok bool
		list, val, ok = list.Pop()
		require.True(t, ok)
		require.Equal(t, size-i, val)
	}

	require.Equal(t, size/2, list.Len())
	last, ok := list.Last()
	require.True(t, ok)
	require.Equal(t, size/2, last)
}

func TestMListGetPanics(t *testing.T) {
	list := MList[int]{}.Push(1).Push(2).Push(3)

	// Valid access
	require.Equal(t, 1, list.Get(0))
	require.Equal(t, 3, list.Get(2))

	// Invalid access should panic
	require.Panics(t, func() { list.Get(-1) })
	require.Panics(t, func() { list.Get(3) })
	require.Panics(t, func() { list.Get(100) })

	// Empty list
	empty := MList[int]{}
	require.Panics(t, func() { empty.Get(0) })
}

func TestMListEmptyOperations(t *testing.T) {
	list := MList[int]{}

	// Empty list properties
	require.Equal(t, 0, list.Len())
	_, ok := list.Last()
	require.False(t, ok)
	_, ok = list.First()
	require.False(t, ok)

	// Pop from empty
	newList, val, ok := list.Pop()
	require.False(t, ok)
	require.Equal(t, 0, val)
	require.Equal(t, list.Len(), newList.Len()) // unchanged
}

func TestMListStringType(t *testing.T) {
	list := MList[string]{}

	words := []string{"hello", "world", "test", "golang"}

	// Push all words
	for _, word := range words {
		list = list.Push(word)
	}

	require.Equal(t, len(words), list.Len())

	// Verify order (oldest to newest)
	for i, word := range words {
		require.Equal(t, word, list.Get(i))
	}

	// Check Last/First
	last, ok := list.Last()
	require.True(t, ok)
	require.Equal(t, "golang", last)
	first, ok := list.First()
	require.True(t, ok)
	require.Equal(t, "hello", first)

	// Pop and verify
	list, val, ok := list.Pop()
	require.True(t, ok)
	require.Equal(t, "golang", val)
	require.Equal(t, len(words)-1, list.Len())
}

func TestMListPersistenceStructuralSharing(t *testing.T) {
	// Test that different versions can coexist
	base := MList[int]{}.Push(1).Push(2)

	// Create two branches
	branch1 := base.Push(10).Push(20)
	branch2 := base.Push(100).Push(200)

	// Base should be unchanged
	require.Equal(t, 2, base.Len())
	require.Equal(t, 1, base.Get(0))
	require.Equal(t, 2, base.Get(1))

	// Branch 1
	require.Equal(t, 4, branch1.Len())
	require.Equal(t, 1, branch1.Get(0))
	require.Equal(t, 2, branch1.Get(1))
	require.Equal(t, 10, branch1.Get(2))
	require.Equal(t, 20, branch1.Get(3))

	// Branch 2
	require.Equal(t, 4, branch2.Len())
	require.Equal(t, 1, branch2.Get(0))
	require.Equal(t, 2, branch2.Get(1))
	require.Equal(t, 100, branch2.Get(2))
	require.Equal(t, 200, branch2.Get(3))
}
