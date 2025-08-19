package go_manual_memory

// import (
// 	"math/bits"
// 	"slices"
// 	"unsafe"

// 	esort "github.com/gabe-lee/go_effect_sort"
// )

// const MIN_CHUNK_SIZE = PAGE_SIZE

// type SpanAllocator struct {
// 	chunks                     [][]byte
// 	free_spans                 []span
// 	free_spans_sorted_by_len   []uint32
// 	free_spans_sorted_by_start []uint32
// }

// func spanEqualLen(slice []span, idx int, val span) bool {
// 	return slice[idx].len == val.len
// }
// func spanGreaterLen(slice []span, idx int, val span) bool {
// 	return slice[idx].len > val.len
// }

// func (s *SpanAllocator) addNewSpan(ss span) {
// }

// func (s *SpanAllocator) findFreeSpan(span_len, align uintptr) (sp span, found bool) {
// 	const _SPAN_BEFORE uint8 = 0b01
// 	const _SPAN_AFTER uint8 = 0b10
// 	const _SPAN_BEFORE_AND_AFTER uint8 = 0b11
// 	const _SPAN_EXACT uint8 = 0b00
// 	span_after := span{len: uint32(span_len)}
// 	free_span_idx := esort.BinaryInsertIndex(s.free_spans, span_after, spanEqualLen, spanGreaterLen)
// 	if free_span_idx >= len(s.free_spans) {
// 		return span_after, false
// 	}
// 	var has_span_before uint8 = 0b00
// 	var has_span_after uint8 = 0b00
// 	var span_before span
// 	free_span := s.free_spans[free_span_idx]
// 	free_span_addr, free_span_align := s.spanAddrAlign(free_span)
// 	free_span_original_len := free_span.len
// 	for free_span_align < align {
// 		next_aligned := (free_span_addr + align - 1) & ^(align - 1)
// 		if next_aligned+span_len > free_span_addr+uintptr(free_span.len) {
// 			free_span_idx += 1
// 			if free_span_idx >= len(s.free_spans) {
// 				return span_after, false
// 			}
// 			free_span := s.free_spans[free_span_idx]
// 			free_span_addr, free_span_align = s.spanAddrAlign(free_span)
// 			free_span_original_len = free_span.len
// 		} else {
// 			delta := uint32(next_aligned - free_span_addr)
// 			span_before = span{
// 				chunk: free_span.chunk,
// 				start: free_span.start,
// 				len:   delta,
// 			}
// 			has_span_before = 0b01
// 			free_span.start += delta
// 			free_span.len -= delta
// 		}
// 	}
// 	leftover := free_span.len - span_after.len
// 	free_span.len = span_after.len
// 	span_after.len = leftover
// 	if span_after.len > 0 {
// 		span_after.chunk = free_span.chunk
// 		span_after.start = free_span.start + free_span.len
// 		has_span_after = 0b10
// 	}
// 	var alter_mode uint8 = has_span_before | has_span_after
// 	switch alter_mode {
// 	case _SPAN_EXACT:
// 		slices.Delete(s.free_spans, free_span_idx, free_span_idx+1)
// 		esort.BinarySearch()
// 	}
// }

// // RawAlloc implements Allocator.
// func (s *SpanAllocator) RawAlloc(len uintptr, align uintptr) (ptr unsafe.Pointer, alloc_len uintptr) {
// 	panic("unimplemented")
// }

// // RawFree implements Allocator.
// func (s *SpanAllocator) RawFree(ptr unsafe.Pointer, len uintptr) {
// 	panic("unimplemented")
// }

// // RawResize implements Allocator.
// func (s *SpanAllocator) RawResize(ptr unsafe.Pointer, old_len uintptr, new_len uintptr, align uintptr) (new_ptr unsafe.Pointer, alloc_len uintptr) {
// 	panic("unimplemented")
// }

// var _ Allocator = (*SpanAllocator)(nil)

// type span struct {
// 	chunk uint32
// 	start uint32
// 	len   uint32
// }

// func (s *SpanAllocator) spanAddrAlign(ss span) (addr, align uintptr) {
// 	cdata := s.chunks[ss.chunk]
// 	caddr := uintptr(unsafe.Pointer(unsafe.SliceData(cdata)))
// 	saddr := caddr + uintptr(ss.start)
// 	return saddr, uintptr(bits.TrailingZeros(uint(saddr)))
// }
