package wasmCompiler

import (
	"compiler/leb128"
	"compiler/wasmCompiler/code"
)

type functionSection struct {
	numFunction            int
	functionTypeSignatures []uint8
}

func newFunctionSection() *functionSection {
	return &functionSection{
		numFunction:            0,
		functionTypeSignatures: make([]uint8, 0),
	}
}

func (s *functionSection) addFunction(typeIndex int) {
	s.functionTypeSignatures = append(s.functionTypeSignatures, leb128.Int32ToULEB128(int32(typeIndex))...)
	s.numFunction++
}

func (s *functionSection) toByteCode() []uint8 {
	return createSection(code.SECTION_FUNCTION, append(leb128.Int32ToULEB128(int32(s.numFunction)), s.functionTypeSignatures...))
}
