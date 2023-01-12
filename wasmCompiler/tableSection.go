package wasmCompiler

import (
	"compiler/leb128"
	"compiler/wasmCompiler/code"
)

type tableSection struct {
	numFunctions int
}

func newTableSection() *tableSection {
	return &tableSection{
		numFunctions: 0,
	}
}

func (s *tableSection) addFunction() int {
	s.numFunctions++
	return s.numFunctions - 1
}

func (s *tableSection) toByteCode() []uint8 {
	if s.numFunctions == 0 {
		return []uint8{}
	}

	byteCode := leb128.Int32ToULEB128(int32(1)) //Number of tables?

	byteCode = append(byteCode, code.ANYFUNC)
	byteCode = append(byteCode, leb128.Int32ToULEB128(int32(1))...) //???

	byteCode = append(byteCode, leb128.Int32ToULEB128(int32(s.numFunctions))...) //Table min values
	byteCode = append(byteCode, leb128.Int32ToULEB128(int32(s.numFunctions))...) // table max values

	return createSection(code.SECTION_TABLE, byteCode)
}

/*


 */
