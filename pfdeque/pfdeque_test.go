package pfdeque

import (
	"fmt"
	"testing"
)

// Helper function to verify deque contents by draining from front
func verifyDequeFromFront[E comparable](t *testing.T, d PFDeque[E], expected []E) {
	t.Helper()

	if d.Len() != len(expected) {
		t.Errorf("deque length = %d, want %d", d.Len(), len(expected))
	}

	for i, want := range expected {
		var val E
		var ok bool
		d, val, ok = d.PopFront()

		if !ok {
			t.Fatalf("PopFront #%d failed", i)
		}

		if val != want {
			t.Errorf("PopFront #%d: got %v, want %v", i, val, want)
		}
	}

	if d.Len() != 0 {
		t.Errorf("deque should be empty after draining, len=%d", d.Len())
	}
}

// Helper function to verify deque contents by draining from back
func verifyDequeFromBack[E comparable](t *testing.T, d PFDeque[E], expected []E) {
	t.Helper()

	if d.Len() != len(expected) {
		t.Errorf("deque length = %d, want %d", d.Len(), len(expected))
	}

	for i, want := range expected {
		var val E
		var ok bool
		d, val, ok = d.PopBack()

		if !ok {
			t.Fatalf("PopBack #%d failed", i)
		}

		if val != want {
			t.Errorf("PopBack #%d: got %v, want %v", i, val, want)
		}
	}

	if d.Len() != 0 {
		t.Errorf("deque should be empty after draining, len=%d", d.Len())
	}
}

func TestEmptyDeque(t *testing.T) {
	var d PFDeque[int]

	if d.Len() != 0 {
		t.Errorf("empty deque length = %d, want 0", d.Len())
	}

	// PopFront on empty
	_, _, ok := d.PopFront()
	if ok {
		t.Errorf("PopFront on empty deque should return ok=false")
	}

	// PopBack on empty
	_, _, ok = d.PopBack()
	if ok {
		t.Errorf("PopBack on empty deque should return ok=false")
	}

	// Front on empty
	_, ok = d.Front()
	if ok {
		t.Errorf("Front on empty deque should return ok=false")
	}

	// Back on empty
	_, ok = d.Back()
	if ok {
		t.Errorf("Back on empty deque should return ok=false")
	}
}

func TestPushFrontSingle(t *testing.T) {
	var d PFDeque[int]

	d = d.PushFront(42)

	if d.Len() != 1 {
		t.Errorf("length = %d, want 1", d.Len())
	}

	val, ok := d.Front()
	if !ok || val != 42 {
		t.Errorf("Front() = (%d, %v), want (42, true)", val, ok)
	}

	val, ok = d.Back()
	if !ok || val != 42 {
		t.Errorf("Back() = (%d, %v), want (42, true)", val, ok)
	}
}

func TestPushBackSingle(t *testing.T) {
	var d PFDeque[int]

	d = d.PushBack(42)

	if d.Len() != 1 {
		t.Errorf("length = %d, want 1", d.Len())
	}

	val, ok := d.Front()
	if !ok || val != 42 {
		t.Errorf("Front() = (%d, %v), want (42, true)", val, ok)
	}

	val, ok = d.Back()
	if !ok || val != 42 {
		t.Errorf("Back() = (%d, %v), want (42, true)", val, ok)
	}
}

func TestPushFrontMultiple(t *testing.T) {
	var d PFDeque[int]

	// Push 1, 2, 3, 4, 5 to front
	for i := 1; i <= 5; i++ {
		d = d.PushFront(i)
	}

	if d.Len() != 5 {
		t.Errorf("length = %d, want 5", d.Len())
	}

	// Front should be 5 (last pushed)
	val, ok := d.Front()
	if !ok || val != 5 {
		t.Errorf("Front() = (%d, %v), want (5, true)", val, ok)
	}

	// Back should be 1 (first pushed)
	val, ok = d.Back()
	if !ok || val != 1 {
		t.Errorf("Back() = (%d, %v), want (1, true)", val, ok)
	}

	// Pop from front should give 5, 4, 3, 2, 1
	expected := []int{5, 4, 3, 2, 1}
	verifyDequeFromFront(t, d, expected)
}

