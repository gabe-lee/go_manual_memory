package go_manual_memory

import "unsafe"

func UnsafeCastPtr[IN any, OUT any](in *IN) *OUT {
	return (*OUT)(unsafe.Pointer(in))
}

func UnsafeCast[IN any, OUT any](in IN) OUT {
	return *(*OUT)(unsafe.Pointer(&in))
}
