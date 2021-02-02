package main

import (
	"os"
	"fmt"
	"strconv"
	"math"
	"math/rand"
	"github.com/jinzhu/copier"

	cxcore "github.com/skycoin/cx/cx"
	"github.com/skycoin/cx/cxgo/actions"
	cxgo "github.com/skycoin/cx/cxgo/cxgo"
	// "github.com/qdm12/reprint"
	// "github.com/barkimedes/go-deepcopy"
	// "github.com/getlantern/deepcopy"
	
	// "github.com/SkycoinProject/skycoin/src/cipher/encoder"
)

// Debug ...
func Debug(args ...interface{}) {
	fmt.Println(args...)
}

// TODOs:
// Only use one variable name for `solution`, `solPrototype`, `sol` when it's of type *cxcore.CXFunction.
// Only use one variable name for `solution` when it's of type [][]string.

// adaptMainFn removes the main function from the main
// package. Then it creates a new main function that will contain a call
// to the solution function.
func adaptSolution(prgrm *cxcore.CXProgram, solution []string) {
	solutionName := solution[len(solution)-1]
	// Ensuring that main pkg exists.
	var mainPkg *cxcore.CXPackage
	mainPkg, err := prgrm.GetPackage(cxcore.MAIN_PKG)
	if err != nil {
		panic(err)
	}
	
	// mainFn, err := mainPkg.GetFunction(cxcore.MAIN_FUNC)
	// if err != nil {
	// 	panic(err)
	// }

	mainFn := cxcore.MakeFunction(cxcore.MAIN_FUNC, "", -1)
	mainFn.Package = mainPkg
	for i, fn := range mainPkg.Functions {
		if fn.Name == cxcore.MAIN_FUNC {
			mainPkg.Functions[i] = mainFn
			break
		}
	}

	mainFn.Expressions = nil
	mainFn.Inputs = nil
	mainFn.Outputs = nil
	
	// idx := -1
	// for i, fn := range pkg.Functions {
	// 	if fn.Name == cxcore.MAIN_FUNC {
	// 		idx = i
	// 		break
	// 	}
	// }
	// _ = idx
	// // Removing main function.
	// pkg.Functions = append(pkg.Functions[:idx], pkg.Functions[idx+1:]...)

	// var solPkg *cxcore.CXPackage
	// solPkg, err = prgrm.GetCurrentPackage()
	// if err != nil {
	// 	panic(err)
	// }
	// pkg, err = prgrm.GetPackage(cxcore.MAIN_PKG)
	// if err != nil {
	// 	panic(err)
	// }
	// mainFn := cxcore.MakeFunction(cxcore.MAIN_FUNC, "", -1)
	// pkg.AddFunction(mainFn)

	var sol *cxcore.CXFunction
	sol, err = mainPkg.GetFunction(solutionName)
	if err != nil {
		panic(err)
	}

	// The size of main function will depend on the number of inputs and outputs.
	mainSize := 0
	
	// Adding inputs to call to solution in main function.
	for _, inp := range sol.Inputs {
		mainFn.AddInput(inp)
		mainSize += inp.TotalSize
	}

	// Adding outputs to call to solution in main function.
	for _, out := range sol.Outputs {
		mainFn.AddInput(out)
		mainSize += out.TotalSize
	}

	// mainInp := MakeArgument("inp", "", -1).AddType("f64")
	// mainInp.Package = mainPkg
	// mainOut := MakeArgument("out", "", -1).AddType("f64")
	// mainOut.Package = mainPkg
	// mainOut.Offset += mainInp.TotalSize
	// mainFn.AddInput(mainInp)
	// mainFn.AddOutput(mainOut)

	

	// // We need to replace the solution function, as it is a pointer.
	// /// It'd maintain a reference to it.
	// newSol := cxcore.MakeFunction(solutionName, sol.FileName, sol.FileLine)
	// // Inputs and outputs can be the same. We only need their offsets.
	// newSol.Inputs = sol.Inputs
	// newSol.Outputs = sol.Outputs
	// newSol.Package = solPkg

	// // We'll need to replace the pointer, not only the object.
	/// So we need the solution index to replace it from solPkg.Functions.
	// solFnIdx := -1
	// for i, fn := range solPkg.Functions {
	// 	if fn.Name == solutionName {
	// 		solFnIdx = i
	// 		break
	// 	}
	// }
	// solPkg.Functions[solFnIdx] = newSol

	expr := cxcore.MakeExpression(sol, "", -1)
	expr.Package = mainPkg
	// expr.AddOutput(mainOut)
	// expr.AddInput(mainInp)

	// Adding inputs to expression which calls solution.
	for _, inp := range sol.Inputs {
		expr.AddInput(inp)
	}

	// Adding outputs to expression which calls solution.
	for _, out := range sol.Outputs {
		expr.AddOutput(out)
	}

	// prnt := cxcore.MakeExpression(cxcore.Natives[cxcore.OpCodes["f64.print"]], "", -1)
	// prnt.Package = mainPkg
	// prnt.AddInput(sol.Outputs[len(sol.Outputs)-1])

	mainFn.AddExpression(expr)
	// mainFn.AddExpression(prnt)
	mainFn.Length = 1
	mainFn.Size = mainSize
}