func TestPushBackMultiple(t *testing.T) {
	var d PFDeque[int]

	// Push 1, 2, 3, 4, 5 to back
	for i := 1; i <= 5; i++ {
		d = d.PushBack(i)
	}

	if d.Len() != 5 {
		t.Errorf("length = %d, want 5", d.Len())
	}

	// Front should be 1 (first pushed)
	val, ok := d.Front()
	if !ok || val != 1 {
		t.Errorf("Front() = (%d, %v), want (1, true)", val, ok)
	}

	// Back should be 5 (last pushed)
	val, ok = d.Back()
	if !ok || val != 5 {
		t.Errorf("Back() = (%d, %v), want (5, true)", val, ok)
	}

	// Pop from front should give 1, 2, 3, 4, 5
	expected := []int{1, 2, 3, 4, 5}
	verifyDequeFromFront(t, d, expected)
}

func TestPopFrontSingle(t *testing.T) {
	var d PFDeque[int]

	d = d.PushFront(42)
	d, val, ok := d.PopFront()

	if !ok || val != 42 {
		t.Errorf("PopFront() = (%d, %v), want (42, true)", val, ok)
	}

	if d.Len() != 0 {
		t.Errorf("length after pop = %d, want 0", d.Len())
	}
}

func TestPopBackSingle(t *testing.T) {
	var d PFDeque[int]

	d = d.PushBack(42)
	d, val, ok := d.PopBack()

	if !ok || val != 42 {
		t.Errorf("PopBack() = (%d, %v), want (42, true)", val, ok)
	}

	if d.Len() != 0 {
		t.Errorf("length after pop = %d, want 0", d.Len())
	}
}

func TestPopFrontMultiple(t *testing.T) {
	var d PFDeque[int]

	// Build: 1, 2, 3, 4, 5 (front to back)
	for i := 1; i <= 5; i++ {
		d = d.PushBack(i)
	}

	// Pop from front
	expected := []int{1, 2, 3, 4, 5}
	for i, want := range expected {
		var val int
		var ok bool
		d, val, ok = d.PopFront()

		if !ok {
			t.Fatalf("PopFront #%d failed", i)
		}

		if val != want {
			t.Errorf("PopFront #%d: got %d, want %d", i, val, want)
		}

		if d.Len() != 5-i-1 {
			t.Errorf("PopFront #%d: length = %d, want %d", i, d.Len(), 5-i-1)
		}
	}
}

func TestPopBackMultiple(t *testing.T) {
	var d PFDeque[int]

	// Build: 1, 2, 3, 4, 5 (front to back)
	for i := 1; i <= 5; i++ {
		d = d.PushBack(i)
	}

	// Pop from back
	expected := []int{5, 4, 3, 2, 1}
	for i, want := range expected {
		var val int
		var ok bool
		d, val, ok = d.PopBack()

		if !ok {
			t.Fatalf("PopBack #%d failed", i)
		}

		if val != want {
			t.Errorf("PopBack #%d: got %d, want %d", i, val, want)
		}

		if d.Len() != 5-i-1 {
			t.Errorf("PopBack #%d: length = %d, want %d", i, d.Len(), 5-i-1)
		}
	}
}

func TestMixedPushFrontBack(t *testing.T) {
	var d PFDeque[int]

	// Push to front: 2, 1
	d = d.PushFront(2)
	d = d.PushFront(1)

	// Push to back: 3, 4
	d = d.PushBack(3)
	d = d.PushBack(4)

	// Deque should be: 1, 2, 3, 4 (front to back)
	if d.Len() != 4 {
		t.Errorf("length = %d, want 4", d.Len())
	}

	val, ok := d.Front()
	if !ok || val != 1 {
		t.Errorf("Front() = (%d, %v), want (1, true)", val, ok)
	}

	val, ok = d.Back()
	if !ok || val != 4 {
		t.Errorf("Back() = (%d, %v), want (4, true)", val, ok)
	}

	// Verify by popping from front
	expected := []int{1, 2, 3, 4}
	verifyDequeFromFront(t, d, expected)
}

