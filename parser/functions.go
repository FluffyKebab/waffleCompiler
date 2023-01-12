package parser

import (
	"compiler/ast"
	"compiler/errors"
	"compiler/token"
	"compiler/types"
	"fmt"
)

// Parses expression that start with !
func parseFunctionExecutionExpression(tokens []token.Token) (ast.ExecuteFunctionExpression, error) {
	if len(tokens) == 0 {
		return ast.ExecuteFunctionExpression{}, errors.NewInternalParserError("No ! at index 0 in tokens given to parseFunctionExecutionExpression")
	}

	if tokens[0].Type != token.EXECUTE_FUNCTION {
		return ast.ExecuteFunctionExpression{}, errors.NewInternalParserError("No ! at index 0 in tokens given to parseFunctionExecutionExpression")
	}

	if len(tokens) == 1 {
		return ast.ExecuteFunctionExpression{}, errors.NewGeneralError(tokens[0].Line, "No expression returning function after function execution symbol")
	}

	output := ast.ExecuteFunctionExpression{}

	//split tokens into the different expression. The first expression return the function to be executed
	split, err := splitTokensByExpression(tokens[1:])
	if err != nil {
		return ast.ExecuteFunctionExpression{}, err
	}

	if len(split) == 0 {
		return ast.ExecuteFunctionExpression{}, errors.NewGeneralError(tokens[0].Line, "No expression returning function after function execution symbol")
	}

	functionParsed, err := parseExpression(split[0])
	if err != nil {
		return ast.ExecuteFunctionExpression{}, err
	}

	output.Function = functionParsed

	for i := 1; i < len(split); i++ {
		curArgumentExpression, err := parseExpression(split[i])
		if err != nil {
			return output, err
		}

		output.Arguments = append(output.Arguments, curArgumentExpression)
	}

	return output, nil
}

func parseFunctionDefinitionExpression(tokens []token.Token) (ast.DefineFunctionExpression, error) {
	outputFunction := ast.DefineFunctionExpression{}

	tokensInFirstParenthesis, tokensInSecondParenthesis, functionBodyTokens, _, hasSpecifiedReturnTypes, err := getFunctionDefinitionExpressionParts(tokens, 0)
	if err != nil {
		return outputFunction, err
	}

	arguments, err := parseListOftVariables(tokensInFirstParenthesis)
	if err != nil {
		return outputFunction, err
	}
	outputFunction.Arguments = arguments

	if hasSpecifiedReturnTypes {
		returnTypes, err := getTypesSeparatedByComma(tokensInSecondParenthesis)
		if err != nil {
			return outputFunction, err
		}
		outputFunction.ReturnTypes = returnTypes
	}

	functionIsOneLine := true
	for i := 0; i < len(functionBodyTokens); i++ {
		if functionBodyTokens[i].Type == token.NEWLINE {
			functionIsOneLine = false
			break
		}
	}

	if !functionIsOneLine && !hasSpecifiedReturnTypes {
		return outputFunction, fmt.Errorf("Function with multiple lines must have specified return types")
	}

	outputFunction.NoReturnTypesSpecified = !hasSpecifiedReturnTypes

	if functionIsOneLine { //Adding return if the function is on one line and there is no return there previously
		if len(functionBodyTokens) >= 1 {
			if functionBodyTokens[0].Type != token.RETURN {
				functionBodyTokens = append([]token.Token{token.New(token.RETURN, token.RETURN, functionBodyTokens[0].Line)}, functionBodyTokens...)
			}
		}
	}

	functionBodyParsed, err := newFunctionParser(functionBodyTokens).Parse()
	if err != nil {
		return outputFunction, err
	}

	outputFunction.FunctionBody = functionBodyParsed.Body
	outputFunction.FunctionType = types.FunctionType{ReturnTypes: outputFunction.ReturnTypes}
	argumentsTypes := make([]types.Type, 0)
	for i := 0; i < len(outputFunction.Arguments); i++ {
		if outputFunction.Arguments[i].Type.String() == types.NONE {
			return outputFunction, fmt.Errorf("Error on line %v: function must have defined argument types", tokens[0].Line)
		}

		argumentsTypes = append(argumentsTypes, outputFunction.Arguments[i].Type)
	}
	outputFunction.FunctionType.ArgumentTypes = argumentsTypes

	return outputFunction, nil
}

