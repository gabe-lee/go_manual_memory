package go_manual_memory

import "unsafe"

const PAGE_SIZE = 8192

type Allocator interface {
	RawAlloc(len, align uintptr) (ptr unsafe.Pointer, alloc_len uintptr)
	RawResizeInPlace(ptr unsafe.Pointer, old_len, new_len uintptr) (newPtr unsafe.Pointer, success bool)
	RawFree(ptr unsafe.Pointer, len uintptr)
}

func Alloc[T any](alloc Allocator, len int) (mem []T) {
	size := unsafe.Sizeof(*new(T))
	align := unsafe.Alignof(*new(T))
	ptr, allocLen := alloc.RawAlloc(size*uintptr(len), align)
	cap := allocLen / size
	mem = unsafe.Slice((*T)(ptr), cap)[:len]
	return
}

func ResizeInPlace[T any](alloc Allocator, mem []T, newLen int) (newMem []T, success bool) {
	size := unsafe.Sizeof(*new(T))
	byteLen := size * uintptr(cap(mem))
	newByteLen := size * uintptr(newLen)
	ptr := unsafe.Pointer(unsafe.SliceData(mem))
	newPtr, success := alloc.RawResizeInPlace(ptr, byteLen, newByteLen)
	newMem = unsafe.Slice((*T)(newPtr), newLen)
	return
}

func ResizeCanMove[T any](alloc Allocator, mem []T, newLen int) (newMem []T) {
	newMem, success := ResizeInPlace(alloc, mem, newLen)
	if success {
		return
	}
	newMem = Alloc[T](alloc, newLen)
	copy(newMem, mem)
	Free(alloc, mem)
	return
}

func Free[T any](alloc Allocator, mem []T) {
	size := unsafe.Sizeof(*new(T))
	byteLen := size * uintptr(cap(mem))
	ptr := unsafe.Pointer(unsafe.SliceData(mem))
	alloc.RawFree(ptr, byteLen)
}

// Create a single-item (scalar) pointer to a new value of type `T`
//
// For many-item (vector) allocations, use the dedicated create function
// for the vector type instead
func Create[T any](alloc Allocator) *T {
	size := unsafe.Sizeof(*new(T))
	align := unsafe.Alignof(*new(T))
	ptr, _ := alloc.RawAlloc(size, align)
	return (*T)(ptr)
}

// Destroy (free) a single-item (scalar) pointer to a value of type `T`
//
// For many-item (vector) de-allocations, use the dedicated method on the
// vector type instead
func Destroy[T any](alloc Allocator, ptr *T) {
	size := unsafe.Sizeof(*ptr)
	alloc.RawFree(unsafe.Pointer(ptr), size)
}
