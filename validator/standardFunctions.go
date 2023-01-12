package validator

import (
	"compiler/token"
	"compiler/types"
)

var standardFunctions = map[string]types.FunctionType{
	"get": {
		ArgumentTypes: []types.Type{
			types.ArrayType{ElementType: types.AnyType{Name: "a"}},
			types.StandardType{Name: token.INT},
		},
		ReturnTypes: []types.Type{
			types.AnyType{Name: "a"},
		},
	},

	"set": {
		ArgumentTypes: []types.Type{
			types.ArrayType{ElementType: types.AnyType{Name: "a"}},
			types.StandardType{Name: token.INT},
			types.AnyType{Name: "a"},
		},
		ReturnTypes: []types.Type{
			types.ArrayType{ElementType: types.AnyType{Name: "a"}},
		},
	},
}
