package wasmCompiler

import (
	"compiler/ast"
	"compiler/symbolTable"
	"compiler/types"
	"fmt"
)

// The function code will be added to the code section and the function name and index will be added to the symbol controller.
func (c *compiler) addGlobalFunction(functionName string, function ast.DefineFunctionExpression) error {

	if _, isDefined, _ := c.symbolController.Resolve(functionName); isDefined {
		return fmt.Errorf("Double declaration in global scope")
	}

	functionType := function.FunctionType
	functionType.TypeIndex = c.typeSection.addType(functionType)

	_, functionIndex := c.symbolController.DefineVariable(functionName, functionType)

	c.funcSection.addFunction(functionType.TypeIndex)
	c.tableSection.addFunction()
	c.elementSection.addFunction(functionIndex)
	c.exportSection.addExport(functionName, functionIndex)

	c.symbolController.PushFunction(c.getFunctionArguments(function.Arguments))
	err := c.compileFunction(function.FunctionBody, functionIndex)
	if err != nil {
		return err
	}

	c.symbolController.PopFunction()

	return nil
}

//Adds function type to type section, function index to function section and function code to code section. Returns the table index / function index and the type index
func (c *compiler) addLocalFunction(function ast.DefineFunctionExpression) (tableIndex int, typeIndex int, e error) {
	functionType := function.FunctionType
	functionIndex := c.symbolController.DefineAnonymousFunction()

	functionType.TypeIndex = c.typeSection.addType(functionType)

	tableIndex = c.tableSection.addFunction()
	c.elementSection.addFunction(functionIndex)
	c.funcSection.addFunction(functionType.TypeIndex)

	c.symbolController.PushFunction(c.getFunctionArguments(function.Arguments))
	err := c.compileFunction(function.FunctionBody, functionIndex)
	if err != nil {
		return -1, -1, err
	}

	c.symbolController.PopFunction()

	return tableIndex, functionType.TypeIndex, nil
}

func (c *compiler) getFunctionArguments(inputArguments []ast.Variable) []symbolTable.Variable {

	outputArguments := make([]symbolTable.Variable, 0)
	for i := 0; i < len(inputArguments); i++ {
		curVariable := symbolTable.Variable{}
		curVariable.Identifier = inputArguments[i].Identifier
		curVariable.Type = inputArguments[i].Type

		curVariableType, isFunction := curVariable.Type.(types.FunctionType)

		if isFunction {
			curVariableType.TypeIndex = c.typeSection.addType(curVariableType)
			curVariable.Type = curVariableType
		}

		outputArguments = append(outputArguments, curVariable)
	}

	return outputArguments
}