func TestMixedPopFrontBack(t *testing.T) {
	var d PFDeque[int]

	// Build: 1, 2, 3, 4, 5
	for i := 1; i <= 5; i++ {
		d = d.PushBack(i)
	}

	// Pop from front (removes 1)
	d, val, ok := d.PopFront()
	if !ok || val != 1 {
		t.Errorf("PopFront() = (%d, %v), want (1, true)", val, ok)
	}

	// Pop from back (removes 5)
	d, val, ok = d.PopBack()
	if !ok || val != 5 {
		t.Errorf("PopBack() = (%d, %v), want (5, true)", val, ok)
	}

	// Should have: 2, 3, 4
	if d.Len() != 3 {
		t.Errorf("length = %d, want 3", d.Len())
	}

	expected := []int{2, 3, 4}
	verifyDequeFromFront(t, d, expected)
}

func TestAlternatePushPop(t *testing.T) {
	var d PFDeque[int]

	// Push to front and back alternately
	d = d.PushFront(2) // 2
	d = d.PushBack(3)  // 2, 3
	d = d.PushFront(1) // 1, 2, 3
	d = d.PushBack(4)  // 1, 2, 3, 4
	d = d.PushFront(0) // 0, 1, 2, 3, 4
	d = d.PushBack(5)  // 0, 1, 2, 3, 4, 5

	if d.Len() != 6 {
		t.Errorf("length = %d, want 6", d.Len())
	}

	// Pop alternately
	d, val, ok := d.PopFront()
	if !ok || val != 0 {
		t.Errorf("PopFront() = (%d, %v), want (0, true)", val, ok)
	}

	d, val, ok = d.PopBack()
	if !ok || val != 5 {
		t.Errorf("PopBack() = (%d, %v), want (5, true)", val, ok)
	}

	// Should have: 1, 2, 3, 4
	expected := []int{1, 2, 3, 4}
	verifyDequeFromFront(t, d, expected)
}

func TestPopFrontFromRightSide(t *testing.T) {
	var d PFDeque[int]

	// Push only to back (right side)
	for i := 1; i <= 5; i++ {
		d = d.PushBack(i)
	}

	// Pop from front should pull from right side when left is empty
	// Right side has: 1, 2, 3, 4, 5 (oldest to newest)
	// PopFront should give 1 (oldest on right = back of right)
	d, val, ok := d.PopFront()
	if !ok || val != 1 {
		t.Errorf("PopFront() = (%d, %v), want (1, true)", val, ok)
	}

	if d.Len() != 4 {
		t.Errorf("length = %d, want 4", d.Len())
	}
}

func TestPopBackFromLeftSide(t *testing.T) {
	var d PFDeque[int]

	// Push only to front (left side)
	for i := 1; i <= 5; i++ {
		d = d.PushFront(i)
	}

	// Pop from back should pull from left side when right is empty
	// Left side has: 5, 4, 3, 2, 1 (newest to oldest)
	// PopBack should give 1 (oldest on left = back of left)
	d, val, ok := d.PopBack()
	if !ok || val != 1 {
		t.Errorf("PopBack() = (%d, %v), want (1, true)", val, ok)
	}

	if d.Len() != 4 {
		t.Errorf("length = %d, want 4", d.Len())
	}
}

func TestFrontBackPeek(t *testing.T) {
	var d PFDeque[int]

	// Build: 1, 2, 3
	d = d.PushBack(1)
	d = d.PushBack(2)
	d = d.PushBack(3)

	// Peek should not modify deque
	originalLen := d.Len()

	val, ok := d.Front()
	if !ok || val != 1 {
		t.Errorf("Front() = (%d, %v), want (1, true)", val, ok)
	}

	if d.Len() != originalLen {
		t.Errorf("Front() modified length")
	}

	val, ok = d.Back()
	if !ok || val != 3 {
		t.Errorf("Back() = (%d, %v), want (3, true)", val, ok)
	}

	if d.Len() != originalLen {
		t.Errorf("Back() modified length")
	}
}

func TestLargeScalePushFront(t *testing.T) {
	var d PFDeque[int]

	// Push 1000 elements to front
	for i := 1; i <= 1000; i++ {
		d = d.PushFront(i)
	}

	if d.Len() != 1000 {
		t.Errorf("length = %d, want 1000", d.Len())
	}

	// Verify order (should be 1000, 999, ..., 1)
	for i := 1000; i >= 1; i-- {
		var val int
		var ok bool
		d, val, ok = d.PopFront()

		if !ok || val != i {
			t.Errorf("PopFront: got (%d, %v), want (%d, true)", val, ok, i)
			break
		}
	}

	if d.Len() != 0 {
		t.Errorf("final length = %d, want 0", d.Len())
	}
}

