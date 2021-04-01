package evolve

import (
	"fmt"
	"io/ioutil"
	"math"
	"runtime"
	"sync"

	"github.com/skycoin/cx-evolves/cmd/maze"
	cxast "github.com/skycoin/cx/cx/ast"
	cxconstants "github.com/skycoin/cx/cx/constants"
)

type EvolveConfig struct {
	MazeBenchmark       bool
	ConstantsBenchmark  bool
	EvensBenchmark      bool
	OddsBenchmark       bool
	PrimesBenchmark     bool
	CompositesBenchmark bool
	RangeBenchmark      bool
	NetworkSimBenchmark bool

	MazeHeight     int
	MazeWidth      int
	RandomMazeSize bool

	NumberOfRounds int

	UpperRange int
	LowerRange int

	EpochLength int
	PlotFitness bool
	SaveAST     bool
	UseAntiLog2 bool
}

// Original Evolve
// func (pop *Population) Evolve() {
// 	errors := make([]float64, pop.PopulationSize)
// 	numIter := pop.Iterations
// 	solProt := pop.FunctionToEvolve
// 	fnToEvolveName := solProt.Name
// 	sPrgrm := cxcore.Serialize(pop.Individuals[0], 0)
// 	targetError := pop.TargetError
// 	inputs := pop.Inputs
// 	outputs := pop.Outputs

// 	// Evolution process.
// 	for c := 0; c < int(numIter); c++ {
// 		// Selection process.
// 		pop1Idx, pop2Idx := tournamentSelection(errors, 0.5, true)
// 		dead1Idx, dead2Idx := tournamentSelection(errors, 0.5, false)

// 		pop1MainPkg, err := pop.Individuals[pop1Idx].GetPackage(cxcore.MAIN_PKG)
// 		if err != nil {
// 			panic(err)
// 		}
// 		parent1, err := pop1MainPkg.GetFunction(fnToEvolveName)
// 		if err != nil {
// 			panic(err)
// 		}

// 		pop2MainPkg, err := pop.Individuals[pop2Idx].GetPackage(cxcore.MAIN_PKG)
// 		if err != nil {
// 			panic(err)
// 		}
// 		parent2, err := pop2MainPkg.GetFunction(fnToEvolveName)
// 		if err != nil {
// 			panic(err)
// 		}

// 		// Crossover process.
// 		crossoverFn := pop.getCrossoverFn()
// 		child1, child2 := crossoverFn(parent1, parent2)
// 		// child1 := parent1
// 		// child2 := parent2

// 		// Mutation process.
// 		_ = sPrgrm
// 		_ = dead1Idx
// 		_ = dead2Idx
// 		_ = child1
// 		_ = child2
// 		randomMutation(pop, sPrgrm)

// 		// Replacing individuals in population.
// 		replaceSolution(pop.Individuals[dead1Idx], fnToEvolveName, child1)
// 		replaceSolution(pop.Individuals[dead2Idx], fnToEvolveName, child2)

// 		// Evaluation process.
// 		for i, _ := range pop.Individuals {
// 			errors[i] = perByteEvaluation(pop.Individuals[i], solProt, inputs, outputs)
// 			if errors[i] <= targetError {
// 				fmt.Printf("\nFound solution:\n\n")
// 				pop.Individuals[i].PrintProgram()
// 				return
// 			}
// 		}

// 		avg := 0.0
// 		for _, err := range errors {
// 			avg += err
// 		}
// 		fmt.Printf("%v\n", float64(avg) / float64(len(errors)))
// 	}
// }

// Used for concurrent output evaluation
var wg = sync.WaitGroup{}

