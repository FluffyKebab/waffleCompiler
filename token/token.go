package token

type Token struct {
	Type    string
	Literal string
	Line    int
}

func New(tokenType, literal string, line int) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
		Line:    line,
	}
}

//Token types & shit
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	NONE    = "NONE"
	NEWLINE = "\n"

	VARIABLE   = "VARIABLE"
	TYPE       = "TYPE"
	ARRAY_TYPE = "[]"
	STRING     = "string"
	FLOAT      = "float"
	INT        = "int"
	BOOL       = "bool"
	TRUE       = "true"
	FALSE      = "false"

	RETURN           = "return"
	IF               = "if"
	ELSE             = "else"
	ASSIGN_VARIABLE  = "="
	FUNCTION_ARROW   = "->"
	EXECUTE_FUNCTION = "!"

	COMMA = ","

	PLUS  = "+"
	MINUS = "-"
	MULT  = "*"
	DIV   = "/"

	EQUAL                 = "=="
	NOT_EQUAL             = "!="
	OR                    = "||"
	AND                   = "&&"
	GREATER_THEN          = ">"
	EQUAL_OR_GREATER_THEN = ">="
	LESS_THEN             = "<"
	EQUAL_OR_LESS_THEN    = "<="

	START_BLOCK       = "{"
	END_BLOCK         = "}"
	LEFT_PARENTHESIS  = "("
	RIGHT_PARENTHESIS = ")"
	START_ARRAY       = "["
	END_ARRAY         = "]"
)

var TypeTokens []string = []string{
	ARRAY_TYPE,
	STRING,
	FLOAT,
	INT,
	BOOL,
}

var TokenLiterals []string = []string{
	RETURN,
	IF,
	ELSE,
	NEWLINE,

	COMMA,
	FUNCTION_ARROW,

	PLUS,
	MINUS,
	MULT,
	DIV,

	EQUAL,
	NOT_EQUAL,
	OR,
	AND,
	EQUAL_OR_GREATER_THEN,
	GREATER_THEN,
	EQUAL_OR_LESS_THEN,
	LESS_THEN,

	EXECUTE_FUNCTION,

	START_BLOCK,
	END_BLOCK,
	LEFT_PARENTHESIS,
	RIGHT_PARENTHESIS,
	START_ARRAY,
	END_ARRAY,

	ASSIGN_VARIABLE,
}

var Operators []string = []string{
	PLUS,
	MINUS,
	MULT,
	DIV,
	EQUAL,
	NOT_EQUAL,
	LESS_THEN,
	GREATER_THEN,
	EQUAL_OR_GREATER_THEN,
	EQUAL_OR_LESS_THEN,
}