func TestLargeScalePushBack(t *testing.T) {
	var d PFDeque[int]

	// Push 1000 elements to back
	for i := 1; i <= 1000; i++ {
		d = d.PushBack(i)
	}

	if d.Len() != 1000 {
		t.Errorf("length = %d, want 1000", d.Len())
	}

	// Verify order (should be 1, 2, ..., 1000)
	for i := 1; i <= 1000; i++ {
		var val int
		var ok bool
		d, val, ok = d.PopFront()

		if !ok || val != i {
			t.Errorf("PopFront: got (%d, %v), want (%d, true)", val, ok, i)
			break
		}
	}

	if d.Len() != 0 {
		t.Errorf("final length = %d, want 0", d.Len())
	}
}

func TestDrainFromBothEnds(t *testing.T) {
	var d PFDeque[int]

	// Build by alternating front and back to use both sides
	for i := 1; i <= 5; i++ {
		d = d.PushBack(i)   // 1, 2, 3, 4, 5 on right
		d = d.PushFront(-i) // -1, -2, -3, -4, -5 on left
	}

	// Deque is: -5, -4, -3, -2, -1, 1, 2, 3, 4, 5

	// Drain from front
	expectedFront := []int{-5, -4, -3, -2, -1}
	for _, want := range expectedFront {
		var val int
		var ok bool
		d, val, ok = d.PopFront()

		if !ok || val != want {
			t.Errorf("PopFront: got (%d, %v), want (%d, true)", val, ok, want)
			break
		}
	}

	// Drain from back
	expectedBack := []int{5, 4, 3, 2, 1}
	for _, want := range expectedBack {
		var val int
		var ok bool
		d, val, ok = d.PopBack()

		if !ok || val != want {
			t.Errorf("PopBack: got (%d, %v), want (%d, true)", val, ok, want)
			break
		}
	}

	if d.Len() != 0 {
		t.Errorf("final length = %d, want 0", d.Len())
	}
}

func TestImmutability(t *testing.T) {
	var d1 PFDeque[int]

	// Build d1
	d1 = d1.PushBack(1)
	d1 = d1.PushBack(2)
	d1 = d1.PushBack(3)

	// Create d2 from d1
	d2 := d1.PushBack(4)

	// d1 should still have 3 elements
	if d1.Len() != 3 {
		t.Errorf("d1 length = %d, want 3", d1.Len())
	}

	// d2 should have 4 elements
	if d2.Len() != 4 {
		t.Errorf("d2 length = %d, want 4", d2.Len())
	}

	// Modify d1 by popping
	d1, _, _ = d1.PopFront()

	// d1 should now have 2 elements
	if d1.Len() != 2 {
		t.Errorf("d1 length after pop = %d, want 2", d1.Len())
	}

	// d2 should still have 4 elements (unchanged)
	if d2.Len() != 4 {
		t.Errorf("d2 length = %d, want 4 (should be unchanged)", d2.Len())
	}
}

func TestWithStrings(t *testing.T) {
	var d PFDeque[string]

	words := []string{"apple", "banana", "cherry", "date", "elderberry"}

	// Push to front
	for _, word := range words {
		d = d.PushFront(word)
	}

	// Should be in reverse order
	expected := []string{"elderberry", "date", "cherry", "banana", "apple"}
	verifyDequeFromFront(t, d, expected)
}

func TestComplexScenario(t *testing.T) {
	var d PFDeque[int]

	// Scenario: simulate a sliding window
	// Add first 5 elements
	for i := 1; i <= 5; i++ {
		d = d.PushBack(i)
	}

	// Sliding window: remove oldest, add newest
	for i := 6; i <= 20; i++ {
		d, _, _ = d.PopFront()
		d = d.PushBack(i)

		if d.Len() != 5 {
			t.Errorf("window size = %d, want 5", d.Len())
		}
	}

	// Final window should be: 16, 17, 18, 19, 20
	expected := []int{16, 17, 18, 19, 20}
	verifyDequeFromFront(t, d, expected)
}

