package evolve

import (
	"fmt"
	"math/rand"
	"strconv"

	copier "github.com/jinzhu/copier"
	cxast "github.com/skycoin/cx/cx/ast"
	cxconstants "github.com/skycoin/cx/cx/constants"
)

// Debug just prints its input arguments using `fmt.Println`.
// It's useful for `grep`ing it and deleting all its instances.
func Debug(args ...interface{}) {
	fmt.Println(args...)
}

func getRandFn(fnSet []*cxast.CXFunction) *cxast.CXFunction {
	return fnSet[rand.Intn(len(fnSet))]
}

func calcFnSize(fn *cxast.CXFunction) (size int) {
	for _, arg := range fn.Inputs {
		size += arg.TotalSize
	}
	for _, arg := range fn.Outputs {
		size += arg.TotalSize
	}
	for _, expr := range fn.Expressions {
		// TODO: We're only considering one output per operator.
		/// Not because of practicality, but because multiple returns in CX are currently buggy anyway.
		if len(expr.Operator.Outputs) > 0 {
			size += expr.Operator.Outputs[0].TotalSize
		}
	}

	return size
}

func getRandInp(expr *cxast.CXExpression) *cxast.CXArgument {
	var rndExprIdx int
	var argToCopy *cxast.CXArgument
	var arg cxast.CXArgument

	// Find available arg options.
	optionsFromInputs, optionsFromExpressions := findArgOptions(expr, expr.Operator.Inputs[0].Type)
	lengthOfOptions := len(optionsFromInputs) + len(optionsFromExpressions)

	// if no options available or if operator is jump, add new i32_LT expression.
	if lengthOfOptions == 0 || expr.Operator.OpCode == cxconstants.OP_JMP {
		// TODO: improve process when there's OP_JMP
		return addNewExpression(expr, cxast.OpCodes["i32.lt"])
	}

	rndExprIdx = rand.Intn(lengthOfOptions)
	gotOptionsFromFunctionInputs := rndExprIdx < len(optionsFromInputs)

	if gotOptionsFromFunctionInputs {
		argToCopy = expr.Function.Inputs[optionsFromInputs[rndExprIdx]]
	} else {
		rndExprIdx -= len(optionsFromInputs)
		argToCopy = expr.Function.Expressions[optionsFromExpressions[rndExprIdx]].Operator.Outputs[0]
	}

	// Making a copy of the argument
	err := copier.Copy(&arg, argToCopy)
	if err != nil {
		panic(err)
	}

	if !gotOptionsFromFunctionInputs {
		determineExpressionOffset(&arg, expr, optionsFromExpressions[rndExprIdx])
		arg.ArgDetails.Name = strconv.Itoa(optionsFromExpressions[rndExprIdx])
	}
	arg.ArgDetails.Package = expr.Function.Package

	return &arg
}

func addNewExpression(expr *cxast.CXExpression, expressionType int) *cxast.CXArgument {
	var rndExprIdx int
	var argToAdd *cxast.CXArgument

	exp := cxast.MakeExpression(cxast.Natives[expressionType], "", -1)
	exp.Operator.Name = cxast.OpNames[expressionType]

	// Add expression's inputs
	for i := 0; i < 2; i++ {
		optionsFromInputs, optionsFromExpressions := findArgOptions(expr, exp.Operator.Inputs[0].Type)
		rndExprIdx = rand.Intn(len(optionsFromInputs) + len(optionsFromExpressions))
		if rndExprIdx < len(optionsFromInputs) {
			argToAdd = expr.Function.Inputs[optionsFromInputs[rndExprIdx]]
		} else {
			rndExprIdx -= len(optionsFromInputs)
			argToAdd = expr.Function.Expressions[optionsFromExpressions[rndExprIdx]].Outputs[0]
		}
		exp.AddInput(argToAdd)
	}

	// Add expression's output
	argOutName := strconv.Itoa(len(expr.Function.Expressions))
	argOut := cxast.MakeField(argOutName, cxconstants.TYPE_BOOL, "", -1)
	argOut.AddType(cxconstants.TypeNames[cxconstants.TYPE_BOOL])
	argOut.ArgDetails.Package = expr.Function.Package
	exp.AddOutput(argOut)
	expr.Function.AddExpression(exp)

	determineExpressionOffset(argOut, expr, len(expr.Function.Expressions))

	return argOut
}