func isFunctionDefinitionExpression(tokens []token.Token, i int) bool {
	hasArguments, i := skipParenthesis(tokens, token.LEFT_PARENTHESIS, token.RIGHT_PARENTHESIS, i)
	if !hasArguments {
		return false
	}

	if i >= len(tokens) {
		return false
	}

	if tokens[i].Type != token.FUNCTION_ARROW {
		return false
	}
	i++

	hasFunctionBody, _ := skipParenthesis(tokens, token.START_BLOCK, token.END_BLOCK, i)
	if hasFunctionBody {
		return true
	}

	hasReturnTypes, i := skipParenthesis(tokens, token.LEFT_PARENTHESIS, token.RIGHT_PARENTHESIS, i)
	if !hasReturnTypes {
		return false
	}

	hasFunctionBody, _ = skipParenthesis(tokens, token.START_BLOCK, token.END_BLOCK, i)
	if hasFunctionBody {
		return true
	}

	return hasFunctionBody
}

func getFunctionDefinitionExpressionParts(tokens []token.Token, startIndex int) (tokensInFirstParenthesis, tokensInSecondParenthesis, functionBodyTokens []token.Token, indexAfterDefinition int, hasSpecifiedReturnTypes bool, e error) {
	tokensInFirstParenthesis, valid, curTokensIndex := getParenthesisContent(tokens, startIndex, token.LEFT_PARENTHESIS, token.RIGHT_PARENTHESIS)
	if !valid {
		return []token.Token{}, []token.Token{}, []token.Token{}, curTokensIndex, false, fmt.Errorf("Internal parser error: tokens given to getFunctionDefinitionExpressionParts not valid as function definition. No valid parenthesis at index 0")
	}

	curTokensIndex++ //skipping the ->

	if curTokensIndex >= len(tokens) {
		return []token.Token{}, []token.Token{}, []token.Token{}, curTokensIndex, false, fmt.Errorf("Internal parser error: tokens given to getFunctionDefinitionExpressionParts not valid as function definition. No tokens after ->")
	}

	if tokens[curTokensIndex].Type == token.START_BLOCK {
		functionBodyTokens, valid, curTokensIndex = getParenthesisContent(tokens, curTokensIndex, token.START_BLOCK, token.END_BLOCK)
		if !valid {
			return []token.Token{}, []token.Token{}, []token.Token{}, curTokensIndex, true, fmt.Errorf("Internal parser error: tokens given to getFunctionDefinitionExpressionParts not valid as function definition. No valid function body")
		}

		return tokensInFirstParenthesis, []token.Token{}, functionBodyTokens, curTokensIndex, false, nil
	}

	tokensInSecondParenthesis, valid, curTokensIndex = getParenthesisContent(tokens, curTokensIndex, token.LEFT_PARENTHESIS, token.RIGHT_PARENTHESIS)
	if !valid {
		return []token.Token{}, []token.Token{}, []token.Token{}, curTokensIndex, false, fmt.Errorf("Internal parser error: tokens given to getFunctionDefinitionExpressionParts not valid as function definition. No valid parenthesis after ->")
	}

	functionBodyTokens, valid, curTokensIndex = getParenthesisContent(tokens, curTokensIndex, token.START_BLOCK, token.END_BLOCK)
	if !valid {
		return []token.Token{}, []token.Token{}, []token.Token{}, curTokensIndex, false, fmt.Errorf("Internal parser error: tokens given to getFunctionDefinitionExpressionParts not valid as function definition. No function body")
	}

	return tokensInFirstParenthesis, tokensInSecondParenthesis, functionBodyTokens, curTokensIndex, true, nil
}

