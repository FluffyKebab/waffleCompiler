package ast

import (
	"compiler/token"
	"compiler/types"
	"encoding/json"
)

type Node interface {
	node()
	GetExpressionReturnType() []types.Type
	GetChildNodes() []Node
}

type Program struct {
	Body BlockStatement
}

func (p Program) node() {}

func (p Program) GetExpressionReturnType() []types.Type {
	return []types.Type{types.StandardType{}}
}

func NewProgram() Program {
	return Program{
		Body: BlockStatement{make([]Node, 0)},
	}
}

func (p Program) ToString() (string, error) {
	treeJSON, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}

	return string(treeJSON), nil
}

type Statement interface {
	statementNode()
}

type AssignmentStatement struct {
	Variables []Variable
	Value     Node
}

func (s AssignmentStatement) node()          {}
func (s AssignmentStatement) statementNode() {}
func (s AssignmentStatement) GetExpressionReturnType() []types.Type {
	return []types.Type{types.StandardType{}}
}
func (s AssignmentStatement) GetChildNodes() []Node {
	return []Node{s.Value}
}

type FunctionStatement struct {
	Expression Node
}

func (p FunctionStatement) node()          {}
func (s FunctionStatement) statementNode() {}
func (s FunctionStatement) GetExpressionReturnType() []types.Type {
	return []types.Type{types.StandardType{}}
}
func (s FunctionStatement) GetChildNodes() []Node {
	return []Node{s.Expression}
}

type ReturnStatement struct {
	Expressions []Node
}

func (p ReturnStatement) node()          {}
func (s ReturnStatement) statementNode() {}
func (s ReturnStatement) GetExpressionReturnType() []types.Type {
	return []types.Type{types.StandardType{}}
}
func (s ReturnStatement) GetChildNodes() []Node {
	return s.Expressions
}

type BlockStatement struct {
	Statements []Node
}

func (p BlockStatement) node()          {}
func (s BlockStatement) statementNode() {}
func (s BlockStatement) GetExpressionReturnType() []types.Type {
	return []types.Type{types.StandardType{}}
}
func (s BlockStatement) GetChildNodes() []Node {
	return s.Statements
}

type Expression interface {
	expressionNode()
}

type Variable struct {
	Identifier string
	Type       types.Type
}

func (p Variable) node()           {}
func (s Variable) expressionNode() {}
func (s Variable) GetExpressionReturnType() []types.Type {
	return []types.Type{s.Type}
}
func (s Variable) GetChildNodes() []Node {
	return []Node{}
}

type FunctionLiteral struct {
	Arguments    []Variable
	ReturnTypes  []types.Type
	FunctionBody BlockStatement
}

func (p FunctionLiteral) node()           {}
func (s FunctionLiteral) expressionNode() {}
func (s FunctionLiteral) GetExpressionReturnType() []types.Type {
	return s.ReturnTypes
}
func (s FunctionLiteral) GetChildNodes() []Node {
	result := []Node{}
	for i := 0; i < len(s.Arguments); i++ {
		result = append(result, s.Arguments[i])
	}
	return append(result, s.FunctionBody)
}

type ExecuteFunctionExpression struct {
	Function    Node
	Arguments   []Node
	ReturnTypes []types.Type
}

func (p ExecuteFunctionExpression) node()           {}
func (s ExecuteFunctionExpression) expressionNode() {}
func (s ExecuteFunctionExpression) GetExpressionReturnType() []types.Type {
	return s.ReturnTypes
}
func (s ExecuteFunctionExpression) GetChildNodes() []Node {
	result := []Node{s.Function}
	for i := 0; i < len(s.Arguments); i++ {
		result = append(result, s.Arguments[i])
	}
	return result
}

type DefineFunctionExpression struct {
	Arguments              []Variable
	ReturnTypes            []types.Type
	FunctionBody           BlockStatement
	FunctionType           types.FunctionType
	NoReturnTypesSpecified bool
}

func (p DefineFunctionExpression) node()           {}
func (s DefineFunctionExpression) expressionNode() {}
func (s DefineFunctionExpression) GetExpressionReturnType() []types.Type {
	return []types.Type{s.FunctionType}
}
func (s DefineFunctionExpression) GetChildNodes() []Node {
	return []Node{s.FunctionBody}
}

type IfExpression struct {
	Condition       Node
	TrueExpression  Node
	FalseExpression Node
	ReturnType      []types.Type
}

func (p IfExpression) node()           {}
func (s IfExpression) expressionNode() {}
func (s IfExpression) GetExpressionReturnType() []types.Type {
	return s.ReturnType
}
func (s IfExpression) GetChildNodes() []Node {
	return []Node{s.Condition, s.TrueExpression, s.FalseExpression}
}

type IntExpression struct {
	Value int32
}

func (p IntExpression) node()           {}
func (s IntExpression) expressionNode() {}
func (s IntExpression) GetExpressionReturnType() []types.Type {
	return []types.Type{types.StandardType{Name: token.INT}}
}
func (s IntExpression) GetChildNodes() []Node {
	return []Node{}
}

type FloatExpression struct {
	Value float64
}

func (p FloatExpression) node()           {}
func (s FloatExpression) expressionNode() {}
func (s FloatExpression) GetExpressionReturnType() []types.Type {
	return []types.Type{types.StandardType{Name: token.FLOAT}}
}
func (s FloatExpression) GetChildNodes() []Node {
	return []Node{}
}

type StringExpression struct {
	Value string
}

func (p StringExpression) node()           {}
func (s StringExpression) expressionNode() {}
func (s StringExpression) GetExpressionReturnType() []types.Type {
	return []types.Type{types.StandardType{Name: token.STRING}}
}
func (s StringExpression) GetChildNodes() []Node {
	return []Node{}
}

type BoolExpression struct {
	Value bool
}

func (p BoolExpression) node()           {}
func (s BoolExpression) expressionNode() {}
func (s BoolExpression) GetExpressionReturnType() []types.Type {
	return []types.Type{types.StandardType{Name: token.BOOL}}
}
func (s BoolExpression) GetChildNodes() []Node {
	return []Node{}
}

type OperatorExpression struct {
	Operator  string
	Type      types.Type
	LeftSide  Node
	RightSide Node
}

func (p OperatorExpression) node()           {}
func (s OperatorExpression) expressionNode() {}
func (s OperatorExpression) GetExpressionReturnType() []types.Type {
	return []types.Type{s.Type}
}
func (s OperatorExpression) GetChildNodes() []Node {
	return []Node{s.LeftSide, s.RightSide}
}

type ArrayExpression struct {
	Type                types.Type
	ElementsExpressions []Node
}

func (p ArrayExpression) node()           {}
func (s ArrayExpression) expressionNode() {}
func (s ArrayExpression) GetExpressionReturnType() []types.Type {
	return []types.Type{s.Type}
}
func (s ArrayExpression) GetChildNodes() []Node {
	return s.ElementsExpressions
}