func getFnBag(prgrm *cxcore.CXProgram, fnBag []string) (fns []*cxcore.CXFunction) {
	pkgName := ""
	fnName := ""
	for i, name := range fnBag {
		if name == "pkg" {
			pkgName = fnBag[i+1]
		}
		if name == "fn" {
			fnName = fnBag[i+1]
		}
		if pkgName != "" && fnName != "" {
			// Then it's a standard library function, like i32.add.
			var fn *cxcore.CXFunction
			if pkgName == cxcore.STDLIB_PKG {
				fn = cxcore.Natives[cxcore.OpCodes[fnName]]
				if fn == nil {
					panic("standard library function not found.")
				}
			} else {
				var err error
				fn, err = prgrm.GetFunction(fnName, pkgName)
				if err != nil {
					panic(err)
				}
			}
			
			fns = append(fns, fn)
			pkgName = ""
			fnName = ""
		}
	}
	return fns
}

func getRandFn(fnBag []*cxcore.CXFunction) *cxcore.CXFunction {
	return fnBag[rand.Intn(len(fnBag))]
}

func getFnArgs(fn *cxcore.CXFunction) (args []*cxcore.CXArgument) {
	for _, arg := range fn.Inputs {
		args = append(args, arg)
	}

	for _, arg := range fn.Outputs {
		args = append(args, arg)
	}

	for _, expr := range fn.Expressions {
		for _, arg := range expr.Inputs {
			args = append(args, arg)
		}

		for _, arg := range expr.Outputs {
			args = append(args, arg)
		}
	}

	return args
}

func calcFnSize(fn *cxcore.CXFunction) (size int) {
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

func getRandInp(fn *cxcore.CXFunction) *cxcore.CXArgument {
	var arg cxcore.CXArgument
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
		arg.Offset += fn.Inputs[c].TotalSize
	}
	for c := 0; c < len(fn.Outputs); c++ {
		arg.Offset += fn.Outputs[c].TotalSize
	}
	for c := 0; c < rndExprIdx; c++ {
		// TODO: We're only considering one output per operator.
		/// Not because of practicality, but because multiple returns in CX are currently buggy anyway.
		arg.Offset += fn.Expressions[c].Operator.Outputs[0].TotalSize
	}

	arg.Package = fn.Package
	arg.Name = strconv.Itoa(rndExprIdx)
	return &arg
}

func getRandOut(fn *cxcore.CXFunction) *cxcore.CXArgument {
	var arg cxcore.CXArgument
	rndExprIdx := rand.Intn(len(fn.Expressions))
	// Making a copy of the argument
	err := copier.Copy(&arg, fn.Expressions[rndExprIdx].Operator.Outputs[0])
	if err != nil {
		panic(err)
	}
	// Determining the offset where the expression should be writing to.
	for c := 0; c < len(fn.Inputs); c++ {
		arg.Offset += fn.Inputs[c].TotalSize
	}
	for c := 0; c < len(fn.Outputs); c++ {
		arg.Offset += fn.Outputs[c].TotalSize
	}
	for c := 0; c < rndExprIdx; c++ {
		// TODO: We're only considering one output per operator.
		/// Not because of practicality, but because multiple returns in CX are currently buggy anyway.
		arg.Offset += fn.Expressions[c].Operator.Outputs[0].TotalSize
	}

	arg.Package = fn.Package
	arg.Name = strconv.Itoa(rndExprIdx)
	return &arg
}

