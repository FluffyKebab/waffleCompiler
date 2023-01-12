package validator

import (
	"compiler/ast"
	"compiler/symbolTable"
	"compiler/token"
	"compiler/types"
	"fmt"
)

func (v *validator) validateExpression(expression ast.Node) (ast.Node, []types.Type, error) {
	switch e := expression.(type) {
	case ast.DefineFunctionExpression:
		return v.validateDefineFunctionExpression(e)
	case ast.IfExpression:
		return v.validateIfExpression(e)
	case ast.ExecuteFunctionExpression:
		return v.validateExecuteFunctionExpression(e)
	case ast.OperatorExpression:
		return v.validateOperatorExpression(e)
	case ast.IntExpression:
		return v.validateLiteral(e, token.INT)
	case ast.FloatExpression:
		return v.validateLiteral(e, token.FLOAT)
	case ast.BoolExpression:
		return v.validateLiteral(e, token.BOOL)
	case ast.StringExpression:
		return v.validateLiteral(e, token.STRING)
	case ast.Variable:
		return v.validateVariable(e)
	case ast.ArrayExpression:
		return v.validateArrayExpression(e)
	}

	return expression, []types.Type{}, fmt.Errorf("Node given to validate expression not valid in expression")
}

func (v *validator) validateDefineFunctionExpression(function ast.DefineFunctionExpression) (ast.DefineFunctionExpression, []types.Type, error) {
	argumentVariables := make([]symbolTable.Variable, 0)
	for i := 0; i < len(function.Arguments); i++ {
		argumentVariables = append(argumentVariables, symbolTable.Variable{
			Identifier: function.Arguments[i].Identifier,
			Type:       function.Arguments[i].Type,
		})
	}

	v.symbolController.PushFunction(argumentVariables)
	validated, returnStatementsExpressionsTypes, err := v.validateBlockStatement(function.FunctionBody, true)
	if err != nil {
		return ast.DefineFunctionExpression{}, []types.Type{}, err
	}
	v.symbolController.PopFunction()
	returnTypes := make([]types.Type, 0)

	if function.NoReturnTypesSpecified {
		if len(returnStatementsExpressionsTypes) != 1 {
			return ast.DefineFunctionExpression{}, []types.Type{}, fmt.Errorf("Function with no specified return type must have one and only one return statement")
		}

		returnTypes = returnStatementsExpressionsTypes[0]
	} else {
		returnTypes = function.ReturnTypes
	}

	function.FunctionBody = validated

	functionType := types.FunctionType{ArgumentTypes: VariablesToTypeList(function.Arguments), ReturnTypes: returnTypes}
	function.FunctionType = functionType
	function.ReturnTypes = functionType.ReturnTypes

	return function, []types.Type{functionType}, nil
}

func (v *validator) validateExecuteFunctionExpression(expression ast.ExecuteFunctionExpression) (ast.ExecuteFunctionExpression, []types.Type, error) {
	argumentReturnTypes := make([]types.Type, 0)

	for i := 0; i < len(expression.Arguments); i++ {
		curArgumentValidated, curArgumentReturnTypes, err := v.validateExpression(expression.Arguments[i])
		if err != nil {
			return ast.ExecuteFunctionExpression{}, []types.Type{}, err
		}

		argumentReturnTypes = append(argumentReturnTypes, curArgumentReturnTypes...)
		expression.Arguments[i] = curArgumentValidated
	}

	functionValidated, functionExpressionTypes, err := v.validateExpression(expression.Function)
	if err != nil {
		return ast.ExecuteFunctionExpression{}, []types.Type{}, err
	}
	expression.Function = functionValidated

	if len(functionExpressionTypes) != 1 {
		return ast.ExecuteFunctionExpression{}, []types.Type{}, fmt.Errorf("No function after function execution symbol")
	}

	functionType, isFunction := functionExpressionTypes[0].(types.FunctionType)
	if !isFunction {
		return ast.ExecuteFunctionExpression{}, []types.Type{}, fmt.Errorf("Expression after function execution symbol does not return function")
	}

	anyTypeToRealType, err := validateFunctionExecutionTypes(argumentReturnTypes, functionType.ArgumentTypes)
	if err != nil {
		return ast.ExecuteFunctionExpression{}, []types.Type{}, err
	}

	returnTypes := make([]types.Type, 0)
	for i := 0; i < len(functionType.ReturnTypes); i++ {
		curReturnType, err := insertAnyTypeRealType(functionType.ReturnTypes[i], anyTypeToRealType)
		if err != nil {
			return ast.ExecuteFunctionExpression{}, []types.Type{}, fmt.Errorf("No function after function execution symbol")
		}
		returnTypes = append(returnTypes, curReturnType)
	}

	expression.ReturnTypes = returnTypes
	return expression, returnTypes, nil
}

