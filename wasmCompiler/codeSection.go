package wasmCompiler

import (
	"compiler/ast"
	"compiler/leb128"
	"compiler/symbolTable"
	"compiler/types"
	"compiler/wasmCompiler/code"
)

type codeSection struct {
	numFunctions int
	functionCode [][]uint8
}

func newCodeSection() *codeSection {
	return &codeSection{
		numFunctions: 0,
		functionCode: make([][]uint8, 0),
	}
}

func (s *codeSection) addFunction(newFunctionCode []uint8, index int) {
	for len(s.functionCode) <= index {
		s.functionCode = append(s.functionCode, nil)
	}

	s.functionCode[index] = encodeVector(newFunctionCode)
	s.numFunctions++
}

func (s *codeSection) toByteCode() []uint8 {
	return createSection(code.SECTION_CODE, append(leb128.Int32ToULEB128(int32(s.numFunctions)), flatten(s.functionCode)...))
}

type functionLocals struct {
	parts []struct {
		partType uint8
		num      int32
	}
}

func newFunctionLocals() *functionLocals {
	return &functionLocals{
		parts: make([]struct {
			partType uint8
			num      int32
		}, 0),
	}
}

func (l *functionLocals) defineLocalVariable(variableType types.Type, variableName string, symbolController *symbolTable.SymbolController) int {
	_, variableIndex := symbolController.DefineVariable(variableName, variableType)

	if len(l.parts) == 0 || l.parts[len(l.parts)-1].partType != variableType.ByteCode() {
		l.parts = append(l.parts, struct {
			partType uint8
			num      int32
		}{partType: variableType.ByteCode(), num: 1})
		return variableIndex
	}

	l.parts[len(l.parts)-1].num++
	return variableIndex
}

func (l *functionLocals) toByteCode() []uint8 {
	output := make([]uint8, 0)
	output = append(output, leb128.Int32ToULEB128(int32(len(l.parts)))...)

	for i := 0; i < len(l.parts); i++ {
		output = append(output, leb128.Int32ToULEB128(l.parts[i].num)...)
		output = append(output, l.parts[i].partType)
	}

	return output
}

func (c *compiler) compileFunction(functionBody ast.BlockStatement, functionIndex int) error {
	bodyByteCode := make([]uint8, 0)

	localVariables := newFunctionLocals()

	for i := 0; i < len(functionBody.Statements); i++ {
		switch s := functionBody.Statements[i].(type) {
		case ast.AssignmentStatement:
			expressionCode, err := c.compileExpression(s.Value, localVariables)
			if err != nil {
				return err
			}

			bodyByteCode = append(bodyByteCode, expressionCode...)

			for i := 0; i < len(s.Variables); i++ {
				variableSymbol, isDefined, _ := c.symbolController.Resolve(s.Variables[i].Identifier)
				variableIndex := int(variableSymbol.Index)

				if !isDefined {
					variableIndex = localVariables.defineLocalVariable(s.Variables[i].Type, s.Variables[i].Identifier, c.symbolController)

				}

				bodyByteCode = append(bodyByteCode, code.LOCAL_SET)
				bodyByteCode = append(bodyByteCode, leb128.Int32ToULEB128(int32(variableIndex))...)
			}
		case ast.ReturnStatement:
			returnExpressionsCode := make([]uint8, 0)

			for i := 0; i < len(s.Expressions); i++ {
				expressionCode, err := c.compileExpression(s.Expressions[i], localVariables)
				if err != nil {
					return err
				}

				returnExpressionsCode = append(returnExpressionsCode, expressionCode...)
			}

			bodyByteCode = append(bodyByteCode, returnExpressionsCode...)
			bodyByteCode = append(bodyByteCode, code.RETURN)
		}
	}

	bodyByteCode = append(bodyByteCode, code.END)

	functionCode := localVariables.toByteCode()
	functionCode = append(functionCode, bodyByteCode...)

	c.codeSection.addFunction(functionCode, functionIndex)

	return nil
}

func flatten(a [][]uint8) []uint8 {
	output := make([]uint8, 0)
	for i := 0; i < len(a); i++ {
		output = append(output, a[i]...)
	}

	return output
}
