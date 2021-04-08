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

func getFunctionSet(prgrm *cxast.CXProgram, fnNames []string) (fns []*cxast.CXFunction) {
	for _, fnName := range fnNames {
		fn := cxast.Natives[cxast.OpCodes[fnName]]
		if fn == nil {
			panic("standard library function not found.")
		}

		fns = append(fns, fn)
	}
	return fns
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
		return addNewExpression(expr, cxconstants.OP_I32_LT)
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
		arg.Name = strconv.Itoa(optionsFromExpressions[rndExprIdx])
	}
	arg.Package = expr.Function.Package

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
	argOut.Package = expr.Function.Package
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
		if inp.Type == argTypeToFind && inp.Name != "" {
			// add index to options from inputs
			optionsFromInputs = append(optionsFromInputs, i)
		}
	}

	// loop in expression outputs
	for i, exp := range expr.Function.Expressions {
		if len(exp.Outputs) > 0 && exp.Outputs[0].Type == argTypeToFind && exp.Outputs[0].Name != "" {
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
	arg.Package = expr.Function.Package
	arg.Name = strconv.Itoa(optionsFromExpressions[rndExprIdx])

	return &arg
}

func determineExpressionOffset(arg *cxast.CXArgument, expr *cxast.CXExpression, indexOfSelectedOption int) {
	// Determining the offset where the expression should be writing to.
	for c := 0; c < len(expr.Function.Inputs); c++ {
		arg.DataSegmentOffset += expr.Function.Inputs[c].TotalSize
	}
	for c := 0; c < len(expr.Function.Outputs); c++ {
		arg.DataSegmentOffset += expr.Function.Outputs[c].TotalSize
	}
	for c := 0; c < indexOfSelectedOption; c++ {
		if len(expr.Function.Expressions[c].Operator.Outputs) > 0 {
			// TODO: We're only considering one output per operator.
			/// Not because of practicality, but because multiple returns in CX are currently buggy anyway.
			arg.DataSegmentOffset += expr.Function.Expressions[c].Operator.Outputs[0].TotalSize
		}
	}
}

// func printData(data [][]byte, typ int) {
// 	switch typ {
// 	case cxcore.TYPE_F64:
// 		for _, datum := range data {
// 			fmt.Printf("%f ", mustDeserializeF64(datum))
// 		}
// 	}
// 	fmt.Printf("\n")
// }

// func mustDeserializeUI32(b []byte) uint32 {
// 	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
// }

// func mustDeserializeUI64(b []byte) uint64 {
// 	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
// 		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
// }

// func mustDeserializeF32(b []byte) float32 {
// 	return math.Float32frombits(mustDeserializeUI32(b))
// }

// func mustDeserializeF64(b []byte) float64 {
// 	return math.Float64frombits(mustDeserializeUI64(b))
// }
