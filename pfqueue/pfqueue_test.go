package pfqueue

import (
	"fmt"
	"testing"
)

// Helper function to verify queue contents by draining
func verifyQueue[E comparable](t *testing.T, q PFQueue[E], expected []E) {
	t.Helper()

	if q.Len() != len(expected) {
		t.Errorf("queue length = %d, want %d", q.Len(), len(expected))
	}

	for i, want := range expected {
		var val E
		var ok bool
		q, val, ok = q.PopBack()

		if !ok {
			t.Fatalf("PopBack #%d failed", i)
		}

		if val != want {
			t.Errorf("PopBack #%d: got %v, want %v", i, val, want)
		}
	}

	if q.Len() != 0 {
		t.Errorf("queue should be empty after draining, len=%d", q.Len())
	}
}

func TestEmptyQueue(t *testing.T) {
	var q PFQueue[int]

	if q.Len() != 0 {
		t.Errorf("empty queue length = %d, want 0", q.Len())
	}

	// PopBack on empty
	_, _, ok := q.PopBack()
	if ok {
		t.Errorf("PopBack on empty queue should return ok=false")
	}

	// Front on empty
	_, ok = q.Front()
	if ok {
		t.Errorf("Front on empty queue should return ok=false")
	}

	// Back on empty
	_, ok = q.Back()
	if ok {
		t.Errorf("Back on empty queue should return ok=false")
	}
}

func TestPushFrontSingle(t *testing.T) {
	var q PFQueue[int]

	q = q.PushFront(42)

	if q.Len() != 1 {
		t.Errorf("length = %d, want 1", q.Len())
	}

	val, ok := q.Front()
	if !ok || val != 42 {
		t.Errorf("Front() = (%d, %v), want (42, true)", val, ok)
	}

	val, ok = q.Back()
	if !ok || val != 42 {
		t.Errorf("Back() = (%d, %v), want (42, true)", val, ok)
	}
}

func TestPushFrontMultiple(t *testing.T) {
	var q PFQueue[int]

	// Enqueue 1, 2, 3, 4, 5
	for i := 1; i <= 5; i++ {
		q = q.PushFront(i)
	}

	if q.Len() != 5 {
		t.Errorf("length = %d, want 5", q.Len())
	}

	// Front should be 5 (last enqueued)
	val, ok := q.Front()
	if !ok || val != 5 {
		t.Errorf("Front() = (%d, %v), want (5, true)", val, ok)
	}

	// Back should be 1 (first enqueued)
	val, ok = q.Back()
	if !ok || val != 1 {
		t.Errorf("Back() = (%d, %v), want (1, true)", val, ok)
	}
}

func TestPopBackSingle(t *testing.T) {
	var q PFQueue[int]

	q = q.PushFront(42)
	q, val, ok := q.PopBack()

	if !ok || val != 42 {
		t.Errorf("PopBack() = (%d, %v), want (42, true)", val, ok)
	}

	if q.Len() != 0 {
		t.Errorf("length after pop = %d, want 0", q.Len())
	}
}

func TestPopBackMultiple(t *testing.T) {
	var q PFQueue[int]

	// Enqueue 1, 2, 3, 4, 5
	for i := 1; i <= 5; i++ {
		q = q.PushFront(i)
	}

	// Dequeue should give FIFO order: 1, 2, 3, 4, 5
	expected := []int{1, 2, 3, 4, 5}
	for i, want := range expected {
		var val int
		var ok bool
		q, val, ok = q.PopBack()

		if !ok {
			t.Fatalf("PopBack #%d failed", i)
		}

		if val != want {
			t.Errorf("PopBack #%d: got %d, want %d", i, val, want)
		}

		if q.Len() != 5-i-1 {
			t.Errorf("PopBack #%d: length = %d, want %d", i, q.Len(), 5-i-1)
		}
	}

	if q.Len() != 0 {
		t.Errorf("final length = %d, want 0", q.Len())
	}
}

