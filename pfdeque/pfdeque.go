package pfdeque

import (
	"github.com/artbin/mlist/segment"
)

type PFDeque[E any] struct {
	left  pfDequeSide[E]
	right pfDequeSide[E]
}

func (m PFDeque[E]) PushFront(elem E) PFDeque[E] {
	return PFDeque[E]{left: m.left.pushFront(elem), right: m.right}
}

func (m PFDeque[E]) PushBack(elem E) PFDeque[E] {
	return PFDeque[E]{left: m.left, right: m.right.pushFront(elem)}
}

func (m PFDeque[E]) PopFront() (PFDeque[E], E, bool) {
	left, elem, ok := m.left.popFront()
	if ok {
		return PFDeque[E]{left: left, right: m.right}, elem, ok
	} else {
		right, elem, ok := m.right.popBack()
		if ok {
			return PFDeque[E]{left: m.left, right: right}, elem, ok
		} else {
			var zero E
			return PFDeque[E]{}, zero, false
		}
	}
}

func (m PFDeque[E]) PopBack() (PFDeque[E], E, bool) {
	right, elem, ok := m.right.popFront()
	if ok {
		return PFDeque[E]{left: m.left, right: right}, elem, ok
	} else {
		left, elem, ok := m.left.popBack()
		if ok {
			return PFDeque[E]{left: left, right: m.right}, elem, ok
		} else {
			var zero E
			return PFDeque[E]{}, zero, false
		}
	}
}

func (m PFDeque[E]) Front() (E, bool) {
	elem, ok := m.left.front()
	if ok {
		return elem, ok
	} else {
		elem, ok := m.right.back()
		if ok {
			return elem, ok
		} else {
			var zero E
			return zero, false
		}
	}
}

func (m PFDeque[E]) Back() (E, bool) {
	elem, ok := m.right.front()
	if ok {
		return elem, ok
	} else {
		elem, ok := m.left.back()
		if ok {
			return elem, ok
		} else {
			var zero E
			return zero, false
		}
	}
}

func (m PFDeque[E]) Len() int {
	return m.left.len + m.right.len
}

type pfDequeSide[E any] struct {
	last         *segment.Segment[E] // Pointer to last segment (contains newest elements)
	first        *segment.Segment[E] // Pointer to first segment (contains oldest elements)
	len          int                 // Total number of elements in the list
	horizonLast  int                 // Index of the newest element in the last segment
	horizonFirst int                 // Index of the oldest element in the first segment
}

func (m pfDequeSide[E]) pushFront(elem E) pfDequeSide[E] {
	// Optimization: call SplitLast only if previous pop operation left horizonLast < len-1.
	// This ensures Merge only sees the visible elements, not hidden ones.
	if m.len > 0 && m.horizonLast != len(m.last.Elems)-1 {
		newLast, first, listlen, horizonLast := segment.SplitLast(m.last, m.first, m.len, m.horizonLast)

		newSide := pfDequeSide[E]{
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

func (m pfDequeSide[E]) merge(elem E) pfDequeSide[E] {
	last, first, listlen := segment.Merge(elem, m.last, m.first, m.len)
	return pfDequeSide[E]{last: last, first: first, len: listlen, horizonLast: len(last.Elems) - 1, horizonFirst: m.horizonFirst}
}

func (m pfDequeSide[E]) popFront() (pfDequeSide[E], E, bool) {
	last, first, listlen, horizonLast, value, ok := segment.PopFront(m.last, m.first, m.len, m.horizonLast)
	return pfDequeSide[E]{last: last, first: first, len: listlen, horizonLast: horizonLast, horizonFirst: m.horizonFirst}, value, ok
}

func (m pfDequeSide[E]) popBack() (pfDequeSide[E], E, bool) {
	last, first, listlen, horizonFirst, value, ok := segment.PopBack(m.last, m.first, m.len, m.horizonFirst)

	// If we have a single segment (last == first), horizonLast needs to be recalculated
	// because PopBack might have compacted the segment, making the old horizonLast invalid.
	// After compaction, the newest element is at the last index of the compacted segment.
	horizonLast := m.horizonLast
	if last != nil && last == first && listlen > 0 {
		horizonLast = len(last.Elems) - 1
	}

	return pfDequeSide[E]{last: last, first: first, len: listlen, horizonLast: horizonLast, horizonFirst: horizonFirst}, value, ok
}

func (m pfDequeSide[E]) front() (E, bool) {
	if m.len == 0 {
		var zero E
		return zero, false
	}

	return m.last.Elems[m.horizonLast], true
}

func (m pfDequeSide[E]) back() (E, bool) {
	if m.len == 0 {
		var zero E
		return zero, false
	}

	return m.first.Elems[m.horizonFirst], true
}
