package evolve

import (
	"fmt"
	"math/rand"
	"strconv"

	copier "github.com/jinzhu/copier"
	cxast "github.com/skycoin/cx/cx/ast"
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

func getRandInp(fn *cxast.CXFunction) *cxast.CXArgument {
	var arg cxast.CXArgument
	// Unlike getRandOut, we need to also consider the function inputs.
	rndExprIdx := rand.Intn(len(fn.Inputs) + len(fn.Expressions))
	// Then we're returning one of fn.Inputs as the input argument.
	if rndExprIdx < len(fn.Inputs) {
		// Making a copy of the operator.
		// Inputs should have already a compiled offset.
		err := copier.Copy(&arg, fn.Inputs[rndExprIdx])
		if err != nil {
			panic(err)
		}
		arg.Package = fn.Package
		return &arg
	}
	// It was not a function input.
	// We need to subtract the number of inputs to rndExprIdx.
	rndExprIdx -= len(fn.Inputs)
	// Making a copy of the argument
	err := copier.Copy(&arg, fn.Expressions[rndExprIdx].Operator.Outputs[0])
	if err != nil {
		panic(err)
	}
	// Determining the offset where the expression should be writing to.
	for c := 0; c < len(fn.Inputs); c++ {
		arg.DataSegmentOffset += fn.Inputs[c].TotalSize
	}
	for c := 0; c < len(fn.Outputs); c++ {
		arg.DataSegmentOffset += fn.Outputs[c].TotalSize
	}
	for c := 0; c < rndExprIdx; c++ {
		// TODO: We're only considering one output per operator.
		/// Not because of practicality, but because multiple returns in CX are currently buggy anyway.
		arg.DataSegmentOffset += fn.Expressions[c].Operator.Outputs[0].TotalSize
	}

	arg.Package = fn.Package
	arg.Name = strconv.Itoa(rndExprIdx)
	return &arg
}

func getRandOut(fn *cxast.CXFunction) *cxast.CXArgument {
	var arg cxast.CXArgument
	rndExprIdx := rand.Intn(len(fn.Expressions))
	// Making a copy of the argument
	err := copier.Copy(&arg, fn.Expressions[rndExprIdx].Operator.Outputs[0])
	if err != nil {
		panic(err)
	}
	// Determining the offset where the expression should be writing to.
	for c := 0; c < len(fn.Inputs); c++ {
		arg.DataSegmentOffset += fn.Inputs[c].TotalSize
	}
	for c := 0; c < len(fn.Outputs); c++ {
		arg.DataSegmentOffset += fn.Outputs[c].TotalSize
	}
	for c := 0; c < rndExprIdx; c++ {
		// TODO: We're only considering one output per operator.
		/// Not because of practicality, but because multiple returns in CX are currently buggy anyway.
		arg.DataSegmentOffset += fn.Expressions[c].Operator.Outputs[0].TotalSize
	}

	arg.Package = fn.Package
	arg.Name = strconv.Itoa(rndExprIdx)
	return &arg
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