// const (
// 	mirrorMutate = iota() // 0
// 	cocoMutate // 1
// 	cucuMutate // 2
// )

// mirrorMutate swaps a gene (*CXExpression) from fn.Expressions (our genome) in a mirror-like manner.
func mirrorMutate(fn *cxcore.CXFunction) {
	randIdx := rand.Intn(len(fn.Expressions))
	tmpExpr := fn.Expressions[randIdx]
	mirrorIdx := len(fn.Expressions) - randIdx - 1
	fn.Expressions[randIdx] = fn.Expressions[mirrorIdx]
	fn.Expressions[mirrorIdx] = tmpExpr
}

func randomMutate(pop []*cxcore.CXProgram, affSol []string, sPrgrm []byte, fns []*cxcore.CXFunction, numExprs int) {
	randIdx := rand.Intn(len(pop))
	pop[randIdx] = cxcore.Deserialize(sPrgrm)
	initSolution(pop[randIdx], affSol, fns, numExprs)
	adaptSolution(pop[randIdx], affSol)
	resetPrgrm(pop[randIdx])
}

func bitflipMutate(fn *cxcore.CXFunction, fnBag []*cxcore.CXFunction) {
	rndExprIdx := rand.Intn(len(fn.Expressions))
	rndFn := getRandFn(fnBag)

	expr := cxcore.MakeExpression(rndFn, "", -1)
	expr.Package = fn.Package
	expr.Inputs = fn.Expressions[rndExprIdx].Inputs
	expr.Outputs = fn.Expressions[rndExprIdx].Outputs

	exprs := make([]*cxcore.CXExpression, len(fn.Expressions))
	for i, ex := range fn.Expressions {
		if i == rndExprIdx {
			exprs[i] = expr
		} else {
			exprs[i] = ex
		}
	}

	// fn.Expressions[rndExprIdx] = expr
	fn.Expressions = exprs
}

func crossover(parent1, parent2 *cxcore.CXFunction) (*cxcore.CXFunction, *cxcore.CXFunction) {
	var child1, child2 cxcore.CXFunction

	cutPoint := rand.Intn(len(parent1.Expressions))

	err := copier.Copy(&child1, *parent1)
	if err != nil {
		panic(err)
	}
	// reprint.FromTo(parent1, &child1)

	// Replacing reference to slice.
	child1.Expressions = make([]*cxcore.CXExpression, len(child1.Expressions))

	// It's okay to keep the same references to expressions, though.
	// We only want to be handling a different slice of `*CXExpression`s.
	for i, expr := range parent1.Expressions {
		child1.Expressions[i] = expr
	}

	err = copier.Copy(&child2, *parent2)
	if err != nil {
		panic(err)
	}
	// reprint.FromTo(parent2, &child2)

	// Replacing expressions as we did for `child1`.
	child2.Expressions = make([]*cxcore.CXExpression, len(child2.Expressions))

	for i, expr := range parent2.Expressions {
		child2.Expressions[i] = expr
	}

	for c := 0; c < cutPoint; c++ {
		child1.Expressions[c] = parent2.Expressions[c]
	}

	for c := 0; c < cutPoint; c++ {
		child2.Expressions[c] = parent1.Expressions[c]
	}

	return &child1, &child2
}

