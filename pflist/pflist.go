package pflist

import (
	"github.com/artbin/mseg/segment"
)

// PFList is a variant of MList that is purely functional (no in-place mutation) and fully persistent.
type PFList[E any] struct {
	last    *segment.Segment[E] // Pointer to last segment (contains newest elements)
	first   *segment.Segment[E] // Pointer to first segment (contains oldest elements)
	len     int                 // Total number of elements in the list
	horizon int                 // Index of the newest element in the last segment
}

// Push returns a new list with elem inserted at the logical front
// (newest position). The operation is purely functional and does not mutate PFList.
//
// Push avoids appending into the current last segment by first splitting the last
// segment into power-of-two sized segments (SplitLast) and creating a
// new last segment, guaranteeing no array writes on shared segments.
func (i PFList[E]) Push(elem E) PFList[E] {
	// Optimization: call SplitLast only if previous operation is Pop.
	if i.len > 0 && i.horizon != len(i.last.Elems)-1 {
		newList, first, listlen, horizon := segment.SplitLast(i.last, i.first, i.len, i.horizon)

		return PFList[E]{last: newList, first: first, len: listlen, horizon: horizon}.merge(elem)
	}

	return i.merge(elem)
}

func (i PFList[E]) merge(elem E) PFList[E] {
	last, first, listlen := segment.Merge(elem, i.last, i.first, i.len)
	return PFList[E]{last: last, first: first, len: listlen, horizon: len(last.Elems) - 1}
}

// Pop removes and returns the logical front element (newest).
// Returns the updated list, the element, and ok=false if empty.
// Amortized O(1); may compact the last segment at half capacity.
func (m PFList[E]) Pop() (PFList[E], E, bool) {
	last, first, listlen, horizon, value, ok := segment.PopFront(m.last, m.first, m.len, m.horizon)
	return PFList[E]{last: last, first: first, len: listlen, horizon: horizon}, value, ok
}

// Last returns the logical front element (newest). Returns false if empty. O(1).
func (i PFList[E]) Last() (E, bool) {
	if i.len == 0 {
		var zero E
		return zero, false
	}

	// Newest element is at the end of the last segment.
	return i.last.Elems[i.horizon], true
}

func (i PFList[E]) Len() int {
	return i.len
}

// Get returns the element at zero-based index counted from the logical back
// (oldest at index 0). Panics if index is out of range [0, Len()). Uses Find.
func (i PFList[E]) Get(index int) E {
	// Find() panics on out-of-range
	elems, offset := i.Find(index)
	return elems[offset]
}

// Find returns the backing slice and in-slice offset for the element at index
// counted from the logical back (oldest = 0). Panics if index is out of range.
// Complexity: O(1) when the element is in the first segment or when the list
// has a single segment; otherwise O(#segments) = O(log n).
func (i PFList[E]) Find(index int) ([]E, int) {
	elems, offset := segment.Find(index, i.last, i.first, i.len, i.horizon)
	return elems, offset
}

// pfListExported is an exported representation of PFList's internal structure
// for debugging, testing, and visualization purposes.
type pfListExported[E any] struct {
	Last    *segment.Segment[E]
	First   *segment.Segment[E]
	Len     int
	Horizon int
}

// PFListExport exports the list to a struct for debugging/testing/visualization.
func PFListExport[E any](i PFList[E]) pfListExported[E] {
	return pfListExported[E]{
		Last:    i.last,
		First:   i.first,
		Len:     i.len,
		Horizon: i.horizon,
	}
}
