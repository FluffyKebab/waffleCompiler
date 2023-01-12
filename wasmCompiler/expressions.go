package wasmCompiler

import (
	"compiler/ast"
	"compiler/leb128"
	"compiler/token"
	"compiler/types"
	"compiler/wasmCompiler/code"
	"fmt"
	"reflect"
)

func (c *compiler) compileExpression(expression ast.Node, functionLocals *functionLocals) ([]uint8, error) {
	byteCode := make([]uint8, 0)

	switch s := expression.(type) {
	case ast.OperatorExpression:
		left, err := c.compileExpression(s.LeftSide, functionLocals)
		if err != nil {
			return []uint8{}, err
		}

		right, err := c.compileExpression(s.RightSide, functionLocals)
		if err != nil {
			return []uint8{}, err
		}

		byteCode = append(byteCode, left...)
		byteCode = append(byteCode, right...)

		operatorCodeIndex, err := getOperatorCode(s.Operator, s.Type.String())
		if err != nil {
			return []byte{}, err
		}

		byteCode = append(byteCode, operatorCodeIndex)

	case ast.ExecuteFunctionExpression:
		for i := 0; i < len(s.Arguments); i++ {
			argumentBytecode, err := c.compileExpression(s.Arguments[i], functionLocals)
			if err != nil {
				return []uint8{}, err
			}

			byteCode = append(byteCode, argumentBytecode...)
		}

		tableIndex := -1
		functionTypeIndex := -1
		isGlobal := true
		var err error

		switch f := s.Function.(type) {
		case ast.DefineFunctionExpression:
			tableIndex, functionTypeIndex, err = c.addLocalFunction(f)
			if err != nil {
				return []uint8{}, err
			}
			isGlobal = false

		case ast.Variable:
			variableSymbol, isDefined, symbolIsGlobal := c.symbolController.Resolve(f.Identifier)
			if !isDefined {
				if isStandardFunction(f.Identifier) {
					argumentTypes := make([]types.Type, 0)
					for i := 0; i < len(s.Arguments); i++ {
						argumentTypes = append(argumentTypes, s.Arguments[i].GetExpressionReturnType()...)
					}

					tableIndex, functionTypeIndex, err = c.getStandardFunctionIndexAndTypeIndex(f.Identifier, argumentTypes)
					break
				}

				return []uint8{}, fmt.Errorf("undefined identifier")
			}

			isGlobal = symbolIsGlobal

			symbolType, isFunction := variableSymbol.Type.(types.FunctionType)
			if !isFunction {
				return []uint8{}, fmt.Errorf("Internal compiler error: type of variable in function given to compile expression not of type function")
			}

			functionTypeIndex = symbolType.TypeIndex
			tableIndex = int(variableSymbol.Index)

		default:
			return []uint8{}, fmt.Errorf("Internal compiler error: expression of type execute function does not operate of function")
		}

		if isGlobal {
			byteCode = append(byteCode, code.I32_CONST)
		} else {
			byteCode = append(byteCode, code.LOCAL_GET)
		}

		byteCode = append(byteCode, leb128.Int32ToULEB128(int32(tableIndex))...)
		byteCode = append(byteCode, callIndirect(functionTypeIndex)...)

	case ast.DefineFunctionExpression:
		functionIndex, _, err := c.addLocalFunction(s)
		if err != nil {
			return []uint8{}, err
		}

		byteCode = append(byteCode, code.I32_CONST)
		byteCode = append(byteCode, leb128.Int32ToULEB128(int32(functionIndex))...)

	case ast.Variable:
		variableSymbol, isDefined, isGlobal := c.symbolController.Resolve(s.Identifier)
		if !isDefined {
			return []uint8{}, fmt.Errorf("undefined identifier")
		}

		if _, isFunction := variableSymbol.Type.(types.FunctionType); isFunction {
			byteCode = append(byteCode, code.I32_CONST)
			byteCode = append(byteCode, leb128.Int32ToULEB128(int32(variableSymbol.Index))...)
			break
		}

		if isGlobal {
			byteCode = append(byteCode, code.GLOBAL_GET)
		} else {
			byteCode = append(byteCode, code.LOCAL_GET)
		}

		byteCode = append(byteCode, leb128.Int32ToULEB128(int32(variableSymbol.Index))...)

	case ast.ArrayExpression:
		arrayType_ := s.Type
		arrayType, ok := arrayType_.(types.ArrayType)
		if !ok {
			return []byte{}, fmt.Errorf("Type in ast.ArrayExpression not array")
		}

		expressionCode, err := c.createArrayCode(arrayType.ElementType, s.ElementsExpressions, functionLocals)
		if err != nil {
			return []byte{}, err
		}

		byteCode = append(byteCode, expressionCode...)

	case ast.IfExpression:
		conditionalCode, err := c.compileExpression(s.Condition, functionLocals)
		if err != nil {
			return []byte{}, err
		}

		trueExpressionCode, err := c.compileExpression(s.TrueExpression, functionLocals)
		if err != nil {
			return []byte{}, err
		}

		falseExpressionCode, err := c.compileExpression(s.FalseExpression, functionLocals)
		if err != nil {
			return []byte{}, err
		}

		ifTypeIndex := c.typeSection.addType(types.FunctionType{
			ArgumentTypes: []types.Type{},
			ReturnTypes:   s.ReturnType,
		})

		expressionCode := make([]uint8, 0)
		expressionCode = append(expressionCode, conditionalCode...)
		expressionCode = append(expressionCode, code.IF)
		expressionCode = append(expressionCode, leb128.Int32ToULEB128(int32(ifTypeIndex))...)
		expressionCode = append(expressionCode, trueExpressionCode...)
		expressionCode = append(expressionCode, code.ELSE)
		expressionCode = append(expressionCode, falseExpressionCode...)
		expressionCode = append(expressionCode, code.END)

		return expressionCode, nil

	case ast.IntExpression:
		byteCode = append(byteCode, code.I32_CONST)
		byteCode = append(byteCode, leb128.Int32ToLEB128(s.Value)...)
	case ast.FloatExpression:
		byteCode = append(byteCode, code.F32_CONST)
		byteCode = append(byteCode, float32ToLittleEndian(float32(s.Value))...)
	case ast.BoolExpression:
		byteCode = append(byteCode, code.I32_CONST)
		if s.Value {
			byteCode = append(byteCode, leb128.Int32ToULEB128(1)...)
		} else {
			byteCode = append(byteCode, leb128.Int32ToULEB128(0)...)
		}
	case ast.StringExpression:
		//TODO :)
	default:
		fmt.Println("Def", reflect.TypeOf(expression))
	}

	return byteCode, nil
}

