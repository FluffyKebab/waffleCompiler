package wasmCompiler

import (
	"compiler/leb128"
	"compiler/wasmCompiler/code"
)

type export struct {
	functionName  string
	functionIndex int32
}

type exportSection struct {
	exports []export
}

func newExportSection() *exportSection {
	return &exportSection{make([]export, 0)}
}

func (s *exportSection) addExport(functionName string, functionIndex int) {
	s.exports = append(s.exports, export{
		functionName:  functionName,
		functionIndex: int32(functionIndex),
	})
}

func (s *exportSection) toByteCode() []uint8 {
	byteCode := leb128.Int32ToULEB128(int32(len(s.exports)))

	for i := 0; i < len(s.exports); i++ {
		byteCode = append(byteCode, stringToByteCode(s.exports[i].functionName)...)
		byteCode = append(byteCode, leb128.Int32ToULEB128(s.exports[i].functionIndex)...)
	}

	return createSection(code.SECTION_EXPORT, byteCode)
}

func stringToByteCode(s string) []uint8 {
	byteCode := leb128.Int32ToULEB128(int32(len(s)))

	byteCode = append(byteCode, []byte(s)...)
	byteCode = append(byteCode, 0)

	return byteCode
}