func TestFIFOBehavior(t *testing.T) {
	var q PFQueue[int]

	// Enqueue 1 to 10
	for i := 1; i <= 10; i++ {
		q = q.PushFront(i)
	}

	// Dequeue should give FIFO order: 1, 2, 3, ..., 10
	for i := 1; i <= 10; i++ {
		var val int
		var ok bool
		q, val, ok = q.PopBack()

		if !ok || val != i {
			t.Errorf("dequeue: got (%d, %v), want (%d, true)", val, ok, i)
		}
	}

	if q.Len() != 0 {
		t.Errorf("queue should be empty, len=%d", q.Len())
	}
}

func TestFrontBackPeek(t *testing.T) {
	var q PFQueue[int]

	// Enqueue 1, 2, 3
	q = q.PushFront(1)
	q = q.PushFront(2)
	q = q.PushFront(3)

	// Peek should not modify queue
	originalLen := q.Len()

	val, ok := q.Front()
	if !ok || val != 3 {
		t.Errorf("Front() = (%d, %v), want (3, true)", val, ok)
	}

	if q.Len() != originalLen {
		t.Errorf("Front() modified length")
	}

	val, ok = q.Back()
	if !ok || val != 1 {
		t.Errorf("Back() = (%d, %v), want (1, true)", val, ok)
	}

	if q.Len() != originalLen {
		t.Errorf("Back() modified length")
	}
}

func TestLargeScale(t *testing.T) {
	var q PFQueue[int]

	// Enqueue 1000 elements
	for i := 1; i <= 1000; i++ {
		q = q.PushFront(i)
	}

	if q.Len() != 1000 {
		t.Errorf("length = %d, want 1000", q.Len())
	}

	// Dequeue all in FIFO order
	for i := 1; i <= 1000; i++ {
		var val int
		var ok bool
		q, val, ok = q.PopBack()

		if !ok || val != i {
			t.Errorf("dequeue #%d: got (%d, %v), want (%d, true)", i, val, ok, i)
			break
		}
	}

	if q.Len() != 0 {
		t.Errorf("final length = %d, want 0", q.Len())
	}
}

func TestImmutability(t *testing.T) {
	var q1 PFQueue[int]

	// Build q1
	q1 = q1.PushFront(1)
	q1 = q1.PushFront(2)
	q1 = q1.PushFront(3)

	// Create q2 from q1
	q2 := q1.PushFront(4)

	// q1 should still have 3 elements
	if q1.Len() != 3 {
		t.Errorf("q1 length = %d, want 3", q1.Len())
	}

	// q2 should have 4 elements
	if q2.Len() != 4 {
		t.Errorf("q2 length = %d, want 4", q2.Len())
	}

	// Modify q1 by popping
	q1, _, _ = q1.PopBack()

	// q1 should now have 2 elements
	if q1.Len() != 2 {
		t.Errorf("q1 length after pop = %d, want 2", q1.Len())
	}

	// q2 should still have 4 elements (unchanged)
	if q2.Len() != 4 {
		t.Errorf("q2 length = %d, want 4 (should be unchanged)", q2.Len())
	}
}

func TestPersistentBehavior(t *testing.T) {
	var q1 PFQueue[int]

	// Build q1: 1, 2, 3, 4
	q1 = q1.PushFront(1)
	q1 = q1.PushFront(2)
	q1 = q1.PushFront(3)
	q1 = q1.PushFront(4)

	// Create multiple versions
	q2 := q1.PushFront(5)    // 1, 2, 3, 4, 5
	q3, _, _ := q1.PopBack() // 2, 3, 4

	// Verify q1 unchanged
	if q1.Len() != 4 {
		t.Errorf("q1 length = %d, want 4", q1.Len())
	}

	// Verify q2: 1, 2, 3, 4, 5
	expected2 := []int{1, 2, 3, 4, 5}
	verifyQueue(t, q2, expected2)

	// Verify q3: 2, 3, 4
	expected3 := []int{2, 3, 4}
	verifyQueue(t, q3, expected3)

	// Verify q1 still: 1, 2, 3, 4
	expected1 := []int{1, 2, 3, 4}
	verifyQueue(t, q1, expected1)
}

