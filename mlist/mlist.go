// This package implements a segmented list of mergeable segments. It supports
// efficient Push/Pop with periodic segment coalescing. Elements are stored in a
// forward-linked chain of array segments.
// Recent operations primarily touch a short coalescing run near the newest end,
// keeping updates amortized O(1) while retaining O(1) Last/First and
// O(log n) random access when the chain remains decomposed, with occasional
// O(n) merges when the list length hits a power of two.
package mlist

import (
	"fmt"

	"github.com/artbin/mlist/segment"
)

// MList implements a segmented, front-oriented list by chaining
// variable-size array segments.
//
// Invariants and conventions:
//   - last points at the segment containing newest elements; first points at
//     the segment containing oldest elements (O(1) First).
//   - Within a segment, elements are laid out from oldest to newest. Across
//     segments, logical order flows from first (oldest) to last (newest). The
//     logical "front" (newest) is at last.Elems[horizon]. After a merge, horizon
//     equals len(last.Elems)-1.
//   - The visible range of the last segment is [0..horizon]. Non-last
//     segments are fully visible. This aligns with the semantics in package
//     segment.
//   - Merge periodically coalesces a run of segments into a single tail segment
//     to keep the structure compact and operations amortized.
type MList[E any] struct {
	last    *segment.Segment[E] // Pointer to last segment (contains newest elements)
	first   *segment.Segment[E] // Pointer to first segment (contains oldest elements)
	len     int                 // Total number of elements in the list
	horizon int                 // Index of the newest element in the last segment
}

func (m MList[E]) merge(elem E) MList[E] {
	last, first, listlen := segment.Merge(elem, m.last, m.first, m.len)
	return MList[E]{last: last, first: first, len: listlen, horizon: len(last.Elems) - 1}
}

// Push inserts elem at the logical front (newest position) and returns the new
// list. The operation may mutate MList.
//
// Complexity: amortized O(1). On certain lengths (power-of-two boundaries),
// push triggers coalescing analogous to carries in a binary counter.
func (m MList[E]) Push(elem E) MList[E] {
	if m.len > 0 && m.horizon != len(m.last.Elems)-1 {
		horizon := m.horizon + 1

		m.last.Elems[horizon] = elem

		return MList[E]{last: m.last, first: m.first, len: m.len + 1, horizon: horizon}
	}

	return m.merge(elem)
}

// Pop removes and returns the logical front element (newest).
// Returns the updated list, the element, and ok=false if empty.
// Amortized O(1); may compact the last segment at half capacity.
func (m MList[E]) Pop() (MList[E], E, bool) {
	last, first, listlen, horizon, value, ok := segment.PopFront(m.last, m.first, m.len, m.horizon)
	return MList[E]{last: last, first: first, len: listlen, horizon: horizon}, value, ok
}

// Last returns the logical front element (newest). Returns false if empty. O(1).
func (m MList[E]) Last() (E, bool) {
	if m.len == 0 {
		var zero E
		return zero, false
	}

	// Newest element is at the end of the last segment.
	return m.last.Elems[m.horizon], true
}

// First returns the logical back element (oldest). Returns false if empty. O(1).
func (m MList[E]) First() (E, bool) {
	if m.len == 0 {
		var zero E
		return zero, false
	}

	// Oldest element resides at the beginning of the first segment.
	return m.first.Elems[0], true
}

// Get returns the element at zero-based index counted from the logical back
// (oldest at index 0). Panics if index is out of range [0, Len()). Uses Find.
func (m MList[E]) Get(index int) E {
	// Find() panics on out-of-range
	elems, offset := m.Find(index)

	return elems[offset]
}

// Find returns the backing slice and in-slice offset for the element at index
// counted from the logical back (oldest = 0). Panics if index is out of range.
// Complexity: O(1) when the element is in the first segment or when the list
// has a single segment; otherwise O(#segments) = O(log n).
func (m MList[E]) Find(index int) ([]E, int) {
	elems, offset := segment.Find(index, m.last, m.first, m.len, m.horizon)

	return elems, offset
}

// Len returns the number of elements in the list. O(1).
func (m MList[E]) Len() int {
	return m.len
}

// MListVerifyInvariants checks internal structural and length invariants of the list.
// Intended for debugging/testing only.
// It returns an error describing the first violation found, or nil if all
// invariants hold.
func MListVerifyInvariants[E any](m MList[E]) error {
	// helper: power-of-two check (>0 and single bit set)
	isPow2 := func(n int) bool {
		return n > 0 && (n&(n-1)) == 0
	}

	// Empty list invariants
	if m.len == 0 {
		if m.last != nil || m.first != nil {
			return fmt.Errorf("empty list must have last==nil and first==nil: last=%v first=%v", m.last, m.first)
		}
		// horizon is irrelevant when empty
		return nil
	}

	// Non-empty list requires valid pointers
	if m.last == nil || m.first == nil {
		return fmt.Errorf("non-empty list must have last and first non-nil: last=%v first=%v", m.last, m.first)
	}

	// Last segment bounds
	if len(m.last.Elems) == 0 {
		return fmt.Errorf("last segment must have non-empty elems")
	}
	if m.horizon < 0 || m.horizon >= len(m.last.Elems) {
		return fmt.Errorf("horizon out of bounds: horizon=%d len(last.elems)=%d", m.horizon, len(m.last.Elems))
	}
	if !isPow2(len(m.last.Elems)) {
		return fmt.Errorf("last segment length must be power of two: got %d", len(m.last.Elems))
	}

	// Detect cycles (Floyd) and later verify first is tail
	slow := m.last
	fast := m.last
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			return fmt.Errorf("cycle detected in segment chain")
		}
	}

	// Traverse to count logical elements and Find tail
	count := 0
	tail := (*segment.Segment[E])(nil)
	for seg := m.last; seg != nil; seg = seg.Next {
		if len(seg.Elems) == 0 {
			return fmt.Errorf("segment with empty elems encountered")
		}
		if !isPow2(len(seg.Elems)) {
			return fmt.Errorf("segment length must be power of two: got %d", len(seg.Elems))
		}
		if seg == m.last {
			count += m.horizon + 1
		} else {
			count += len(seg.Elems)
		}
		if seg.Next == nil {
			tail = seg
		}
	}

	if tail != m.first {
		return fmt.Errorf("first pointer mismatch: expected tail=%p got first=%p", tail, m.first)
	}
	if m.first.Next != nil {
		return fmt.Errorf("first.next must be nil")
	}
	if count != m.len {
		return fmt.Errorf("length mismatch: fields len=%d but counted=%d", m.len, count)
	}

	return nil
}

// mList is a mutable struct for debugging/testing/visualization.
type mList[E any] struct {
	Last    *segment.Segment[E]
	First   *segment.Segment[E]
	Len     int
	Horizon int
}

// MListExport exports the list to a mutable struct for debugging/testing/visualization.
func MListExport[E any](m MList[E]) mList[E] {
	return mList[E]{
		Last:    m.last,
		First:   m.first,
		Len:     m.len,
		Horizon: m.horizon,
	}
}
