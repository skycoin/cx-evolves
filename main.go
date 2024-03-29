package main

import (
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	evolve "github.com/skycoin/cx-evolves/evolve"
	cxmutation "github.com/skycoin/cx-evolves/mutation"
	cxprobability "github.com/skycoin/cx-evolves/probability"
	cxtasks "github.com/skycoin/cx-evolves/tasks"
	cxast "github.com/skycoin/cx/cx/ast"
	cxconstants "github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
	cxactions "github.com/skycoin/cx/cxparser/actions"
	cxparsing "github.com/skycoin/cx/cxparser/cxparsing"
)

// Maze and Output Configuration
var (
	taskName    string
	taskVersion int

	// Maze Config
	mazeWidth      int
	mazeHeight     int
	randomMazeSize bool

	numberOfRounds int

	constantsTarget int

	upperRange int
	lowerRange int

	epochLength int
	plotFitness bool
	saveAST     bool
	logFitness  bool

	workersAvailable int

	randomSearch bool

	selectRankCutoff bool
)

// Evolve Configuration
var (
	expressionsCount int
	populationSize   int
	iterations       int
	functionToEvolve string

	// What functions from CX standard library can we use to create expressions in the programs.
	functionSetNames = []string{"i32.jmpeq", "i32.jmpuneq", "i32.jmpgt", "i32.jmpgteq", "i32.jmplt", "i32.jmplteq", "i32.jmpzero", "i32.jmpnotzero", "jmp", "nop", "i32.add", "i32.mul", "i32.sub", "i32.eq", "i32.uneq", "i32.gt", "i32.gteq", "i32.lt", "i32.lteq", "bool.not", "bool.or", "bool.and", "bool.uneq", "bool.eq", "i32.neg", "i32.abs", "i32.bitand", "i32.bitor", "i32.bitxor", "i32.bitclear", "i32.bitshl", "i32.bitshr", "i32.max", "i32.min", "i32.rand"}
	// Missing
	// "i32.div", "i32.mod",

	// What's the input signature of the programs being evolved.
	inputSignature []string

	// What's the output signature of the programs being evolved.
	outputSignature []string
)

func InitialProgram() *cxast.CXProgram {
	// Creating the initial CX program.
	prgrm := cxast.MakeProgram()
	prgrm.SetCurrentCxProgram()
	cxactions.SelectProgram(prgrm)

	// Adding `main` package.
	mainPkg := cxast.MakePackage(cxconstants.MAIN_PKG)
	prgrm.AddPackage(mainPkg)

	// Adding `main` function to `main` package.
	mainFn := cxast.MakeFunction(cxconstants.MAIN_FUNC, "", -1)
	mainFn.Package = mainPkg
	mainPkg.AddFunction(mainFn)

	// Adding function to evolve (`FunctionToEvolve`).
	toEvolveFn := cxast.MakeFunction(functionToEvolve, "", -1)
	mainPkg.AddFunction(toEvolveFn)

	// Adding input signature to function to evolve (`FunctionToEvolve`).
	for _, inpType := range inputSignature {
		var dataType types.Code
		if inpType == "i32" {
			dataType = types.I32
		}

		inp := cxast.MakeArgument(cxactions.MakeGenSym("evo_inp"), "", -1).AddType(dataType)
		inp.AddPackage(mainPkg)
		toEvolveFn.AddInput(inp)
	}

	// Adding output signature to function to evolve (`FunctionToEvolve`).
	for _, outType := range outputSignature {
		var dataType types.Code
		if outType == "i32" {
			dataType = types.I32
		}
		out := cxast.MakeArgument(cxactions.MakeGenSym("evo_out"), "", -1).AddType(dataType)
		out.AddPackage(mainPkg)
		toEvolveFn.AddOutput(out)
	}

	// Creating an init function for the CX program.
	cxparsing.AddInitFunction(prgrm)

	return prgrm
}

func setInputOutputSignature() {
	// Set input and output signature based on what to benchmark
	if cxtasks.IsMazeTask(taskName) {
		inputSignature = []string{"i32", "i32", "i32", "i32", "i32", "i32", "i32", "i32", "i32", "i32", "i32", "i32", "i32"}
		outputSignature = []string{"i32"}
	}
	if cxtasks.IsConstantsTask(taskName) ||
		cxtasks.IsEvensTask(taskName) ||
		cxtasks.IsOddsTask(taskName) ||
		cxtasks.IsPrimesTask(taskName) ||
		cxtasks.IsCompositesTask(taskName) ||
		cxtasks.IsRangeTask(taskName) ||
		cxtasks.IsNetworkSimulatorTask(taskName) {
		inputSignature = []string{"i32"}
		outputSignature = []string{"i32"}
	}
}