func TestEnqueueDequeuePattern(t *testing.T) {
	var q PFQueue[int]

	// Pattern: enqueue 2, dequeue 1, repeat
	for i := 1; i <= 20; i++ {
		q = q.PushFront(i * 2)
		q = q.PushFront(i*2 + 1)

		if q.Len() > 0 {
			q, _, _ = q.PopBack()
		}
	}

	// Should have 20 elements remaining
	if q.Len() != 20 {
		t.Errorf("length = %d, want 20", q.Len())
	}

	// Drain and verify no panics
	count := 0
	for q.Len() > 0 {
		q, _, _ = q.PopBack()
		count++
		if count > 100 {
			t.Fatalf("infinite loop detected")
		}
	}

	if count != 20 {
		t.Errorf("drained %d elements, want 20", count)
	}
}

func TestWithStrings(t *testing.T) {
	var q PFQueue[string]

	words := []string{"apple", "banana", "cherry", "date"}

	// Enqueue words
	for _, word := range words {
		q = q.PushFront(word)
	}

	// Dequeue in FIFO order
	for i, want := range words {
		var val string
		var ok bool
		q, val, ok = q.PopBack()

		if !ok || val != want {
			t.Errorf("dequeue #%d: got (%s, %v), want (%s, true)", i, val, ok, want)
		}
	}
}

func TestFrontAfterOperations(t *testing.T) {
	var q PFQueue[int]

	// Enqueue 1, 2, 3
	q = q.PushFront(1)
	q = q.PushFront(2)
	q = q.PushFront(3)

	// Front should be 3
	front, _ := q.Front()
	if front != 3 {
		t.Errorf("Front = %d, want 3", front)
	}

	// Dequeue (removes 1)
	q, _, _ = q.PopBack()

	// Front should still be 3
	front, _ = q.Front()
	if front != 3 {
		t.Errorf("Front after dequeue = %d, want 3", front)
	}

	// Back should be 2
	back, _ := q.Back()
	if back != 2 {
		t.Errorf("Back after dequeue = %d, want 2", back)
	}
}

func TestBackAfterOperations(t *testing.T) {
	var q PFQueue[int]

	// Enqueue 1, 2, 3
	q = q.PushFront(1)
	q = q.PushFront(2)
	q = q.PushFront(3)

	// Back should be 1
	back, _ := q.Back()
	if back != 1 {
		t.Errorf("Back = %d, want 1", back)
	}

	// Dequeue twice (removes 1 and 2)
	q, _, _ = q.PopBack()
	q, _, _ = q.PopBack()

	// Back should now be 3
	back, _ = q.Back()
	if back != 3 {
		t.Errorf("Back after dequeues = %d, want 3", back)
	}

	// Front should also be 3 (single element)
	front, _ := q.Front()
	if front != 3 {
		t.Errorf("Front after dequeues = %d, want 3", front)
	}
}

func TestEmptyAfterOperations(t *testing.T) {
	var q PFQueue[int]

	// Build and drain multiple times
	for cycle := 0; cycle < 3; cycle++ {
		// Build
		for i := 1; i <= 10; i++ {
			q = q.PushFront(i)
		}

		// Drain
		for q.Len() > 0 {
			q, _, _ = q.PopBack()
		}

		// Verify empty
		if q.Len() != 0 {
			t.Errorf("cycle %d: length = %d, want 0", cycle, q.Len())
		}

		_, ok := q.Front()
		if ok {
			t.Errorf("cycle %d: Front() on empty should return ok=false", cycle)
		}

		_, ok = q.Back()
		if ok {
			t.Errorf("cycle %d: Back() on empty should return ok=false", cycle)
		}
	}
}

