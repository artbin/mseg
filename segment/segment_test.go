package segment

import (
	"testing"
)

// Helper function to verify segment chain structure
func verifyChain[E comparable](t *testing.T, last *Segment[E], first *Segment[E], expectedLen int, horizon int, expectedElems []E) {
	t.Helper()

	// Collect all elements from newest to oldest
	var collected []E
	if last != nil {
		// Add visible elements from last segment
		for i := horizon; i >= 0; i-- {
			collected = append(collected, last.Elems[i])
		}

		// Add elements from remaining segments
		for seg := last.Next; seg != nil; seg = seg.Next {
			for i := len(seg.Elems) - 1; i >= 0; i-- {
				collected = append(collected, seg.Elems[i])
			}
		}
	}

	if len(collected) != expectedLen {
		t.Errorf("collected %d elements, expected %d", len(collected), expectedLen)
	}

	if len(collected) != len(expectedElems) {
		t.Errorf("collected length %d != expected length %d", len(collected), len(expectedElems))
		return
	}

	for i := range collected {
		if collected[i] != expectedElems[i] {
			t.Errorf("element at index %d: got %v, want %v", i, collected[i], expectedElems[i])
		}
	}

	// Verify first pointer correctness
	if last != nil && first != nil {
		// Walk to the end to find actual first
		actualFirst := last
		for actualFirst.Next != nil {
			actualFirst = actualFirst.Next
		}
		if actualFirst != first {
			t.Errorf("first pointer mismatch: got %p, walking from last gives %p", first, actualFirst)
		}
	}
}

func TestMergeEmpty(t *testing.T) {
	last, first, listlen := Merge(42, nil, nil, 0)

	if listlen != 1 {
		t.Errorf("listlen = %d, want 1", listlen)
	}

	if last == nil {
		t.Fatal("last is nil")
	}

	if first != last {
		t.Errorf("first != last for single element")
	}

	if len(last.Elems) != 1 || last.Elems[0] != 42 {
		t.Errorf("last.Elems = %v, want [42]", last.Elems)
	}

	if last.Next != nil {
		t.Errorf("last.Next should be nil")
	}

	horizon := len(last.Elems) - 1
	verifyChain(t, last, first, listlen, horizon, []int{42})
}

func TestMergeTwoElements(t *testing.T) {
	// Add first element
	last, first, listlen := Merge(1, nil, nil, 0)
	horizon := len(last.Elems) - 1

	// Add second element - should coalesce (length 2 is power of two)
	last, first, listlen = Merge(2, last, first, listlen)
	horizon = len(last.Elems) - 1

	if listlen != 2 {
		t.Errorf("listlen = %d, want 2", listlen)
	}

	if last == nil {
		t.Fatal("last is nil")
	}

	if first != last {
		t.Errorf("first != last after coalescing to power of two")
	}

	if len(last.Elems) != 2 {
		t.Errorf("last.Elems length = %d, want 2", len(last.Elems))
	}

	verifyChain(t, last, first, listlen, horizon, []int{2, 1})
}

func TestMergePowerOfTwo(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	// Add 8 elements (power of two)
	for i := 1; i <= 8; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1

	// At power of two, everything should coalesce into one segment
	if last != first {
		t.Errorf("expected full coalescing at power of two")
	}

	if len(last.Elems) != 8 {
		t.Errorf("last.Elems length = %d, want 8", len(last.Elems))
	}

	if last.Next != nil {
		t.Errorf("last.Next should be nil after full coalescing")
	}

	expected := []int{8, 7, 6, 5, 4, 3, 2, 1}
	verifyChain(t, last, first, listlen, horizon, expected)
}

func TestMergeNonPowerOfTwo(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	// Add 7 elements (not power of two)
	// Binary: 111 = segments of size 4, 2, 1
	for i := 1; i <= 7; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1

	if listlen != 7 {
		t.Errorf("listlen = %d, want 7", listlen)
	}

	// Count segments
	segmentCount := 0
	for seg := last; seg != nil; seg = seg.Next {
		segmentCount++
	}

	if segmentCount != 3 {
		t.Errorf("segment count = %d, want 3 for binary 111", segmentCount)
	}

	expected := []int{7, 6, 5, 4, 3, 2, 1}
	verifyChain(t, last, first, listlen, horizon, expected)
}

