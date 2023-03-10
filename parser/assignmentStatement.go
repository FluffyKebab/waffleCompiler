package parser

import (
	"compiler/ast"
	"compiler/errors"
	"compiler/token"
	"compiler/types"
)

// A slice of tokens should be parsed as a assignment statement if the slice contains an = token
func isAssignmentStatement(tokensOnLine []token.Token) bool {
	for i := 0; i < len(tokensOnLine); i++ {
		if tokensOnLine[i].Type == token.ASSIGN_VARIABLE {
			return true
		}
	}
	return false
}

func (p *parser) parseAssignmentStatement(tokensBeforeAssignmentToken []token.Token) (ast.AssignmentStatement, error) {
	statement := ast.AssignmentStatement{}

	variables, err := parseListOftVariables(tokensBeforeAssignmentToken)
	if err != nil {
		return statement, err
	}

	expressionTokens, err := p.GetAllTokensInExpression()
	if err != nil {
		return statement, err
	}

	expression, err := parseExpression(expressionTokens)
	if err != nil {
		return statement, err
	}

	statement.Variables = variables
	statement.Value = expression

	return statement, nil
}

func parseListOftVariables(variableTokens []token.Token) ([]ast.Variable, error) {
	variables := make([]ast.Variable, 0)

	for i := 0; i < len(variableTokens); i++ {
		if variableTokens[i].Type == token.COMMA {
			continue
		}

		if variableTokens[i].Type == token.VARIABLE {
			curVariable := ast.Variable{
				Identifier: variableTokens[i].Literal,
			}

			variableType, indexAfter, isValidType, err := parseTypeLiteral(variableTokens, i+1)
			if err != nil {
				return []ast.Variable{}, err
			}

			if isValidType {
				curVariable.Type = variableType
				i = indexAfter - 1
			} else {
				curVariable.Type = types.StandardType{Name: types.NONE}
			}

			variables = append(variables, curVariable)
			continue
		}

		return variables, errors.NewSyntaxErrorUnexpectedToken(variableTokens[i].Line, variableTokens[i].Type, "identifier")
	}

	return variables, nil
}
