package evolve

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/skycoin/cx-evolves/cmd/maze"
	cxcore "github.com/skycoin/cx/cx"
)

type EvolveConfig struct {
	MazeHeight     int
	MazeWidth      int
	Epochs         int
	PlotFitness    bool
	SaveAST        bool
	RandomMazeSize bool
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

// Used for concurrent maze moves evaluation
var wg = sync.WaitGroup{}

func (pop *Population) Evolve(cfg EvolveConfig) {
	var histoValues []float64
	var averageValues []float64
	var mostFit []float64
	var game maze.Game
	var saveDirectory string

	moves := make([]float64, pop.PopulationSize)
	numIter := pop.Iterations
	solProt := pop.FunctionToEvolve
	fnToEvolveName := solProt.Name
	sPrgrm := cxcore.Serialize(pop.Individuals[0], 0)

	if cfg.Epochs == 0 {
		cfg.Epochs = 1
	}

	if cfg.PlotFitness || cfg.SaveAST {
		saveDirectory = getSaveDirectory(&cfg)

		// create directory
		_ = os.Mkdir(saveDirectory, 0700)

		if cfg.SaveAST {
			_ = os.Mkdir(saveDirectory+"AST/", 0700)
		}
	}

	// Evolution process.
	for c := 0; c < int(numIter); c++ {
		// Maze Creation
		if c%cfg.Epochs == 0 || c == 0 {
			if cfg.RandomMazeSize {
				setRandomMazeSize(&cfg)
			}

			game = maze.Game{}
			game.Init(cfg.MazeWidth, cfg.MazeHeight)
		}

		// Selection process.
		pop1Idx, pop2Idx := tournamentSelection(moves, 0.5, true)
		dead1Idx, dead2Idx := tournamentSelection(moves, 0.5, false)

		pop1MainPkg, err := pop.Individuals[pop1Idx].GetPackage(cxcore.MAIN_PKG)
		if err != nil {
			panic(err)
		}
		parent1, err := pop1MainPkg.GetFunction(fnToEvolveName)
		if err != nil {
			panic(err)
		}

		pop2MainPkg, err := pop.Individuals[pop2Idx].GetPackage(cxcore.MAIN_PKG)
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
				moves[j] = mazeMovesEvaluation(pop.Individuals[j], solProt, game)
				fmt.Printf("moves of program[%v]:%v\n", j, moves[j])
			}(i)
		}
		wg.Wait()

		var total int = 0
		var fittestIndex int = 0
		var fittest int = int(moves[0])
		for z := 0; z < len(moves); z++ {
			fitness := int(moves[z])
			total = total + fitness

			// Get Best fitness per generation
			if fitness < fittest {
				fittest = fitness
				fittestIndex = z
			}

			// Add fitness for histogram
			histoValues = append(histoValues, float64(fitness))
		}

		ave := total / pop.PopulationSize

		// Add average values for Average fitness graph
		averageValues = append(averageValues, float64(ave))

		// Add fittest values for Fittest per generation graph
		mostFit = append(mostFit, float64(fittest))

		if cfg.SaveAST {
			// Save best ASTs per generation
			saveASTDirectory := saveDirectory + "AST/"
			astName := fmt.Sprintf("Generation_%v", c)
			pop.Individuals[fittestIndex].PrintProgram()
			if err := ioutil.WriteFile(saveASTDirectory+astName+".ast", []byte(fmt.Sprintf("%v", pop.Individuals[fittestIndex])), 0644); err != nil {
				panic(err)
			}
		}
	}

	if cfg.PlotFitness {
		pointsPlot(averageValues, "Generation Number", "Ave Fitness", "Average Fitness Of Individuals In Generation N", saveDirectory+"AverageFitness.png")
		pointsPlot(mostFit, "Generation Number", "Fitness", "Fittest Per Generation N", saveDirectory+"FittestPerGeneration.png")
		histogramPlot(histoValues, "Fitness Distribution of all programs across all generations", saveDirectory+"Histogram.png")
	}
}

func getSaveDirectory(cfg *EvolveConfig) string {
	// Unixtime-Maze-2x2
	mazeSize := fmt.Sprintf("%vx%v", cfg.MazeWidth, cfg.MazeHeight)
	if cfg.RandomMazeSize {
		mazeSize = "random"
	}

	return fmt.Sprintf("./Results/%v-%v-%v/", time.Now().Unix(), "Maze", mazeSize)
}

func setRandomMazeSize(cfg *EvolveConfig) {
	rand.Seed(time.Now().Unix())
	randOptions := []int{2, 3, 4, 5, 6, 7, 8}
	size := randOptions[rand.Int()%len(randOptions)]
	cfg.MazeWidth = size
	cfg.MazeHeight = size
}
