package pflist

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMListFullyPersistent(t *testing.T) {
	list := PFList[int]{} // fully persistent mode

	// Test that old versions remain unchanged
	v1 := list.Push(1)
	v2 := v1.Push(2)
	v3 := v2.Push(3)

	// Check v1 is unchanged
	require.Equal(t, 1, v1.Len())
	val, ok := v1.Last()
	require.True(t, ok)
	require.Equal(t, 1, val)

	// Check v2 is unchanged
	require.Equal(t, 2, v2.Len())
	val, ok = v2.Last()
	require.True(t, ok)
	require.Equal(t, 2, val)

	// Check v3
	require.Equal(t, 3, v3.Len())
	val, ok = v3.Last()
	require.True(t, ok)
	require.Equal(t, 3, val)
}

func TestEmptyList(t *testing.T) {
	list := PFList[int]{}

	// Empty list should have length 0
	require.Equal(t, 0, list.Len())

	// Last on empty list should return false
	val, ok := list.Last()
	require.False(t, ok)
	require.Equal(t, 0, val)

	// Pop on empty list should return false
	newList, val, ok := list.Pop()
	require.False(t, ok)
	require.Equal(t, 0, val)
	require.Equal(t, 0, newList.Len())
}

func TestSingleElement(t *testing.T) {
	list := PFList[int]{}
	list = list.Push(42)

	// Check length
	require.Equal(t, 1, list.Len())

	// Check Last
	val, ok := list.Last()
	require.True(t, ok)
	require.Equal(t, 42, val)

	// Pop the element
	list2, val, ok := list.Pop()
	require.True(t, ok)
	require.Equal(t, 42, val)
	require.Equal(t, 0, list2.Len())

	// Original list should be unchanged
	require.Equal(t, 1, list.Len())
	val, ok = list.Last()
	require.True(t, ok)
	require.Equal(t, 42, val)
}

func TestPushPop(t *testing.T) {
	list := PFList[int]{}

	// Push several elements
	list = list.Push(1)
	list = list.Push(2)
	list = list.Push(3)
	list = list.Push(4)

	require.Equal(t, 4, list.Len())

	// Pop them and verify LIFO order
	var val int
	var ok bool

	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, 4, val)
	require.Equal(t, 3, list.Len())

	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, 3, val)
	require.Equal(t, 2, list.Len())

	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, 2, val)
	require.Equal(t, 1, list.Len())

	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, 1, val)
	require.Equal(t, 0, list.Len())

	// One more pop should return false
	list, val, ok = list.Pop()
	require.False(t, ok)
	require.Equal(t, 0, list.Len())
}

func TestPopThenPush(t *testing.T) {
	// This tests the SplitLast optimization path in Push
	list := PFList[int]{}

	// Build up a list
	for i := 0; i < 10; i++ {
		list = list.Push(i)
	}
	require.Equal(t, 10, list.Len())

	// Pop one element (this changes horizon)
	list, val, ok := list.Pop()
	require.True(t, ok)
	require.Equal(t, 9, val)
	require.Equal(t, 9, list.Len())

	// Push another element (this should trigger SplitLast path)
	list = list.Push(100)
	require.Equal(t, 10, list.Len())

	val, ok = list.Last()
	require.True(t, ok)
	require.Equal(t, 100, val)

	// Pop and verify order
	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, 100, val)
}

func TestLargeSequence(t *testing.T) {
	list := PFList[int]{}
	n := 1000

	// Push many elements
	for i := 0; i < n; i++ {
		list = list.Push(i)
	}
	require.Equal(t, n, list.Len())

	// Check last element
	val, ok := list.Last()
	require.True(t, ok)
	require.Equal(t, n-1, val)

	// Pop all elements and verify order
	for i := n - 1; i >= 0; i-- {
		var ok bool
		list, val, ok = list.Pop()
		require.True(t, ok)
		require.Equal(t, i, val)
		require.Equal(t, i, list.Len())
	}

	require.Equal(t, 0, list.Len())
}

func TestInterleavedOperations(t *testing.T) {
	list := PFList[int]{}

	// Interleave push and pop operations
	list = list.Push(1)
	list = list.Push(2)

	var val int
	var ok bool

	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, 2, val)

	list = list.Push(3)
	list = list.Push(4)

	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, 4, val)

	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, 3, val)

	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, 1, val)

	require.Equal(t, 0, list.Len())
}

