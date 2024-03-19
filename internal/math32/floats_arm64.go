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
func _dot_product_neon(a unsafe.Pointer, b unsafe.Pointer, n uintptr, ret unsafe.Pointer)

//go:noescape
func _squared_l2_neon(a, b unsafe.Pointer, n uintptr, result unsafe.Pointer)

func dot(a, b []float32) float32 {
	switch {
	case useNEON:
		var ret float32

		if len(a) > 0 {
			_dot_product_neon(unsafe.Pointer(&a[0]), unsafe.Pointer(&b[0]), uintptr(len(a)), unsafe.Pointer(&ret))
		}

		return ret
	default:
		return dotGeneric(a, b)
	}
}

func squaredL2(a, b []float32) float32 {
	switch {
	case useNEON:
		var ret float32

		if len(a) > 0 {
			_squared_l2_neon(unsafe.Pointer(&a[0]), unsafe.Pointer(&b[0]), uintptr(len(a)), unsafe.Pointer(&ret))
		}

		return ret
	default:
		return squaredL2Generic(a, b)
	}
}
