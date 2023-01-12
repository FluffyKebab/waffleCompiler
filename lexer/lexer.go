package lexer

import (
	"compiler/token"
	"compiler/types"
	"fmt"
	"unicode"
)

type Lexer struct {
	input        []rune
	readPosition int
	curLine      int
}

func New(input string) Lexer {
	return Lexer{input: []rune(input)}
}

func (l *Lexer) NextToken() token.Token {

	for unicode.IsSpace(l.getCurChar()) || l.readPosition >= len(l.input) {
		if l.readPosition >= len(l.input) {
			return l.newToken(token.EOF, "", l.curLine)
		}

		if l.getCurChar() == '\n' {
			l.curLine++
			return l.newToken(token.NEWLINE, token.NEWLINE, l.curLine-1)
		}

		l.readPosition++
	}

	for _, tokenType := range token.TokenLiterals {
		if l.tokenIs(tokenType) {
			return l.newToken(tokenType, tokenType, l.curLine)
		}
	}

	for _, variableType := range types.ValidTypes {
		if l.tokenIs(variableType) {
			return l.newToken(token.TYPE, variableType, l.curLine)
		}
	}

	for _, boolean := range []string{token.TRUE, token.FALSE} {
		if l.tokenIs(boolean) {
			return l.newToken(token.BOOL, boolean, l.curLine)
		}
	}

	if l.curCharIsInt() {
		return l.readNumber()
	}

	if l.getCurChar() == '"' {
		l.readPosition++
		return l.readString()
	}

	if l.curCharIsValidInVariable() {
		return l.readVariable()
	}

	fmt.Printf("ilegal: %d \n", l.getCurChar())

	return l.newToken(token.ILLEGAL, string(l.getCurChar()), l.curLine)
}

func (l *Lexer) newToken(tokenType, ch string, curLine int) token.Token {
	l.readPosition += len(ch)
	return token.New(tokenType, ch, curLine)
}

func (l *Lexer) tokenIs(token string) bool {
	for i := 0; i < len(token); i++ {
		if i+l.readPosition >= len(l.input) {
			return false
		}
		if []rune(token)[i] != l.input[i+l.readPosition] {
			return false
		}
	}

	return true
}

func (l *Lexer) getCurChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) curCharIsInt() bool {
	return unicode.IsDigit(l.getCurChar())
}

func (l *Lexer) curCharIsValidInVariable() bool {
	char := l.getCurChar()
	return unicode.IsDigit(char) || unicode.IsLetter(char)
}

func (l *Lexer) readString() token.Token {
	l.readPosition++
	stringLiteral := ""

	for {
		curChar := l.getCurChar()
		l.readPosition++

		if curChar == 0 || curChar == '"' { // 0 is end of string
			break
		}

		if curChar == '\n' {
			l.curLine++
		}

		stringLiteral += string(curChar)
	}

	return token.New(token.STRING, stringLiteral, l.curLine)
}

func (l *Lexer) readVariable() token.Token {
	variableLiteral := ""

	for {
		if !l.curCharIsValidInVariable() {
			break
		}

		variableLiteral += string(l.getCurChar())
		l.readPosition++
	}

	return token.New(token.VARIABLE, variableLiteral, l.curLine)
}

func (l *Lexer) readNumber() token.Token {
	numberLiteral := ""
	numberType := token.INT

	for {
		if l.getCurChar() == '.' {
			l.readPosition++
			if numberType == token.FLOAT {
				continue
			}

			numberLiteral += "."
			numberType = token.FLOAT
			continue
		}

		if !unicode.IsDigit(l.getCurChar()) {
			break
		}

		numberLiteral += string(l.getCurChar())
		l.readPosition++
	}

	return token.New(numberType, numberLiteral, l.curLine)
}
