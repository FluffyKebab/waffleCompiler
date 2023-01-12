package symbolTable

import (
	"compiler/types"
)

type SymbolController struct {
	NumGlobals    int32
	NumFunctions  int32
	globalScope   *symbolTable
	functionScope *symbolTableStack
}

func NewSymbolController() *SymbolController {
	return &SymbolController{
		globalScope: &symbolTable{
			store:          make(map[string]Symbol),
			numDefinitions: 0,
		},
		functionScope: &symbolTableStack{
			stack:        make([]*symbolTable, 0),
			stackPointer: -1,
		},
	}
}

func (s *SymbolController) DefineAnonymousFunction() int {
	s.NumFunctions++
	return int(s.NumFunctions) - 1
}

//Defines variable in global scope if no symbol tables on function stack and in current local scope if there exists a function symbol table on the function stack
func (s *SymbolController) DefineVariable(variableName string, variableType types.Type) (Symbol, int) {
	currentLocalSymbolTable, isInFunction := s.functionScope.getCur()
	if isInFunction {
		return currentLocalSymbolTable.Define(variableName, variableType, currentLocalSymbolTable.numDefinitions)
	}

	switch variableType.(type) {
	case types.StandardType:
		symbol, _ := s.globalScope.Define(variableName, variableType, s.NumGlobals)
		s.NumGlobals++
		return symbol, int(s.NumGlobals) - 1
	case types.FunctionType:
		symbol, _ := s.globalScope.Define(variableName, variableType, s.NumFunctions)
		s.NumFunctions++
		return symbol, int(s.NumFunctions) - 1
	}

	return Symbol{}, -1
}

func (s *SymbolController) PushFunction(arguments []Variable) {
	s.functionScope.push(newFunctionSymbolTable(arguments))
}

func (s *SymbolController) PopFunction() {
	s.functionScope.pop()
}

func (s *SymbolController) Resolve(variableName string) (symbol Symbol, exists bool, isGlobal bool) {
	functionScope, isInFunction := s.functionScope.getCur()
	if isInFunction {
		symbol, exists := functionScope.Resolve(variableName)
		if exists {
			return symbol, true, false
		}
	}

	symbol, exists = s.globalScope.Resolve(variableName)
	if exists {
		return symbol, true, true
	}

	return Symbol{}, false, false
}

type Symbol struct {
	Name  string
	Type  types.Type
	Index int32
}

type symbolTable struct {
	store          map[string]Symbol
	numDefinitions int32
}

type Variable struct {
	Identifier string
	Type       types.Type
}

func newFunctionSymbolTable(arguments []Variable) *symbolTable {
	table := &symbolTable{
		store:          make(map[string]Symbol),
		numDefinitions: 0,
	}

	for i := 0; i < len(arguments); i++ {
		table.Define(arguments[i].Identifier, arguments[i].Type, table.numDefinitions)
	}

	return table
}

func (s *symbolTable) Define(name string, symbolType types.Type, index int32) (Symbol, int) {
	symbol := Symbol{Name: name, Index: int32(index), Type: symbolType}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol, int(index)
}

func (s *symbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := s.store[name]
	return symbol, ok
}

type symbolTableStack struct {
	stack        []*symbolTable
	stackPointer int
}

func (s *symbolTableStack) getCur() (*symbolTable, bool) {
	if s.stackPointer >= len(s.stack) || s.stackPointer < 0 {
		return &symbolTable{}, false
	}

	return s.stack[s.stackPointer], true
}

func (s *symbolTableStack) push(table *symbolTable) {
	s.stackPointer++
	if s.stackPointer >= len(s.stack) {
		s.stack = append(s.stack, table)
		return
	}

	s.stack[s.stackPointer] = table
}

func (s *symbolTableStack) pop() {
	s.stackPointer--
}
