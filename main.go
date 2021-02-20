package main

import (
	"math/rand"
	"time"
	
	evolve "github.com/skycoin/cx-evolves/evolve"
	cxgo "github.com/skycoin/cx/cxgo/cxgo"
	actions "github.com/skycoin/cx/cxgo/actions"
	cxcore "github.com/skycoin/cx/cx"
	encoder "github.com/skycoin/skycoin/src/cipher/encoder"
)

var expressionsCount = 4
var populationSize = 100
var iterations = 10000
var targetError = 0.1
var functionToEvolve = "polynomialFitting"
var functionSetNames = []string{"f64.add", "f64.mul", "f64.sub", "f64.div", "f64.neg", "f64.neg", "f64.abs", "f64.pow", "f64.cos", "f64.sin", "f64.acos", "f64.asin", "f64.sqrt", "f64.log"}
var crossoverFunction = evolve.CrossoverSinglePoint
var evaluationFunction = evolve.EvaluationPerByte
var inputSignature = []string{"f64", "f64"}
var outputSignature = []string{"f64"}

func InitialProgram() *cxcore.CXProgram {
	// Creating the initial CX program.
	prgrm := cxcore.MakeProgram()
	prgrm.SelectProgram()
	actions.SelectProgram(prgrm)

	// Adding `main` package.
	mainPkg := cxcore.MakePackage(cxcore.MAIN_PKG)
	prgrm.AddPackage(mainPkg)

	// Adding `main` function to `main` package.
	mainFn := cxcore.MakeFunction(cxcore.MAIN_FUNC, "", -1)
	mainFn.Package = mainPkg
	mainPkg.AddFunction(mainFn)

	// Adding function to evolve (`FunctionToEvolve`).
	toEvolveFn := cxcore.MakeFunction(functionToEvolve, "", -1)
	mainPkg.AddFunction(toEvolveFn)

	// Adding input signature to function to evolve (`FunctionToEvolve`).
	for _, inpType := range inputSignature {
		inp := cxcore.MakeArgument(cxcore.MakeGenSym("evo_inp"), "", -1).AddType(inpType)
		inp.AddPackage(mainPkg)
		toEvolveFn.AddInput(inp)
	}

	// Adding output signature to function to evolve (`FunctionToEvolve`).
	for _, outType := range outputSignature {
		out := cxcore.MakeArgument(cxcore.MakeGenSym("evo_out"), "", -1).AddType(outType)
		out.AddPackage(mainPkg)
		toEvolveFn.AddOutput(out)
	}

	cxgo.AddInitFunction(prgrm)

	return prgrm
}

func polynomial(inp1 float64, inp2 float64) float64 {
	return inp1*inp1 + inp2*inp2
}

func polyDataPoints(paramCount, sampleSize int) ([][]byte, [][]byte) {
	inputs := make([][]byte, paramCount)
	outputs := make([][]byte, 1)
	for c := 0; c < paramCount; c++ {
		for i := 0; i < sampleSize; i++ {
			inputs[c] = append(inputs[c], encoder.Serialize(float64(i))...)
		}
	}
	for i := 0; i < sampleSize; i++ {
		outputs[0] = append(outputs[0], encoder.Serialize(polynomial(float64(i), float64(i)))...)
	}

	return inputs, outputs
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	initPrgrm := InitialProgram()
	// initPrgrm.PrintProgram()
	// initPrgrm.RunCompiled(0, nil)

	sampleSize := 100
	paramCount := 2

	inputs, outputs := polyDataPoints(paramCount, sampleSize)

	pop := evolve.MakePopulation(populationSize)
	
	pop.SetIterations(iterations)
	pop.SetExpressionsCount(expressionsCount)
	pop.SetTargetError(targetError)
	pop.SetInputs(inputs)
	pop.SetOutputs(outputs)
	
	pop.InitIndividuals(initPrgrm)
	pop.InitFunctionSet(functionSetNames)
	pop.InitFunctionsToEvolve(functionToEvolve)
	

	pop.Evolve()
	
	// evolve.Evolve(initPrgrm, functionSetNames, functionToEvolve, populationSize, expressionsCount, iterations, targetError, inputs, outputs)
}
