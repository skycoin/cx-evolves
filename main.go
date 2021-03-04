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

// How many expressions a program can have.
var expressionsCount = 4

// How many programs will our population have.
var populationSize = 100

// How many iterations need to pass until the algorithm stops evolving.
var iterations = 10000

// If the algorithm reaches this error, the evolutionary process stops.
var targetError = 0.1

// The name of the function to be evolved in the programs. It can be any name.
var functionToEvolve = "polynomialFitting"

// What functions from CX standard library can we use to create expressions in the programs.
var functionSetNames = []string{"f64.add", "f64.mul", "f64.sub", "f64.div", "f64.neg", "f64.neg", "f64.abs", "f64.pow", "f64.cos", "f64.sin", "f64.acos", "f64.asin", "f64.sqrt", "f64.log"}

// What function (evolve/crossover.go) will we use to perform crossover.
var crossoverFunction = evolve.CrossoverSinglePoint

// What function (evolve/evaluation.go) will we use to evaluate individuals.
var evaluationFunction = evolve.EvaluationPerByte

// What's the input signature of the programs being evolved.
var inputSignature = []string{"f64", "f64"}

// What's the output signature of the programs being evolved.
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

	// Creating an init function for the CX program.
	cxgo.AddInitFunction(prgrm)

	return prgrm
}


// polynomial is used to create a data model for the programs to evolve.
// This can be changed to whatever you want.
func polynomial(inp1 float64, inp2 float64) float64 {
	return inp1*inp1 + inp2*inp2
}

// ployDataPoints uses `polynomial` to create the data model.
// This can be changed to whatever you want. The important thing is to generate
// some data represented by slices of type [][]byte.
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
	// Setting seed so results vary every time we run the example.
	rand.Seed(time.Now().UTC().UnixNano())

	// We create an initial CX program, with a
	initPrgrm := InitialProgram()

	// How big will our data model be (how many data points in the dataset).
	sampleSize := 100
	// How many inputs in the function to be evolved.
	paramCount := 2
	// Generating the datasets.
	inputs, outputs := polyDataPoints(paramCount, sampleSize)

	// Generate a population.
	pop := evolve.MakePopulation(populationSize)

	// Configuring the population. The method calls are self-explanatory.
	pop.SetIterations(iterations)
	pop.SetExpressionsCount(expressionsCount)
	pop.SetTargetError(targetError)
	pop.SetInputs(inputs)
	pop.SetOutputs(outputs)
	
	pop.InitIndividuals(initPrgrm)
	pop.InitFunctionSet(functionSetNames)
	pop.InitFunctionsToEvolve(functionToEvolve)
	
	// Evolving the population. The errors between the real and simulated data will be printed to standard output.
	pop.Evolve()
}
