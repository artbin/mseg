package pfqueue

import (
	"github.com/artbin/mlist/segment"
)

type PFQueue[E any] struct {
	last         *segment.Segment[E] // Pointer to last segment (contains newest elements)
	first        *segment.Segment[E] // Pointer to first segment (contains oldest elements)
	len          int                 // Total number of elements in the list
	horizonLast  int                 // Index of the newest element in the last segment
	horizonFirst int                 // Index of the oldest element in the first segment
}

func (m PFQueue[E]) PushFront(elem E) PFQueue[E] {
	// Optimization: call SplitLast only if previous pop operation left horizonLast < len-1.
	// This ensures Merge only sees the visible elements, not hidden ones.
	if m.len > 0 && m.horizonLast != len(m.last.Elems)-1 {
		newLast, first, listlen, horizonLast := segment.SplitLast(m.last, m.first, m.len, m.horizonLast)

		newSide := PFQueue[E]{
			last:         newLast,
			first:        first,
			len:          listlen,
			horizonLast:  horizonLast,
			horizonFirst: m.horizonFirst,
		}
		return newSide.merge(elem)
	}

	return m.merge(elem)
}

func (m PFQueue[E]) merge(elem E) PFQueue[E] {
	last, first, listlen := segment.Merge(elem, m.last, m.first, m.len)
	return PFQueue[E]{last: last, first: first, len: listlen, horizonLast: len(last.Elems) - 1, horizonFirst: m.horizonFirst}
}

func (m PFQueue[E]) PopBack() (PFQueue[E], E, bool) {
	last, first, listlen, horizonFirst, value, ok := segment.PopBack(m.last, m.first, m.len, m.horizonFirst)

	// If we have a single segment (last == first), horizonLast needs to be recalculated
	// because PopBack might have compacted the segment, making the old horizonLast invalid.
	// After compaction, the newest element is at the last index of the compacted segment.
	horizonLast := m.horizonLast
	if last != nil && last == first && listlen > 0 {
		horizonLast = len(last.Elems) - 1
	}

	return PFQueue[E]{last: last, first: first, len: listlen, horizonLast: horizonLast, horizonFirst: horizonFirst}, value, ok
}

func (m PFQueue[E]) Front() (E, bool) {
	if m.len == 0 {
		var zero E
		return zero, false
	}

	return m.last.Elems[m.horizonLast], true
}

func (m PFQueue[E]) Back() (E, bool) {
	if m.len == 0 {
		var zero E
		return zero, false
	}

	return m.first.Elems[m.horizonFirst], true
}

func (m PFQueue[E]) Len() int {
	return m.len
}