func TestLengthTracking(t *testing.T) {
	var q PFQueue[int]

	lengths := []int{}

	// Track length changes
	lengths = append(lengths, q.Len()) // 0

	q = q.PushFront(1)
	lengths = append(lengths, q.Len()) // 1

	q = q.PushFront(2)
	lengths = append(lengths, q.Len()) // 2

	q = q.PushFront(3)
	lengths = append(lengths, q.Len()) // 3

	q, _, _ = q.PopBack()
	lengths = append(lengths, q.Len()) // 2

	q, _, _ = q.PopBack()
	lengths = append(lengths, q.Len()) // 1

	q, _, _ = q.PopBack()
	lengths = append(lengths, q.Len()) // 0

	expected := []int{0, 1, 2, 3, 2, 1, 0}
	for i, want := range expected {
		if lengths[i] != want {
			t.Errorf("length[%d] = %d, want %d", i, lengths[i], want)
		}
	}
}

func TestPushAfterPopOptimization(t *testing.T) {
	// This test exercises the SplitLast optimization in PushFront
	var q PFQueue[int]

	// Enqueue 1 to 10
	for i := 1; i <= 10; i++ {
		q = q.PushFront(i)
	}

	// Dequeue 5 elements (this leaves horizonFirst > 0)
	for i := 0; i < 5; i++ {
		q, _, _ = q.PopBack()
	}

	// Now enqueue more - this should trigger the SplitLast optimization
	q = q.PushFront(100)
	q = q.PushFront(200)

	// Verify correctness
	if q.Len() != 7 {
		t.Errorf("length = %d, want 7", q.Len())
	}

	// Dequeue and verify order: 6, 7, 8, 9, 10, 100, 200
	expected := []int{6, 7, 8, 9, 10, 100, 200}
	verifyQueue(t, q, expected)
}

func TestAlternatingEnqueueDequeue(t *testing.T) {
	var q PFQueue[int]

	// Simulate sliding window: enqueue, then dequeue oldest
	for i := 1; i <= 20; i++ {
		q = q.PushFront(i)

		// Keep window size at 5
		if i > 5 {
			q, _, _ = q.PopBack()
		}
	}

	// Should have last 5 elements: 16, 17, 18, 19, 20
	if q.Len() != 5 {
		t.Errorf("length = %d, want 5", q.Len())
	}

	expected := []int{16, 17, 18, 19, 20}
	verifyQueue(t, q, expected)
}

func TestSingleElementOperations(t *testing.T) {
	var q PFQueue[int]

	q = q.PushFront(42)

	// Front and Back should be same
	front, ok1 := q.Front()
	back, ok2 := q.Back()

	if !ok1 || !ok2 {
		t.Errorf("Front/Back on single element should succeed")
	}

	if front != 42 || back != 42 {
		t.Errorf("Front=%d, Back=%d, want both=42", front, back)
	}

	// Pop the element
	q, val, ok := q.PopBack()

	if !ok || val != 42 {
		t.Errorf("PopBack: got (%d, %v), want (42, true)", val, ok)
	}

	// Now empty
	_, ok = q.Front()
	if ok {
		t.Errorf("Front on empty after pop should fail")
	}
}

func TestMultipleVersions(t *testing.T) {
	var base PFQueue[int]

	// Create base queue
	base = base.PushFront(1)
	base = base.PushFront(2)

	// Create multiple derived versions
	v1 := base.PushFront(3)
	v2 := base.PushFront(4)
	v3, _, _ := base.PopBack()

	// Verify all versions
	if base.Len() != 2 {
		t.Errorf("base length = %d, want 2", base.Len())
	}

	if v1.Len() != 3 {
		t.Errorf("v1 length = %d, want 3", v1.Len())
	}

	if v2.Len() != 3 {
		t.Errorf("v2 length = %d, want 3", v2.Len())
	}

	if v3.Len() != 1 {
		t.Errorf("v3 length = %d, want 1", v3.Len())
	}

	// Verify contents don't interfere
	front1, _ := v1.Front()
	front2, _ := v2.Front()

	if front1 != 3 {
		t.Errorf("v1 front = %d, want 3", front1)
	}

	if front2 != 4 {
		t.Errorf("v2 front = %d, want 4", front2)
	}
}