func initSolution(prgrm *cxcore.CXProgram, solution []string, fns []*cxcore.CXFunction, numExprs int) {
	solutionName := solution[len(solution)-1]

	pkg, err := prgrm.GetPackage(cxcore.MAIN_PKG)
	if err != nil {
		panic(err)
	}

	var newPkg cxcore.CXPackage
	copier.Copy(&newPkg, *pkg)
	pkgs := make([]*cxcore.CXPackage, len(prgrm.Packages))
	for i, _ := range pkgs {
		pkgs[i] = prgrm.Packages[i]
	}
	prgrm.Packages = pkgs

	for i, pkg := range prgrm.Packages {
		if pkg.Name == cxcore.MAIN_PKG {
			prgrm.Packages[i] = &newPkg
			break
		}
	}
	
	fn, err := prgrm.GetFunction(solutionName, cxcore.MAIN_PKG)
	if err != nil {
		panic(err)
	}

	var newFn cxcore.CXFunction
	newFn.Name = fn.Name
	newFn.Inputs = fn.Inputs
	newFn.Outputs = fn.Outputs
	newFn.Package = fn.Package
	// copier.Copy(&newFn, *fn)

	tmpFns := make([]*cxcore.CXFunction, len(newPkg.Functions))
	for i, _ := range tmpFns {
		tmpFns[i] = newPkg.Functions[i]
	}
	newPkg.Functions = tmpFns

	for i, fn := range newPkg.Functions {
		if fn.Name == solutionName {
			newPkg.Functions[i] = &newFn
			break
		}
	}
	
	preExistingExpressions := len(newFn.Expressions)
	// Checking if we need to add more expressions.
	for i := 0; i < numExprs - preExistingExpressions; i++ {
		op := getRandFn(fns)
		expr := cxcore.MakeExpression(op, "", -1)
		for c := 0; c < len(op.Inputs); c++ {
			expr.Inputs = append(expr.Inputs, getRandInp(&newFn))
		}
		// We need to add the expression at this point, so we
		// can consider this expression's output as a
		// possibility to assign stuff.
		newFn.Expressions = append(newFn.Expressions, expr)
		// Adding last expression, so output must be fn's output.
		if i == numExprs - preExistingExpressions - 1 {
			expr.Outputs = append(expr.Outputs, newFn.Outputs[0])
		} else {
			for c := 0; c < len(op.Outputs); c++ {
				expr.Outputs = append(expr.Outputs, getRandOut(&newFn))
			}
		}
	}
	newFn.Size = calcFnSize(&newFn)
	newFn.Length = numExprs
}

// injectMainInputs injects `inps` at the beginning of `prgrm`'s memory,
// which should always represent the memory sent to the first expression contained
// in `prgrm`'s `main`'s function.
func injectMainInputs(prgrm *cxcore.CXProgram, inps []byte) {
	for i := 0; i < len(inps); i++ {
		prgrm.Memory[i] = inps[i]
	}
}

func extractMainOutputs(prgrm *cxcore.CXProgram, solPrototype *cxcore.CXFunction) [][]byte {
	outputs := make([][]byte, len(solPrototype.Outputs))
	for c := 0; c < len(solPrototype.Outputs); c++ {
		size := solPrototype.Outputs[c].TotalSize
		off := solPrototype.Outputs[0].Offset
		outputs[c] = prgrm.Memory[off:off+size]
	}

	return outputs
}

func resetPrgrm(prgrm *cxcore.CXProgram) {
	// Creating a copy of `prgrm`'s memory.
	mem := make([]byte, len(prgrm.Memory))
	copy(mem, prgrm.Memory)
	// Replacing `prgrm.Memory` with its copy, so individuals don't share the same memory.
	prgrm.Memory = mem
	
	prgrm.CallCounter = 0
	prgrm.StackPointer = 0
	prgrm.CallStack = make([]cxcore.CXCall, cxcore.CALLSTACK_SIZE)
	prgrm.Terminated = false
	// minHeapSize := minHeapSize()
	// prgrm.Memory = make([]byte, STACK_SIZE+minHeapSize)
}

func mae(real, sim []float64) float64 {
	var sum float64
	for c := 0; c < len(real); c++ {
		sum += math.Abs(real[c] - sim[c])
	}
	return sum / float64(len(real))
}