func TestBuildFromBothEnds(t *testing.T) {
	var d PFDeque[int]

	// Build by alternating front and back
	// Back: 1, 3, 5, 7, 9
	// Front: 0, -2, -4, -6, -8
	for i := 0; i < 5; i++ {
		d = d.PushBack(i*2 + 1) // 1, 3, 5, 7, 9
		d = d.PushFront(-i * 2) // 0, -2, -4, -6, -8
	}

	// Should be: -8, -6, -4, -2, 0, 1, 3, 5, 7, 9
	if d.Len() != 10 {
		t.Errorf("length = %d, want 10", d.Len())
	}

	expected := []int{-8, -6, -4, -2, 0, 1, 3, 5, 7, 9}
	verifyDequeFromFront(t, d, expected)
}

func TestQueueBehavior(t *testing.T) {
	// Use deque as a queue (FIFO): PushBack, PopFront
	var d PFDeque[int]

	// Enqueue 1 to 5
	for i := 1; i <= 5; i++ {
		d = d.PushBack(i)
	}

	// Dequeue should give 1, 2, 3, 4, 5
	expected := []int{1, 2, 3, 4, 5}
	for i, want := range expected {
		var val int
		var ok bool
		d, val, ok = d.PopFront()

		if !ok || val != want {
			t.Errorf("dequeue #%d: got (%d, %v), want (%d, true)", i, val, ok, want)
		}
	}
}

func TestStackBehavior(t *testing.T) {
	// Use deque as a stack (LIFO): PushFront, PopFront
	var d PFDeque[int]

	// Push 1 to 5
	for i := 1; i <= 5; i++ {
		d = d.PushFront(i)
	}

	// Pop should give 5, 4, 3, 2, 1
	expected := []int{5, 4, 3, 2, 1}
	for i, want := range expected {
		var val int
		var ok bool
		d, val, ok = d.PopFront()

		if !ok || val != want {
			t.Errorf("pop #%d: got (%d, %v), want (%d, true)", i, val, ok, want)
		}
	}
}

func TestEmptyAfterOperations(t *testing.T) {
	var d PFDeque[int]

	// Build and drain multiple times
	for cycle := 0; cycle < 3; cycle++ {
		// Build
		for i := 1; i <= 10; i++ {
			d = d.PushBack(i)
		}

		// Drain
		for d.Len() > 0 {
			d, _, _ = d.PopFront()
		}

		// Verify empty
		if d.Len() != 0 {
			t.Errorf("cycle %d: length = %d, want 0", cycle, d.Len())
		}

		_, ok := d.Front()
		if ok {
			t.Errorf("cycle %d: Front() on empty should return ok=false", cycle)
		}

		_, ok = d.Back()
		if ok {
			t.Errorf("cycle %d: Back() on empty should return ok=false", cycle)
		}
	}
}

func TestPersistentBehavior(t *testing.T) {
	var d1 PFDeque[int]

	// Build d1 using both ends to avoid single-segment dual-horizon issues
	d1 = d1.PushFront(2)
	d1 = d1.PushFront(1)
	d1 = d1.PushBack(3)

	// d1 is now: 1, 2, 3 (1,2 on left, 3 on right)

	// Create multiple versions
	d2 := d1.PushBack(4)      // 1, 2, 3, 4
	d3 := d1.PushFront(0)     // 0, 1, 2, 3
	d4, _, _ := d1.PopFront() // 2, 3

	// Verify d1 unchanged
	if d1.Len() != 3 {
		t.Errorf("d1 length = %d, want 3", d1.Len())
	}

	// Verify d2: 1, 2, 3, 4
	expected2 := []int{1, 2, 3, 4}
	verifyDequeFromFront(t, d2, expected2)

	// Verify d3: 0, 1, 2, 3
	expected3 := []int{0, 1, 2, 3}
	verifyDequeFromFront(t, d3, expected3)

	// Verify d4: 2, 3
	expected4 := []int{2, 3}
	verifyDequeFromFront(t, d4, expected4)

	// Verify d1 still: 1, 2, 3
	expected1 := []int{1, 2, 3}
	verifyDequeFromFront(t, d1, expected1)
}

