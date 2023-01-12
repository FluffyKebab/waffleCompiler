package wasmCompiler

import (
	"compiler/leb128"
	"compiler/wasmCompiler/code"
)

type memorySection struct {
	size int
}

func newMemorySection(size int) *memorySection {
	return &memorySection{
		size: size,
	}
}

func (s *memorySection) toByteCode() []byte {
	byteCode := leb128.Int32ToULEB128(1) //1 storing number of memories
	byteCode = append(byteCode, newLimit(s.size, s.size)...)
	return createSection(code.SECTION_MEMORY, byteCode)
}

func newLimit(min, max int) []byte {
	byteCode := []byte{code.LIMIT_MIN_MAX}
	byteCode = append(byteCode, leb128.Int32ToULEB128(int32(min))...)
	byteCode = append(byteCode, leb128.Int32ToULEB128(int32(max))...)

	return byteCode
}