func TestPersistenceWithPop(t *testing.T) {
	list := PFList[int]{}

	// Create multiple versions with push
	v1 := list.Push(1)
	v2 := v1.Push(2)
	v3 := v2.Push(3)

	// Pop from v3 to create v4
	v4, val, ok := v3.Pop()
	require.True(t, ok)
	require.Equal(t, 3, val)
	require.Equal(t, 2, v4.Len())

	// v3 should still have 3 elements
	require.Equal(t, 3, v3.Len())
	val, ok = v3.Last()
	require.True(t, ok)
	require.Equal(t, 3, val)

	// v4 should have 2 elements with 2 as last
	val, ok = v4.Last()
	require.True(t, ok)
	require.Equal(t, 2, val)

	// Push to v4 to create v5
	v5 := v4.Push(100)
	require.Equal(t, 3, v5.Len())
	val, ok = v5.Last()
	require.True(t, ok)
	require.Equal(t, 100, val)

	// All previous versions should be unchanged
	require.Equal(t, 1, v1.Len())
	require.Equal(t, 2, v2.Len())
	require.Equal(t, 3, v3.Len())
	require.Equal(t, 2, v4.Len())
}

func TestMultipleBranches(t *testing.T) {
	// Test creating multiple branches from the same version
	base := PFList[string]{}
	base = base.Push("a")
	base = base.Push("b")

	// Create branch 1
	branch1 := base.Push("c1")
	branch1 = branch1.Push("d1")

	// Create branch 2 from base
	branch2 := base.Push("c2")
	branch2 = branch2.Push("d2")

	// Verify base is unchanged
	require.Equal(t, 2, base.Len())
	val, ok := base.Last()
	require.True(t, ok)
	require.Equal(t, "b", val)

	// Verify branch1
	require.Equal(t, 4, branch1.Len())
	val, ok = branch1.Last()
	require.True(t, ok)
	require.Equal(t, "d1", val)

	// Verify branch2
	require.Equal(t, 4, branch2.Len())
	val, ok = branch2.Last()
	require.True(t, ok)
	require.Equal(t, "d2", val)

	// Pop from branch1 and verify branch2 is unaffected
	branch1, val, ok = branch1.Pop()
	require.True(t, ok)
	require.Equal(t, "d1", val)

	val, ok = branch2.Last()
	require.True(t, ok)
	require.Equal(t, "d2", val)
}

func TestStringValues(t *testing.T) {
	list := PFList[string]{}

	list = list.Push("hello")
	list = list.Push("world")
	list = list.Push("!")

	require.Equal(t, 3, list.Len())

	val, ok := list.Last()
	require.True(t, ok)
	require.Equal(t, "!", val)

	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, "!", val)

	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, "world", val)

	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, "hello", val)

	require.Equal(t, 0, list.Len())
}

func TestStructValues(t *testing.T) {
	type Point struct {
		X, Y int
	}

	list := PFList[Point]{}

	list = list.Push(Point{1, 2})
	list = list.Push(Point{3, 4})
	list = list.Push(Point{5, 6})

	require.Equal(t, 3, list.Len())

	val, ok := list.Last()
	require.True(t, ok)
	require.Equal(t, Point{5, 6}, val)

	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, Point{5, 6}, val)
	require.Equal(t, 2, list.Len())
}

func TestRepeatedPopPush(t *testing.T) {
	// Test repeated pop-push cycles to exercise the SplitLast path
	list := PFList[int]{}

	// Build initial list
	for i := 0; i < 20; i++ {
		list = list.Push(i)
	}

	// Repeatedly pop and push
	for i := 0; i < 50; i++ {
		var val int
		var ok bool
		list, val, ok = list.Pop()
		require.True(t, ok)
		require.Equal(t, 19, list.Len())

		list = list.Push(val + 1000)
		require.Equal(t, 20, list.Len())
	}

	// List should still be functional
	require.Equal(t, 20, list.Len())
	val, ok := list.Last()
	require.True(t, ok)
	require.NotEqual(t, 0, val)
}

func TestZeroValues(t *testing.T) {
	list := PFList[int]{}

	// Push zero values
	list = list.Push(0)
	list = list.Push(0)
	list = list.Push(0)

	require.Equal(t, 3, list.Len())

	// Pop and verify we get zeros back
	var val int
	var ok bool

	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, 0, val)

	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, 0, val)

	list, val, ok = list.Pop()
	require.True(t, ok)
	require.Equal(t, 0, val)

	require.Equal(t, 0, list.Len())
}
