package wasmCompiler

import (
	"compiler/leb128"
	"compiler/types"
	"compiler/wasmCompiler/code"
	"fmt"
)

type typeSection struct {
	typesByteCode   []uint8
	numTypes        int
	typeToTypeIndex map[string]int
}

func newTypeSection() *typeSection {
	return &typeSection{
		typeToTypeIndex: make(map[string]int),
		numTypes:        0,
		typesByteCode:   make([]uint8, 0),
	}
}

//Gets the index of already defined type. If type is not defined the function gives an error.
func (s *typeSection) getFunctionTypeIndex(functionType types.FunctionType) (int, error) {
	if typeIndex, typeIsDeclared := s.typeToTypeIndex[functionType.String()]; typeIsDeclared {
		return typeIndex, nil
	}

	return -1, fmt.Errorf("Internal compiler error: function type (%s) given to get function type index is not defined", functionType.String())
}

//Makes sure the type is in the type section byte code and returns the type index
func (s *typeSection) addType(functionType types.FunctionType) int {
	if typeIndex, typeIsDeclared := s.typeToTypeIndex[functionType.String()]; typeIsDeclared {
		return typeIndex
	}

	typeIndex := s.numTypes
	s.numTypes++
	s.typeToTypeIndex[functionType.String()] = typeIndex

	s.typesByteCode = append(s.typesByteCode, code.FUNC)
	s.typesByteCode = append(s.typesByteCode, encodeVector(listOfTypesToByteCodeOfTypes(functionType.ArgumentTypes))...)
	s.typesByteCode = append(s.typesByteCode, encodeVector(listOfTypesToByteCodeOfTypes(functionType.ReturnTypes))...)

	return typeIndex
}

func (s *typeSection) toByteCode() []byte {
	return createSection(code.SECTION_TYPE, append(leb128.Int32ToULEB128(int32(s.numTypes)), s.typesByteCode...))
}

func listOfTypesToByteCodeOfTypes(inputTypes []types.Type) []uint8 {
	output := make([]uint8, 0)

	for i := 0; i < len(inputTypes); i++ {
		output = append(output, inputTypes[i].ByteCode())
	}

	return output
}
