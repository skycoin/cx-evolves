package evolve

import (
	copier "github.com/jinzhu/copier"
	cxcore "github.com/skycoin/cx/cx"
)

// adaptSolution removes the main function from the main
// package. Then it creates a new main function that will contain a call
// to the solution function.
func adaptSolution(prgrm *cxcore.CXProgram, fnToEvolve *cxcore.CXFunction) {
	// Ensuring that main pkg exists.
	var mainPkg *cxcore.CXPackage
	mainPkg, err := prgrm.GetPackage(cxcore.MAIN_PKG)
	if err != nil {
		panic(err)
	}

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

	var sol *cxcore.CXFunction
	sol, err = mainPkg.GetFunction(fnToEvolve.Name)
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

	mainFn.AddExpression(expr)
	mainFn.Length = 1
	mainFn.Size = mainSize
}

func initSolution(prgrm *cxcore.CXProgram, fnToEvolve *cxcore.CXFunction, fns []*cxcore.CXFunction, numExprs int) {
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

	fn, err := prgrm.GetFunction(fnToEvolve.Name, cxcore.MAIN_PKG)
	if err != nil {
		panic(err)
	}

	var newFn cxcore.CXFunction
	newFn.Name = fn.Name
	newFn.Inputs = fn.Inputs
	newFn.Outputs = fn.Outputs
	newFn.Package = fn.Package
	
	solutionName := fn.Name

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
		expr.Package = &newPkg
		expr.Function = &newFn
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

func replaceSolution(ind *cxcore.CXProgram, solutionName string, sol *cxcore.CXFunction) {
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