func main() {
	EvolveApp := &cli.App{
		Name:    "Evolve with Maze",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "Task Name",
				Value:       "maze",
				Aliases:     []string{"task"},
				Usage:       "Name of task to benchmark",
				Destination: &taskName,
			},
			&cli.IntFlag{
				Name:        "TaskVersion",
				Aliases:     []string{"task-version"},
				Usage:       "version of task",
				Value:       1,
				Destination: &taskVersion,
			},
			&cli.BoolFlag{
				Name:        "log 2 for fitness",
				Aliases:     []string{"use-log-fitness"},
				Usage:       "set true if fitness will be log2",
				Destination: &logFitness,
			},
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
				Name:        "rounds number",
				Aliases:     []string{"rounds-total"},
				Usage:       "number of rounds for numbers benchmarking",
				Value:       6,
				Destination: &numberOfRounds,
			},
			&cli.IntFlag{
				Name:        "target constants",
				Aliases:     []string{"constants-target"},
				Usage:       "target number for constants benchmarking",
				Value:       -1,
				Destination: &constantsTarget,
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
				Name:        "Epoch Length",
				Aliases:     []string{"epoch-length"},
				Usage:       "Maze changes for every N generations, if set 0 then maze changes every generations",
				Value:       1,
				Destination: &epochLength,
			},
			&cli.IntFlag{
				Name:        "Workers Available",
				Aliases:     []string{"workers"},
				Usage:       "Number of CX Programs workers available/deployed",
				Value:       1,
				Destination: &workersAvailable,
			},
			&cli.BoolFlag{
				Name:        "Random Maze Size",
				Aliases:     []string{"random-maze-size"},
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
			&cli.BoolFlag{
				Name:        "RandomSearch",
				Aliases:     []string{"random-search"},
				Usage:       "set true to have no mutation on individuals",
				Destination: &randomSearch,
			},
			&cli.BoolFlag{
				Name:        "SelectRankCutoff",
				Aliases:     []string{"select-rank-cutoff"},
				Usage:       "set true if selection is select, rank, and cutoff",
				Destination: &selectRankCutoff,
			},
		},
		Action: func(c *cli.Context) error {
			Evolve()
			return nil
		},
	}

	err := EvolveApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Evolve() {
	// Setting seed so results vary every time we run the example.
	rand.Seed(time.Now().UTC().UnixNano())

	setInputOutputSignature()
	functionToEvolve = taskName

	// We create an initial CX program, with a
	initPrgrm := InitialProgram()

	// Initialize point mutation operators
	cxmutation.RegisterMutationOperators()

	// Initialize point operator probability
	pointOpFns := cxmutation.GetAllMutationOperatorFunctionSet()

	pointMutationOperatorCDF := cxprobability.NewProbability(cxprobability.GetEqualDensity(len(pointOpFns)))
	mutationCrossoverCDF := cxprobability.NewProbability([]float32{1, 1, 98})

	// Generate a population.
	pop := evolve.MakePopulation(populationSize)

	// Configuring the population. The method calls are self-explanatory.
	pop.SetIterations(iterations)
	pop.SetExpressionsCount(expressionsCount)

	pop.InitIndividuals(initPrgrm)
	pop.InitFunctionSet(functionSetNames)
	pop.InitFunctionsToEvolve(functionToEvolve)

	// Evolving the population. The errors between the real and simulated data will be printed to standard output.
	pop.Evolve(evolve.EvolveConfig{
		TaskName: taskName,
		Version:  taskVersion,

		MazeWidth:  mazeWidth,
		MazeHeight: mazeHeight,

		NumberOfRounds: numberOfRounds,

		ConstantsTarget: constantsTarget,

		UpperRange: upperRange,
		LowerRange: lowerRange,

		EpochLength:    epochLength,
		PlotFitness:    plotFitness,
		SaveAST:        saveAST,
		RandomMazeSize: randomMazeSize,
		UseAntiLog2:    logFitness,

		WorkersAvailable: workersAvailable,

		RandomSearch: randomSearch,

		SelectRankCutoff: selectRankCutoff,

		PointMutationOperatorCDF: pointMutationOperatorCDF,
		MutationCrossoverCDF:     mutationCrossoverCDF,
	})
}
