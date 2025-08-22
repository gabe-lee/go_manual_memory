package go_manual_memory

// import "unsafe"

// // This allocator ALWAYS returns addresses
// type MiniSlabAllocator struct {
// 	parentAlloc       Allocator
// 	slabIds           Slice[uint32]
// 	addresses         Slice[uintptr]
// 	maxSlabSize       uint32
// 	maxSlabSizeBits   uint32
// 	maxSlabCount      uint32
// 	slabIdMask        uint32
// 	extendedSlabCount uint32
// }

// // Expand implements MiniAllocator.
// func (m *MiniSlabAllocator) Expand(addr MiniAddr) unsafe.Pointer {
// 	panic("unimplemented")
// }

// // RawAlloc implements MiniAllocator.
// func (m *MiniSlabAllocator) RawAlloc(len uint32, align uint32) (addr MiniAddr, alloc_len uint32) {
// 	uptr, ulen := m.parentAlloc.RawAlloc(uintptr(len), uintptr(align))
// 	uaddr := uintptr(uptr)
// 	lowAddr :=
// }

// // RawFree implements MiniAllocator.
// func (m *MiniSlabAllocator) RawFree(addr MiniAddr, len uint32) {
// 	panic("unimplemented")
// }

// // RawResizeInPlace implements MiniAllocator.
// func (m *MiniSlabAllocator) RawResizeInPlace(addr MiniAddr, old_len uint32, new_len uint32) (newAddr MiniAddr, success bool) {
// 	panic("unimplemented")
// }

// var _ MiniAllocator = (*MiniSlabAllocator)(nil)
