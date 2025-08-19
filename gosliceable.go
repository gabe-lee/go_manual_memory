package go_manual_memory

type GoSlicable[T any] interface {
	// Return a Golang slice `[]T` representing the underlying memory of the type
	GoSlice() []T
}