// evalInd evaluates the solution function contained in `ind`.
func evalInd(ind *cxcore.CXProgram, solPrototype *cxcore.CXFunction, inputs [][]byte, outputs [][]byte) float64 {
	var tmp *cxcore.CXProgram
	tmp = cxcore.PROGRAM
	cxcore.PROGRAM = ind

	// TODO: We're calculating the error in here.
	/// Migrate to functions when we have other fitness functions.

	inpFullByteSize := 0
	for c := 0; c < len(solPrototype.Inputs); c++ {
		inpFullByteSize += solPrototype.Inputs[c].TotalSize
	}

	var sum float64

	// `numElts` represents the number of elements per input array calculated by the inputs function.
	// All the inputs represent arrays of the same size, regardless of element type
	// (for example, 10 `i32`s and 10 `f64`s). So it is safe to assume that
	// looping over `inputs[0]` will make us loop over all `inputs` from 1 to N.
	numElts := len(inputs[0]) / solPrototype.Inputs[0].TotalSize
	
	for i := 0; i < numElts; {
		// Now we'll loop over each of the `inputs`.
		/// We want to extract the `i`th element from each of the `inputs`.
		/// For example, if we are sending two arrays (inputs), a [10]i32 and a [10]f64,
		/// we want to extract the `i`th i32 and the `i`th f64 and send those two inputs to the solution.

		// We'll store the `i`th inputs on `inps`.
		inps := make([]byte, inpFullByteSize)
		// `inpsOff` helps us keep track of what byte in `inps` we can write to.
		inpsOff := 0
		
		for c := 0; c < len(inputs); c++ {
			// The size of the input.
			inpSize := solPrototype.Inputs[c].TotalSize
			// The bytes representing the input.
			inp := inputs[c][inpSize*i:inpSize*(i+1)]

			// Copying the input `b`ytes.
			for b := 0; b < len(inp); b++ {
				inps[inpsOff+b] = inp[b]
			}

			// Updating offset.
			inpsOff += inpSize
		}

		// Updating how many `b`ytes we read from `inputs[0]`.
		// b += solPrototype.Inputs[0].TotalSize

		// Injecting the input bytes `inps` to program `ind`.
		injectMainInputs(ind, inps)
		
		// Running program `ind`.
		ind.RunCompiled(0, nil)

		// Extracting outputs processed by `solPrototype`.
		simOuts := extractMainOutputs(ind, solPrototype)

		// Comparing real vs simulated outputs (error).
		for o := 0; o < len(solPrototype.Outputs); o++ {
			outSize := solPrototype.Outputs[o].TotalSize
			for b := 0; b < len(simOuts[o]); b++ {
				// Comparing byte by byte.
				sum += math.Abs(float64(outputs[o][i*outSize+b] - simOuts[o][b]))
			}
		}
		i++
	}

	cxcore.PROGRAM = tmp
	return sum
}

func getLowErrorIdxs(errors []float64, chance float32) (int, int) {
	idx := 0
	secondIdx := 0
	lowest := errors[0]
	for i, err := range errors {
		if err <= lowest && rand.Float32() <= chance {
			lowest = err
			secondIdx = idx
			idx = i
		}
	}
	return idx, secondIdx
}

func getHighErrorIdxs(errors []float64, chance float32) (int, int) {
	idx := 0
	secondIdx := 0
	highest := errors[0]
	for i, err := range errors {
		if err >= highest && rand.Float32() <= chance {
			highest = err
			secondIdx = idx
			idx = i
		}
	}
	return idx, secondIdx
}

func rouletteSelection(errors []float64, isMinimizing bool) int {
	// var prevProb float64
	var errSum float64
	for _, err := range errors {
		errSum += err
	}
	value := rand.Float64() * errSum
	for i, err := range errors {
		value -= err
		if value <= 0 {
			return i
		}
	}

	// _ = prevProb

	return 0
}

func replaceSolution(ind *cxcore.CXProgram, solution []string, sol *cxcore.CXFunction) {
	solutionName := solution[len(solution)-1]
	mainPkg, err := ind.GetPackage(cxcore.MAIN_PKG)
	if err != nil {
		panic(err)
	}
	for i, fn := range mainPkg.Functions {
		if fn.Name == solutionName {
			// mainPkg.Functions[i] = sol
			// We need to replace expression by expression, otherwise we'll
			// end up with duplicated pointers all over the population.
			for j, _ := range sol.Expressions {
				mainPkg.Functions[i].Expressions[j] = sol.Expressions[j]
			}
		}
	}
	mainFn, err := mainPkg.GetFunction(cxcore.MAIN_FUNC)
	if err != nil {
		panic(err)
	}
	mainFn.Expressions[0].Operator = sol
}