func TestMergeLarge(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	// Add 100 elements
	for i := 1; i <= 100; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1

	if listlen != 100 {
		t.Errorf("listlen = %d, want 100", listlen)
	}

	// Verify all elements are present
	expected := make([]int, 100)
	for i := 0; i < 100; i++ {
		expected[i] = 100 - i
	}

	verifyChain(t, last, first, listlen, horizon, expected)
}

func TestPopFrontEmpty(t *testing.T) {
	last, first, listlen, horizon, val, ok := PopFront[int](nil, nil, 0, 0)

	if ok {
		t.Errorf("PopFront on empty list should return ok=false")
	}

	if val != 0 {
		t.Errorf("PopFront on empty list should return zero value")
	}

	if last != nil || first != nil {
		t.Errorf("PopFront on empty list should return nil pointers")
	}

	if listlen != 0 {
		t.Errorf("PopFront on empty list should return listlen=0")
	}

	if horizon != 0 {
		t.Errorf("PopFront on empty list should return horizon=0")
	}
}

func TestPopFrontSingleElement(t *testing.T) {
	last, first, listlen := Merge(42, nil, nil, 0)
	horizon := len(last.Elems) - 1

	last, first, listlen, horizon, val, ok := PopFront(last, first, listlen, horizon)

	if !ok {
		t.Fatal("PopFront should succeed")
	}

	if val != 42 {
		t.Errorf("popped value = %d, want 42", val)
	}

	if listlen != 0 {
		t.Errorf("listlen = %d, want 0", listlen)
	}

	if last != nil || first != nil {
		t.Errorf("list should be empty after popping single element")
	}
}

func TestPopFrontMultiple(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	// Build list: 1, 2, 3, 4, 5
	for i := 1; i <= 5; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1

	// Pop front (newest first): should get 5, 4, 3, 2, 1
	expected := []int{5, 4, 3, 2, 1}
	for i, want := range expected {
		var val int
		var ok bool
		last, first, listlen, horizon, val, ok = PopFront(last, first, listlen, horizon)

		if !ok {
			t.Fatalf("PopFront #%d failed", i)
		}

		if val != want {
			t.Errorf("PopFront #%d: got %d, want %d", i, val, want)
		}

		if listlen != 5-i-1 {
			t.Errorf("PopFront #%d: listlen = %d, want %d", i, listlen, 5-i-1)
		}
	}

	// Should be empty now
	if listlen != 0 {
		t.Errorf("final listlen = %d, want 0", listlen)
	}

	if last != nil || first != nil {
		t.Errorf("list should be empty after popping all elements")
	}
}

func TestPopFrontCompaction(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	// Build a list large enough to trigger compaction
	for i := 1; i <= 16; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1
	initialCap := cap(last.Elems)

	// Pop elements to trigger compaction (when horizon == cap/2)
	popsNeeded := horizon - (initialCap >> 1)
	for i := 0; i < popsNeeded; i++ {
		last, first, listlen, horizon, _, _ = PopFront(last, first, listlen, horizon)
	}

	// Next pop should trigger compaction
	capBeforeCompaction := cap(last.Elems)
	last, first, listlen, horizon, _, ok := PopFront(last, first, listlen, horizon)

	if !ok {
		t.Fatal("PopFront should succeed")
	}

	if cap(last.Elems) >= capBeforeCompaction {
		t.Errorf("expected compaction to reduce capacity, before=%d, after=%d", capBeforeCompaction, cap(last.Elems))
	}

	// Verify remaining elements are still correct
	remaining := make([]int, listlen)
	for i := 0; i < listlen; i++ {
		remaining[i] = listlen - i
	}

	verifyChain(t, last, first, listlen, horizon, remaining)
}

func TestPopBackEmpty(t *testing.T) {
	_, _, listlen, _, val, ok := PopBack[int](nil, nil, 0, 0)

	if ok {
		t.Errorf("PopBack on empty list should return ok=false")
	}

	if val != 0 {
		t.Errorf("PopBack on empty list should return zero value")
	}

	if listlen != 0 {
		t.Errorf("PopBack on empty list should return listlen=0")
	}
}

