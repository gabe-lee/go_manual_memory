package go_manual_memory

import (
	"slices"
	"unsafe"

	esort "github.com/gabe-lee/go_effect_sort"
	ll "github.com/gabe-lee/go_list_like"
)

// This allocator simply uses Golang's normal allocation strategy to create new memory,
// and keeps a cached reference of each memory slice until it is freed, allowing
// the user to store the returned memory in any format desired without fear of
// loss to the garbage collector
type GoAllocator struct {
	slices  [][]byte
	adapter ll.SliceAdapter[[]byte]
}

func NewGoAllocator() *GoAllocator {
	a := GoAllocator{
		slices: make([][]byte, 0),
	}
	a.adapter = ll.New(&a.slices)
	return &a
}

func memSlicesSameAddr(slice ll.SliceLike[[]byte], idx int, val []byte) bool {
	valAddr := uintptr(unsafe.Pointer(unsafe.SliceData(val)))
	mem := ll.Get(slice, idx)
	idxAddr := uintptr(unsafe.Pointer(unsafe.SliceData(mem)))
	idxEnd := idxAddr + uintptr(len(mem))
	return idxAddr <= valAddr && valAddr < idxEnd
}
func memSliceGreaterAddr(slice ll.SliceLike[[]byte], idx int, val []byte) bool {
	valPtr := uintptr(unsafe.Pointer(unsafe.SliceData(val)))
	mem := ll.Get(slice, idx)
	idxPtr := uintptr(unsafe.Pointer(unsafe.SliceData(mem)))
	return idxPtr > valPtr
}

// RawAlloc implements Allocator.
func (g *GoAllocator) RawAlloc(len uintptr, align uintptr) (ptr unsafe.Pointer, alloc_len uintptr) {
	padding := align - 1
	mem := make([]byte, len+align)
	alloc_len = uintptr(cap(mem))
	addr := uintptr(unsafe.Pointer(unsafe.SliceData(mem)))
	alignedAddr := (addr + padding) & ^padding
	delta := alignedAddr - addr
	ptr = unsafe.Pointer(uintptr(unsafe.Pointer(unsafe.SliceData(mem))) + delta)
	alloc_len -= delta
	esort.Sorted_Insert(g.adapter, mem, memSlicesSameAddr, memSliceGreaterAddr, esort.MoveNoSideEffect)
	return
}

// RawFree implements Allocator.
func (g *GoAllocator) RawFree(ptr unsafe.Pointer, len uintptr) {
	mem := unsafe.Slice((*byte)(ptr), len)
	idx, found := esort.Sorted_Search(g.adapter, mem, memSlicesSameAddr, memSliceGreaterAddr)
	if found {
		g.slices = slices.Delete(g.slices, idx, idx+1)
	}
}

// RawResize implements Allocator.
func (g *GoAllocator) RawResizeInPlace(ptr unsafe.Pointer, old_len uintptr, new_len uintptr) (newPtr unsafe.Pointer, success bool) {
	if new_len <= old_len {
		return ptr, true
	}
	mem := unsafe.Slice((*byte)(ptr), old_len)
	idx, found := esort.Sorted_Search(g.adapter, mem, memSlicesSameAddr, memSliceGreaterAddr)
	if found {
		foundMem := g.slices[idx]
		foundMemAddr := uintptr(unsafe.Pointer(unsafe.SliceData(foundMem)))
		ptrAddr := uintptr(ptr)
		delta := ptrAddr - foundMemAddr
		foundCap := uintptr(cap(foundMem))
		foundSpace := foundCap - delta
		if foundSpace >= new_len {
			return ptr, true
		}
	}
	return ptr, false
}

var _ Allocator = (*GoAllocator)(nil)
