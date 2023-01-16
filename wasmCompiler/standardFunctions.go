package wasmCompiler

import (
	"compiler/readWasm"
	"compiler/token"
	"compiler/types"
	"fmt"
)

type typeAndFuncIndex struct {
	typeIndex int
	funcIndex int
}

type standardFunctions struct {
	standardFunctionIndexes map[string]typeAndFuncIndex
}

func (c *compiler) getStandardFunctionIndexTypeIndexAndExtraArguments(name string, arguments []types.Type) (int, int, []byte, error) {
	realFunctionName, err := getStandardFunctionRealName(name, arguments)
	if err != nil {
		return 0, 0, []byte{}, err
	}

	extraArguments, err := getStandardFunctionExtraArguments(name, arguments)
	if err != nil {
		return 0, 0, []byte{}, err
	}

	if indexes, isImported := c.standardFunctions.standardFunctionIndexes[realFunctionName]; isImported {
		return indexes.funcIndex, indexes.typeIndex, extraArguments, nil
	}

	funcIndex, typeIndex, err := c.importStandardFunction(realFunctionName)
	return funcIndex, typeIndex, extraArguments, err
}

//the memory handler must be added before all other function because multiple standard functions written in wasm reference these functions indexes with 0, 1 and 2
func (c *compiler) importMemoryHandler() error {
	_, _, err := c.importStandardFunction("allocate")
	_, _, err = c.importStandardFunction("deAllocate")
	_, _, err = c.importStandardFunction("array")
	return err
}

//Returns func index and type index
func (c *compiler) importStandardFunction(functionName string) (int, int, error) {
	for i := 0; i < len(standardFunctionsData); i++ {
		if standardFunctionsData[i].name != functionName {
			continue
		}

		functionCode, err := readWasm.GetFuncFromFile(standardFunctionsData[i].fileName, standardFunctionsData[i].funcIndex)
		if err != nil {
			return 0, 0, fmt.Errorf("Internal compiler error: Error getting standard function %v from file %v: %v", standardFunctionsData[i].funcIndex, standardFunctionsData[i].fileName, err.Error())
		}

		funcIndex := c.symbolController.DefineAnonymousFunction()
		typeIndex := c.typeSection.addType(standardFunctionsData[i].funcType)

		c.funcSection.addFunction(typeIndex)
		c.tableSection.addFunction()
		c.elementSection.addFunction(funcIndex)
		c.codeSection.addFunction(functionCode, funcIndex)
		c.standardFunctions.standardFunctionIndexes[functionName] = typeAndFuncIndex{funcIndex: funcIndex, typeIndex: typeIndex}

		return funcIndex, typeIndex, nil
	}

	return 0, 0, fmt.Errorf("Internal compiler error: Standard function with name %s not found in standard function data", functionName)
}

func getStandardFunctionRealName(functionName string, functionArguments []types.Type) (string, error) {
	for _, functionNameNotDependingOnArgumentsTypes := range []string{"array", "allocate", "deAllocate", "length", "take"} {
		if functionNameNotDependingOnArgumentsTypes == functionName {
			return functionName, nil
		}
	}

	if functionName == "get" || functionName == "set" {
		if len(functionArguments) < 1 {
			return "", fmt.Errorf("Error in validation process: wrong amount of arguments in %s call", functionName)
		}

		typePrefix, err := getArrayTypePrefix(functionArguments[0])
		return typePrefix + functionName, err
	}

	return "", fmt.Errorf("Internal compiler error: getting real name of %s not implemented", functionName)
}

func getStandardFunctionExtraArguments(functionName string, functionArguments []types.Type) ([]byte, error) {
	switch functionName {
	case "take":
		if len(functionArguments) != 2 {
			return []byte{}, fmt.Errorf("Error in validation process: wrong amount of arguments to take ")
		}

		sizeOfElementsInArray, err := getArrayTypeElementSize(functionArguments[1])
		if err != nil {
			return []byte{}, err
		}

		return addConst(sizeOfElementsInArray), nil
	}

	return []byte{}, nil
}

func getArrayTypePrefix(inputType types.Type) (string, error) {
	arrayType, isArrayType := inputType.(types.ArrayType)
	if !isArrayType {
		return "", fmt.Errorf("Internal compiler error: argument inputType in getArrayTypePrefix not of arrayType")
	}

	switch t := arrayType.ElementType.(type) {
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
			return "", fmt.Errorf("Type %s given to getArrayTypePrefix not supported", inputType)
		}

	case types.ArrayType:
		return "i32", nil
	case types.FunctionType:
		return "i32", nil
	}

	return "", fmt.Errorf("Type %s given to getArrayTypePrefix not supported", inputType)
}

func getArrayTypeElementSize(inputType types.Type) (int, error) {
	arrayType, isArrayType := inputType.(types.ArrayType)
	if !isArrayType {
		return 0, fmt.Errorf("Internal compiler error: argument inputType in getArrayTypeElementSize not of arrayType")
	}

	switch t := arrayType.ElementType.(type) {
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
			return 0, fmt.Errorf("Type %s given to getArrayTypeElementSize not supported", arrayType)
		}

	case types.ArrayType:
		return 4, nil
	case types.FunctionType:
		return 4, nil
	}

	return 0, fmt.Errorf("Type %s given to getArrayTypeElementSize not supported", arrayType)
}