func TestPopBackSingleElement(t *testing.T) {
	last, first, listlen := Merge(42, nil, nil, 0)
	horizon := 0 // For PopBack, horizon starts at 0 (oldest index in first)

	last, first, listlen, horizon, val, ok := PopBack(last, first, listlen, horizon)

	if !ok {
		t.Fatal("PopBack should succeed")
	}

	if val != 42 {
		t.Errorf("popped value = %d, want 42", val)
	}

	if listlen != 0 {
		t.Errorf("listlen = %d, want 0", listlen)
	}
}

func TestPopBackMultiple(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	// Build list: 1, 2, 3, 4, 5 (1 is oldest, 5 is newest)
	for i := 1; i <= 5; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := 0 // For PopBack, horizon starts at 0 (oldest index in first)

	// Pop back (oldest first): should get 1, 2, 3, 4, 5
	expected := []int{1, 2, 3, 4, 5}
	for i, want := range expected {
		var val int
		var ok bool
		last, first, listlen, horizon, val, ok = PopBack(last, first, listlen, horizon)

		if !ok {
			t.Fatalf("PopBack #%d failed", i)
		}

		if val != want {
			t.Errorf("PopBack #%d: got %d, want %d", i, val, want)
		}

		if listlen != 5-i-1 {
			t.Errorf("PopBack #%d: listlen = %d, want %d", i, listlen, 5-i-1)
		}
	}

	// Should be empty now
	if listlen != 0 {
		t.Errorf("final listlen = %d, want 0", listlen)
	}
}

func TestSplitLastEmpty(t *testing.T) {
	last, first, listlen, horizon := SplitLast[int](nil, nil, 0, 0)

	if last != nil || first != nil {
		t.Errorf("SplitLast on nil should return nil")
	}

	if listlen != 0 || horizon != 0 {
		t.Errorf("SplitLast on nil should preserve zero values")
	}
}

func TestSplitLastSingleElement(t *testing.T) {
	last, first, listlen := Merge(42, nil, nil, 0)
	horizon := len(last.Elems) - 1

	origLast := last
	last, first, listlen, horizon = SplitLast(last, first, listlen, horizon)

	// Single element should not split
	if last != origLast {
		t.Errorf("SplitLast on single element should not change last pointer")
	}

	verifyChain(t, last, first, listlen, horizon, []int{42})
}

func TestSplitLastPowerOfTwo(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	// Build list with 8 elements (power of two)
	for i := 1; i <= 8; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1

	// SplitLast should split into segments of size 1, 2, 4
	last, first, listlen, horizon = SplitLast(last, first, listlen, horizon)

	// Count segments
	segmentCount := 0
	sizes := []int{}
	for seg := last; seg != nil; seg = seg.Next {
		segmentCount++
		sizes = append(sizes, len(seg.Elems))
	}

	// 8 = 0b1000 has popcount 1, but we have 8 elements to split
	// Actually 8 visible elements = horizon 7, newHorizon = 8
	// 8 = 0b1000, so popcount = 1, one segment of size 8
	// Wait, let me recalculate...
	// If horizon = 7, newHorizon = 8
	// 8 in binary is 1000, popcount = 1
	// So we get one segment of size 8

	expectedSegments := 1
	if segmentCount != expectedSegments {
		t.Logf("segment sizes: %v", sizes)
		t.Errorf("segment count = %d, want %d", segmentCount, expectedSegments)
	}

	// Verify all elements are present
	expected := []int{8, 7, 6, 5, 4, 3, 2, 1}
	verifyChain(t, last, first, listlen, horizon, expected)
}

func TestSplitLastNonPowerOfTwo(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	// Build list with 7 elements
	for i := 1; i <= 7; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1

	// After merge, this could be multiple segments already
	// Force into single segment for testing
	elems := make([]int, 7)
	for i := 0; i < 7; i++ {
		elems[i] = i + 1
	}
	last = &Segment[int]{Elems: elems}
	first = last
	horizon = 6

	// SplitLast with horizon=6 means 7 visible elements
	// 7 = 0b111, popcount = 3, so segments of size 1, 2, 4
	last, first, listlen, horizon = SplitLast(last, first, listlen, horizon)

	// Count segments
	segmentCount := 0
	sizes := []int{}
	for seg := last; seg != nil; seg = seg.Next {
		segmentCount++
		sizes = append(sizes, len(seg.Elems))
	}

	expectedSegments := 3
	expectedSizes := []int{1, 2, 4}
	if segmentCount != expectedSegments {
		t.Errorf("segment count = %d, want %d. Sizes: %v", segmentCount, expectedSegments, sizes)
	}

	if len(sizes) == len(expectedSizes) {
		for i := range sizes {
			if sizes[i] != expectedSizes[i] {
				t.Errorf("segment %d size = %d, want %d", i, sizes[i], expectedSizes[i])
			}
		}
	}

	// Verify all elements are present (newest to oldest)
	expected := []int{7, 6, 5, 4, 3, 2, 1}
	verifyChain(t, last, first, listlen, horizon, expected)
}

