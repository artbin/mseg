package marray

import (
	"github.com/artbin/mseg/mlist"
)

// MArray is a mutable, dynamic-array-style facade over MList.
//
// It exposes a conventional API with in-place updates (via pointer receiver)
// by replacing the internal persistent value on structural edits. The Set
// method performs an in-place element overwrite in the underlying segment for
// O(1) updates after locating the segment. This intentionally sacrifices the
// persistence guarantee for element updates when using this wrapper.
//
// Complexity characteristics (amortized unless noted):
//   - Push/Pop: O(1) amortized; may trigger internal merges
//   - Last/First: O(1)
//   - Get: O(1) if the index falls inside the first segment; worst-case O(log n)
//     due to the power-of-two segment scheme
//   - Set: O(1) after locating the segment (the Find itself follows Get’s bounds)
//
// Note: Because Set mutates the shared slice within a persistent segment,
// any other references to the same MList (or structures sharing
// those segments) will observe the updated value.
type MArray[E any] struct {
	internal mlist.MList[E]
}

// Push inserts elem at the logical front (newest position).
func (m *MArray[E]) Push(elem E) {
	m.internal = m.internal.Push(elem)
}

// Pop removes and returns the logical front element.
// Returns None if the list is empty.
func (m *MArray[E]) Pop() (E, bool) {
	internal, elem, ok := m.internal.Pop()
	if !ok {
		var zero E
		return zero, false
	}

	m.internal = internal

	return elem, true
}

// Last returns the logical front element (newest) if present.
func (m *MArray[E]) Last() (E, bool) {
	return m.internal.Last()
}

// First returns the logical back element (oldest) if present.
func (m *MArray[E]) First() (E, bool) {
	return m.internal.First()
}

// Get returns the element at index measured from the logical back (oldest at 0).
// Complexity: O(1) if in first; worst-case O(log n) overall.
func (m *MArray[E]) Get(index int) E {
	return m.internal.Get(index)
}

// Set overwrites the element at index with elem.
//
// This locates the underlying segment and assigns directly (in-place mutation)
// into the segment’s slice.
func (m *MArray[E]) Set(index int, elem E) {
	elems, offset := m.internal.Find(index)

	elems[offset] = elem
}

// Len returns the number of elements.
func (m *MArray[E]) Len() int {
	return m.internal.Len()
}
