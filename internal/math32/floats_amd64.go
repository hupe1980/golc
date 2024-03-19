//go:build amd64 && !noasm

package math32

import (
	"unsafe"

	"golang.org/x/sys/cpu"
)

func init() {
	useAVX = cpu.X86.HasAVX
}

//go:noescape
func _dot_product_avx(a, b unsafe.Pointer, n uintptr, result unsafe.Pointer)

//go:noescape
func _squared_l2_avx(a, b unsafe.Pointer, n uintptr, result unsafe.Pointer)

func dot(a, b []float32) float32 {
	switch {
	case useAVX:
		var ret float32

		if len(a) > 0 {
			_dot_product_avx(unsafe.Pointer(&a[0]), unsafe.Pointer(&b[0]), uintptr(len(a)), unsafe.Pointer(&ret))
		}

		return ret
	default:
		return dotGeneric(a, b)
	}
}

func squaredL2(a, b []float32) float32 {
	switch {
	case useAVX:
		var ret float32

		if len(a) > 0 {
			_squared_l2_avx(unsafe.Pointer(&a[0]), unsafe.Pointer(&b[0]), uintptr(len(a)), unsafe.Pointer(&ret))
		}

		return ret
	default:
		return squaredL2Generic(a, b)
	}
}
