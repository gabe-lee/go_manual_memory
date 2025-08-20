package go_manual_memory

import (
	"fmt"
	"unsafe"

	ll "github.com/gabe-lee/go_list_like"
)

type Slice[T any] struct {
	ptr *T
	len uint32
	cap uint32
}

func sliceFromSlice[T any](slice []T) Slice[T] {
	return Slice[T]{
		ptr: unsafe.SliceData(slice),
		len: uint32(len(slice)),
		cap: uint32(cap(slice)),
	}
}

// Copies the data from provided Golang slice into a new `Slice[T]`
// using the provided `Allocator`
func CreateSliceCopyFrom[T any](data []T, alloc Allocator) Slice[T] {
	slice := Alloc[T](alloc, len(data))
	copy(slice, data)
	return sliceFromSlice(slice)
}

// Cretes a new `Slice[T]` with specified length, using provided `Allocator`
func CreateSlice[T any](sliceLen int, alloc Allocator) Slice[T] {
	slice := Alloc[T](alloc, sliceLen)
	return sliceFromSlice(slice)
}

// Cretes a new `Slice[T]` with specified capacity (length 0), using provided `Allocator`
func CreateEmptySlice[T any](sliceCap int, alloc Allocator) Slice[T] {
	slice := Alloc[T](alloc, sliceCap)
	return sliceFromSlice(slice[:0])
}

// Copies the data from this `Slice[T]` into a new `Slice[T]`
// using the provided `Allocator`
func (s Slice[T]) Clone(alloc Allocator) Slice[T] {
	this_goslice := s.GoSlice()
	new_slice := CreateSlice[T](int(s.cap), alloc)
	new_slice.len = s.len
	new_goslice := new_slice.GoSlice()
	copy(new_goslice[:s.len], this_goslice[:s.len])
	return new_slice
}

// Convert this `Slice[T]` into a `List[T]`, handing off ownership
// of the data to the new list
func (s *Slice[T]) ToList(alloc Allocator) List[T] {
	l := List[T]{
		ptr:   s.ptr,
		len:   s.len,
		cap:   s.cap,
		alloc: alloc,
	}
	s.ptr = nil
	s.len = 0
	s.cap = 0
	return l
}

// Destroy this `Slice[T]`, returning the memory to the provided `Allocator`
//
// The caller MUST ensure the provided `Allocator` is the exact one originally used to
// allocate this slice
func (s *Slice[T]) Destroy(alloc Allocator) {
	slice := s.GoSlice()
	Free(alloc, slice)
	s.ptr = nil
	s.cap = 0
	s.len = 0
}

// Return the length of the slice
//
// Anologous to `len(slice)`
func (s Slice[T]) Len() int {
	return int(s.len)
}

// Return the capacity of the slice
//
// Anologous to `cap(slice)`
func (s Slice[T]) Cap() int {
	return int(s.cap)
}

// Return a pointer to the value at index `idx` in this slice
func (s Slice[T]) GetPtr(idx int) *T {
	return &unsafe.Slice(s.ptr, s.len)[idx]
}

var _ ll.SliceLike[byte] = (*Slice[byte])(nil)

// Return a sub-slice of the original slice that cannot be freed
//
// The original `Slice[T]` retains ownership of the data, and any changes to values in the sub-slice
// will be reflected in the original slice and any other overlapping sub-slices or Golang
// slices (`[]T`)
//
// Analogous to `slice[start:end]`
func (s Slice[T]) SubSlice(start, end int) SubSlice[T] {
	if end > int(s.len) {
		panic(fmt.Sprintf("fatal: go_manual_memory: Slice[T].SubSlice(): end index %d is greater than len %d", end, s.len))
	}
	goslice := s.GoSlice()
	len := uint32(end - start)
	sub_ptr := unsafe.SliceData(goslice[start:end])
	return SubSlice[T]{
		ptr: sub_ptr,
		len: len,
		cap: len,
	}
}

// Return a sub-slice of the original slice that cannot be freed
//
// The original `Slice[T]` retains ownership of the data, and any changes to values in the sub-slice
// will be reflected in the original slice and any other overlapping sub-slices or Golang
// slices (`[]T`)
//
// Analogous to `slice[:]`
func (s Slice[T]) WholeSlice() SubSlice[T] {
	return s.SubSlice(0, int(s.len))
}

