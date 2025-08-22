package go_manual_memory

import (
	"unsafe"

	ll "github.com/gabe-lee/go_list_like"
)

type List[T any] struct {
	ptr   *T
	len   uint32
	cap   uint32
	alloc Allocator
}

// Copies the data from provided Golang slice into a new `List[T]`
// using the provided `Allocator`
func CreateListCopyFrom[T any](data []T, alloc Allocator) List[T] {
	slice := CreateSliceCopyFrom(data, alloc)
	return slice.ToList(alloc)
}

// Cretes a new `List[T]` with specified length, using provided `Allocator`
func CreateList[T any](listLen int, alloc Allocator) List[T] {
	slice := CreateSlice[T](listLen, alloc)
	return slice.ToList(alloc)
}

// Convert this `List[T]` into a `Slice[T]`, handing off ownership
// of the data to the new slice
func (l *List[T]) ToSlice() Slice[T] {
	s := l.AsSlice()
	l.ptr = nil
	l.len = 0
	l.cap = 0
	l.alloc = nil
	return s
}

// Convert this `List[T]` into a `Slice[T]`, WITHOUT handing off ownership
// of the data to the new slice
//
// If the new slice is freed, this list becomes invalid
func (l *List[T]) AsSlice() Slice[T] {
	return Slice[T]{
		ptr: l.ptr,
		len: l.len,
		cap: l.cap,
	}
}

// Copies the data from this `List[T]` into a new `List[T]`
// using the cached `Allocator`
func (l *List[T]) Clone() List[T] {
	slice := l.AsSlice()
	newSlice := slice.Clone(l.alloc)
	return newSlice.ToList(l.alloc)
}

// Destroy this `List[T]`, returning the memory to the cached `Allocator`
//
// The caller MUST ensure the cached `Allocator` is the exact one originally used to
// allocate the ORIGINAL slice/list
func (l *List[T]) Destroy() {
	alloc := l.alloc
	s := l.ToSlice()
	s.Destroy(alloc)
}

// Return the length of the list
//
// Anologous to `len(slice)`
func (l *List[T]) Len() int {
	return int(l.len)
}

// Return the length of the list
//
// Anologous to `len(slice)`
func (l *List[T]) Cap() int {
	return int(l.cap)
}

// Return a pointer to the value at index `idx` in this slice
func (l *List[T]) GetPtr(idx int) *T {
	return &unsafe.Slice(l.ptr, l.len)[idx]
}

// Grow or shrink the list length, resizing/reallocating if neccessary
func (l *List[T]) OffsetLen(delta int) {
	if delta < 0 {
		l.len -= uint32(-delta)
		return
	}
	space := l.cap - l.len
	if delta <= int(space) {
		l.len += uint32(delta)
		return
	}
	newSlice := ResizeCanMove(l.alloc, l.GoSlice(), l.Len()+delta)
	newSSLice := sliceFromSlice(newSlice)
	*l = newSSLice.ToList(l.alloc)
}

var _ ll.ListLike[byte] = (*List[byte])(nil)

// Return a sub-slice of the original list's data that cannot be freed
//
// The original `List[T]` retains ownership of the data, and any changes to values in the sub-slice
// will be reflected in the original list and any other overlapping sub-slices or Golang
// slices (`[]T`)
//
// Analogous to `slice[start:end]`
func (l List[T]) SubSlice(start, end int) SubSlice[T] {
	slice := l.AsSlice()
	return slice.SubSlice(start, end)
}

// Return a sub-slice of the original list's data that cannot be freed
//
// The original `List[T]` retains ownership of the data, and any changes to values in the sub-slice
// will be reflected in the original list and any other overlapping sub-slices or Golang
// slices (`[]T`)
//
// Analogous to `slice[:]`
func (l List[T]) WholeSlice() SubSlice[T] {
	slice := l.AsSlice()
	return slice.WholeSlice()
}

// Return a sub-slice of the original list's data that cannot be freed
//
// The original `List[T]` retains ownership of the data, and any changes to values in the sub-slice
// will be reflected in the original list and any other overlapping sub-slices or Golang
// slices (`[]T`)
//
// Analogous to `slice[start:]`
func (l List[T]) EndSlice(start int) SubSlice[T] {
	slice := l.AsSlice()
	return slice.EndSlice(start)
}

// Return a sub-slice of the original list's data that cannot be freed
//
// The original `List[T]` retains ownership of the data, and any changes to values in the sub-slice
// will be reflected in the original list and any other overlapping sub-slices or Golang
// slices (`[]T`)
//
// Analogous to `slice[:end]`
func (l List[T]) StartSlice(end int) SubSlice[T] {
	slice := l.AsSlice()
	return slice.StartSlice(end)
}

// Return the `List[T]` transformed into a standard Golang `[]T` slice
//
// The `List[T]` retains ownership of the data, and any changes to values in the slice
// will be reflected in the original list and any other overlapping sub-slices or Golang
// slices (`[]T`)
func (l List[T]) GoSlice() []T {
	slice := l.AsSlice()
	return slice.GoSlice()
}