func TestFindFirstSegment(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	// Build list with several elements
	for i := 1; i <= 10; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1

	// Find element in first segment (oldest segment)
	// Index 0 should be the oldest element (1)
	slice, offset := Find(0, last, first, listlen, horizon)

	if slice[offset] != 1 {
		t.Errorf("Find(0) = %d, want 1", slice[offset])
	}
}

func TestFindLastSegment(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	// Build list
	for i := 1; i <= 10; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1

	// Find element in last segment (newest segment)
	// Index listlen-1 should be the newest element (10)
	slice, offset := Find(listlen-1, last, first, listlen, horizon)

	if slice[offset] != 10 {
		t.Errorf("Find(%d) = %d, want 10", listlen-1, slice[offset])
	}
}

func TestFindAllElements(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	// Build list
	for i := 1; i <= 20; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1

	// Find all elements by index (counting from oldest=0)
	for i := 0; i < listlen; i++ {
		slice, offset := Find(i, last, first, listlen, horizon)
		expected := i + 1 // Element values are 1, 2, 3, ..., 20

		if slice[offset] != expected {
			t.Errorf("Find(%d) = %d, want %d", i, slice[offset], expected)
		}
	}
}

func TestFindOutOfBounds(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	for i := 1; i <= 5; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1

	// Test negative index
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Find(-1) should panic")
		}
	}()
	Find(-1, last, first, listlen, horizon)
}

func TestFindOutOfBoundsHigh(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	for i := 1; i <= 5; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1

	// Test index >= listlen
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Find(listlen) should panic")
		}
	}()
	Find(listlen, last, first, listlen, horizon)
}

func TestMixedOperations(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	// Build initial list: 1, 2, 3, 4
	for i := 1; i <= 4; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	frontHorizon := len(last.Elems) - 1

	// Pop from front (removes 4)
	last, first, listlen, frontHorizon, val, ok := PopFront(last, first, listlen, frontHorizon)
	if !ok || val != 4 {
		t.Errorf("PopFront: got (%d, %v), want (4, true)", val, ok)
	}

	if listlen != 3 {
		t.Errorf("after PopFront: listlen = %d, want 3", listlen)
	}

	// Before merging after PopFront, materialize the horizon via SplitLast
	// This ensures Merge sees only the visible elements
	last, first, listlen, frontHorizon = SplitLast(last, first, listlen, frontHorizon)

	// Add more elements (5 and 6)
	last, first, listlen = Merge(5, last, first, listlen)
	last, first, listlen = Merge(6, last, first, listlen)

	if listlen != 5 {
		t.Errorf("after merges: listlen = %d, want 5", listlen)
	}

	frontHorizon = len(last.Elems) - 1

	// Verify remaining elements (newest to oldest): 6, 5, 3, 2, 1
	expected := []int{6, 5, 3, 2, 1}
	verifyChain(t, last, first, listlen, frontHorizon, expected)

	// Now pop everything from front to verify correctness
	for i, want := range expected {
		last, first, listlen, frontHorizon, val, ok = PopFront(last, first, listlen, frontHorizon)
		if !ok {
			t.Fatalf("PopFront #%d failed", i)
		}
		if val != want {
			t.Errorf("PopFront #%d: got %d, want %d", i, val, want)
		}
	}

	if listlen != 0 {
		t.Errorf("final listlen = %d, want 0", listlen)
	}
}

func TestSegmentInvariantBinaryCounter(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	// Add many elements and verify segment count stays O(log n)
	for i := 1; i <= 100; i++ {
		last, first, listlen = Merge(i, last, first, listlen)

		// Count segments
		segmentCount := 0
		for seg := last; seg != nil; seg = seg.Next {
			segmentCount++
		}

		// Segment count should be at most floor(log2(n)) + 1
		// For n=100, log2(100) ≈ 6.64, so at most 7 segments
		maxSegments := 0
		n := i
		for n > 0 {
			maxSegments++
			n >>= 1
		}

		if segmentCount > maxSegments {
			t.Errorf("after %d merges: segment count = %d, exceeds max %d", i, segmentCount, maxSegments)
		}
	}
}