func findArgOptions(expr *cxast.CXExpression, argTypeToFind int) ([]int, []int) {
	var optionsFromInputs []int
	var optionsFromExpressions []int

	// loop in inputs
	for i, inp := range expr.Function.Inputs {
		if inp.Type == argTypeToFind && inp.ArgDetails.Name != "" {
			// add index to options from inputs
			optionsFromInputs = append(optionsFromInputs, i)
		}
	}

	// loop in expression outputs
	for i, exp := range expr.Function.Expressions {
		if len(exp.Outputs) > 0 && exp.Outputs[0].Type == argTypeToFind && exp.Outputs[0].ArgDetails.Name != "" {
			// add index to options from inputs
			optionsFromExpressions = append(optionsFromExpressions, i)
		}
	}
	return optionsFromInputs, optionsFromExpressions
}

func getRandOut(expr *cxast.CXExpression) *cxast.CXArgument {
	var arg cxast.CXArgument
	var optionsFromExpressions []int

	for i, exp := range expr.Function.Expressions {
		if len(exp.Operator.Outputs) > 0 && exp.Operator.Outputs[0].Type == expr.Operator.Outputs[0].Type {
			optionsFromExpressions = append(optionsFromExpressions, i)
		}
	}

	rndExprIdx := rand.Intn(len(optionsFromExpressions))
	// Making a copy of the argument
	err := copier.Copy(&arg, expr.Function.Expressions[optionsFromExpressions[rndExprIdx]].Operator.Outputs[0])
	if err != nil {
		panic(err)
	}

	determineExpressionOffset(&arg, expr, optionsFromExpressions[rndExprIdx])
	arg.ArgDetails.Package = expr.Function.Package
	arg.ArgDetails.Name = strconv.Itoa(optionsFromExpressions[rndExprIdx])

	return &arg
}

func determineExpressionOffset(arg *cxast.CXArgument, expr *cxast.CXExpression, indexOfSelectedOption int) {
	// Determining the offset where the expression should be writing to.
	for c := 0; c < len(expr.Function.Inputs); c++ {
		arg.Offset += expr.Function.Inputs[c].TotalSize
	}
	for c := 0; c < len(expr.Function.Outputs); c++ {
		arg.Offset += expr.Function.Outputs[c].TotalSize
	}
	for c := 0; c < indexOfSelectedOption; c++ {
		if len(expr.Function.Expressions[c].Operator.Outputs) > 0 {
			// TODO: We're only considering one output per operator.
			/// Not because of practicality, but because multiple returns in CX are currently buggy anyway.
			arg.Offset += expr.Function.Expressions[c].Operator.Outputs[0].TotalSize
		}
	}
}

func GetFunctionSet(fnNames []string) (fns []*cxast.CXFunction) {
	for _, fnName := range fnNames {
		fn := cxast.Natives[cxast.OpCodes[fnName]]
		if fn == nil {
			panic("standard library function not found.")
		}

		fns = append(fns, fn)
	}
	return fns
}

func GenerateRandomExpressions(inputFn *cxast.CXFunction, inputPkg *cxast.CXPackage, fns []*cxast.CXFunction, numExprs int) {
	preExistingExpressions := len(inputFn.Expressions)
	// Checking if we need to add more expressions.
	for i := 0; i < numExprs-preExistingExpressions; i++ {
		op := getRandFn(fns)
		// Last expression output must be the same as function output.
		if i == (numExprs-preExistingExpressions)-1 && len(op.Outputs) > 0 && len(inputFn.Outputs) > 0 {
			for len(op.Outputs) == 0 || op.Outputs[0].Type != inputFn.Outputs[0].Type {
				op = getRandFn(fns)
			}
		}

		expr := cxast.MakeExpression(op, "", -1)
		expr.Package = inputPkg
		expr.Function = inputFn
		for c := 0; c < len(op.Inputs); c++ {
			expr.Inputs = append(expr.Inputs, getRandInp(expr))
		}

		// if operator is jmp, add then and else lines
		if op.OpCode == cxconstants.OP_JMP {
			lineNumOptions := numExprs - len(expr.Function.Expressions)
			if lineNumOptions < 0 {
				lineNumOptions = (lineNumOptions * -1) - 2
			}
			randThenLineIndex := 0
			if lineNumOptions > 0 {
				randThenLineIndex = rand.Intn(lineNumOptions)
			}

			expr.ThenLines = 1
			expr.ElseLines = randThenLineIndex
		}

		// We need to add the expression at this point, so we
		// can consider this expression's output as a
		// possibility to assign stuff.
		inputFn.Expressions = append(inputFn.Expressions, expr)

		// Adding last expression, so output must be fn's output.
		if i == numExprs-preExistingExpressions-1 {
			expr.Outputs = append(expr.Outputs, inputFn.Outputs[0])
		} else {
			for c := 0; c < len(op.Outputs); c++ {
				expr.Outputs = append(expr.Outputs, getRandOut(expr))
			}
		}
	}
	inputFn.Size = calcFnSize(inputFn)
	inputFn.Length = numExprs
}
