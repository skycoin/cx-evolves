package mutation

import (
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
	cxast "github.com/skycoin/cx/cx/ast"
	cxastapi "github.com/skycoin/cx/cx/astapi"
)

type MutationHandler func(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument)

var (
	PointMutationOperators map[int]MutationHandler
	MutationOpNames        map[int]string
	MutationOpCodes        map[string]int
)

const (
	MOP_INSERT_RAND_I8_AS__I32_LIT = iota + 1
	MOP_INSERT_RAND_I16_AS_I32_LIT
	MOP_INSERT_RAND__I32_LIT
	MOP_HALF_I32_LIT
	MOP_DOUBLE_I32_LIT
	MOP_ZERO_I32_LIT
	MOP_ADD_ONE_I32_LIT
	MOP_ADD_RAND_I32_LIT
	MOP_SUB_ONE_I32_LIT
	MOP_SUB_RAND_I32_LIT
	MOP_BIT_OR_I32_LIT
	MOP_BIT_AND_I32_LIT
	MOP_BIT_XOR_I32_LIT
	MOP_OR_I32_LIT
	MOP_AND_I32_LIT
	MOP_XOR_I32_LIT
	MOP_BIT_ROTATE_LEFT_I32_LIT
	MOP_BIT_ROTATE_RIGHT_I32_LIT
	MOP_SHIFT_BIT_LEFT_I32_LIT
	MOP_SHIFT_BIT_RIGHT_I32_LIT
)

// RegisterMutationOperator
func RegisterMutationOperator(code int, name string, handler MutationHandler) {
	// Check if duplicate
	if PointMutationOperators[code] != nil {
		panic(fmt.Sprintf("duplicate opcode %d : '%s' width '%s'.\n", code, name, MutationOpNames[code]))
	}

	PointMutationOperators[code] = handler
	MutationOpNames[code] = name
	MutationOpCodes[name] = code
}

// GetCompatiblePositionForOperator returns list of line numbers where the operator can be inserted to.
func GetCompatiblePositionForOperator(cxprogram *cxast.CXProgram, fnName, operatorName string) ([]int, error) {
	var lines []int

	// Get operator function
	operatorFn := cxast.Natives[cxast.OpCodes[operatorName]]
	if operatorFn == nil {
		return []int{}, errors.New("standard library function not found")
	}

	// Check if operatorFn has an output or not.
	hasOutput := false
	if len(operatorFn.Outputs) > 0 {
		hasOutput = true
	}

	fn, err := cxastapi.FindFunction(cxprogram, fnName)
	if err != nil {
		return []int{}, errors.New("function not found")
	}

	for i, expr := range fn.Expressions {
		if !hasOutput {
			lines = append(lines, i)
			continue
		}

		for _, arg := range expr.Inputs {
			if arg.Type == operatorFn.Outputs[0].Type {
				lines = append(lines, i)
				break
			}
		}
	}

	return lines, nil
}

// ReplaceArgInput replaces an expression's input with the argToPut.
func ReplaceArgInput(expr *cxast.CXExpression, argIndex int, argToPut *cxast.CXArgument) error {
	// Check if arg index is valid
	if (len(expr.Inputs)-1) < argIndex || argIndex < 0 {
		return errors.New("invalid arg index")
	}

	// Check if arg type in argIndex is same as arg type of argToPut
	if expr.Inputs[argIndex].Type != argToPut.Type {
		return errors.New("arg types are not the same")
	}

	var arg cxast.CXArgument
	// Making a copy of the argument
	err := copier.Copy(&arg, argToPut)
	if err != nil {
		return err
	}

	expr.Inputs[argIndex] = &arg

	return nil
}
