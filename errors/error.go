package errors

import "fmt"

type GeneralError struct {
	line    int
	message string
}

func (e GeneralError) Error() string {
	return fmt.Sprintf("Error on line %d: %s", e.line, e.message)
}

func NewGeneralError(line int, message string) GeneralError {
	return GeneralError{
		line:    line,
		message: message,
	}
}

type SyntaxErrorUnexpectedToken struct {
	line        int
	tokenGotten string
	expected    string
}

func (e SyntaxErrorUnexpectedToken) Error() string {
	return fmt.Sprintf("Error on line %d: unexpected %s, expected %s", e.line, e.tokenGotten, e.expected)
}

func NewSyntaxErrorUnexpectedToken(line int, tokenGotten, expected string) SyntaxErrorUnexpectedToken {
	return SyntaxErrorUnexpectedToken{
		line:        line,
		tokenGotten: tokenGotten,
		expected:    expected,
	}
}

type SyntaxErrorInvalidToken struct {
	line        int
	tokenGotten string
}

func (e SyntaxErrorInvalidToken) Error() string {
	return fmt.Sprintf("Error on line %d: invalid token %s", e.line, e.tokenGotten)
}

func NewSyntaxErrorInvalidToken(line int, tokenGotten string) SyntaxErrorInvalidToken {
	return SyntaxErrorInvalidToken{
		line:        line,
		tokenGotten: tokenGotten,
	}
}

type InternalParserError struct {
	description string
}

func (e InternalParserError) Error() string {
	return fmt.Sprintf("Internal parser error: %s", e.description)
}

func NewInternalParserError(description string) InternalParserError {
	return InternalParserError{
		description: description,
	}
}

const GetParenthesisContentError string = "Error parsing expression: tokens from i given to getParenthesisContent not valid parenthesis"
const ParseFunctionTypeLiteralError string = "Error parsing type literal: tokens from i given to parseFunctionTypeLiteral not as function type"