func insertAnyTypeRealType(returnType types.Type, anyTypeIdentifierToRealType map[string]types.Type) (types.Type, error) {
	if returnTypeAnyType, isAnyType := returnType.(types.AnyType); isAnyType {
		anyTypeRealType, ok := anyTypeIdentifierToRealType[returnTypeAnyType.Name]
		if !ok {
			return types.StandardType{}, fmt.Errorf("Any type with identifier %v not used in function definition", returnTypeAnyType.Name)
		}

		return anyTypeRealType, nil
	}

	if returnTypeArrayType, isArrayType := returnType.(types.ArrayType); isArrayType {
		arrayTypeElementRealType, err := insertAnyTypeRealType(returnTypeArrayType.ElementType, anyTypeIdentifierToRealType)
		if err != nil {
			return types.StandardType{}, err
		}

		return types.ArrayType{ElementType: arrayTypeElementRealType}, nil
	}

	return returnType, nil
}

func validateFunctionExecutionTypes(argumentTypes []types.Type, expectedArgumentTypes []types.Type) (map[string]types.Type, error) {
	if len(argumentTypes) != len(expectedArgumentTypes) {
		return map[string]types.Type{}, fmt.Errorf("Amount of arguments given to function does not match expected number of arguments")
	}

	anyTypeIdentifierToRealType := make(map[string]types.Type)

	for i := 0; i < len(argumentTypes); i++ {
		isEquivalent, newAnyTypeIdentifierToRealType := isActualTypeEquivalentToExpectedType(argumentTypes[i], expectedArgumentTypes[i])
		if !isEquivalent {
			return anyTypeIdentifierToRealType, fmt.Errorf("Argument %v does not match expected argument %v in function. Expected type: %v. Actual type: %v", i, i, argumentTypes[i].String(), expectedArgumentTypes[i].String())
		}

		for anyTypeIdentifier, curAnyTypeRealType := range newAnyTypeIdentifierToRealType {
			curAnyTypeIdentifierToRealTypeValue, hasAnyTypePreviously := anyTypeIdentifierToRealType[anyTypeIdentifier]
			if !hasAnyTypePreviously {
				anyTypeIdentifierToRealType[anyTypeIdentifier] = curAnyTypeRealType
				continue
			}

			if curAnyTypeIdentifierToRealTypeValue.String() != curAnyTypeRealType.String() {
				return anyTypeIdentifierToRealType, fmt.Errorf("Argument %v does not match expected argument %v in function. Expected type: %v. Actual type: %v", i, i, argumentTypes[i].String(), expectedArgumentTypes[i].String())
			}
		}
	}

	return anyTypeIdentifierToRealType, nil
}

func isActualTypeEquivalentToExpectedType(actualType, expectedType types.Type) (bool, map[string]types.Type) {
	if expectedTypeAnyType, expectedTypeIsAnyType := expectedType.(types.AnyType); expectedTypeIsAnyType {
		return true, map[string]types.Type{expectedTypeAnyType.Name: actualType}
	}

	if expectedTypeArrayType, expectedTypeIsArrayType := expectedType.(types.ArrayType); expectedTypeIsArrayType {
		actualTypeArraytype, actualTypeIsArrayType := actualType.(types.ArrayType)
		if !actualTypeIsArrayType {
			return false, make(map[string]types.Type)
		}

		return isActualTypeEquivalentToExpectedType(actualTypeArraytype.ElementType, expectedTypeArrayType.ElementType)
	}

	return actualType.String() == expectedType.String(), make(map[string]types.Type)
}