func TestStressTest(t *testing.T) {
	var q PFQueue[int]

	// Mix enqueue and dequeue operations
	for i := 0; i < 100; i++ {
		// Enqueue
		q = q.PushFront(i)

		// Occasionally dequeue
		if i%3 == 0 && q.Len() > 1 {
			q, _, _ = q.PopBack()
		}
	}

	// Should have elements remaining
	if q.Len() == 0 {
		t.Errorf("queue should not be empty")
	}

	// Drain completely
	count := 0
	for q.Len() > 0 {
		q, _, _ = q.PopBack()
		count++
		if count > 200 {
			t.Fatalf("infinite loop detected")
		}
	}
}

func TestPopBackCompaction(t *testing.T) {
	var q PFQueue[int]

	// Build queue with enough elements to trigger compaction
	for i := 1; i <= 16; i++ {
		q = q.PushFront(i)
	}

	// Dequeue multiple times to trigger compaction in PopBack
	for i := 0; i < 8; i++ {
		var val int
		var ok bool
		q, val, ok = q.PopBack()

		if !ok {
			t.Fatalf("PopBack #%d failed", i)
		}

		if val != i+1 {
			t.Errorf("PopBack #%d: got %d, want %d", i, val, i+1)
		}
	}

	// Verify remaining elements
	if q.Len() != 8 {
		t.Errorf("length = %d, want 8", q.Len())
	}

	// Verify correct order of remaining elements: 9 to 16
	for i := 9; i <= 16; i++ {
		var val int
		var ok bool
		q, val, ok = q.PopBack()

		if !ok || val != i {
			t.Errorf("dequeue: got (%d, %v), want (%d, true)", val, ok, i)
		}
	}
}

func TestZeroValue(t *testing.T) {
	// Test that zero value is a valid empty queue
	var q PFQueue[int]

	if q.Len() != 0 {
		t.Errorf("zero value length = %d, want 0", q.Len())
	}

	// Should be able to use immediately
	q = q.PushFront(1)

	if q.Len() != 1 {
		t.Errorf("length after push = %d, want 1", q.Len())
	}

	val, ok := q.Front()
	if !ok || val != 1 {
		t.Errorf("Front() = (%d, %v), want (1, true)", val, ok)
	}
}

func TestPushFrontWithoutSplit(t *testing.T) {
	// Test PushFront path that doesn't need SplitLast
	var q PFQueue[int]

	// First push - no split needed
	q = q.PushFront(1)
	if q.Len() != 1 {
		t.Errorf("length = %d, want 1", q.Len())
	}

	// More pushes without any pops - no split needed
	q = q.PushFront(2)
	q = q.PushFront(3)
	q = q.PushFront(4)

	if q.Len() != 4 {
		t.Errorf("length = %d, want 4", q.Len())
	}

	// Verify FIFO order
	expected := []int{1, 2, 3, 4}
	verifyQueue(t, q, expected)
}

func TestPushFrontWithSplitNeeded(t *testing.T) {
	// Test PushFront path that requires SplitLast
	var q PFQueue[int]

	// Build initial queue
	for i := 1; i <= 8; i++ {
		q = q.PushFront(i)
	}

	// Pop several to create horizonFirst > 0
	for i := 0; i < 4; i++ {
		q, _, _ = q.PopBack()
	}

	// At this point, horizonLast should be < len(last.Elems)-1
	// Next push should trigger SplitLast

	q = q.PushFront(99)

	if q.Len() != 5 {
		t.Errorf("length = %d, want 5", q.Len())
	}

	// Verify order: 5, 6, 7, 8, 99
	expected := []int{5, 6, 7, 8, 99}
	verifyQueue(t, q, expected)
}

