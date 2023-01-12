package wasmCompiler

import (
	"compiler/leb128"
	"compiler/wasmCompiler/code"
)

type elementSection struct {
	functions []int
}

func newElementSection() *elementSection {
	return &elementSection{functions: make([]int, 0)}
}

func (s *elementSection) addFunction(functionIndex int) {
	s.functions = append(s.functions, functionIndex)
}

func (s *elementSection) toByteCode() []uint8 {
	if len(s.functions) == 0 {
		return []uint8{}
	}

	byteCode := leb128.Int32ToULEB128(int32(1))                     //Number of elements
	byteCode = append(byteCode, leb128.Int32ToULEB128(int32(0))...) //Element index

	byteCode = append(byteCode, code.I32_CONST) //Expression
	byteCode = append(byteCode, leb128.Int32ToULEB128(int32(0))...)
	byteCode = append(byteCode, code.END)

	byteCode = append(byteCode, leb128.Int32ToULEB128(int32(len(s.functions)))...)

	for i := 0; i < len(s.functions); i++ {
		byteCode = append(byteCode, leb128.Int32ToULEB128(int32(s.functions[i]))...)
	}

	return createSection(code.SECTION_ELEMENT, byteCode)
}

/*
Example of element section:

9 	Section element
8
1 	num elements
0 	element index
65	i32 const
0	const value
11	end
2	num functions
0	function index
1	function index
*/