func TestLengthTracking(t *testing.T) {
	var d PFDeque[int]

	lengths := []int{}

	// Track length changes
	lengths = append(lengths, d.Len()) // 0

	d = d.PushFront(1)
	lengths = append(lengths, d.Len()) // 1

	d = d.PushBack(2)
	lengths = append(lengths, d.Len()) // 2

	d = d.PushFront(0)
	lengths = append(lengths, d.Len()) // 3

	d, _, _ = d.PopFront()
	lengths = append(lengths, d.Len()) // 2

	d, _, _ = d.PopBack()
	lengths = append(lengths, d.Len()) // 1

	d, _, _ = d.PopFront()
	lengths = append(lengths, d.Len()) // 0

	expected := []int{0, 1, 2, 3, 2, 1, 0}
	for i, want := range expected {
		if lengths[i] != want {
			t.Errorf("length[%d] = %d, want %d", i, lengths[i], want)
		}
	}
}

func TestFrontBackAfterMixedOps(t *testing.T) {
	var d PFDeque[int]

	// Push 1, 2, 3 to back
	d = d.PushBack(1)
	d = d.PushBack(2)
	d = d.PushBack(3)

	// Check Front and Back
	front, _ := d.Front()
	back, _ := d.Back()
	if front != 1 || back != 3 {
		t.Errorf("before pops: Front=%d, Back=%d, want Front=1, Back=3", front, back)
	}

	// Pop from front (removes 1)
	d, _, _ = d.PopFront()

	// Check Front and Back
	front, _ = d.Front()
	back, _ = d.Back()
	if front != 2 || back != 3 {
		t.Errorf("after pop front: Front=%d, Back=%d, want Front=2, Back=3", front, back)
	}

	// Pop from back (removes 3)
	d, _, _ = d.PopBack()

	// Check Front and Back
	front, _ = d.Front()
	back, _ = d.Back()
	if front != 2 || back != 2 {
		t.Errorf("after pop back: Front=%d, Back=%d, want Front=2, Back=2", front, back)
	}

	// Pop last element
	d, _, _ = d.PopFront()

	// Check Front and Back on empty
	_, ok1 := d.Front()
	_, ok2 := d.Back()
	if ok1 || ok2 {
		t.Errorf("after draining: Front ok=%v, Back ok=%v, want both false", ok1, ok2)
	}
}

func TestPopFromBackWhenLeftEmpty(t *testing.T) {
	var d PFDeque[int]

	// Push only to right (back)
	for i := 1; i <= 5; i++ {
		d = d.PushBack(i)
	}

	// Pop from back should work (right side has elements)
	d, val, ok := d.PopBack()
	if !ok || val != 5 {
		t.Errorf("PopBack() = (%d, %v), want (5, true)", val, ok)
	}

	// Continue popping from back
	expected := []int{4, 3, 2, 1}
	verifyDequeFromBack(t, d, expected)
}

func TestPopFromFrontWhenRightEmpty(t *testing.T) {
	var d PFDeque[int]

	// Push only to left (front)
	for i := 1; i <= 5; i++ {
		d = d.PushFront(i)
	}

	// Pop from front should work (left side has elements)
	d, val, ok := d.PopFront()
	if !ok || val != 5 {
		t.Errorf("PopFront() = (%d, %v), want (5, true)", val, ok)
	}

	// Continue popping from front
	expected := []int{4, 3, 2, 1}
	verifyDequeFromFront(t, d, expected)
}

func TestFrontFromRightSide(t *testing.T) {
	var d PFDeque[int]

	// Push only to back (right side)
	d = d.PushBack(1)
	d = d.PushBack(2)
	d = d.PushBack(3)

	// Front should pull from right's back (oldest element)
	val, ok := d.Front()
	if !ok || val != 1 {
		t.Errorf("Front() = (%d, %v), want (1, true)", val, ok)
	}
}

func TestBackFromLeftSide(t *testing.T) {
	var d PFDeque[int]

	// Push only to front (left side)
	d = d.PushFront(3)
	d = d.PushFront(2)
	d = d.PushFront(1)

	// Back should pull from left's back (oldest element)
	val, ok := d.Back()
	if !ok || val != 3 {
		t.Errorf("Back() = (%d, %v), want (3, true)", val, ok)
	}
}