func (v *validator) validateOperatorExpression(expression ast.OperatorExpression) (ast.OperatorExpression, []types.Type, error) {
	leftSideValidated, leftTypes, err := v.validateExpression(expression.LeftSide)
	if err != nil {
		return ast.OperatorExpression{}, []types.Type{}, err
	}

	rightSideValidated, rightTypes, err := v.validateExpression(expression.RightSide)
	if err != nil {
		return ast.OperatorExpression{}, []types.Type{}, err
	}

	expression.LeftSide = leftSideValidated
	expression.RightSide = rightSideValidated

	if len(leftTypes) != 1 || len(rightTypes) != 1 {
		return ast.OperatorExpression{}, []types.Type{}, generateOperatorError(expression.Operator, leftTypes, rightTypes)
	}

	leftType, isStandardType := leftTypes[0].(types.StandardType)
	if !isStandardType {
		return ast.OperatorExpression{}, []types.Type{}, generateOperatorError(expression.Operator, leftTypes, rightTypes)
	}

	rightType, isStandardType := rightTypes[0].(types.StandardType)
	if !isStandardType {
		return ast.OperatorExpression{}, []types.Type{}, generateOperatorError(expression.Operator, leftTypes, rightTypes)
	}

	if leftType.String() != rightType.String() {
		return ast.OperatorExpression{}, []types.Type{}, generateOperatorError(expression.Operator, leftTypes, rightTypes)
	}

	if expression.Operator == token.AND || expression.Operator == token.OR {
		if leftType.Name != token.BOOL {
			return ast.OperatorExpression{}, []types.Type{}, generateOperatorError(expression.Operator, leftTypes, rightTypes)
		}

		expression.Type = leftTypes[0]
		return expression, []types.Type{types.StandardType{Name: token.BOOL}}, nil
	}

	if expression.Operator == token.EQUAL || expression.Operator == token.NOT_EQUAL {
		if leftType.Name == token.INT || leftType.Name == token.FLOAT || leftType.Name == token.BOOL {
			expression.Type = leftType
			return expression, []types.Type{types.StandardType{Name: token.BOOL}}, nil
		}

		return ast.OperatorExpression{}, []types.Type{}, generateOperatorError(expression.Operator, leftTypes, rightTypes)
	}

	if isInList(expression.Operator, []string{token.PLUS, token.MINUS, token.DIV, token.MULT}) {
		if leftType.Name == token.INT || leftType.Name == token.FLOAT {
			expression.Type = leftType
			return expression, []types.Type{types.StandardType{Name: leftType.Name}}, nil
		}

		return ast.OperatorExpression{}, []types.Type{}, generateOperatorError(expression.Operator, leftTypes, rightTypes)
	}

	if isInList(expression.Operator, []string{token.GREATER_THEN, token.LESS_THEN, token.EQUAL_OR_GREATER_THEN, token.EQUAL_OR_LESS_THEN}) {
		if leftType.Name == token.INT || leftType.Name == token.FLOAT {
			expression.Type = leftType
			return expression, []types.Type{types.StandardType{Name: token.BOOL}}, nil
		}

		return ast.OperatorExpression{}, []types.Type{}, generateOperatorError(expression.Operator, leftTypes, rightTypes)
	}

	return expression, []types.Type{leftType}, fmt.Errorf("Operator in operator expression not valid")
}

func (v *validator) validateLiteral(expression ast.Node, literalType string) (ast.Node, []types.Type, error) {
	return expression, []types.Type{types.StandardType{Name: literalType}}, nil
}