func TestSplitLastWithRemainder(t *testing.T) {
	// Create a scenario where last has a Next pointer
	var last, first *Segment[int]
	listlen := 0

	// Build initial list
	for i := 1; i <= 5; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	// Manually create a chain with specific structure
	// Last segment with 4 elements, then an older segment
	older := &Segment[int]{Elems: []int{1}}
	last = &Segment[int]{
		Elems: []int{2, 3, 4, 5},
		Next:  older,
	}
	first = older
	listlen = 5
	horizon := 3 // 4 visible elements in last

	// Split last (4 elements = 0b100, one segment)
	last, first, listlen, horizon = SplitLast(last, first, listlen, horizon)

	// Verify the older segment is still linked
	foundOlder := false
	for seg := last; seg != nil; seg = seg.Next {
		if seg == older {
			foundOlder = true
			break
		}
	}

	if !foundOlder {
		t.Errorf("SplitLast should preserve the chain beyond the split segment")
	}

	// Verify all elements
	expected := []int{5, 4, 3, 2, 1}
	verifyChain(t, last, first, listlen, horizon, expected)
}

func TestPopBackSingleSegmentMultipleElements(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	// Build a list that coalesces into one segment (power of two)
	for i := 1; i <= 4; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	if last != first {
		t.Fatalf("expected single segment for power of two")
	}

	horizon := 0 // For PopBack

	// Pop all elements from back
	for i := 1; i <= 4; i++ {
		var val int
		var ok bool
		last, first, listlen, horizon, val, ok = PopBack(last, first, listlen, horizon)

		if !ok {
			t.Fatalf("PopBack #%d failed", i)
		}

		if val != i {
			t.Errorf("PopBack #%d: got %d, want %d", i, val, i)
		}
	}

	if listlen != 0 {
		t.Errorf("listlen = %d, want 0 after popping all", listlen)
	}
}

func TestPopFrontAdvanceToNextSegment(t *testing.T) {
	// Create a specific structure with multiple segments
	seg2 := &Segment[int]{Elems: []int{1, 2}}
	seg1 := &Segment[int]{Elems: []int{3}, Next: seg2}

	last := seg1
	first := seg2
	listlen := 3
	horizon := 0 // Only one visible element in last

	// Pop the only element in last segment
	last, first, listlen, horizon, val, ok := PopFront(last, first, listlen, horizon)

	if !ok || val != 3 {
		t.Errorf("PopFront: got (%d, %v), want (3, true)", val, ok)
	}

	// Should advance to next segment
	if last != seg2 {
		t.Errorf("last should advance to next segment")
	}

	// Horizon should point to the last element of the new last segment
	if horizon != len(last.Elems)-1 {
		t.Errorf("horizon = %d, want %d", horizon, len(last.Elems)-1)
	}
}

func TestLargeScaleMergeAndPop(t *testing.T) {
	var last, first *Segment[int]
	listlen := 0

	// Add 1000 elements
	for i := 1; i <= 1000; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	if listlen != 1000 {
		t.Errorf("listlen = %d, want 1000", listlen)
	}

	horizon := len(last.Elems) - 1

	// Pop all from front
	for i := 1000; i >= 1; i-- {
		var val int
		var ok bool
		last, first, listlen, horizon, val, ok = PopFront(last, first, listlen, horizon)

		if !ok {
			t.Fatalf("PopFront at element %d failed", i)
		}

		if val != i {
			t.Errorf("PopFront: got %d, want %d", val, i)
		}
	}

	if listlen != 0 || last != nil || first != nil {
		t.Errorf("list should be empty after popping all elements")
	}
}

func TestCompactionEdgeCases(t *testing.T) {
	// Test PopFront compaction when last == first
	last := &Segment[int]{Elems: make([]int, 8)}
	for i := 0; i < 8; i++ {
		last.Elems[i] = i + 1
	}
	first := last
	listlen := 8
	horizon := 7

	// Pop until compaction threshold
	for horizon > len(last.Elems)>>1 {
		last, first, listlen, horizon, _, _ = PopFront(last, first, listlen, horizon)
	}

	// Next pop triggers compaction, and last == first
	oldLast := last
	last, first, listlen, horizon, _, _ = PopFront(last, first, listlen, horizon)

	if last == oldLast {
		t.Errorf("compaction should create a new segment")
	}

	// After compaction with last == first initially, first should also update
	if last != first {
		t.Errorf("first should equal last after compaction of single segment")
	}
}

