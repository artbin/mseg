MArray vs classic dynamic array (e.g., Go slice)

MArray is part of the same segment-based family as MList/PFList/PFQueue/PFDeque: it reuses the power-of-two segmented representation, but exposes a mutable, array-like API (`Get`/`Set`) on top.

### Comparison with classic dynamic array

- Purpose and mutability: dynamic arrays are mutable contiguous buffers optimized for append-at-back and O(1) random access; MArray is a persistent, segment-based list optimized for front operations and version sharing.
- Front operations: dynamic array Push/Pop are O(n) due to shifting; MArray Push/Pop are amortized O(1) via head segments and periodic merging.
- Back operations: dynamic array Append/PopBack are amortized O(1); MArray provides O(1) First access but is front-oriented (no append-at-back API in the core structure).
- Random access: dynamic array O(1); MArray O(1) when the index lies in the tail segment and O(log n) in the worst case overall due to at most one segment per power-of-two size. When the total number of elements is a power of two, MArray coalesces into a single segment and random access is O(1).
- Memory layout & locality: dynamic arrays store elements contiguously (excellent cache locality). MArray stores elements in a few segments and periodically merges a head prefix, reducing pointer chasing relative to naive chunked lists but still not as cache-perfect as a single contiguous array.
- Resizing and copying costs: dynamic arrays occasionally reallocate and copy O(n) elements on growth; front insertions shift O(n). MArray tends to copy only a small merged head prefix after front updates; unmerged suffix segments are shared persistently.
- Persistence and snapshots: dynamic arrays require full copying for immutable snapshots; MArray provides structural sharing—creating new versions is O(1) plus the cost of modified/merged segments.
- Concurrency: dynamic arrays need synchronization or copy-on-write for readers; MArray’s immutable versions are naturally safe for concurrent readers.
- Middle insert/delete: both are O(n) work in the general case; neither structure targets efficient arbitrary middle edits.

### Advantages of MArray over classic dynamic array

- O(1) amortized front operations (push/pop), avoiding O(n) element shifts on front-heavy workloads.
- Persistence with structural sharing: cheap snapshots and branches without duplicating the entire sequence.
- Reduced copying under front-heavy mutation: merges copy only a compact head prefix rather than the whole buffer.
- Stable O(1) Last/First access; First is constant-time via an explicit tail pointer.
- Random access that is often fast in practice (O(1) for tail, O(1) when length is a power of two) and O(log n) worst-case otherwise.
- Naturally read-safe for concurrent consumers thanks to immutable prior versions.

### When a classic dynamic array is preferable

- Workloads dominated by random access and in-place updates across the entire range.
- Back-append-heavy streams where persistence is not required and maximum cache locality is desired.
- Tight numeric loops where contiguous memory and SIMD-friendly access patterns dominate performance.

### Use cases

- **Front-heavy buffer with occasional indexed access**: use `Push`/`Pop` as the primary update path, and `Get` for reads; the segmented layout avoids O(n) front-shifts typical of a slice.
- **Mutable “array-like” view over MList**: use `Set(i, v)` when you want in-place overwrites (e.g. patching recent entries) and you do not need persistence for element updates.
- **Read-both-ends without a linked list**: O(1) `Last` (newest) and O(1) `First` (oldest) while avoiding `container/list`’s per-node allocation overhead; updates remain front-oriented.
