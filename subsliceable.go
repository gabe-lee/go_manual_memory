package go_manual_memory

type SubSliceable[T any] interface {
	// Return a `SubSlice[T]` of the data analogous to `slice[:]`
	WholeSlice() SubSlice[T]
	// Return a `SubSlice[T]` of the data analogous to `slice[start:end]`
	SubSlice(start uint32, end uint32) SubSlice[T]
	// Return a `SubSlice[T]` of the data analogous to `slice[start:]`
	EndSlice(start uint32) SubSlice[T]
	// Return a `SubSlice[T]` of the data analogous to `slice[:end]`
	StartSlice(end uint32) SubSlice[T]
}