func TestMixedPushPopLargeScale(t *testing.T) {
	var d PFDeque[int]

	// Complex pattern: push to both ends, pop from both ends
	for i := 0; i < 100; i++ {
		d = d.PushFront(i * 2)
		d = d.PushBack(i*2 + 1)

		if i%3 == 0 && d.Len() > 0 {
			d, _, _ = d.PopFront()
		}

		if i%5 == 0 && d.Len() > 0 {
			d, _, _ = d.PopBack()
		}
	}

	// Should have some elements remaining
	if d.Len() == 0 {
		t.Errorf("deque should not be empty after mixed operations")
	}

	// Drain and verify no panics
	count := 0
	for d.Len() > 0 {
		d, _, _ = d.PopFront()
		count++
		if count > 1000 {
			t.Fatalf("infinite loop detected")
		}
	}
}

func TestPushAfterPopOptimization(t *testing.T) {
	// This test specifically exercises the SplitLast optimization in pushFront
	var d PFDeque[int]

	// Build a deque on the left side
	for i := 1; i <= 10; i++ {
		d = d.PushFront(i)
	}

	// Pop some elements from front (this leaves horizonLast < len(last.Elems)-1)
	for i := 0; i < 5; i++ {
		d, _, _ = d.PopFront()
	}

	// Now push again - this should trigger the SplitLast optimization
	d = d.PushFront(100)
	d = d.PushFront(200)

	// Verify correctness
	if d.Len() != 7 {
		t.Errorf("length = %d, want 7", d.Len())
	}

	// Pop and verify order: 200, 100, 5, 4, 3, 2, 1
	expected := []int{200, 100, 5, 4, 3, 2, 1}
	for i, want := range expected {
		var val int
		var ok bool
		d, val, ok = d.PopFront()

		if !ok {
			t.Fatalf("PopFront #%d failed", i)
		}

		if val != want {
			t.Errorf("PopFront #%d: got %d, want %d", i, val, want)
		}
	}

	if d.Len() != 0 {
		t.Errorf("final length = %d, want 0", d.Len())
	}
}

func TestPushBackAfterPopOptimization(t *testing.T) {
	// Test the SplitLast optimization for right side
	var d PFDeque[int]

	// Build a deque on the right side
	for i := 1; i <= 10; i++ {
		d = d.PushBack(i)
	}

	// Pop some elements from front (which pops from right's back)
	for i := 0; i < 5; i++ {
		d, _, _ = d.PopFront()
	}

	// Now push to back again - this should trigger the SplitLast optimization
	d = d.PushBack(100)
	d = d.PushBack(200)

	// Verify correctness
	if d.Len() != 7 {
		t.Errorf("length = %d, want 7", d.Len())
	}

	// Pop from front and verify order: 6, 7, 8, 9, 10, 100, 200
	expected := []int{6, 7, 8, 9, 10, 100, 200}
	for i, want := range expected {
		var val int
		var ok bool
		d, val, ok = d.PopFront()

		if !ok {
			t.Fatalf("PopFront #%d failed", i)
		}

		if val != want {
			t.Errorf("PopFront #%d: got %d, want %d", i, val, want)
		}
	}

	if d.Len() != 0 {
		t.Errorf("final length = %d, want 0", d.Len())
	}
}

// Benchmarks

func BenchmarkPFDequePushFront(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				d := PFDeque[int]{}
				for j := 0; j < size; j++ {
					d = d.PushFront(j)
				}
			}
		})
	}
}

func BenchmarkPFDequePushBack(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				d := PFDeque[int]{}
				for j := 0; j < size; j++ {
					d = d.PushBack(j)
				}
			}
		})
	}
}

func BenchmarkPFDequePopFront(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		setup := PFDeque[int]{}
		for j := 0; j < size; j++ {
			setup = setup.PushFront(j)
		}

		b.Run(fmt.Sprintf("n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				d := setup
				for j := 0; j < size; j++ {
					var val int
					d, val, _ = d.PopFront()
					_ = val
				}
			}
		})
	}
}

func BenchmarkPFDequePopBack(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		setup := PFDeque[int]{}
		for j := 0; j < size; j++ {
			setup = setup.PushFront(j)
		}

		b.Run(fmt.Sprintf("n=%d", size), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				d := setup
				for j := 0; j < size; j++ {
					var val int
					d, val, _ = d.PopBack()
					_ = val
				}
			}
		})
	}
}