func (v *validator) validateVariable(expression ast.Variable) (ast.Variable, []types.Type, error) {
	variableSymbol, exists, _ := v.symbolController.Resolve(expression.Identifier)
	if !exists {
		if standardFunctionType, isStandardFunction := standardFunctions[expression.Identifier]; isStandardFunction {
			expression.Type = standardFunctionType
			return expression, []types.Type{standardFunctionType}, nil
		}

		return ast.Variable{}, []types.Type{}, fmt.Errorf("Identifier " + expression.Identifier + " is not defined")
	}
	expression.Type = variableSymbol.Type
	return expression, []types.Type{variableSymbol.Type}, nil
}

func (v *validator) validateArrayExpression(expression ast.ArrayExpression) (ast.ArrayExpression, []types.Type, error) {
	if len(expression.ElementsExpressions) == 0 {
		return expression, []types.Type{}, fmt.Errorf("Array literal with zero expressions not valid. Use make")
	}

	var arrayElementsType types.Type = nil

	for i := 0; i < len(expression.ElementsExpressions); i++ {
		curElementValidated, curElementType, err := v.validateExpression(expression.ElementsExpressions[i])
		if err != nil {
			return expression, []types.Type{}, err
		}

		if len(curElementType) != 1 {
			return expression, []types.Type{}, fmt.Errorf("Array expressions must return single value")
		}

		expression.ElementsExpressions[i] = curElementValidated

		if arrayElementsType == nil {
			arrayElementsType = curElementType[0]
			continue
		}

		if curElementType[0].String() != arrayElementsType.String() {
			return expression, []types.Type{}, fmt.Errorf("Can not add element of type %s to array of type []%s", curElementType[0].String(), arrayElementsType.String())
		}
	}

	expression.Type = types.ArrayType{ElementType: arrayElementsType}
	return expression, []types.Type{types.ArrayType{ElementType: arrayElementsType}}, nil
}

func (v *validator) validateIfExpression(expression ast.IfExpression) (ast.IfExpression, []types.Type, error) {
	conditionValidated, conditionReturnTypes, err := v.validateExpression(expression.Condition)
	if err != nil {
		return ast.IfExpression{}, []types.Type{}, err
	}

	trueValidated, trueExpressionReturnTypes, err := v.validateExpression(expression.TrueExpression)
	if err != nil {
		return ast.IfExpression{}, []types.Type{}, err
	}

	falseValidated, falseExpressionReturnTypes, err := v.validateExpression(expression.FalseExpression)
	if err != nil {
		return ast.IfExpression{}, []types.Type{}, err
	}

	if len(conditionReturnTypes) != 1 {
		return ast.IfExpression{}, []types.Type{}, fmt.Errorf("Condition expression must return one value in if expression")
	}

	if conditionReturnTypes[0].String() != token.BOOL {
		return ast.IfExpression{}, []types.Type{}, fmt.Errorf("Condition expression must return a value of type bool in if expression")
	}

	if len(trueExpressionReturnTypes) != len(falseExpressionReturnTypes) {
		return ast.IfExpression{}, []types.Type{}, fmt.Errorf("Return values of true and false expressions in if expression must be the same ")
	}

	for i := 0; i < len(trueExpressionReturnTypes); i++ {
		if trueExpressionReturnTypes[i].String() != falseExpressionReturnTypes[i].String() {
			return ast.IfExpression{}, []types.Type{}, fmt.Errorf("Return values of true and false expressions in if expression must be the same ")
		}
	}

	expression.ReturnType = trueExpressionReturnTypes
	expression.Condition = conditionValidated
	expression.TrueExpression = trueValidated
	expression.FalseExpression = falseValidated

	return expression, expression.ReturnType, nil
}

func VariablesToTypeList(variables []ast.Variable) []types.Type {
	variablesType := make([]types.Type, 0)
	for i := 0; i < len(variables); i++ {
		variablesType = append(variablesType, variables[i].Type)
	}
	return variablesType
}