func (pop *Population) Evolve(cfg EvolveConfig) {
	var histoValues []float64
	var averageValues []float64
	var mostFit []float64
	var game maze.Game
	var saveDirectory string

	output := make([]float64, pop.PopulationSize)
	numIter := pop.Iterations
	solProt := pop.FunctionToEvolve
	fnToEvolveName := solProt.Name
	sPrgrm := cxast.SerializeCXProgram(pop.Individuals[0], true)

	setEpochLength(&cfg)
	saveDirectory = makeDirectory(&cfg)

	// Evolution process.
	for c := 0; c < int(numIter); c++ {
		// Maze Creation if Maze Benchmark
		if cfg.MazeBenchmark {
			generateNewMaze(c, &cfg, &game)
		}

		// Selection process.
		pop1Idx, pop2Idx := tournamentSelection(output, 0.5, true)
		dead1Idx, dead2Idx := tournamentSelection(output, 0.5, false)

		pop1MainPkg, err := pop.Individuals[pop1Idx].GetPackage(cxconstants.MAIN_PKG)
		if err != nil {
			panic(err)
		}
		parent1, err := pop1MainPkg.GetFunction(fnToEvolveName)
		if err != nil {
			panic(err)
		}

		pop2MainPkg, err := pop.Individuals[pop2Idx].GetPackage(cxconstants.MAIN_PKG)
		if err != nil {
			panic(err)
		}
		parent2, err := pop2MainPkg.GetFunction(fnToEvolveName)
		if err != nil {
			panic(err)
		}

		// Crossover process.
		crossoverFn := pop.getCrossoverFn()
		child1, child2 := crossoverFn(parent1, parent2)

		// Mutation process.
		_ = sPrgrm
		_ = dead1Idx
		_ = dead2Idx
		_ = child1
		_ = child2
		randomMutation(pop, sPrgrm)

		// Replacing individuals in population.
		replaceSolution(pop.Individuals[dead1Idx], fnToEvolveName, child1)
		replaceSolution(pop.Individuals[dead2Idx], fnToEvolveName, child2)

		runtime.GOMAXPROCS(4)
		// Evaluation process.
		for i := range pop.Individuals {
			wg.Add(1)
			go func(j int) {
				output[j] = RunBenchmark(pop.Individuals[j], solProt, &cfg, &game)
				wg.Done()
				fmt.Printf("output of program[%v]:%v\n", j, output[j])
			}(i)
		}
		wg.Wait()

		var total float64 = 0
		var fittestIndex int = 0
		var fittest float64 = output[0]
		for z := 0; z < len(output); z++ {
			fitness := output[z]
			total = total + fitness

			// Get Best fitness per generation
			if fitness < fittest {
				fittest = fitness
				fittestIndex = z
			}

			// Add fitness for histogram
			histoValues = append(histoValues, float64(fitness))
		}

		ave := total / float64(pop.PopulationSize)

		if cfg.UseAntiLog2 {
			ave = math.Pow(2, ave)
			fittest = math.Pow(2, fittest)
		}

		// Add average values for Average fitness graph
		averageValues = append(averageValues, ave)

		// Add fittest values for Fittest per generation graph
		mostFit = append(mostFit, fittest)

		// if cfg.SaveAST {
		// 	// Save best ASTs per generation
		// 	saveASTDirectory := saveDirectory + "AST/"
		// 	astName := fmt.Sprintf("Generation_%v", c)

		// 	// Save as human-readable string .txt format
		// 	astAsString := []byte(cxast.ToString(pop.Individuals[fittestIndex]))
		// 	if err := ioutil.WriteFile(saveASTDirectory+astName+".txt", astAsString, 0644); err != nil {
		// 		panic(err)
		// 	}

		// 	// Save as serialized bytes.
		// 	astInBytes := cxast.SerializeCXProgram(pop.Individuals[fittestIndex], false)
		// 	if err := ioutil.WriteFile(saveASTDirectory+astName+"_serialized"+".ast", []byte(astInBytes), 0644); err != nil {
		// 		panic(err)
		// 	}
		// }
		if cfg.SaveAST {
			// Save best ASTs per generation
			saveASTDirectory := saveDirectory + "AST/"
			astName := fmt.Sprintf("Generation_%v", c)

			// Save as human-readable string .txt format.
			astInBytes := []byte(cxast.ToString(pop.Individuals[fittestIndex]))
			if err := ioutil.WriteFile(saveASTDirectory+astName+".txt", astInBytes, 0644); err != nil {
				panic(err)
			}

			// // Save as serialized bytes.
			// astInBytess := cxast.SerializeCXProgram(pop.Individuals[fittestIndex], false)
			// if err := ioutil.WriteFile(saveASTDirectory+astName+"_serialized"+".ast", []byte(astInBytess), 0644); err != nil {
			// 	panic(err)
			// }

			// dat, _ := ioutil.ReadFile(saveASTDirectory + astName + "_serialized" + ".ast")
			// deserProg := cxast.Deserialize(dat)
			// fmt.Println("Deserialized ast")
			// deserProg.PrintProgram()
		}
	}

	if cfg.PlotFitness {
		saveGraphs(averageValues, mostFit, histoValues, saveDirectory)
	}
}

func RunBenchmark(cxprogram *cxast.CXProgram, solProt *cxast.CXFunction, cfg *EvolveConfig, game *maze.Game) float64 {
	var intOut float64
	if cfg.MazeBenchmark {
		intOut = mazeMovesEvaluation(cxprogram, solProt, *game)
	}
	if cfg.ConstantsBenchmark {
		intOut = perByteEvaluation_Constants(cxprogram, solProt, cfg.NumberOfRounds)
	}

	if cfg.EvensBenchmark {
		intOut = perByteEvaluation_Evens(cxprogram, solProt, cfg.NumberOfRounds)
	}

	if cfg.OddsBenchmark {
		intOut = perByteEvaluation_Odds(cxprogram, solProt, cfg.NumberOfRounds)
	}

	if cfg.PrimesBenchmark {
		intOut = perByteEvaluation_Primes(cxprogram, solProt, cfg.NumberOfRounds)
	}

	if cfg.CompositesBenchmark {
		intOut = perByteEvaluation_Composites(cxprogram, solProt, cfg.NumberOfRounds)
	}

	if cfg.RangeBenchmark {
		intOut = perByteEvaluation_Range(cxprogram, solProt, cfg.NumberOfRounds, cfg.UpperRange, cfg.LowerRange)
	}

	if cfg.NetworkSimBenchmark {
		intOut = perByteEvaluation_NetworkSim(cxprogram, solProt, cfg.NumberOfRounds)
	}
	return intOut
}
