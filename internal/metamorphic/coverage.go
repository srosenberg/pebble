package metamorphic

import (
	"fmt"
	"math/bits"
	"strconv"
	"unsafe"
)

type Slice struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

func coverage() []byte {
	addr := unsafe.Pointer(&counters)
	size := uintptr(unsafe.Pointer(&ecounters)) - uintptr(addr)

	var res []byte
	*(*Slice)(unsafe.Pointer(&res)) = Slice{
		Data: addr,
		Len:  int(size),
		Cap:  int(size),
	}
	return res
}

func countBits(cov []byte) int {
	n := 0
	for _, c := range cov {
		n += bits.OnesCount8(c)
	}
	return n
}

func countBytes(cov []byte) int {
	n := 0
	for _, c := range cov {
		if bits.OnesCount8(c) > 0 {
			n += 1
		}
	}
	return n
}

func NumEdges() int {
	return len(coverage())
}

//go:linkname counters internal/fuzz._counters
var counters [0]byte

//go:linkname ecounters internal/fuzz._ecounters
var ecounters [0]byte

func BitCoverage() string {
	return fmt.Sprintf("Bit Coverage: %d\n", countBits(coverage()))
}

func EdgeCoverage() string {
	return fmt.Sprintf("Edge Coverage: %d\n", countBytes(coverage()))
}

func EdgeCoverageVerbose() string {
	res := ""
	for i, e := range coverage() {
		if bits.OnesCount8(e) > 0 {
			res += strconv.Itoa(i)
			res += ","
		}
	}
	return "covered edges: " + res
}
