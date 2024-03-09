//go:build arm64 && !noasm

package math32

import (
	"unsafe"

	"golang.org/x/sys/cpu"
)

func init() {
	useNEON = cpu.ARM64.HasASIMD
}

//go:noescape
func vdotNEON(a unsafe.Pointer, b unsafe.Pointer, n uintptr, ret unsafe.Pointer)

func dot(a, b []float32) float32 {
	switch {
	case useNEON:
		var ret float32

		if len(a) > 0 {
			vdotNEON(unsafe.Pointer(&a[0]), unsafe.Pointer(&b[0]), uintptr(len(a)), unsafe.Pointer(&ret))
		}

		return ret
	default:
		return dotGeneric(a, b)
	}
}
