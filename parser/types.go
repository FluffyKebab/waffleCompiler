package parser

//Functions helping with parsing type literals

import (
	"compiler/errors"
	"compiler/token"
	"compiler/types"
)

//Returns the parsed type, index after type, bool storing if is valid type and error if a syntax error is found
func parseTypeLiteral(tokens []token.Token, i int) (types.Type, int, bool, error) {
	if standardType, indexAfter, isStandardType, err := parseStandardTypeLiteral(tokens, i); err != nil || isStandardType {
		return standardType, indexAfter, isStandardType, err
	}

	if arrayType, indexAfter, isArrayType, err := parseArrayTypeLiteral(tokens, i); err != nil || isArrayType {
		return arrayType, indexAfter, isArrayType, err
	}

	if functionType, indexAfter, isFunctionType, err := parseFunctionTypeLiteral(tokens, i); err != nil || isFunctionType {
		return functionType, indexAfter, isFunctionType, err
	}

	return types.StandardType{}, i, false, nil
}

func parseStandardTypeLiteral(tokens []token.Token, i int) (types.Type, int, bool, error) {
	if i >= len(tokens) {
		return types.StandardType{}, i, false, nil
	}

	for _, standardType := range types.ValidTypes {
		if tokens[i].Literal == standardType {
			return types.StandardType{Name: standardType}, i + 1, true, nil
		}
	}

	return types.StandardType{}, i, false, nil
}

func parseArrayTypeLiteral(tokens []token.Token, i int) (types.Type, int, bool, error) {
	if i >= len(tokens) {
		return types.StandardType{}, i, false, nil
	}

	if !(tokens[i].Type == token.TYPE && tokens[i].Literal == token.ARRAY_TYPE) {
		return types.StandardType{}, i, false, nil
	}

	arrayElementType, indexAfter, isvalidType, err := parseTypeLiteral(tokens, i+1)
	if err != nil {
		return types.StandardType{}, i, false, err
	}

	if !isvalidType {
		return types.StandardType{}, i, false, errors.NewGeneralError(tokens[i].Line, "valid type after [] is expected")
	}

	return types.ArrayType{ElementType: arrayElementType}, indexAfter, true, nil
}

func parseFunctionTypeLiteral(tokens []token.Token, curTokensIndex int) (types.FunctionType, int, bool, error) {
	outputFunctionType := types.FunctionType{}

	tokensInFirstParenthesis, isValidParenthesis, curTokensIndex := getParenthesisContent(tokens, curTokensIndex, token.LEFT_PARENTHESIS, token.RIGHT_PARENTHESIS)
	if !isValidParenthesis {
		return types.FunctionType{}, curTokensIndex, false, nil
	}

	if tokens[curTokensIndex].Type != token.FUNCTION_ARROW {
		return types.FunctionType{}, curTokensIndex, false, nil
	}
	curTokensIndex++

	tokensInSecondParenthesis, isValidParenthesis, curTokensIndex := getParenthesisContent(tokens, curTokensIndex, token.LEFT_PARENTHESIS, token.RIGHT_PARENTHESIS)
	if !isValidParenthesis {
		return types.FunctionType{}, curTokensIndex, false, nil
	}

	argumentTypes, err := getTypesSeparatedByComma(tokensInFirstParenthesis)
	if err != nil {
		return outputFunctionType, curTokensIndex, false, err
	}

	returnTypes, err := getTypesSeparatedByComma(tokensInSecondParenthesis)
	if err != nil {
		return outputFunctionType, curTokensIndex, false, err
	}

	outputFunctionType.ArgumentTypes = argumentTypes
	outputFunctionType.ReturnTypes = returnTypes

	return outputFunctionType, curTokensIndex, true, nil
}

func getTypesSeparatedByComma(tokens []token.Token) ([]types.Type, error) {
	outputTypes := make([]types.Type, 0)

	for i := 0; true; {
		curType, indexAfter, valid, err := parseTypeLiteral(tokens, i)
		if err != nil {
			return []types.Type{}, err
		}

		if !valid {
			return []types.Type{}, errors.NewGeneralError(tokens[0].Line, "Function type not valid ")
		}

		outputTypes = append(outputTypes, curType)

		if indexAfter >= len(tokens) {
			break
		}

		if tokens[indexAfter].Type != token.COMMA {
			return []types.Type{}, errors.NewGeneralError(tokens[0].Line, "Commas between types in function type is expected")
		}

		i = indexAfter + 1
	}

	return outputTypes, nil
}

