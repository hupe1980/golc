//go:build amd64 && !noasm

package math32

import (
	"golang.org/x/sys/cpu"
)

func init() {
	useAVX512 = cpu.X86.HasAVX512
}

func dot(a, b []float32) float32 {
	switch {
	case useAVX512:
		return dotGeneric(a, b) // TODO
	default:
		return dotGeneric(a, b)
	}
}
