// Package segment provides primitives for a segmented list.
//
// Model
// The structure is a singly linked chain of segments ordered from newest to
// oldest. Pointer Next moves toward older segments. Each segment holds a
// contiguous run of elements in Elems laid out oldest→newest.
//
// Cursors
// Operations are expressed in terms of four cursors carried by the caller:
//   - last: pointer to the newest segment
//   - first: pointer to the oldest segment
//   - listlen: total logical element count
//   - horizon: index of the newest element within last.Elems (visible length of
//     last is horizon+1) or the oldest element within first.Elems (visible length of
//     first is horizon+1)
//
// Horizon
// In PopFront, Find, SplitLast functions horizon indexes the newest element of the last segment; only the
// visible range [0..horizon] within last.Elems is visible. All non-last segments are
// considered fully visible (indices [0..len-1]). Callers should treat horizon conceptually as “newest index in
// last” and pass through the returned value between operations.
// In PopBack function horizon indexes the oldest element of the first segment; only the
// visible range [horizon..len-1] within first.Elems is visible. All non-first segments are
// considered fully visible (indices [0..len-1]). Callers should treat horizon conceptually as “oldest index in
// first” and pass through the returned value between operations.
//
// Complexity summary (Big-O)
// Symbols: n = total logical elements; s = number of segments (≤ floor(log2(n))+1);
//
//	h = visible length in the last segment (horizon+1); r = remaining elements
//	in the first segment after a back-pop.
//
// | Operation    | Worst-case                          | Amortized/Typical                            | Notes                                                            |
// |--------------|-------------------------------------|----------------------------------------------|------------------------------------------------------------------|
// | Merge (Push) | O(n) when fully coalescing; O(s)    | O(1)                                         | Binary-counter-like carries; partial coalescing                  |
// | PopFront     | O(h) on compaction                  | O(1)                                         | Compact when horizon == cap(last)/2; keep visible range [0..h-1] |
// | PopBack      | O(s + r) on compaction/removal      | O(1)                                         | Scan predecessor when first size==1; compact at half             |
// | SplitLast    | O(popcount(h)) = O(log n)           | O(log n)                                     | Reslice-only; no element copies                                  |
// | Find(i)      | O(s)                                | O(1) if in first or when s==1; else O(log n) | First-segment fast path                                          |
package segment

import "math/bits"

// Segment is a node in a segmented list.
// Next points toward older segments (away from the newest).
// Elems holds a contiguous run of elements in logical order (oldest to newest).
//
// Let n be the logical element count and s the number of segments. Invariants:
//   - s is O(log n) and specifically s ≤ floor(log2(n)) + 1.
//   - When n is a power of two, s = 1 (the list fully coalesces).
//
// These bounds underlie the per-operation complexities documented below.
type Segment[E any] struct {
	Next  *Segment[E]
	Elems []E
}

// Merge appends elem as the newest element and partially coalesces segments
// near the tail. It returns the updated last segment, first segment (which may
// change if everything coalesces), and the new list length.
// The algorithm coalesces a run of segments starting at last while either the
// total list length or the cumulative run size remains a power of two. This
// keeps the structure compact and yields amortized O(1) appends without
// destroying persistence. After merge, horizon indexes the newest element in
// the returned last segment (len(last.Elems)-1).
//
// Complexity:
//   - Worst case: O(n) time when newLen is a power of two (the entire list
//     coalesces into one segment), plus O(s) segment headers visited.
//   - Typical case: O(k) where k is the number of elements copied into the new
//     last segment for the current coalescing run; s visited is O(log n).
//   - Amortized: O(1) per push (proof sketch: merges mirror carries in a binary
//     counter; each element moves O(1) times across successive pushes, so total
//     copying over N pushes is O(N)).
//   - Space: O(k) for the new backing slice of the coalesced tail; older
//     segments remain for persistence and are reclaimed by GC when unreachable.
func Merge[E any](elem E, last *Segment[E], first *Segment[E], listlen int) (*Segment[E], *Segment[E], int) {
	newLen := listlen + 1

	// Empty list fast path: allocate a single-element segment.
	if last == nil {
		last := &Segment[E]{Next: nil, Elems: make([]E, 1)}
		last.Elems[0] = elem
		return last, last, newLen
	}

	newLenIsPowerOfTwo := (newLen & (newLen - 1)) == 0

	// Start with the new element already accounted for.
	elemsLen := uint(1)
	numSegments := 0
	// If we include zero existing segments, the remainder is the current last.
	nextAfterCoalesced := last

	// Identify the coalescing run (starting after the new element) to coalesce.
	for segment := last; segment != nil; segment = segment.Next {
		trial := elemsLen + uint(len(segment.Elems))
		elemsLenIsPowerOfTwo := (trial & (trial - 1)) == 0
		// Stop when both the total list length and the cumulative coalescing run size
		// cease to be powers of two. This limits the amount of copying while
		// ensuring periodic compaction.
		if !newLenIsPowerOfTwo && !elemsLenIsPowerOfTwo {
			break
		}
		// Include this segment in the coalescing run
		elemsLen = trial
		numSegments++
		nextAfterCoalesced = segment.Next
	}

	// Allocate space for the coalescing run elements.
	elems := make([]E, elemsLen)

	if numSegments > 0 {
		// Copy from oldest to newest without an intermediate slice by
		// computing target offsets from the total coalescing run length.
		totalMerged := int(elemsLen) - 1 // reserve last slot for the new element
		copied := 0
		remaining := 0
		for segment := last; remaining < numSegments && segment != nil; segment = segment.Next {
			segmentLen := len(segment.Elems)
			start := totalMerged - copied - segmentLen
			copy(elems[start:], segment.Elems)
			copied += segmentLen
			remaining++
		}

		// Place the new element (newest) at the end.
		elems[totalMerged] = elem
	} else {
		// Only the new element is in the coalescing run.
		elems[0] = elem
	}

	// Link new last segment to the first unmerged segment (if any), preserving the first segment
	// of the list. The coalesced elements become the new last segment payload.
	newLast := &Segment[E]{Next: nextAfterCoalesced, Elems: elems}

	// If we merged all segments, the last segment is also the first segment.
	if newLast.Next == nil {
		first = newLast
	}

	return newLast, first, newLen
}