func TestPopBackCompactionSingleSegment(t *testing.T) {
	// Test PopBack compaction when last == first
	last := &Segment[int]{Elems: []int{1, 2, 3, 4, 5, 6, 7, 8}}
	first := last
	listlen := 8
	horizon := 0

	// Pop until we reach compaction threshold (horizon == len/2)
	for i := 0; i < 3; i++ {
		last, first, listlen, horizon, _, _ = PopBack(last, first, listlen, horizon)
	}

	// Next pop should trigger compaction
	oldFirst := first
	last, first, listlen, horizon, val, ok := PopBack(last, first, listlen, horizon)

	if !ok || val != 4 {
		t.Errorf("PopBack: got (%d, %v), want (4, true)", val, ok)
	}

	if first == oldFirst {
		t.Errorf("compaction should create a new first segment")
	}

	// After compaction, verify remaining elements
	if listlen != 4 {
		t.Errorf("listlen = %d, want 4", listlen)
	}

	// Continue popping and verify values
	expected := []int{5, 6, 7, 8}
	for _, want := range expected {
		last, first, listlen, horizon, val, ok = PopBack(last, first, listlen, horizon)
		if !ok || val != want {
			t.Errorf("PopBack: got (%d, %v), want (%d, true)", val, ok, want)
		}
	}

	if listlen != 0 {
		t.Errorf("final listlen = %d, want 0", listlen)
	}
}

func TestPopBackCompactionMultiSegment(t *testing.T) {
	// Build a list with multiple segments
	var last, first *Segment[int]
	listlen := 0

	for i := 1; i <= 7; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	// This creates segments with binary decomposition: 7 = 4 + 2 + 1
	// Now pop from back to trigger compaction in the first segment
	horizon := 0

	// Pop elements until compaction triggers
	for i := 1; i <= 7; i++ {
		var val int
		var ok bool
		last, first, listlen, horizon, val, ok = PopBack(last, first, listlen, horizon)
		if !ok {
			t.Fatalf("PopBack #%d failed", i)
		}
		if val != i {
			t.Errorf("PopBack #%d: got %d, want %d", i, val, i)
		}
	}

	if listlen != 0 {
		t.Errorf("listlen should be 0, got %d", listlen)
	}
}

func TestMergePowerOfTwoSequence(t *testing.T) {
	// Test that merging up to each power of two results in single segment
	for pow := 0; pow <= 10; pow++ {
		n := 1 << pow
		var last, first *Segment[int]
		listlen := 0

		for i := 1; i <= n; i++ {
			last, first, listlen = Merge(i, last, first, listlen)
		}

		if last != first {
			t.Errorf("n=%d: expected single segment after merging power of two", n)
		}

		if len(last.Elems) != n {
			t.Errorf("n=%d: segment length = %d, want %d", n, len(last.Elems), n)
		}

		if last.Next != nil {
			t.Errorf("n=%d: last.Next should be nil", n)
		}
	}
}

func TestFindSingleSegmentWithHorizon(t *testing.T) {
	// Build a single segment and then pop some elements
	var last, first *Segment[int]
	listlen := 0

	for i := 1; i <= 8; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1

	// Pop 3 elements from front
	for i := 0; i < 3; i++ {
		last, first, listlen, horizon, _, _ = PopFront(last, first, listlen, horizon)
	}

	// Now listlen=5, elements are 1,2,3,4,5 (oldest to newest)
	// Index 0 should be 1
	slice, offset := Find(0, last, first, listlen, horizon)
	if slice[offset] != 1 {
		t.Errorf("Find(0) = %d, want 1", slice[offset])
	}

	// Index 4 should be 5
	slice, offset = Find(4, last, first, listlen, horizon)
	if slice[offset] != 5 {
		t.Errorf("Find(4) = %d, want 5", slice[offset])
	}

	// Index 2 should be 3
	slice, offset = Find(2, last, first, listlen, horizon)
	if slice[offset] != 3 {
		t.Errorf("Find(2) = %d, want 3", slice[offset])
	}
}

