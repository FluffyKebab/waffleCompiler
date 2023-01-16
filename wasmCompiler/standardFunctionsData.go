package wasmCompiler

import (
	"compiler/token"
	"compiler/types"
)

type standardFunctionDataElement struct {
	fileName  string
	funcIndex int
	funcType  types.FunctionType
	name      string
}

var standardFunctionsData []standardFunctionDataElement = []standardFunctionDataElement{
	{
		name: "allocate",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		fileName:  "./builtInsCode/memoryManagement.wasm",
		funcIndex: 0,
	},
	{
		name: "deAllocate",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{},
		},
		fileName:  "./builtInsCode/memoryManagement.wasm",
		funcIndex: 1,
	},
	{
		name: "array",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		fileName:  "./builtInsCode/memoryManagement.wasm",
		funcIndex: 2,
	},
	{
		name: "take",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		fileName:  "./builtInsCode/memoryManagement.wasm",
		funcIndex: 3,
	},

	{
		name: "i32get",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		fileName:  "./builtInsCode/setterAndGetters.wasm",
		funcIndex: 0,
	},
	{
		name: "i8get",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		fileName:  "./builtInsCode/setterAndGetters.wasm",
		funcIndex: 1,
	},
	{
		name: "f32get",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.FLOAT}},
		},
		fileName:  "./builtInsCode/setterAndGetters.wasm",
		funcIndex: 2,
	},

	{
		name: "i32set",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		fileName:  "./builtInsCode/setterAndGetters.wasm",
		funcIndex: 3,
	},
	{
		name: "i8set",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		fileName:  "./builtInsCode/setterAndGetters.wasm",
		funcIndex: 4,
	},
	{
		name: "f32set",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}, types.StandardType{Name: token.INT}, types.StandardType{Name: token.FLOAT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		fileName:  "./builtInsCode/setterAndGetters.wasm",
		funcIndex: 5,
	},
	{
		name: "length",
		funcType: types.FunctionType{
			ArgumentTypes: []types.Type{types.StandardType{Name: token.INT}},
			ReturnTypes:   []types.Type{types.StandardType{Name: token.INT}},
		},
		fileName:  "./builtInsCode/arrayFunctions.wasm",
		funcIndex: 0,
	},
}

var isOpenStandardFunction = map[string]bool{
	"set":    true,
	"get":    true,
	"length": true,
	"take":   true,
}