// PopFront removes and returns the logical front (newest) element.
// Returns: last, first, newLen, newHorizon, value, ok.
// ok is false when the list is empty.
//
// Complexity:
//   - O(1) for the common path (move horizon or advance to next segment).
//   - O(h) when compaction triggers at half capacity, where h ≈ visible size in
//     the last segment after the pop (one-time copy). Amortized O(1) overall
//     because each halving copies elements at most once per capacity drop.
func PopFront[E any](last *Segment[E], first *Segment[E], listlen int, horizon int) (*Segment[E], *Segment[E], int, int, E, bool) {
	if listlen == 0 {
		var v E
		return last, first, listlen, horizon, v, false
	}

	// The logical front (newest) resides at last.elems[horizon].
	elem := last.Elems[horizon]

	// nextHorizon will become the horizon of the returned list.
	nextHorizon := 0

	if horizon == 0 {
		// Case 1: we popped the last visible element from the current last segment.
		// Advance last to the next segment.
		last = last.Next

		if last == nil {
			// No more segments: list becomes empty.
			first = nil
		} else {
			// New horizon points to the newest element in the next last segment.
			nextHorizon = len(last.Elems) - 1
		}
	} else if horizon == (len(last.Elems) >> 1) {
		// Case 2: compaction threshold hit (horizon == capacity/2 for the last segment).
		// Compact the last segment by retaining only the visible range that remains after
		// the pop: [0..horizon-1]. Allocate a right-sized slice of length horizon,
		// copy that visible range, and rebuild last to reduce memory and improve locality.
		newElems := make([]E, horizon)
		copy(newElems, last.Elems)
		newLast := &Segment[E]{Next: last.Next, Elems: newElems}

		if last == first {
			// If last segment was also first segment, update first segment to the new compacted segment.
			first = newLast
		}

		last = newLast

		// After compaction, the newest element sits at the end of the new last segment slice.
		nextHorizon = len(last.Elems) - 1
	} else {
		// Case 3: normal fast-path. Just move the horizon left by one.
		nextHorizon = horizon - 1
	}

	return last, first, listlen - 1, nextHorizon, elem, true
}

// PopBack removes and returns the logical back (oldest) element.
// Returns: last, first, newLen, newHorizon, value, ok.
// ok is false when the list is empty.
//
// Complexity:
//   - O(1) when the oldest segment has more than one remaining element and no
//     compaction is due (just increment horizon).
//   - O(s) to locate the predecessor when the oldest segment size is 1,
//     where s is the segment count (s = O(log n)).
//   - O(r + s) when compaction triggers at half capacity of the oldest segment,
//     where r is the remaining element count copied into the compacted segment.
//   - Amortized behavior across a draining sequence is O(1) per pop: scans
//     happen once per segment removal and capacity shrinks geometrically.
func PopBack[E any](last *Segment[E], first *Segment[E], listlen int, horizon int) (*Segment[E], *Segment[E], int, int, E, bool) {
	if listlen == 0 {
		var zero E
		return last, first, listlen, horizon, zero, false
	}

	newListLen := listlen - 1

	value := first.Elems[horizon]

	// Move the horizon right by one; if that fully consumes the current first segment,
	// move the first segment pointer to its predecessor and reset horizon.
	nextHorizon := horizon + 1
	if len(first.Elems) == 1 {
		if last == first {
			return last, first, newListLen, 0, value, true
		} else {
			prev := last
			for prev != nil && prev.Next != nil && prev.Next != first {
				prev = prev.Next
			}

			return last, prev, newListLen, 0, value, true
		}
	} else if nextHorizon == len(first.Elems)>>1 {
		// Compact the first segment after crossing half capacity.
		remain := first.Elems[nextHorizon:]

		firstElems := make([]E, len(remain))
		copy(firstElems, remain)
		newFirst := &Segment[E]{Elems: firstElems, Next: nil}

		// If last == first, then we only have one segment
		if last == first {
			return newFirst, newFirst, newListLen, 0, value, true
		}

		// Find the predecessor of first and link it to the new compacted first
		prev := last
		for prev != nil && prev.Next != nil && prev.Next != first {
			prev = prev.Next
		}

		if prev != nil {
			prev.Next = newFirst
		}

		return last, newFirst, newListLen, 0, value, true
	}

	return last, first, newListLen, nextHorizon, value, true
}