func TestPopBackSingleSegmentCompaction(t *testing.T) {
	// Test PopBack compaction when last == first
	var q PFQueue[int]

	// Build a queue that will coalesce into single segment
	for i := 1; i <= 8; i++ {
		q = q.PushFront(i)
	}

	// Pop several times to trigger compaction in single segment
	for i := 1; i <= 4; i++ {
		var val int
		var ok bool
		q, val, ok = q.PopBack()

		if !ok || val != i {
			t.Errorf("PopBack #%d: got (%d, %v), want (%d, true)", i, val, ok, i)
		}
	}

	// Continue popping - should handle horizonLast correctly
	for i := 5; i <= 8; i++ {
		var val int
		var ok bool
		q, val, ok = q.PopBack()

		if !ok || val != i {
			t.Errorf("PopBack #%d: got (%d, %v), want (%d, true)", i, val, ok, i)
		}
	}

	if q.Len() != 0 {
		t.Errorf("final length = %d, want 0", q.Len())
	}
}

func TestPushFrontEmptyQueue(t *testing.T) {
	// Test pushing to empty queue (special case in SplitLast check)
	var q PFQueue[int]

	// len == 0, so SplitLast optimization shouldn't trigger
	q = q.PushFront(42)

	if q.Len() != 1 {
		t.Errorf("length = %d, want 1", q.Len())
	}

	val, ok := q.Front()
	if !ok || val != 42 {
		t.Errorf("Front() = (%d, %v), want (42, true)", val, ok)
	}
}

func TestMergeHelper(t *testing.T) {
	// Indirectly test the merge helper through PushFront
	var q PFQueue[int]

	// Multiple pushes should all use merge
	for i := 1; i <= 10; i++ {
		q = q.PushFront(i)
	}

	// Verify all elements are present
	if q.Len() != 10 {
		t.Errorf("length = %d, want 10", q.Len())
	}

	expected := make([]int, 10)
	for i := 0; i < 10; i++ {
		expected[i] = i + 1
	}

	verifyQueue(t, q, expected)
}

func TestSplitLastOptimizationPath(t *testing.T) {
	// Very specific test to trigger the SplitLast optimization
	var q PFQueue[int]

	// Build a queue with 16 elements (power of 2, will coalesce)
	for i := 1; i <= 16; i++ {
		q = q.PushFront(i)
	}

	// Pop 8 elements from back
	for i := 0; i < 8; i++ {
		q, _, _ = q.PopBack()
	}

	// At this point:
	// - m.len > 0 (we have 8 elements)
	// - m.horizonLast might not equal len(m.last.Elems)-1 if there was compaction
	// Now push should trigger SplitLast

	q = q.PushFront(100)
	q = q.PushFront(101)

	// Verify correctness: should have 9, 10, 11, 12, 13, 14, 15, 16, 100, 101
	expected := []int{9, 10, 11, 12, 13, 14, 15, 16, 100, 101}
	verifyQueue(t, q, expected)
}

func TestSplitLastConditionCoverage(t *testing.T) {
	// Test to ensure we cover the if condition properly
	var q PFQueue[int]

	// Build initial queue
	q = q.PushFront(1)
	q = q.PushFront(2)
	q = q.PushFront(3)
	q = q.PushFront(4)

	// Pop one to potentially create horizonFirst > 0
	q, _, _ = q.PopBack()

	// This push might need SplitLast
	q = q.PushFront(5)

	// Verify: should have 2, 3, 4, 5
	expected := []int{2, 3, 4, 5}
	verifyQueue(t, q, expected)
}

// Benchmarks

func BenchmarkPFQueuePushFront(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				q := PFQueue[int]{}
				for j := 0; j < size; j++ {
					q = q.PushFront(j)
				}
			}
		})
	}
}

func BenchmarkPFQueuePopBack(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		setup := PFQueue[int]{}
		for j := 0; j < size; j++ {
			setup = setup.PushFront(j)
		}

		b.Run(fmt.Sprintf("n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				q := setup
				for j := 0; j < size; j++ {
					var val int
					q, val, _ = q.PopBack()
					_ = val
				}
			}
		})
	}
}
