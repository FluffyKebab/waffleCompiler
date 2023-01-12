package parser

//Functions helping with parsing type literals

import (
	"compiler/errors"
	"compiler/token"
	"compiler/types"
	"fmt"
)

//Gives error if tokens given form i not valid as function type literal and returns index after function type
func parseFunctionTypeLiteral(tokens []token.Token, curTokensIndex int) (types.FunctionType, int, error) {
	outputFunctionType := types.FunctionType{}

	tokensInFirstParenthesis, isValidParenthesis, curTokensIndex := getParenthesisContent(tokens, curTokensIndex, token.LEFT_PARENTHESIS, token.RIGHT_PARENTHESIS)
	if !isValidParenthesis {
		return types.FunctionType{}, curTokensIndex, fmt.Errorf("Internal parser error: tokens given to parse function type literal not parseable as function. No valid parenthesis at curTokens index")
	}

	curTokensIndex++ //skipping the ->

	tokensInSecondParenthesis, isValidParenthesis, curTokensIndex := getParenthesisContent(tokens, curTokensIndex, token.LEFT_PARENTHESIS, token.RIGHT_PARENTHESIS)
	if !isValidParenthesis {
		return types.FunctionType{}, curTokensIndex, fmt.Errorf("Internal parser error: tokens given to parse function type literal not parseable as function. No valid parenthesis after ->")
	}

	argumentTypes, err := getTypesSeparatedByComma(tokensInFirstParenthesis)
	if err != nil {
		return outputFunctionType, curTokensIndex, err
	}

	returnTypes, err := getTypesSeparatedByComma(tokensInSecondParenthesis)
	if err != nil {
		return outputFunctionType, curTokensIndex, err
	}

	outputFunctionType.ArgumentTypes = argumentTypes
	outputFunctionType.ReturnTypes = returnTypes

	return outputFunctionType, curTokensIndex, nil
}

func getTypesSeparatedByComma(tokens []token.Token) ([]types.Type, error) {
	commaSeparatedTokens := splitTokenSliceByComma(tokens) // Split slice into multiple slices that contain tokens between comma

	outputTypes := make([]types.Type, 0)

	for i := 0; i < len(commaSeparatedTokens); i++ {
		if len(commaSeparatedTokens[i]) == 0 {
			continue
		}

		if isStandardTypeLiteral(commaSeparatedTokens[i], 0) {
			if len(commaSeparatedTokens[i]) > 1 {
				return outputTypes, errors.NewSyntaxErrorInvalidToken(commaSeparatedTokens[i][1].Line, commaSeparatedTokens[i][1].Literal)
			}

			outputTypes = append(outputTypes, types.StandardType{
				Name: commaSeparatedTokens[i][0].Literal,
			})

			continue
		}

		if isFunctionTypeLiteral(commaSeparatedTokens[i], 0) {
			functionType, _, err := parseFunctionTypeLiteral(commaSeparatedTokens[i], 0)
			if err != nil {
				return outputTypes, err
			}

			outputTypes = append(outputTypes, functionType)

			continue
		}

		return outputTypes, errors.NewGeneralError(tokens[0].Line, "Tokens not parsable as type literal in function type")
	}

	return outputTypes, nil
}

func isValidAsTypeLiteral(tokens []token.Token, i int) bool {
	return isStandardTypeLiteral(tokens, i) || isFunctionTypeLiteral(tokens, i)
}

// The start of a list of tokens is a functionType if (...) -> (...). Returns
func isFunctionTypeLiteral(tokens []token.Token, i int) bool {
	isFunctionType, _ := skipFunctionTypeLiteral(tokens, i)
	return isFunctionType
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

func isStandardTypeLiteral(tokens []token.Token, i int) bool {
	if i >= len(tokens) {
		return false
	}

	if tokens[i].Type == token.TYPE {
		return true
	}

	return false
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