// SplitLast splits the last segment's visible range [0..horizon] into segments
// sized by the power-of-two decomposition (1,2,4,...).
// Newest elements remain in the returned last segment; older chunks follow via
// Next pointers. Elements are not copied; slices are re-sliced views.
// No-op when the visible length <= 1.
//
// Complexity: O(t) time and space where t = popcount(horizon+1).
// Since t ≤ floor(log2(horizon+1)) + 1, this is O(log n) in n = list length.
func SplitLast[E any](last *Segment[E], first *Segment[E], listlen int, horizon int) (*Segment[E], *Segment[E], int, int) {
	if last == nil {
		return last, first, listlen, horizon
	}

	newHorizon := horizon + 1
	if newHorizon <= 1 {
		return last, first, listlen, horizon
	}

	// Preserve the remainder of the chain after the current last segment. We'll append it
	// after the split chain is constructed.
	oldNext := last.Next

	newLast := (*Segment[E])(nil)
	splitFirst := (*Segment[E])(nil)

	// Preallocate exactly the number of segments we will produce based on the
	// population count of newHorizon (power-of-two decomposition).
	numSegments := bits.OnesCount(uint(newHorizon))
	segments := make([]Segment[E], numSegments)
	segmentIndex := 0
	// Position one past the newest logical element (newHorizon).
	pos := newHorizon
	// Decompose newHorizon into powers of two by iterating set bits from least
	// significant to most significant so we produce segments sized 1,2,4,...
	// Example: horizon=6 (newHorizon=7 => 0b111) yields slices [6:7],[4:6],[0:4].
	// Example: horizon=5 (newHorizon=6 => 0b110) yields slices [4:6],[0:4].
	for bit := 0; (1 << bit) <= newHorizon; bit++ {
		if (newHorizon>>bit)&1 == 0 {
			continue
		}

		size := 1 << bit
		pos -= size

		elems := last.Elems[pos : pos+size]

		segment := &segments[segmentIndex]
		segmentIndex++
		segment.Elems = elems

		if newLast == nil {
			// First chunk corresponds to the newest suffix; it must become the
			// last segment so that Last remains at last.elems[horizon].
			newLast = segment
			splitFirst = segment
		} else {
			// Append older chunks to the end of the split chain to preserve
			// newest-to-oldest ordering in the segment chain.
			splitFirst.Next = segment
			splitFirst = segment
		}
	}

	// Link the split chain to the remainder of the list.
	splitFirst.Next = oldNext

	// Determine the resulting first pointer.
	if splitFirst.Next == nil {
		first = splitFirst
	}

	return newLast, first, listlen, len(newLast.Elems) - 1
}

// Find returns the backing slice and in-slice offset for the element at index
// counted from the logical back (oldest = 0). Panics if index is out of range.
// Fast path checks the first segment; otherwise it walks forward from last,
// accumulating lengths until the index lands in a segment.
// Complexity:
//   - O(1) if the element is in the first (oldest) segment.
//   - O(1) when the list length is a power of two (single segment after merge).
//   - O(s) in general, where s is the number of segments (O(log n)).
func Find[E any](index int, last *Segment[E], first *Segment[E], listlen int, horizon int) ([]E, int) {
	if index < 0 || index >= listlen {
		panic("index out of range")
	}

	// First segment fast path: when the index falls in the first segment. If the list
	// consists of a single segment (last == first), only elements up to
	// horizon are logically present.
	if first != nil {
		segmentLen := 0

		if last == first {
			segmentLen = horizon + 1
		} else {
			segmentLen = len(first.Elems)
		}

		if index < segmentLen {
			return first.Elems, index
		}
	}

	// Start with a negative offset such that as we add segment lengths while
	// walking forward from last segment, the target segment is where offset >= 0.
	offset := index - listlen

	for segment := last; segment != nil; segment = segment.Next {
		segmentLen := 0

		if segment == last {
			segmentLen = horizon + 1
		} else {
			segmentLen = len(segment.Elems)
		}

		// Accumulate the number of elements in each segment until we pass the
		// desired logical index measured from the back.
		offset += segmentLen

		if offset >= 0 {
			// Found the segment containing the desired index; offset is now the
			// in-segment index.
			return segment.Elems, offset
		}
	}

	// Unreachable due to bounds check.
	panic("unreachable")
}
