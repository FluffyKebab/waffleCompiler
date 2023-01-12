package parser

import (
	"compiler/ast"
	"compiler/errors"
	"compiler/token"
	"compiler/types"
	"fmt"
	"strconv"
)

func (p *parser) GetAllTokensInExpression() ([]token.Token, error) {
	functionDepth := 0
	outputTokens := make([]token.Token, 0)

	for true {
		if p.curToken.Type == token.EOF {
			break
		}

		if p.curToken.Type == token.START_BLOCK {
			functionDepth++
		}

		if p.curToken.Type == token.END_BLOCK {
			if functionDepth == 0 {
				return outputTokens, errors.NewSyntaxErrorInvalidToken(p.curToken.Line, token.END_BLOCK)
			}

			functionDepth--
		}

		if functionDepth == 0 && p.curToken.Type == token.NEWLINE {
			return outputTokens, nil
		}

		outputTokens = append(outputTokens, p.curToken)
		p.NextToken()

	}

	if functionDepth != 0 {
		return outputTokens, errors.NewSyntaxErrorUnexpectedToken(p.curToken.Line, "end of file", "end of function")
	}

	return outputTokens, nil
}

func parseExpression(tokens []token.Token) (ast.Node, error) {

	if len(tokens) == 0 {
		return ast.StringExpression{}, fmt.Errorf("NO TOKEN")
	}

	if len(tokens) == 1 {
		if tokens[0].Type == token.INT {
			tokenValue, err := strconv.Atoi(tokens[0].Literal)
			if err != nil {
				return ast.StringExpression{}, errors.NewInternalParserError("Token with int type not parsable as int: " + err.Error())
			}

			return ast.IntExpression{
				Value: int32(tokenValue),
			}, nil
		}

		if tokens[0].Type == token.FLOAT {
			tokenValue, err := strconv.ParseFloat(tokens[0].Literal, 64)
			if err != nil {
				return ast.StringExpression{}, errors.NewInternalParserError("Token with int type not parsable as int: " + err.Error())
			}

			return ast.FloatExpression{
				Value: tokenValue,
			}, nil
		}

		if tokens[0].Type == token.STRING {
			return ast.StringExpression{
				Value: tokens[0].Literal,
			}, nil
		}

		if tokens[0].Type == token.BOOL {
			tokenValue := false
			if tokens[0].Literal == token.TRUE {
				tokenValue = true
			} else if tokens[0].Literal == token.FALSE {
				tokenValue = false
			} else {
				return ast.BoolExpression{}, errors.NewInternalParserError("Token with type bool not parsable as bool")
			}

			return ast.BoolExpression{
				Value: tokenValue,
			}, nil
		}

		if tokens[0].Type == token.VARIABLE {
			return ast.Variable{
				Identifier: tokens[0].Literal,
				Type:       types.StandardType{Name: types.NONE},
			}, nil
		}

		return ast.IntExpression{}, errors.NewSyntaxErrorInvalidToken(tokens[0].Line, tokens[0].Literal)
	}

	// If expression has no operators and starts with ! the next part of the expression must be a variable or function literal
	if tokens[0].Type == token.EXECUTE_FUNCTION {
		return parseFunctionExecutionExpression(tokens)
	}

	if tokens[0].Type == token.IF {
		return parseIfExpression(tokens)
	}

	//Find operator to be executed last
	// 	leftmost comparative operator that is not in parenthesis
	// 	leftmost + or - that is not in parenthesis
	// 	leftmost * or / that is not in parenthesis

	operator, pos, err := findLeftmostTokenOfType([]string{
		token.OR,
		token.EQUAL,
		token.NOT_EQUAL,
		token.EQUAL_OR_GREATER_THEN,
		token.EQUAL_OR_LESS_THEN,
	}, tokens, true)

	if pos != -1 {
		return createOperatorExpression(tokens, operator, pos)
	}

	operator, pos, err = findLeftmostTokenOfType([]string{
		token.PLUS,
		token.MINUS,
	}, tokens, true)

	if pos != -1 {
		return createOperatorExpression(tokens, operator, pos)
	}

	operator, pos, err = findLeftmostTokenOfType([]string{
		token.MULT,
		token.DIV,
	}, tokens, true)

	if pos != -1 {
		return createOperatorExpression(tokens, operator, pos)
	}

	if err != nil {
		return ast.IntExpression{}, err
	}

	if isFunctionDefinitionExpression(tokens, 0) {
		return parseFunctionDefinitionExpression(tokens)
	}

	if tokens[0].Type == token.LEFT_PARENTHESIS {
		return parseParenthesisExpression(tokens)
	}

	if tokens[0].Type == token.START_ARRAY {
		return parseArrayExpression(tokens)
	}

	tokensString := ""
	for i := 0; i < len(tokens); i++ {
		tokensString += tokens[i].Literal + " "
	}

	return ast.StringExpression{}, fmt.Errorf("Error on line %v: Unable to parse as expression: %s", tokens[0].Line, tokensString)
}