func callIndirect(functionTypeIndex int) []byte {
	byteCode := []byte{code.CALL_INDIRECT}
	byteCode = append(byteCode, leb128.Int32ToULEB128(int32(functionTypeIndex))...)
	byteCode = append(byteCode, leb128.Int32ToULEB128(int32(0))...)

	return byteCode
}

func getOperatorCode(operatorType, argumentsType string) (byte, error) {
	if !(argumentsType == token.INT || argumentsType == token.FLOAT || argumentsType == token.BOOL) {
		return 0, fmt.Errorf("Type %s not supported in getOperatorCode ", operatorType)
	}

	switch operatorType {
	case token.PLUS:
		if argumentsType == token.INT {
			return code.I32_ADD, nil
		} else {
			return code.F32_ADD, nil
		}

	case token.MINUS:
		if argumentsType == token.INT {
			return code.I32_SUB, nil
		} else {
			return code.F32_SUB, nil
		}

	case token.MULT:
		if argumentsType == token.INT {
			return code.I32_MUL, nil
		} else {
			return code.F32_MUL, nil
		}

	case token.DIV:
		if argumentsType == token.INT {
			return code.I32_DIV_S, nil
		} else {
			return code.F32_DIV, nil
		}

	case token.EQUAL:
		if argumentsType == token.INT {
			return code.I32_EQ, nil
		} else {
			return code.F32_EQ, nil
		}

	case token.GREATER_THEN:
		if argumentsType == token.INT {
			return code.I32_GT_S, nil
		} else {
			return code.F32_GT, nil
		}

	case token.LESS_THEN:
		if argumentsType == token.INT {
			return code.I32_LT_S, nil
		} else {
			return code.F32_LT, nil
		}

	case token.EQUAL_OR_LESS_THEN:
		if argumentsType == token.INT {
			return code.I32_LE_S, nil
		} else {
			return code.F32_LE, nil
		}

	case token.EQUAL_OR_GREATER_THEN:
		if argumentsType == token.INT {
			return code.I32_GE_S, nil
		} else {
			return code.F32_GE, nil
		}

	case token.AND:
		return code.I32_AND, nil

	case token.OR:
		return code.I32_OR, nil
	}

	return 0, fmt.Errorf("unknown operator %s", operatorType)
}