func removeTokenFromTokenSlice(tokens []token.Token, tokenToRemoveLiteral string) []token.Token {
	outputTokens := make([]token.Token, 0)

	for i := 0; i < len(tokens); i++ {
		if tokens[i].Literal == tokenToRemoveLiteral {
			continue
		}

		outputTokens = append(outputTokens, tokens[i])
	}

	return outputTokens
}

//Splits tokens into expressions. [10 + 1 10 + 3] -> [[10 + 1] [10 + 3]], [10 !f 10 10+1] -> [[10] [!f 10 10 + 1]]
func splitTokensByExpression(tokens []token.Token) ([][]token.Token, error) {
	split := make([][]token.Token, 0)
	curExpression := make([]token.Token, 0)
	lastWasBinaryOperator := true

	for i := 0; i < len(tokens); i++ {
		if isOperator(tokens[i]) {
			lastWasBinaryOperator = true
			curExpression = append(curExpression, tokens[i])
			continue
		}

		if !lastWasBinaryOperator {
			split = append(split, curExpression)
			curExpression = make([]token.Token, 0)
		}

		lastWasBinaryOperator = false

		if tokens[i].Type == token.EXECUTE_FUNCTION {
			curExpression = append(curExpression, tokens[i:]...)
			break
		}

		if isFunctionDefinitionExpression(tokens, i) { //If a list of tokens can be parsed as a function type literal it can also be parsed as function definition expression
			indexBeforeFunction := i
			_, _, _, indexAfterFunction, _, err := getFunctionDefinitionExpressionParts(tokens, i)
			if err != nil {
				return split, err
			}

			curExpression = append(curExpression, tokens[indexBeforeFunction:indexAfterFunction]...)
			i = indexAfterFunction - 1
			continue
		}

		tokensInArray, isValidArray, indexAfterArray := getParenthesisContent(tokens, i, token.START_ARRAY, token.END_ARRAY)
		if isValidArray {
			curExpression = append(curExpression, addParenthesis(tokensInArray, token.START_ARRAY, token.END_ARRAY, tokens[i].Line)...)
			i = indexAfterArray - 1
			continue
		}

		tokensInParenthesis, isValidParenthesis, indexAfterParenthesis := getParenthesisContent(tokens, i, token.LEFT_PARENTHESIS, token.RIGHT_PARENTHESIS)
		if isValidParenthesis {
			curExpression = append(curExpression, addParenthesis(tokensInParenthesis, token.LEFT_PARENTHESIS, token.RIGHT_PARENTHESIS, tokens[i].Line)...)
			i = indexAfterParenthesis - 1
			continue
		}

		curExpression = append(curExpression, tokens[i])
	}

	split = append(split, curExpression)

	return split, nil
}

func addParenthesis(parenthesisContent []token.Token, leftParenthesis, rightParenthesis string, line int) []token.Token {
	output := make([]token.Token, 0)
	output = append(output, token.New(leftParenthesis, leftParenthesis, line))
	output = append(output, parenthesisContent...)
	output = append(output, token.New(rightParenthesis, rightParenthesis, line))

	return output
}

func isOperator(t token.Token) bool {
	for i := 0; i < len(token.Operators); i++ {
		if t.Literal == token.Operators[i] {
			return true
		}
	}

	return false
}

// adder = _ + _
//
// map (3+_) _ . range _

// !(adder 15 _ . adder) 12 12
// !adder 15 (!adder 12 12)

//func parseFunctionTypeLiteral(tokens []token.Token, curTokensIndex int) (ast.FunctionType, int, error) {

//funksjon 34 34
// !

/*
	b = (a int, b int) -> (int) { a * 5 + b  }

	k = !b 10 10

	fn = (a int, b int) -> (int) {
		c = 6
		d = 10
		return a * c + a * d
	}

	k = !fn 10 40
*/
