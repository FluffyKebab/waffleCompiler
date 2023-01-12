package wasmCompiler

import (
	"compiler/token"
	"compiler/types"
)

type standardFunction struct {
	fileName  string
	funcIndex int
	funcType  types.FunctionType
	name      string
	open      bool
}

var standardFunctions []standardFunction = []standardFunction{
	{
		name: "allocate",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		open:      false,
		fileName:  "./builtInsCode/memoryManagement.wasm",
		funcIndex: 0,
	},
	{
		name: "deAllocate",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{},
		},
		open:      false,
		fileName:  "./builtInsCode/memoryManagement.wasm",
		funcIndex: 1,
	},
	{
		name: "array",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		open:      false,
		fileName:  "./builtInsCode/memoryManagement.wasm",
		funcIndex: 2,
	},
	{
		name: "checkArraySize",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{},
		},
		open:      false,
		fileName:  "./builtInsCode/memoryManagement.wasm",
		funcIndex: 3,
	},
	{
		name: "i32set",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		open:      false,
		fileName:  "./builtInsCode/memoryManagement.wasm",
		funcIndex: 4,
	},
	{
		name: "i8set",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		open:      false,
		fileName:  "./builtInsCode/memoryManagement.wasm",
		funcIndex: 5,
	},
	{
		name: "f32set",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}, types.StandardType{Name: token.FLOAT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		open:      false,
		fileName:  "./builtInsCode/memoryManagement.wasm",
		funcIndex: 6,
	},
	{
		name: "i32get",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		open:      false,
		fileName:  "./builtInsCode/memoryManagement.wasm",
		funcIndex: 7,
	},
	{
		name: "i8get",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		open:      false,
		fileName:  "./builtInsCode/memoryManagement.wasm",
		funcIndex: 8,
	},
	{
		name: "f32get",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.FLOAT}},
		},
		open:      false,
		fileName:  "./builtInsCode/memoryManagement.wasm",
		funcIndex: 9,
	},
}

var standardFunctionsIndex = map[string]int{
	"allocate":   0,
	"deAllocate": 1,
	"array":      2,
	"i32set":     4,
	"i8set":      5,
	"f32set":     6,
	"i32get":     7,
	"i8get":      8,
	"f32get":     9,
}