func parseIfExpression(tokens []token.Token) (ast.IfExpression, error) {
	if len(tokens) == 0 {
		return ast.IfExpression{}, errors.NewInternalParserError("Length of tokens given to parseIfExpression is 0, must be more at least 1")
	}

	_, elsePos, err := findLeftmostTokenOfType([]string{token.ELSE}, tokens[1:], false)
	if err != nil {
		return ast.IfExpression{}, err
	}

	if elsePos == -1 {
		return ast.IfExpression{}, fmt.Errorf("No else expression found after if expression. All if expressions must have an else expression after them")
	}

	elsePos++ //Adding one because the search started at index 1

	outputIfExpression := ast.IfExpression{}

	outputIfExpression.FalseExpression, err = parseExpression(tokens[elsePos+1:])
	if err != nil {
		return ast.IfExpression{}, err
	}

	ifExpressionToken := tokens[1:elsePos]
	ifExpressionSplit, err := splitTokensByExpression(ifExpressionToken)
	if err != nil {
		return ast.IfExpression{}, err
	}

	if len(ifExpressionSplit) != 2 {
		return ast.IfExpression{}, fmt.Errorf("Amount of expressions between if and else token is not two. There must be two expression between the if and else token. The first returning a boolean value.")
	}

	outputIfExpression.Condition, err = parseExpression(ifExpressionSplit[0])
	if err != nil {
		return ast.IfExpression{}, err
	}

	outputIfExpression.TrueExpression, err = parseExpression(ifExpressionSplit[1])
	if err != nil {
		return ast.IfExpression{}, err
	}

	return outputIfExpression, nil
}

func parseArrayExpression(tokens []token.Token) (ast.ArrayExpression, error) {
	outputExpression := ast.ArrayExpression{ElementsExpressions: make([]ast.Node, 0), Type: types.StandardType{Name: types.NONE}}

	arrayContentTokens, isArrayExpression, _ := getParenthesisContent(tokens, 0, token.START_ARRAY, token.END_ARRAY)
	if !isArrayExpression {
		if len(tokens) == 1 {
			return ast.ArrayExpression{}, errors.NewSyntaxErrorInvalidToken(tokens[0].Line, tokens[0].Literal)
		}

		return ast.ArrayExpression{}, errors.NewSyntaxErrorUnexpectedToken(tokens[len(tokens)-1].Line, tokens[len(tokens)-1].Literal, token.END_ARRAY)
	}

	tokensInExpressions := splitTokenSliceByComma(arrayContentTokens)
	for i := 0; i < len(tokensInExpressions); i++ {
		curExpression, err := parseExpression(tokensInExpressions[i])
		if err != nil {
			return outputExpression, err
		}

		outputExpression.ElementsExpressions = append(outputExpression.ElementsExpressions, curExpression)
	}

	return outputExpression, nil
}

//Returns -1 if no operator is not found
func findLeftmostTokenOfType(tokensToFind []string, tokens []token.Token, stopAtIfAndFunctionExecution bool) (string, int, error) {
	parenthesisDepth := 0
	functionDepth := 0
	arrayDepth := 0
	leftMostToken := ""
	leftMostTokenPos := -1

	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == token.START_BLOCK {
			functionDepth++
			continue
		}

		if tokens[i].Type == token.END_BLOCK {
			functionDepth--
			continue
		}

		if tokens[i].Type == token.LEFT_PARENTHESIS {
			parenthesisDepth++
			continue
		}

		if tokens[i].Type == token.RIGHT_PARENTHESIS {
			parenthesisDepth--
			continue
		}

		if tokens[i].Type == token.START_ARRAY {
			arrayDepth++
			continue
		}

		if tokens[i].Type == token.END_ARRAY {
			arrayDepth--
			continue
		}

		if parenthesisDepth != 0 || functionDepth != 0 || arrayDepth != 0 {
			continue
		}

		if stopAtIfAndFunctionExecution {
			if tokens[i].Type == token.EXECUTE_FUNCTION {
				break
			}

			if tokens[i].Type == token.IF {
				break
			}
		}

		for j := 0; j < len(tokensToFind); j++ {
			if tokensToFind[j] == tokens[i].Type {
				leftMostToken = tokensToFind[j]
				leftMostTokenPos = i
			}
		}
	}

	if parenthesisDepth != 0 {
		return "", -1, errors.NewGeneralError(tokens[0].Line, "Amount of left and right parenthesizes in expression does not match")
	}

	if functionDepth != 0 {
		return "", -1, errors.NewGeneralError(tokens[0].Line, "Amount of { and } in expression does not match")
	}

	if arrayDepth != 0 {
		return "", -1, errors.NewGeneralError(tokens[0].Line, "Amount of [ and ] in expression does not match")
	}

	return leftMostToken, leftMostTokenPos, nil
}

func createOperatorExpression(tokens []token.Token, operator string, operatorPos int) (ast.OperatorExpression, error) {
	leftExpression, err := parseExpression(tokens[0:operatorPos])
	if err != nil {
		return ast.OperatorExpression{}, err
	}

	rightExpression, err := parseExpression(tokens[operatorPos+1:])
	if err != nil {
		return ast.OperatorExpression{}, err
	}

	return ast.OperatorExpression{
		LeftSide:  leftExpression,
		RightSide: rightExpression,
		Operator:  operator,
	}, nil
}

func parseParenthesisExpression(tokens []token.Token) (ast.Node, error) {
	if tokens[len(tokens)-1].Type != token.RIGHT_PARENTHESIS {
		return ast.IntExpression{}, errors.NewGeneralError(tokens[0].Line, "Expected ) before end of expression")
	}

	return parseExpression(tokens[1 : len(tokens)-1])
}

// v = 5 + 3 ))
// 5 + 6 ) + 3