func skipFunctionTypeLiteral(tokens []token.Token, i int) (bool, int) {
	isValidParenthesis, i := skipParenthesis(tokens, token.LEFT_PARENTHESIS, token.RIGHT_PARENTHESIS, i)

	if !isValidParenthesis {
		return false, -1
	}

	if i >= len(tokens) {
		return false, -1
	}

	if tokens[i].Type != token.FUNCTION_ARROW {
		return false, -1
	}

	isValidParenthesis, i = skipParenthesis(tokens, token.LEFT_PARENTHESIS, token.RIGHT_PARENTHESIS, i+1)
	if !isValidParenthesis {
		return false, -1
	}

	return true, i
}

func skipParenthesis(tokens []token.Token, leftParenthesis, rightParenthesis string, i int) (bool, int) {
	if i >= len(tokens) {
		return false, i
	}

	if tokens[i].Type != leftParenthesis {
		return false, i
	}

	i++
	depth := 1

	for i < len(tokens) {
		if tokens[i].Type == leftParenthesis {
			depth++
		}

		if tokens[i].Type == rightParenthesis {
			depth--
		}

		i++

		if depth == 0 {
			return true, i
		}
	}

	return false, i
}

//Gets all tokens inside a parenthesis pair starting at index i. Returns index after parenthesis pair and error and bool if not valid parenthesis pair starting at index i of tokens slice.
func getParenthesisContent(tokens []token.Token, i int, leftParenthesis, rightParenthesis string) ([]token.Token, bool, int) {
	if i >= len(tokens) {
		return tokens, false, i
	}

	if tokens[i].Type != leftParenthesis {
		return tokens, false, i
	}

	i++
	depth := 1
	tokensInParenthesis := make([]token.Token, 0)

	for i < len(tokens) {
		if tokens[i].Type == leftParenthesis {
			depth++
		}

		if tokens[i].Type == rightParenthesis {
			depth--
		}

		if depth == 0 {
			return tokensInParenthesis, true, i + 1
		}

		tokensInParenthesis = append(tokensInParenthesis, tokens[i])
		i++
	}

	return tokens, false, i
}

func splitTokenSliceByComma(tokens []token.Token) [][]token.Token {
	output := make([][]token.Token, 0)
	curSlice := make([]token.Token, 0)

	for i := 0; i < len(tokens); i++ {
		parenthesisContent, isParenthesisStart, indexAfter := getParenthesisContent(tokens, i, token.LEFT_PARENTHESIS, token.RIGHT_PARENTHESIS)
		if isParenthesisStart {
			curSlice = append(curSlice, addParenthesis(parenthesisContent, token.LEFT_PARENTHESIS, token.RIGHT_PARENTHESIS, tokens[i].Line)...)
			i = indexAfter - 1
			continue
		}

		functionContent, isFunctionStart, indexAfter := getParenthesisContent(tokens, i, token.START_BLOCK, token.END_BLOCK)
		if isFunctionStart {
			curSlice = append(curSlice, addParenthesis(functionContent, token.START_BLOCK, token.END_BLOCK, tokens[i].Line)...)
			i = indexAfter - 1
			continue
		}

		arrayContent, isArray, indexAfter := getParenthesisContent(tokens, i, token.START_ARRAY, token.END_ARRAY)
		if isArray {
			curSlice = append(curSlice, addParenthesis(arrayContent, token.START_ARRAY, token.END_ARRAY, tokens[i].Line)...)
			i = indexAfter - 1
			continue
		}

		if tokens[i].Type == token.COMMA {
			output = append(output, curSlice)
			curSlice = make([]token.Token, 0)
			continue
		}

		curSlice = append(curSlice, tokens[i])
	}

	return append(output, curSlice)
}
