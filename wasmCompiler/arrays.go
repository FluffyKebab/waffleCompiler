package wasmCompiler

import (
	"compiler/ast"
	"compiler/leb128"
	"compiler/token"
	"compiler/types"
	"compiler/wasmCompiler/code"
)

func (c *compiler) createArrayCode(arrayElementType types.Type, arrayElementsExpression []ast.Node, functionLocals *functionLocals) ([]byte, error) {
	outputCode := make([]byte, 0)

	elementSizeInBytes, err := getArrayTypeElementSize(types.ArrayType{ElementType: arrayElementType})
	if err != nil {
		return []byte{}, err
	}

	arrayFunctionIndex, arrayFunctionTypeIndex, _, err := c.getStandardFunctionIndexTypeIndexAndExtraArguments("array", []types.Type{})

	outputCode = append(outputCode, addConst(len(arrayElementsExpression))...)
	outputCode = append(outputCode, addConst(elementSizeInBytes)...)
	outputCode = append(outputCode, addConst(arrayFunctionIndex)...)
	outputCode = append(outputCode, callIndirect(arrayFunctionTypeIndex)...)

	arrayVariableIndex := functionLocals.defineLocalVariable(types.StandardType{Name: token.INT}, "", c.symbolController)
	outputCode = append(outputCode, code.LOCAL_SET)
	outputCode = append(outputCode, leb128.Int32ToULEB128((int32(arrayVariableIndex)))...)

	for i := 0; i < len(arrayElementsExpression); i++ {
		expressionCode, err := c.compileExpression(arrayElementsExpression[i], functionLocals)
		if err != nil {
			return []byte{}, err
		}

		indexCode := addConst(i)

		variableCode := make([]byte, 0)
		variableCode = append(variableCode, code.LOCAL_GET)
		variableCode = append(variableCode, leb128.Int32ToULEB128(int32(arrayVariableIndex))...)

		curSetterCode, err := c.createSetArrayCode(arrayElementType, variableCode, expressionCode, indexCode)
		if err != nil {
			return []byte{}, err
		}

		curSetterCode = append(curSetterCode, code.DROP)
		outputCode = append(outputCode, curSetterCode...)
	}

	outputCode = append(outputCode, code.LOCAL_GET)
	outputCode = append(outputCode, leb128.Int32ToULEB128((int32(arrayVariableIndex)))...)

	return outputCode, nil
}

func (c *compiler) createSetArrayCode(arrayVariableType types.Type, arrayVariableCode, elementExpressionCode, indexExpressionCode []byte) ([]byte, error) {
	outputCode := make([]byte, 0)

	setterFunctionIndex, setterTypeIndex, _, err := c.getStandardFunctionIndexTypeIndexAndExtraArguments("set", []types.Type{types.ArrayType{ElementType: arrayVariableType}})
	if err != nil {
		return []byte{}, err
	}

	outputCode = append(outputCode, arrayVariableCode...)
	outputCode = append(outputCode, indexExpressionCode...)
	outputCode = append(outputCode, elementExpressionCode...)
	outputCode = append(outputCode, addConst(setterFunctionIndex)...)
	outputCode = append(outputCode, callIndirect(setterTypeIndex)...)

	return outputCode, nil
}

func addConst(constValue int) []byte {
	outputCode := []byte{code.I32_CONST}
	outputCode = append(outputCode, leb128.Int32ToULEB128(int32(constValue))...)
	return outputCode
}
