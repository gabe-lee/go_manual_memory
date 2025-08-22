package go_manual_memory

// import "unsafe"

// type MiniAddr uint32

// type MiniAllocator interface {
// 	RawAlloc(len, align uint32) (addr MiniAddr, alloc_len uint32)
// 	RawResizeInPlace(addr MiniAddr, old_len, new_len uint32) (newAddr MiniAddr, success bool)
// 	RawFree(addr MiniAddr, len uint32)
// 	Expand(addr MiniAddr) unsafe.Pointer
// }