// getAffOutput gets the appropriate element from the CX program (a function or an argument),
// evaluates it and then returns a [][]byte that represents the value of either the function that
// was evaluated or the value of a CXArgument, such as a variable. This function also returns
// a Boolean that indicates its caller if it's a function (true) or an argument (false).
func getAffOutput(prgrm *cxcore.CXProgram, sol *cxcore.CXFunction, aff []string, inputs [][]byte) ([][]byte, bool) {
	isFunction := false
	pkgName := ""
	fnName := ""
	argName := ""
	for i, name := range aff {
		if name == "pkg" {
			pkgName = aff[i+1]
		}
		if name == "fn" {
			fnName = aff[i+1]
		}
		if name == "arg" {
			argName = aff[i+1]
		}
	}

	var outputs [][]byte
	if fnName != "" {
		isFunction = true

		fn, err := prgrm.GetFunction(fnName, pkgName)
		if err != nil {
			panic(err)
		}

		if inputs != nil {
			// `numElts` represents the number of elements per input array calculated by the inputs function.
			// All the inputs represent arrays of the same size, regardless of element type
			// (for example, 10 `i32`s and 10 `f64`s). So it is safe to assume that
			// looping over `inputs[0]` will make us loop over all `inputs` 1-N.
			numElts := len(inputs[0]) / sol.Inputs[0].TotalSize

			outputs = make([][]byte, len(fn.Outputs))

			for i := 0; i < numElts; {
				// `inps` will hold the ith slice element of each input in `inputs`.
				inps := make([][]byte, len(inputs))
				for c, _ := range inputs {
					inpSize := sol.Inputs[c].TotalSize
					inps[c] = inputs[c][inpSize*i:inpSize*(i+1)]
				}
				outs := prgrm.Callback(fn, inps)

				for j, out := range outs {
					outputs[j] = append(outputs[j], out...)
				}
				i++
			}
		} else {
			outputs = prgrm.Callback(fn, nil)
		}
	}

	if argName != "" {
		cxcore.Debug("argName", argName)
	}

	return outputs, isFunction
}

func getSolution(prgrm *cxcore.CXProgram, aff []string) *cxcore.CXFunction {
	pkgName := ""
	fnName := ""
	
	for i, name := range aff {
		if name == "pkg" {
			pkgName = aff[i+1]
		}
		if name == "fn" {
			fnName = aff[i+1]
		}
	}
	
	if fnName == "" {
		return nil
	}

	var err error
	var fn *cxcore.CXFunction
	if pkgName == "" {
		// Then user is telling us to look in cxcore.MAIN_PKG.
		fn, err = prgrm.GetFunction(fnName, cxcore.MAIN_PKG)
	} else {
		fn, err = prgrm.GetFunction(fnName, pkgName)
	}

	if err != nil {
		panic(err)
	}
	return fn
}

