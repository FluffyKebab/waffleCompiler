package validator

import (
	"compiler/ast"
	"compiler/symbolTable"
	"compiler/types"
	"fmt"
)

func Validate(syntaxTree ast.Program) (ast.Program, error) {
	v := validator{symbolController: symbolTable.NewSymbolController()}
	return v.validate(syntaxTree)
}

type validator struct {
	symbolController *symbolTable.SymbolController
}

func (v *validator) validate(syntaxTre ast.Program) (ast.Program, error) {
	validated, returnStatementsReturnTypes, err := v.validateBlockStatement(syntaxTre.Body, false)
	if err != nil {
		return ast.Program{}, err
	}

	if len(returnStatementsReturnTypes) != 0 {
		return ast.Program{}, fmt.Errorf("Return statement in global scope")
	}

	syntaxTre.Body = validated
	return syntaxTre, nil
}

func (v *validator) validateBlockStatement(block ast.BlockStatement, isFunction bool) (ast.BlockStatement, [][]types.Type, error) {
	returnStatementsReturnTypes := make([][]types.Type, 0)

	for i := 0; i < len(block.Statements); i++ {
		switch s := block.Statements[i].(type) {
		case ast.AssignmentStatement:
			functionIsRecursive := false
			if funcDefinitionExpression, isFunctionDefinitionExpression := s.Value.(ast.DefineFunctionExpression); isFunctionDefinitionExpression {
				if isUsingRecursion(funcDefinitionExpression, s.Variables[0].Identifier) {
					functionIsRecursive = true
					if funcDefinitionExpression.NoReturnTypesSpecified {
						return ast.BlockStatement{}, returnStatementsReturnTypes, fmt.Errorf("Function definition using recursion must have specified return types")
					}

					if len(s.Variables) != 1 {
						return ast.BlockStatement{}, returnStatementsReturnTypes, fmt.Errorf("Number of expression return types does not match number of variables in assignment statement")
					}

					v.addVariableToSymbolController(s.Variables[0].Identifier, s.Variables[0].Type, funcDefinitionExpression.FunctionType)
				}
			}

			validated, expressionReturnTypes, err := v.validateExpression(s.Value)
			if err != nil {
				return ast.BlockStatement{}, returnStatementsReturnTypes, err
			}

			s.Value = validated
			block.Statements[i] = s

			if functionIsRecursive {
				continue
			}

			if len(s.Variables) != len(expressionReturnTypes) {
				return ast.BlockStatement{}, returnStatementsReturnTypes, fmt.Errorf("Number of expression return types does not match number of variables in assignment statement")
			}

			for i := 0; i < len(s.Variables); i++ {
				err := v.addVariableToSymbolController(s.Variables[i].Identifier, s.Variables[i].Type, expressionReturnTypes[i])
				if err != nil {
					return ast.BlockStatement{}, returnStatementsReturnTypes, err
				}
			}

		case ast.ReturnStatement:
			if !isFunction {
				return ast.BlockStatement{}, returnStatementsReturnTypes, fmt.Errorf("Return statement in global scope")
			}

			returnExpressionsTypes := make([]types.Type, 0)

			for i := 0; i < len(s.Expressions); i++ {
				validated, curExpressionType, err := v.validateExpression(s.Expressions[i])
				if err != nil {
					return ast.BlockStatement{}, returnStatementsReturnTypes, err
				}

				s.Expressions[i] = validated

				returnExpressionsTypes = append(returnExpressionsTypes, curExpressionType...)
			}

			returnStatementsReturnTypes = append(returnStatementsReturnTypes, returnExpressionsTypes)
		}
	}

	return block, returnStatementsReturnTypes, nil
}

func areListsMatching(list1, list2 []types.Type) bool {
	for i := 0; i < len(list1); i++ {
		if i >= len(list2) {
			return false
		}

		if list1[i].String() != list2[i].String() {
			return false
		}
	}

	return true
}

func generateOperatorError(operator string, leftTypes, rightTypes []types.Type) error {
	errorMessage := "use of operator " + operator + " on "
	for i := 0; i < len(leftTypes); i++ {
		errorMessage += leftTypes[i].String()
		if i+1 != len(leftTypes) {
			errorMessage += ", "
		}
	}

	errorMessage += " and "
	for i := 0; i < len(rightTypes); i++ {
		errorMessage += rightTypes[i].String()
		if i+1 != len(rightTypes) {
			errorMessage += ", "
		}
	}

	errorMessage += " not supported"

	return fmt.Errorf(errorMessage)
}

func isInList(s string, sList []string) bool {
	for i := 0; i < len(sList); i++ {
		if s == sList[i] {
			return true
		}
	}

	return false
}

func (v *validator) addVariableToSymbolController(variableName string, variableType, expressionReturnType types.Type) error {
	variableSymbol, alreadyDefined, isGlobal := v.symbolController.Resolve(variableName)
	if alreadyDefined {
		if isGlobal {
			return fmt.Errorf("Attempt at mutating global variable")
		}

		if variableSymbol.Type.String() != variableType.String() {
			return fmt.Errorf("Attempt at changing variable type")
		}

		return nil
	}

	if variableType.String() != types.NONE {
		if variableType.String() != expressionReturnType.String() {
			return fmt.Errorf("Given variable type does not match return type from expression")
		}
	} else {
		variableType = expressionReturnType
	}

	v.symbolController.DefineVariable(variableName, expressionReturnType)
	return nil
}

func isUsingRecursion(expression ast.Node, functionName string) bool {
	if variable, isVariable := expression.(ast.Variable); isVariable {
		return variable.Identifier == functionName
	}

	childExpressions := expression.GetChildNodes()
	for i := 0; i < len(childExpressions); i++ {
		if isUsingRecursion(childExpressions[i], functionName) {
			return true
		}
	}

	return false
}
