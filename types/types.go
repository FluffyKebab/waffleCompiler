package types

import (
	"compiler/token"
	"compiler/wasmCompiler/code"
	"fmt"
)

type Type interface {
	node()
	String() string
	ByteCode() uint8
}

type StandardType struct {
	Name string
}

func (t StandardType) node() {}

func (t StandardType) String() string {
	return t.Name
}

func (t StandardType) ByteCode() uint8 {
	switch t.Name {
	case token.INT:
		return code.I32
	case token.FLOAT:
		return code.F32
	case token.BOOL:
		return code.I32
	default:
		fmt.Printf("Warning: standard type %s to byte code not defined \n", t.Name)
		return 0
	}
}

type FunctionType struct {
	TypeIndex     int
	ArgumentTypes []Type
	ReturnTypes   []Type
}

func NewFunctionType() FunctionType {
	return FunctionType{
		ArgumentTypes: make([]Type, 0),
		ReturnTypes:   make([]Type, 0),
		TypeIndex:     -1,
	}
}

func (t FunctionType) node() {}

func (t FunctionType) String() string {
	output := "func ("
	for i := 0; i < len(t.ArgumentTypes); i++ {
		output += t.ArgumentTypes[i].String()
		if i+1 != len(t.ArgumentTypes) {
			output += ", "
		}
	}

	output += ") -> ("
	for i := 0; i < len(t.ReturnTypes); i++ {
		output += t.ReturnTypes[i].String()
		if i+1 != len(t.ReturnTypes) {
			output += ", "
		}
	}

	return output + ")"
}

func (t FunctionType) ByteCode() uint8 {
	return code.I32
}

type ArrayType struct {
	ElementType Type
}

func (p ArrayType) node() {}

func (t ArrayType) String() string {
	return "[]" + t.ElementType.String()
}

func (t ArrayType) ByteCode() uint8 {
	return code.I32
}

type AnyType struct {
	Name string
}

func (p AnyType) node() {}

func (t AnyType) String() string {
	return "any" + t.Name
}

func (t AnyType) ByteCode() uint8 {
	fmt.Println("Warning can not turn anyType into byteCode")
	return 0
}

const (
	INT    = "int"
	FLOAT  = "float"
	STRING = "string"
	BOOL   = "bool"
	NONE   = "NONE"
)

var ValidTypes []string = []string{
	INT,
	FLOAT,
	STRING,
	BOOL,
}