func mustDeserializeUI32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func mustDeserializeUI64(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func mustDeserializeF32(b []byte) float32 {
	return math.Float32frombits(mustDeserializeUI32(b))
}

func mustDeserializeF64(b []byte) float64 {
	return math.Float64frombits(mustDeserializeUI64(b))
}

func printData(data [][]byte, typ int) {
	switch typ {
	case cxcore.TYPE_F64:
		for _, datum := range data {
			fmt.Printf("%f ", mustDeserializeF64(datum))
		}
	}
	fmt.Printf("\n")
}

func opEvolve(prgrm *cxcore.CXProgram) {
	expr := prgrm.GetExpr()
	fp := prgrm.GetFramePointer()

	inp1, inp2, inp3, inp4, inp5, inp6, inp7, inp8, inp9 := expr.Inputs[0], expr.Inputs[1], expr.Inputs[2], expr.Inputs[3], expr.Inputs[4], expr.Inputs[5], expr.Inputs[6], expr.Inputs[7], expr.Inputs[8]

	affSol := cxcore.GetInferActions(inp1, fp)
	fnBag := cxcore.GetInferActions(inp2, fp)

	affInps := cxcore.GetInferActions(inp3, fp)
	affOuts := cxcore.GetInferActions(inp4, fp)
	eval := cxcore.GetInferActions(inp5, fp)

	// Solution prototype. Used only to get solution's signature and *nothing* more.
	// Why nothing more? Because it's not the cxcore.CXFunction we're *actually* evolving.
	solProt := getSolution(prgrm, affSol)

	inps, isInpsFn := getAffOutput(prgrm, solProt, affInps, nil)
	outs, isOutsFn := getAffOutput(prgrm, solProt, affOuts, inps)

	// cxcore.Debug("inps:")
	// printData(inps, cxcore.TYPE_F64)
	// cxcore.Debug("outs:")
	// printData(outs, cxcore.TYPE_F64)

 	// inps := ReadSliceBytes(fp, inp3, inp3.Type)
	// outs := ReadSliceBytes(fp, inp4, inp4.Type)
	// eval := ReadSliceBytes(fp, inp4, inp4.Type)
	
	numExprs := int(cxcore.ReadI32(fp, inp6))
	numIter := cxcore.ReadI32(fp, inp7)
	numPop := cxcore.ReadI32(fp, inp8)
	eps := cxcore.ReadF64(fp, inp9)

	fns := getFnBag(prgrm, fnBag)

	// TODO: Delete these.
	_ = isInpsFn
	// _ = affOuts
	_ = isOutsFn
	// _ = affSol
	// _ = inps
	// _ = outs
	_ = eval
	// _ = numExprs
	// _ = numIter
	// _ = numPop
	// _ = eps
	// _ = fns

	// Serializing root CX program to create copies of it.
	sPrgrm := cxcore.Serialize(prgrm, 0)
	// Initializing population.
	pop := make([]*cxcore.CXProgram, numPop)
	errors := make([]float64, numPop)

	for i, _ := range pop {
		// err := copier.Copy(&pop[i], prgrm)
		// if err != nil {
		// 	panic(err)
		// }

		// reprint.FromTo(prgrm.Packages, pop[i].Packages)
		// reprint.FromTo(&(prgrm.Memory), &(pop[i].Memory))

		// var copy []byte
		// reprint.FromTo(&sPrgrm, &copy)
		// dsPrgrm := cxcore.Deserialize(copy)
		// pop[i] = *dsPrgrm
		// reprint.FromTo(dsPrgrm, &pop[i])
		pop[i] = cxcore.Deserialize(sPrgrm)

		// pop[i].Packages = make([]*cxcore.CXPackage, len(prgrm.Packages))
		// for c, pkg := range prgrm.Packages {
		// 	reprint.FromTo(&(*pkg), &(*(pop[i].Packages[c])))
		// 	cxcore.Debug("huh")
		// }

		// cxcore.Debug("meow", prgrm.Packages, pop[i].Packages)
		// popCopy, err := deepcopy.Anything(prgrm)
		// if err != nil {
		// 	panic(err)
		// }
		// pop[i] = popCopy.(*cxcore.CXProgram)
		
		// Initialize solution with random expressions.
		initSolution(pop[i], affSol, fns, numExprs)
		adaptSolution(pop[i], affSol)

		resetPrgrm(pop[i])

		// Evaluating solution.
		errors[i] = evalInd(pop[i], solProt, inps, outs)

		// pop[i].PrintProgram()

		// cxcore.Debug("originMem", prgrm.Memory[0:32])
	}

	// fmt.Printf("errors: %v\n", errors)

	// Crossover.
	for c := 0; c < int(numIter); c++ {
		pop1Idx, pop2Idx := getLowErrorIdxs(errors, 0.5)
		dead1Idx, dead2Idx := getHighErrorIdxs(errors, 0.5)
		// pop1Idx := rand.Intn(int(numPop))
		// pop2Idx := rand.Intn(int(numPop))
		// dead1Idx := rand.Intn(int(numPop))
		// dead2Idx := rand.Intn(int(numPop))
		// pop1 := pop[pop1Idx]
		// pop2 := pop[pop2Idx]

		pop1MainPkg, err := pop[pop1Idx].GetPackage(cxcore.MAIN_PKG)
		if err != nil {
			panic(err)
		}
		parent1, err := pop1MainPkg.GetFunction(affSol[len(affSol)-1])
		if err != nil {
			panic(err)
		}

		pop2MainPkg, err := pop[pop2Idx].GetPackage(cxcore.MAIN_PKG)
		if err != nil {
			panic(err)
		}
		parent2, err := pop2MainPkg.GetFunction(affSol[len(affSol)-1])
		if err != nil {
			panic(err)
		}

		child1, child2 := crossover(parent1, parent2)
		// child1, child2 := parent1, parent2

		// child1, child2 := crossover(parent1, parent2)
		
		// if rand.Float32() < 1.1 {
		// 	// mutateFn(child1, fns)
		// 	mirrorMutate(child1)
		// }
		// mirrorMutate(child1)
		// bitflipMutate(child1, fns)
		// if rand.Float32() < 1.1 {
		// 	// mutateFn(child2, fns)
		// 	mirrorMutate(child2)
		// }
		// mirrorMutate(child2)
		// bitflipMutate(child2, fns)

		randomMutate(pop, affSol, sPrgrm, fns, numExprs)

		// cxcore.Debug("parent1")
		// pop[pop1Idx].PrintProgram()
		// cxcore.Debug("parent2")
		// pop[pop2Idx].PrintProgram()

		// cxcore.Debug("cross", pop1Idx, pop2Idx)
		// cxcore.Debug("deads", dead1Idx, dead2Idx)

		replaceSolution(pop[dead1Idx], affSol, child1)
		replaceSolution(pop[dead2Idx], affSol, child2)

		// cxcore.Debug("child1")
		// pop[dead1Idx].PrintProgram()
		// cxcore.Debug("child2")
		// pop[dead2Idx].PrintProgram()

		// panic("stop")

		for i, _ := range pop {
			errors[i] = evalInd(pop[i], solProt, inps, outs)
			if errors[i] <= eps {
				fmt.Printf("Found affSol. Bot #%d", i)
				fmt.Printf("errors: %v\n", errors)
				pop[i].PrintProgram()
				return
			}
		}

		// When all errors are the same, print programs and panic()
		// same := true
		// val := errors[0]
		// for _, error := range errors {
		// 	if error != val {
		// 		same = false
		// 		break
		// 	}
		// }
		// if same {
		// 	cxcore.Debug("Same")
		// 	for i, _ := range pop {
		// 		pop[i].PrintProgram()
		// 	}
		// 	panic("All same")
		// }

		// fmt.Printf("errors: %v\n", errors)
		avg := 0.0
		for _, err := range errors {
			avg += err
		}
		fmt.Printf("avg. error: %v\n", float64(avg) / float64(len(errors)))
	}
}

func main() {
	// Registering this library as a CX library.
	cxcore.RegisterPackage("evolve")
	cxcore.Op(cxcore.GetOpCodeCount(), "evolve.evolve", opEvolve, cxcore.In(cxcore.Slice(cxcore.TYPE_AFF), cxcore.Slice(cxcore.TYPE_AFF), cxcore.Slice(cxcore.TYPE_AFF), cxcore.Slice(cxcore.TYPE_AFF), cxcore.Slice(cxcore.TYPE_AFF), cxcore.AI32, cxcore.AI32, cxcore.AI32, cxcore.AF64), nil)

	// Creating a new CX program.
	actions.PRGRM = cxcore.MakeProgram()

	// Reading flags.
	options := defaultCmdFlags()
	parseFlags(&options, os.Args[1:])
	args := commandLine.Args()

	// Reading source code and parsing it to a valid CX program.
	_, sourceCode, fileNames := cxcore.ParseArgsForCX(args, true)
	cxgo.ParseSourceCode(sourceCode, fileNames)
	cxgo.AddInitFunction(actions.PRGRM)

	actions.PRGRM.RunCompiled(0, nil)
}
