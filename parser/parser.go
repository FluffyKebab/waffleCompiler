package parser

import (
	"compiler/ast"
	"compiler/errors"
	"compiler/lexer"
	"compiler/token"
)

func Parse(input string) (ast.Program, error) {
	lexer := lexer.New(input)
	parser := new(lexer)
	return parser.Parse()
}

type tokenGiver interface {
	nextToken() token.Token
	goBack([]token.Token)
}

type programTokenGiver struct {
	goBackTokens []token.Token
	l            *lexer.Lexer
}

func (tokenGiver programTokenGiver) nextToken() token.Token {
	if len(tokenGiver.goBackTokens) != 0 {
		curToken := tokenGiver.goBackTokens[0]
		if len(tokenGiver.goBackTokens) == 1 {
			tokenGiver.goBackTokens = make([]token.Token, 0)
			return curToken
		}

		tokenGiver.goBackTokens = tokenGiver.goBackTokens[1:]
		return curToken
	}

	return tokenGiver.l.NextToken()
}

func (tokenGiver programTokenGiver) goBack(tokens []token.Token) {
	tokenGiver.goBackTokens = append(tokens, tokenGiver.goBackTokens...)
}

type functionTokenGiver struct {
	tokens []token.Token
	curPos int
}

func (tokenGiver *functionTokenGiver) nextToken() token.Token {
	tokenGiver.curPos++
	if tokenGiver.curPos >= len(tokenGiver.tokens) {
		return token.Token{
			Type:    token.EOF,
			Literal: token.EOF,
			Line:    0,
		}
	}

	return tokenGiver.tokens[tokenGiver.curPos]
}

func (tokenGiver *functionTokenGiver) goBack(tokens []token.Token) {
	tokenGiver.curPos -= len(tokens)
}

func newFunctionTokenGiver(tokens []token.Token) *functionTokenGiver {
	return &functionTokenGiver{
		tokens: tokens,
		curPos: -1,
	}
}

type parser struct {
	l        tokenGiver
	curToken token.Token
}

func (p *parser) NextToken() {
	p.curToken = p.l.nextToken()
}

func new(l lexer.Lexer) *parser {
	return &parser{
		l:        programTokenGiver{l: &l},
		curToken: token.Token{},
	}
}

func newFunctionParser(tokens []token.Token) *parser {
	return &parser{
		l:        newFunctionTokenGiver(tokens),
		curToken: token.Token{},
	}
}

func (p *parser) Parse() (ast.Program, error) {
	ast := ast.NewProgram()

	for true {
		p.NextToken()
		if p.curToken.Type == token.EOF {
			break
		}

		if p.curToken.Type == token.NEWLINE {
			continue
		}

		statements, err := p.parseStatement(ast.Body)
		if err != nil {
			return ast, err
		}

		ast.Body = statements
	}

	return ast, nil
}

func (p *parser) parseStatement(statementParent ast.BlockStatement) (ast.BlockStatement, error) {
	if p.curToken.Type == token.VARIABLE {
		tokensOnLineBeforeStartBlock := p.getTokensBeforeToken([]string{token.ASSIGN_VARIABLE, token.NEWLINE})

		if p.curToken.Type == token.ASSIGN_VARIABLE {
			p.NextToken() //Skip the = token
			statement, err := p.parseAssignmentStatement(tokensOnLineBeforeStartBlock)
			if err != nil {
				return ast.BlockStatement{}, err
			}

			statementParent.Statements = append(statementParent.Statements, statement)
			return statementParent, nil
		}

		return ast.BlockStatement{}, errors.NewSyntaxErrorUnexpectedToken(p.curToken.Line, p.curToken.Literal, "token valid at start of statement")
	}

	if p.curToken.Type == token.RETURN {
		p.NextToken() //skip return token
		expressionsTokens, err := p.GetAllTokensInExpression()
		if err != nil {
			return ast.BlockStatement{}, err
		}

		expressionsTokensSplit := splitTokenSliceByComma(expressionsTokens)
		if err != nil {
			return ast.BlockStatement{}, err
		}

		returnExpressions := make([]ast.Node, 0)

		for i := 0; i < len(expressionsTokensSplit); i++ {
			returnExpression, err := parseExpression(expressionsTokensSplit[i])
			if err != nil {
				return ast.BlockStatement{}, err
			}

			returnExpressions = append(returnExpressions, returnExpression)
		}

		statementParent.Statements = append(statementParent.Statements, ast.ReturnStatement{
			Expressions: returnExpressions,
		})

		return statementParent, nil
	}

	return ast.BlockStatement{}, errors.NewSyntaxErrorUnexpectedToken(p.curToken.Line, p.curToken.Literal, "token valid at start of statement")
}

func (p *parser) getTokensBeforeToken(stopTokenLiterals []string) []token.Token {
	tokens := make([]token.Token, 0)

	for true {
		if p.curToken.Type == token.EOF {
			break
		}

		for i := 0; i < len(stopTokenLiterals); i++ {
			if p.curToken.Literal == stopTokenLiterals[i] {
				return tokens
			}
		}

		tokens = append(tokens, p.curToken)

		p.NextToken()
	}

	return tokens
}

//Etter Ast er generert sett in riktig type til variablene

/*
f10 = !( (x int) -> (int) {
	x = x*5
	x = x +5
	return x
} ) 10
*/
