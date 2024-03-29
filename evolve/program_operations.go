package evolve

import (
	copier "github.com/jinzhu/copier"
	cxast "github.com/skycoin/cx/cx/ast"
	cxconstants "github.com/skycoin/cx/cx/constants"
	cxgenerator "github.com/skycoin/cx/cx/generator"
	"github.com/skycoin/cx/cx/types"
)

// adaptSolution removes the main function from the main
// package. Then it creates a new main function that will contain a call
// to the solution function.
func adaptSolution(prgrm *cxast.CXProgram, fnToEvolve *cxast.CXFunction) {
	// Ensuring that main pkg exists.
	var mainPkg *cxast.CXPackage
	mainPkg, err := prgrm.GetPackage(cxconstants.MAIN_PKG)
	if err != nil {
		panic(err)
	}

	mainFn := cxast.MakeFunction(cxconstants.MAIN_FUNC, "", -1)
	mainFn.Package = mainPkg
	for i, fn := range mainPkg.Functions {
		if fn.Name == cxconstants.MAIN_FUNC {
			mainPkg.Functions[i] = mainFn
			break
		}
	}

	mainFn.Expressions = nil
	mainFn.Inputs = nil
	mainFn.Outputs = nil

	var sol *cxast.CXFunction
	sol, err = mainPkg.GetFunction(fnToEvolve.Name)
	if err != nil {
		panic(err)
	}

	// The size of main function will depend on the number of inputs and outputs.
	var mainSize types.Pointer = 0

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

	expr := cxast.MakeExpression(sol, "", -1)
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
	mainFn.LineCount = 1
	mainFn.Size = mainSize
}

func initSolution(prgrm *cxast.CXProgram, fnToEvolve *cxast.CXFunction, fns []*cxast.CXFunction, numExprs int) {
	pkg, err := prgrm.GetPackage(cxconstants.MAIN_PKG)
	if err != nil {
		panic(err)
	}

	var newPkg cxast.CXPackage
	copier.Copy(&newPkg, *pkg)

	prgrm.Packages[cxconstants.MAIN_PKG] = &newPkg

	fn, err := prgrm.GetFunction(fnToEvolve.Name, cxconstants.MAIN_PKG)
	if err != nil {
		panic(err)
	}

	var newFn cxast.CXFunction
	newFn.Name = fn.Name
	newFn.Inputs = fn.Inputs
	newFn.Outputs = fn.Outputs
	newFn.Package = fn.Package

	solutionName := fn.Name
	newPkg.Functions[solutionName] = &newFn

	cxgenerator.GenerateRandomExpressions(&newFn, &newPkg, fns, numExprs)
}

func resetPrgrm(prgrm *cxast.CXProgram) {
	// Creating a copy of `prgrm`'s memory.
	mem := make([]byte, len(prgrm.Memory))
	copy(mem, prgrm.Memory)
	// Replacing `prgrm.Memory` with its copy, so individuals don't share the same memory.
	prgrm.Memory = mem

	prgrm.CallCounter = 0
	prgrm.Stack.Pointer = 0
	prgrm.CallStack = make([]cxast.CXCall, cxconstants.CALLSTACK_SIZE)
	prgrm.Terminated = false
	// minHeapSize := minHeapSize()
	// prgrm.Memory = make([]byte, STACK_SIZE+minHeapSize)
}

func replaceSolution(ind *cxast.CXProgram, solutionName string, sol *cxast.CXFunction) {
	mainPkg, err := ind.GetPackage(cxconstants.MAIN_PKG)
	if err != nil {
		panic(err)
	}
	for i, fn := range mainPkg.Functions {
		if fn.Name == solutionName {
			// mainPkg.Functions[i] = sol
			// We need to replace expression by expression, otherwise we'll
			// end up with duplicated pointers all over the population.

			var replaceRange int
			replaceRange = len(mainPkg.Functions[i].Expressions)
			if len(sol.Expressions) < len(mainPkg.Functions[i].Expressions) {
				replaceRange = len(sol.Expressions)
			}

			for j := 0; j < replaceRange; j++ {
				mainPkg.Functions[i].Expressions[j] = sol.Expressions[j]
			}
		}
	}
	mainFn, err := mainPkg.GetFunction(cxconstants.MAIN_FUNC)
	if err != nil {
		panic(err)
	}
	mainFn.Expressions[0].Operator = sol
}

func ReplaceIndividualWithRandomExpressions(prgrm *cxast.CXProgram, fnToEvolve *cxast.CXFunction, fns []*cxast.CXFunction, numExprs int) {
	initSolution(prgrm, fnToEvolve, fns, numExprs)
	adaptSolution(prgrm, fnToEvolve)
	resetPrgrm(prgrm)
}
