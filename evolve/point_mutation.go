package evolve

import (
	"errors"

	copier "github.com/jinzhu/copier"
	cxast "github.com/skycoin/cx/cx/ast"
	cxastapi "github.com/skycoin/cx/cx/astapi"
)

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