func TestMergeAfterCompleteEmpty(t *testing.T) {
	// Build list, drain it completely, then rebuild
	var last, first *Segment[int]
	listlen := 0

	// Build
	for i := 1; i <= 5; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1

	// Drain completely
	for listlen > 0 {
		last, first, listlen, horizon, _, _ = PopFront(last, first, listlen, horizon)
	}

	if last != nil || first != nil {
		t.Errorf("after complete drain: pointers should be nil")
	}

	// Rebuild
	for i := 10; i <= 15; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	if listlen != 6 {
		t.Errorf("listlen = %d, want 6", listlen)
	}

	horizon = len(last.Elems) - 1
	expected := []int{15, 14, 13, 12, 11, 10}
	verifyChain(t, last, first, listlen, horizon, expected)
}

func TestSplitLastMultipleTimes(t *testing.T) {
	// Build a list and split multiple times
	var last, first *Segment[int]
	listlen := 0

	for i := 1; i <= 16; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1

	// Split twice
	last, first, listlen, horizon = SplitLast(last, first, listlen, horizon)
	last, first, listlen, horizon = SplitLast(last, first, listlen, horizon)

	// Should still have all elements
	if listlen != 16 {
		t.Errorf("listlen = %d, want 16", listlen)
	}

	expected := make([]int, 16)
	for i := 0; i < 16; i++ {
		expected[i] = 16 - i
	}

	verifyChain(t, last, first, listlen, horizon, expected)
}

func TestPopFrontUntilAdvanceSegment(t *testing.T) {
	// Build a list with multiple segments, pop until advancing to next segment
	seg2 := &Segment[int]{Elems: []int{1, 2, 3}}
	seg1 := &Segment[int]{Elems: []int{4, 5}, Next: seg2}

	last := seg1
	first := seg2
	listlen := 5
	horizon := 1 // Two elements visible in last: [4, 5]

	// Pop once (removes 5)
	last, first, listlen, horizon, val, ok := PopFront(last, first, listlen, horizon)
	if !ok || val != 5 {
		t.Errorf("PopFront: got (%d, %v), want (5, true)", val, ok)
	}
	if horizon != 0 {
		t.Errorf("horizon = %d, want 0", horizon)
	}

	// Pop again (removes 4, should advance to next segment)
	last, first, listlen, horizon, val, ok = PopFront(last, first, listlen, horizon)
	if !ok || val != 4 {
		t.Errorf("PopFront: got (%d, %v), want (4, true)", val, ok)
	}

	// Should have advanced to seg2
	if last != seg2 {
		t.Errorf("should have advanced to next segment")
	}

	if horizon != len(last.Elems)-1 {
		t.Errorf("horizon = %d, want %d", horizon, len(last.Elems)-1)
	}

	// Continue popping
	expectedVals := []int{3, 2, 1}
	for i, want := range expectedVals {
		last, first, listlen, horizon, val, ok = PopFront(last, first, listlen, horizon)
		if !ok || val != want {
			t.Errorf("PopFront #%d: got (%d, %v), want (%d, true)", i, val, ok, want)
		}
	}

	if listlen != 0 {
		t.Errorf("listlen = %d, want 0", listlen)
	}
}

func TestFindInMultiSegmentList(t *testing.T) {
	// Build a specific multi-segment structure
	seg3 := &Segment[int]{Elems: []int{1, 2}, Next: nil}
	seg2 := &Segment[int]{Elems: []int{3, 4, 5, 6}, Next: seg3}
	seg1 := &Segment[int]{Elems: []int{7, 8}, Next: seg2}

	last := seg1
	first := seg3
	listlen := 8
	horizon := 1

	// Find each element
	expectedVals := []int{1, 2, 3, 4, 5, 6, 7, 8}
	for i, want := range expectedVals {
		slice, offset := Find(i, last, first, listlen, horizon)
		if slice[offset] != want {
			t.Errorf("Find(%d) = %d, want %d", i, slice[offset], want)
		}
	}
}

