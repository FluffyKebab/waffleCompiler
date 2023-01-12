package wasmCompiler

import (
	"compiler/ast"
	"compiler/leb128"
	"compiler/readWasm"
	"compiler/token"
	"compiler/types"
	"compiler/wasmCompiler/code"
	"fmt"
)

func (c *compiler) importStandardFunctions() error {
	for i := 0; i < len(standardFunctions); i++ {
		functionCode, err := readWasm.GetFuncFromFile(standardFunctions[i].fileName, standardFunctions[i].funcIndex)
		if err != nil {
			return fmt.Errorf("Error getting standard function %v from file %v: %v", standardFunctions[i].funcIndex, standardFunctions[i].fileName, err.Error())
		}

		var funcIndex int
		if standardFunctions[i].open {
			_, funcIndex = c.symbolController.DefineVariable(standardFunctions[i].name, standardFunctions[i].funcType)
		} else {
			funcIndex = c.symbolController.DefineAnonymousFunction()
		}

		typeIndex := c.typeSection.addType(standardFunctions[i].funcType)

		c.funcSection.addFunction(typeIndex)
		c.tableSection.addFunction()
		c.elementSection.addFunction(funcIndex)
		c.codeSection.addFunction(functionCode, funcIndex)
	}

	return nil
}

func (c *compiler) createArrayCode(arrayElementType types.Type, arrayElementsExpression []ast.Node, functionLocals *functionLocals) ([]byte, error) {
	outputCode := make([]byte, 0)

	elementSizeInBytes, err := getArrayTypeElementSize(arrayElementType)
	if err != nil {
		return []byte{}, err
	}

	outputCode = append(outputCode, addConst(len(arrayElementsExpression))...)
	outputCode = append(outputCode, addConst(elementSizeInBytes)...)
	outputCode = append(outputCode, addConst(standardFunctionsIndex["array"])...)

	funcTypeIndex, err := c.getStandardFunctionTypeIndex("array")
	if err != nil {
		return []byte{}, err
	}
	outputCode = append(outputCode, callIndirect(funcTypeIndex)...)

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

func (c *compiler) getStandardFunctionTypeIndex(funcName string) (int, error) {
	return c.typeSection.getFunctionTypeIndex(standardFunctions[standardFunctionsIndex[funcName]].funcType)
}

func (c *compiler) createSetArrayCode(arrayVariableType types.Type, arrayVariableCode, elementExpressionCode, indexExpressionCode []byte) ([]byte, error) {
	outputCode := make([]byte, 0)

	setterFunctionIndex, setterTypeIndex, err := c.getSetOrGetArraySetFunctionIndexAndTypeIndex("set", arrayVariableType)
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

func isStandardFunction(functionName string) bool {
	switch functionName {
	case "get":
		return true
	case "set":
		return true
	default:
		return false
	}
}

//Return type index, table index and an error
func (c *compiler) getStandardFunctionIndexAndTypeIndex(functionName string, argumentTypes []types.Type) (int, int, error) {

	if functionName == "get" || functionName == "set" {
		if len(argumentTypes) < 1 {
			return 0, 0, fmt.Errorf("Error in validation process: wrong amount of arguments in %s call", functionName)
		}

		arrayType, firstIsArrayType := argumentTypes[0].(types.ArrayType)
		if !firstIsArrayType {
			return 0, 0, fmt.Errorf("Error in validation process: first argument of %s is not array type", functionName)
		}

		return c.getSetOrGetArraySetFunctionIndexAndTypeIndex(functionName, arrayType.ElementType)
	}

	return 0, 0, fmt.Errorf("Function name %s given to getStandardFunctionIndex not supported", functionName)
}

func getStandardFunctionIndex(functionName string) (int, error) {
	functionIndex, ok := standardFunctionsIndex[functionName]
	if !ok {
		return 0, fmt.Errorf("Function %v not found in standardFunctionIndex.", functionName)
	}

	return functionIndex, nil
}

func (c *compiler) getSetOrGetArraySetFunctionIndexAndTypeIndex(functionName string, arrayElementType types.Type) (int, int, error) {
	typePrefix, err := getArrayTypePrefix(arrayElementType)
	if err != nil {
		return 0, 0, err
	}

	realFunctionName := typePrefix + functionName

	functionIndex, ok := standardFunctionsIndex[realFunctionName]
	if !ok {
		return 0, 0, fmt.Errorf("Error running getArraySetFunctionIndex: Function %vget not found in standardFunctionIndex.", typePrefix)
	}

	typeIndex, err := c.getStandardFunctionTypeIndex(realFunctionName)
	if err != nil {
		return 0, 0, err
	}

	return functionIndex, typeIndex, nil
}

func getArrayTypePrefix(arrayElementType types.Type) (string, error) {
	switch t := arrayElementType.(type) {
	case types.StandardType:
		switch t.Name {
		case token.INT:
			return "i32", nil
		case token.FLOAT:
			return "f32", nil
		case token.BOOL:
			return "i8", nil
		case token.STRING:
			return "i32", nil
		default:
			return "", fmt.Errorf("Type %s given to getArrayTypePrefix not supported", arrayElementType.String())
		}

	case types.ArrayType:
		return "i32", nil
	case types.FunctionType:
		return "i32", nil
	}

	return "", fmt.Errorf("Type %s given to getArrayTypePrefix not supported", arrayElementType.String())
}

func getArrayTypeElementSize(arrayType types.Type) (int, error) {
	switch t := arrayType.(type) {
	case types.StandardType:
		switch t.Name {
		case token.INT:
			return 4, nil
		case token.FLOAT:
			return 4, nil
		case token.BOOL:
			return 1, nil
		case token.STRING:
			return 4, nil
		default:
			return 0, fmt.Errorf("Type %s given to getArrayTypeElementSize not supported", arrayType.String())
		}

	case types.ArrayType:
		return 4, nil
	case types.FunctionType:
		return 4, nil
	}

	return 0, fmt.Errorf("Type %s given to getArrayTypeElementSize not supported", arrayType.String())
}

func addConst(constValue int) []byte {
	outputCode := []byte{code.I32_CONST}
	outputCode = append(outputCode, leb128.Int32ToULEB128(int32(constValue))...)
	return outputCode
}
