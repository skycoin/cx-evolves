package main

import (
	"math/rand"
	"os"
	"time"

	evolve "github.com/skycoin/cx-evolves/evolve"
	cxcore "github.com/skycoin/cx/cx"
	actions "github.com/skycoin/cx/cxgo/actions"
	cxgo "github.com/skycoin/cx/cxgo/cxgo"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Maze and Output Configuration
var (
	mazeWidth      int
	mazeHeight     int
	epochsCount    int
	plotFitness    bool
	saveAST        bool
	randomMazeSize bool
)

// Evolve Configuration
var (
	expressionsCount int
	populationSize   int
	iterations       int
	functionToEvolve string

	// What functions from CX standard library can we use to create expressions in the programs.
	functionSetNames = []string{"i32.add", "i32.mul", "i32.sub", "i32.neg", "i32.abs", "i32.bitand", "i32.bitor", "i32.bitxor", "i32.bitclear", "i32.bitshl", "i32.bitshr", "i32.max", "i32.min", "i32.rand"}
	// Missing
	// ,"i32.mod"
	// ,"i32.div"
	// ,"i32.gt"
	// ,"i32.gteq"
	// ,"i32.lt"
	// ,"i32.lteq"
	// ,"i32.eq"
	// ,"i32.uneq"

	// If the algorithm reaches this error, the evolutionary process stops.
	// var targetError = 0.1

	// What function (evolve/crossover.go) will we use to perform crossover.
	crossoverFunction = evolve.CrossoverSinglePoint

	// What function (evolve/evaluation.go) will we use to evaluate individuals.
	evaluationFunction = evolve.EvaluationPerByte

	// What's the input signature of the programs being evolved.
	inputSignature = []string{"i32", "i32", "i32", "i32", "i32", "i32", "i32", "i32", "i32", "i32", "i32", "i32", "i32"}

	// What's the output signature of the programs being evolved.
	outputSignature = []string{"i32"}
)

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

// // polynomial is used to create a data model for the programs to evolve.
// // This can be changed to whatever you want.
// func polynomial(inp1 float64, inp2 float64) float64 {
// 	return inp1*inp1 + inp2*inp2
// }

// // ployDataPoints uses `polynomial` to create the data model.
// // This can be changed to whatever you want. The important thing is to generate
// // some data represented by slices of type [][]byte.
// func polyDataPoints(paramCount, sampleSize int) ([][]byte, [][]byte) {
// 	inputs := make([][]byte, paramCount)
// 	outputs := make([][]byte, 1)
// 	for c := 0; c < paramCount; c++ {
// 		for i := 0; i < sampleSize; i++ {
// 			inputs[c] = append(inputs[c], encoder.Serialize(float64(i))...)
// 		}
// 	}
// 	for i := 0; i < sampleSize; i++ {
// 		outputs[0] = append(outputs[0], encoder.Serialize(polynomial(float64(i), float64(i)))...)
// 	}

// 	return inputs, outputs
// }

func main() {
	EvolveApp := &cli.App{
		Name:    "Evolve with Maze",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "width",
				Aliases:     []string{"W"},
				Usage:       "width of the generated maze",
				Value:       2,
				Destination: &mazeWidth,
			},
			&cli.IntFlag{
				Name:        "height",
				Aliases:     []string{"H"},
				Usage:       "height of the generated maze",
				Value:       2,
				Destination: &mazeHeight,
			},
			&cli.IntFlag{
				Name:        "Population Size",
				Aliases:     []string{"population"},
				Usage:       "population size",
				Value:       300,
				Destination: &populationSize,
			},
			&cli.IntFlag{
				Name:        "Generations",
				Aliases:     []string{"generations"},
				Usage:       "Number of generations",
				Value:       500,
				Destination: &iterations,
			},
			&cli.IntFlag{
				Name:        "Expression Count",
				Aliases:     []string{"expressions"},
				Usage:       "Number of expressions a program can have",
				Value:       30,
				Destination: &expressionsCount,
			},
			&cli.IntFlag{
				Name:        "Epochs",
				Aliases:     []string{"epochs"},
				Usage:       "Maze changes for every N generations, if set 0 then maze changes every generations",
				Value:       1,
				Destination: &epochsCount,
			},
			&cli.StringFlag{
				Name:        "Generated Program Name",
				Value:       "MazeRunner",
				Aliases:     []string{"name"},
				Usage:       "Name of program to evolve",
				Destination: &functionToEvolve,
			},
			&cli.BoolFlag{
				Name:        "Random Maze Size",
				Aliases:     []string{"random"},
				Usage:       "set true if generated mazes will be random from NxN 2,3,4,5,6,7, or 8",
				Destination: &randomMazeSize,
			},
			&cli.BoolFlag{
				Name:        "Plot Fitness Graphs",
				Aliases:     []string{"graphs"},
				Usage:       "set true if fitness graphs should be plotted",
				Destination: &plotFitness,
			},
			&cli.BoolFlag{
				Name:        "Save ASTs",
				Aliases:     []string{"ast"},
				Usage:       "set true if best ASTs per generation should be saved to a file",
				Destination: &saveAST,
			},
		},
		Action: func(c *cli.Context) error {
			EvolveWithMaze()
			return nil
		},
	}

	err := EvolveApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func EvolveWithMaze() {
	// Setting seed so results vary every time we run the example.
	rand.Seed(time.Now().UTC().UnixNano())

	// We create an initial CX program, with a
	initPrgrm := InitialProgram()

	// How big will our data model be (how many data points in the dataset).
	// sampleSize := 100
	// How many inputs in the function to be evolved.
	// paramCount := 1
	// Generating the datasets.
	// inputs, outputs := polyDataPoints(paramCount, sampleSize)

	// Generate a population.
	pop := evolve.MakePopulation(populationSize)

	// Configuring the population. The method calls are self-explanatory.
	pop.SetIterations(iterations)
	pop.SetExpressionsCount(expressionsCount)
	// pop.SetTargetError(targetError)
	// pop.SetInputs(inputs)
	// pop.SetOutputs(outputs)

	pop.InitIndividuals(initPrgrm)
	pop.InitFunctionSet(functionSetNames)
	pop.InitFunctionsToEvolve(functionToEvolve)

	// Evolving the population. The errors between the real and simulated data will be printed to standard output.
	pop.Evolve(evolve.EvolveConfig{
		MazeWidth:      mazeWidth,
		MazeHeight:     mazeHeight,
		Epochs:         epochsCount,
		PlotFitness:    plotFitness,
		SaveAST:        saveAST,
		RandomMazeSize: randomMazeSize,
	})
}
