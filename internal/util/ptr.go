package util

// AddrOrNil returns nil if x is the zero value for T,
// or &x otherwise.
func AddrOrNil[T comparable](x T) *T {
	var z T
	if x == z {
		return nil
	}

	return &x
}