func TestPopBackRemovingLastSegment(t *testing.T) {
	// Create a list where popping will remove the entire first segment
	seg2 := &Segment[int]{Elems: []int{10}}
	seg1 := &Segment[int]{Elems: []int{20, 21}, Next: seg2}

	last := seg1
	first := seg2
	listlen := 3
	backHorizon := 0

	// Pop the only element in first segment
	last, first, listlen, backHorizon, val, ok := PopBack(last, first, listlen, backHorizon)

	if !ok || val != 10 {
		t.Errorf("PopBack: got (%d, %v), want (10, true)", val, ok)
	}

	// First should now be seg1
	if first != seg1 {
		t.Errorf("first should be seg1 after removing old first")
	}

	if listlen != 2 {
		t.Errorf("listlen = %d, want 2", listlen)
	}

	// Verify remaining by popping from back
	last, first, listlen, backHorizon, val, ok = PopBack(last, first, listlen, backHorizon)
	if !ok || val != 20 {
		t.Errorf("PopBack: got (%d, %v), want (20, true)", val, ok)
	}

	last, first, listlen, backHorizon, val, ok = PopBack(last, first, listlen, backHorizon)
	if !ok || val != 21 {
		t.Errorf("PopBack: got (%d, %v), want (21, true)", val, ok)
	}

	if listlen != 0 {
		t.Errorf("listlen = %d, want 0", listlen)
	}
}

func TestAlternatingPushPop(t *testing.T) {
	// Test pattern: build a list, pop some, add more, pop some more
	var last, first *Segment[int]
	listlen := 0

	// Build initial list
	for i := 1; i <= 10; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	horizon := len(last.Elems) - 1

	// Pop 5 elements from front
	for i := 0; i < 5; i++ {
		last, first, listlen, horizon, _, _ = PopFront(last, first, listlen, horizon)
	}

	// Now have elements 1-5, need to split before merging
	last, first, listlen, horizon = SplitLast(last, first, listlen, horizon)

	// Add 5 more elements
	for i := 11; i <= 15; i++ {
		last, first, listlen = Merge(i, last, first, listlen)
	}

	// Should have 10 elements: 1, 2, 3, 4, 5, 11, 12, 13, 14, 15
	if listlen != 10 {
		t.Errorf("listlen = %d, want 10", listlen)
	}

	horizon = len(last.Elems) - 1

	// Pop all and verify order (newest to oldest)
	expected := []int{15, 14, 13, 12, 11, 5, 4, 3, 2, 1}
	for i, want := range expected {
		var val int
		var ok bool
		last, first, listlen, horizon, val, ok = PopFront(last, first, listlen, horizon)
		if !ok {
			t.Fatalf("PopFront #%d failed", i)
		}
		if val != want {
			t.Errorf("PopFront #%d: got %d, want %d", i, val, want)
		}
	}

	if listlen != 0 {
		t.Errorf("final listlen = %d, want 0", listlen)
	}
}

func TestMergeWithDifferentTypes(t *testing.T) {
	// Test with string type
	var last, first *Segment[string]
	listlen := 0

	words := []string{"apple", "banana", "cherry", "date", "elderberry"}
	for _, word := range words {
		last, first, listlen = Merge(word, last, first, listlen)
	}

	if listlen != 5 {
		t.Errorf("listlen = %d, want 5", listlen)
	}

	horizon := len(last.Elems) - 1

	// Pop and verify
	for i := len(words) - 1; i >= 0; i-- {
		var val string
		var ok bool
		last, first, listlen, horizon, val, ok = PopFront(last, first, listlen, horizon)
		if !ok || val != words[i] {
			t.Errorf("PopFront: got (%s, %v), want (%s, true)", val, ok, words[i])
		}
	}
}

func TestSegmentCountBounds(t *testing.T) {
	// Verify segment count stays within O(log n) bounds
	testSizes := []int{10, 50, 100, 200, 500, 1000}

	for _, n := range testSizes {
		var last, first *Segment[int]
		listlen := 0

		for i := 1; i <= n; i++ {
			last, first, listlen = Merge(i, last, first, listlen)
		}

		// Count segments
		segmentCount := 0
		for seg := last; seg != nil; seg = seg.Next {
			segmentCount++
		}

		// Calculate max allowed segments: floor(log2(n)) + 1
		maxSegments := 0
		temp := n
		for temp > 0 {
			maxSegments++
			temp >>= 1
		}

		if segmentCount > maxSegments {
			t.Errorf("n=%d: segment count = %d, exceeds bound %d", n, segmentCount, maxSegments)
		}
	}
}