// Return a sub-slice of the original slice that cannot be freed
//
// The original `Slice[T]` retains ownership of the data, and any changes to values in the sub-slice
// will be reflected in the original slice and any other overlapping sub-slices or Golang
// slices (`[]T`)
//
// Analogous to `slice[start:]`
func (s Slice[T]) EndSlice(start int) SubSlice[T] {
	return s.SubSlice(start, int(s.len))
}

// Return a sub-slice of the original slice that cannot be freed
//
// The original `Slice[T]` retains ownership of the data, and any changes to values in the sub-slice
// will be reflected in the original slice and any other overlapping sub-slices or Golang
// slices (`[]T`)
//
// Analogous to `slice[:end]`
func (s Slice[T]) StartSlice(end int) SubSlice[T] {
	return s.SubSlice(0, end)
}

// Return the `Slice[T]` transformed into a standard Golang `[]T` slice
//
// The `Slice[T]` retains ownership of the data, and any changes to values in the sub-slice
// will be reflected in the original slice and any other overlapping sub-slices or Golang
// slices (`[]T`)
func (s Slice[T]) GoSlice() []T {
	return unsafe.Slice(s.ptr, s.cap)[:s.len]
}

type SubSlice[T any] struct {
	ptr *T
	len uint32
	cap uint32
}

func (ss SubSlice[T]) Len() int {
	return int(ss.len)
}

func (ss SubSlice[T]) GetPtr(idx int) *T {
	return &unsafe.Slice(ss.ptr, ss.len)[idx]
}

// Return a sub-slice of this sub-slice that cannot be freed
//
// The original `Slice[T]` retains ownership of the data, and any changes to values in the sub-slice
// will be reflected in the original slice and any other overlapping sub-slices or Golang
// slices (`[]T`)
//
// Analogous to `slice[start:end]`
func (ss SubSlice[T]) SubSlice(start, end int) SubSlice[T] {
	if end > int(ss.len) {
		panic(fmt.Sprintf("fatal: go_manual_memory: SubSlice[T].SubSlice(): end index %d is greater than len %d", end, ss.len))
	}
	goslice := ss.GoSlice()
	len := uint32(end - start)
	sub_ptr := unsafe.SliceData(goslice[start:end])
	return SubSlice[T]{
		ptr: sub_ptr,
		len: len,
		cap: len,
	}
}

// Return a sub-slice of this sub-slice that cannot be freed
//
// The original `Slice[T]` retains ownership of the data, and any changes to values in the sub-slice
// will be reflected in the original slice and any other overlapping sub-slices or Golang
// slices (`[]T`)
//
// Analogous to `slice[:]`
func (ss SubSlice[T]) WholeSlice() SubSlice[T] {
	return ss.SubSlice(0, int(ss.len))
}

// Return a sub-slice of this sub-slice that cannot be freed
//
// The original `Slice[T]` retains ownership of the data, and any changes to values in the sub-slice
// will be reflected in the original slice and any other overlapping sub-slices or Golang
// slices (`[]T`)
//
// Analogous to `slice[start:]`
func (ss SubSlice[T]) EndSlice(start int) SubSlice[T] {
	return ss.SubSlice(start, int(ss.len))
}

// Return a sub-slice of this sub-slice that cannot be freed
//
// The original `Slice[T]` retains ownership of the data, and any changes to values in the sub-slice
// will be reflected in the original slice and any other overlapping sub-slices or Golang
// slices (`[]T`)
//
// Analogous to `slice[:end]`
func (ss SubSlice[T]) StartSlice(end int) SubSlice[T] {
	return ss.SubSlice(0, end)
}

// Return the `SubSlice[T]` transformed into a standard Golang `[]T` slice
//
// The original `Slice[T]` retains ownership of the data, and any changes to values in the sub-slice
// will be reflected in the original slice and any other overlapping sub-slices or Golang
// slices (`[]T`)
func (ss SubSlice[T]) GoSlice() []T {
	return unsafe.Slice(ss.ptr, ss.cap)[:ss.len]
}

var _ ll.SliceLike[byte] = SubSlice[byte]{}
