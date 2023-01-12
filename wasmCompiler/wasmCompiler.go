package wasmCompiler

import (
	"compiler/ast"
	"compiler/leb128"
	"compiler/symbolTable"
	"compiler/validator"
	"compiler/wasmCompiler/code"
	"encoding/binary"
	"fmt"
	"math"
)

func Compile(syntaxTree ast.Program) ([]byte, error) {
	c := &compiler{
		typeSection:      newTypeSection(),
		tableSection:     newTableSection(),
		elementSection:   newElementSection(),
		funcSection:      newFunctionSection(),
		codeSection:      newCodeSection(),
		exportSection:    newExportSection(),
		memorySection:    newMemorySection(1),
		symbolController: symbolTable.NewSymbolController(),
	}

	err := c.compile(syntaxTree)
	if err != nil {
		return []byte{}, err
	}

	return c.toByteCode(), nil
}

type compiler struct {
	typeSection      *typeSection
	tableSection     *tableSection
	elementSection   *elementSection
	funcSection      *functionSection
	codeSection      *codeSection
	exportSection    *exportSection
	memorySection    *memorySection
	symbolController *symbolTable.SymbolController
}

func (c *compiler) compile(syntaxTree ast.Program) error {
	validated, err := validator.Validate(syntaxTree)
	if err != nil {
		return err
	}

	err = c.importStandardFunctions()
	if err != nil {
		return err
	}

	for i := 0; i < len(validated.Body.Statements); i++ {
		curStatement := validated.Body.Statements[i]
		assignStatement, ok := curStatement.(ast.AssignmentStatement)
		if !ok {
			return fmt.Errorf("Only function declaration valid in global scope")
		}

		functionDeclaration, ok := assignStatement.Value.(ast.DefineFunctionExpression)
		if !ok {
			return fmt.Errorf("Only function declaration valid in global scope")
		}

		err := c.addGlobalFunction(assignStatement.Variables[0].Identifier, functionDeclaration)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *compiler) toByteCode() []byte {
	result := make([]byte, 0)

	result = append(result, code.MagicModuleHeader...)
	result = append(result, code.ModuleVersion...)
	result = append(result, c.typeSection.toByteCode()...)
	result = append(result, c.funcSection.toByteCode()...)
	result = append(result, c.tableSection.toByteCode()...)
	result = append(result, c.memorySection.toByteCode()...)
	result = append(result, c.exportSection.toByteCode()...)
	result = append(result, c.elementSection.toByteCode()...)
	result = append(result, c.codeSection.toByteCode()...)

	return result
}

func encodeVector(data []uint8) []uint8 {
	return append(leb128.Int32ToULEB128(int32(len(data))), data...)
}

func createSection(sectionType uint8, data []uint8) []uint8 {
	return append([]uint8{sectionType}, encodeVector(data)...)
}

func float32ToLittleEndian(f float32) []byte {
	bits := math.Float32bits(f)
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, bits)
	return buf
}
