package evolve

import (
	"fmt"
	"io/ioutil"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/skycoin/cx-evolves/cxexecutes/worker"
	workerclient "github.com/skycoin/cx-evolves/cxexecutes/worker/client"
	cxplotter "github.com/skycoin/cx-evolves/plotter"
	cxast "github.com/skycoin/cx/cx/ast"
	cxconstants "github.com/skycoin/cx/cx/constants"
)

// Used for concurrent output evaluation
var wg = sync.WaitGroup{}

func (pop *Population) Evolve(cfg EvolveConfig) {
	var histoValues []float64
	var averageValues []float64
	var mostFit []float64
	var availPorts []int
	var saveDirectory string

	output := make([]float64, pop.PopulationSize)
	numIter := pop.Iterations
	solProt := pop.FunctionToEvolve
	fnToEvolveName := solProt.Name
	sPrgrm := cxast.SerializeCXProgramV2(pop.Individuals[0], true, true)

	setEpochLength(&cfg)
	saveDirectory = makeDirectory(&cfg)

	// Make worker ports channel
	availWorkers := worker.GetAvailableWorkers(cfg.WorkersAvailable)
	availPorts = append(availPorts, availWorkers...)
	availPortsCh := make(chan int, len(availPorts))
	go func() {
		for _, val := range availPorts {
			availPortsCh <- val
		}
	}()

	// Evolution process.
	for c := 0; c < int(numIter); c++ {
		startTimeGeneration := time.Now()

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

		// Point Mutation
		// pointMutation(pop)

		// Replacing individuals in population.
		replaceSolution(pop.Individuals[dead1Idx], fnToEvolveName, child1)
		replaceSolution(pop.Individuals[dead2Idx], fnToEvolveName, child2)

		if cfg.MazeBenchmark {
			cfg.RandSeed = generateNewSeed(c, cfg)
		}

		runtime.GOMAXPROCS(48)
		// Evaluation process.
		for i := range pop.Individuals {
			wg.Add(1)
			go func(j, genCount int, cfg EvolveConfig) {
				defer wg.Done()

				// Get worker port number from
				// avail ports channel.
				currPortNum := <-availPortsCh
				cfg.WorkerPortNum = currPortNum

				// pop.Individuals[j].PrintProgram()
				output[j], err = RunBenchmark(pop.Individuals[j], solProt, cfg)
				if err != nil {
					fmt.Printf("err=%v", err)
					output[j] = float64(math.MaxInt32)
				}

				// Append back the worker port number used so that
				// it can be used by another go routine.
				availPortsCh <- currPortNum

				// fmt.Printf("output of program[%v]:%v\n", j, output[j])
			}(i, c, cfg)
		}
		wg.Wait()

		var fittestIndex int = 0
		err = UpdateGraphValues(output, &fittestIndex, &histoValues, &mostFit, &averageValues, &cfg, pop.PopulationSize)
		if err != nil {
			panic(err)
		}

		if cfg.SaveAST || c == numIter-1 {
			err := SaveAST(pop.Individuals[fittestIndex], saveDirectory, c)
			if err != nil {
				panic(err)
			}
		}

		if (cfg.MazeBenchmark && c != 0 && c%cfg.EpochLength == 0) || (cfg.ConstantsBenchmark && c != 0 && c%100 == 0) {
			graphTitle := fmt.Sprintf("Average Fitness Of Individuals (%v)", getBenchmarkName(&cfg))
			cxplotter.PointsPlot(cxplotter.PointsPlotCfg{
				Values:       averageValues,
				Xlabel:       averageXLabel,
				Ylabel:       averageYLabel,
				Title:        graphTitle,
				SaveLocation: saveDirectory + fmt.Sprintf("Generation_%v_", c) + "AverageFitness.png",
			})
		}

		fmt.Printf("Time to finish generation[%v]=%v\n", c, time.Since(startTimeGeneration))
	}

	if cfg.PlotFitness {
		saveGraphs(averageValues, mostFit, histoValues, saveDirectory, getBenchmarkName(&cfg))
	}
}

func RunBenchmark(cxprogram *cxast.CXProgram, solProt *cxast.CXFunction, cfg EvolveConfig) (output float64, err error) {
	var result worker.Result
	var TaskName string
	var VersionNum int = 1

	if cfg.MazeBenchmark {
		TaskName = "maze"
	}

	if cfg.ConstantsBenchmark {
		TaskName = "constants"
	}

	if cfg.EvensBenchmark {
		TaskName = "evens"
	}

	if cfg.OddsBenchmark {
		TaskName = "odds"
	}

	if cfg.PrimesBenchmark {
		TaskName = "primes"
	}

	if cfg.CompositesBenchmark {
		TaskName = "composites"
	}

	if cfg.RangeBenchmark {
		TaskName = "range"
	}

	if cfg.NetworkSimBenchmark {
		TaskName = "network_simulator"
	}

	taskCfg := setTaskParams(cfg)
	workerAddr := fmt.Sprintf(":%v", cfg.WorkerPortNum)
	workerclient.CallWorker(
		workerclient.CallWorkerConfig{
			Task:    TaskName,
			Version: VersionNum,
			Program: cxprogram,
			SolProt: solProt,
			TaskCfg: taskCfg,
		},
		workerAddr,
		&result,
	)

	output = result.Output
	return output, nil
}

func SaveAST(cxprogram *cxast.CXProgram, saveDir string, generationNum int) error {
	// Save best ASTs per generation
	saveASTDirectory := saveDir + "AST/"
	astName := fmt.Sprintf("Generation_%v", generationNum)

	// Save as human-readable string .txt format
	astAsString := []byte(cxast.ToString(cxprogram))
	if err := ioutil.WriteFile(saveASTDirectory+astName+".txt", astAsString, 0644); err != nil {
		return err
	}

	// Save as serialized bytes.
	astInBytes := cxast.SerializeCXProgramV2(cxprogram, true, false)
	if err := ioutil.WriteFile(saveASTDirectory+astName+"_serialized"+".ast", []byte(astInBytes), 0644); err != nil {
		return err
	}

	return nil
}

func UpdateGraphValues(output []float64, fittestIndex *int, histoValues, mostFit, averageValues *[]float64, cfg *EvolveConfig, popuSize int) error {
	var total float64 = 0
	var fittest float64 = output[0]
	for z := 0; z < len(output); z++ {
		fitness := output[z]
		total = total + fitness

		// Get Best fitness per generation
		if fitness < fittest {
			fittest = fitness
			*fittestIndex = z
		}

		// Add fitness for histogram
		*histoValues = append(*histoValues, float64(fitness))
	}

	ave := total / float64(popuSize)
	fmt.Printf("Average score=%v\n", ave)
	if cfg.UseAntiLog2 {
		ave = math.Pow(2, ave)
		fittest = math.Pow(2, fittest)
	}

	// Add average values for Average fitness graph
	*averageValues = append(*averageValues, ave)

	// Add fittest values for Fittest per generation graph
	*mostFit = append(*mostFit, fittest)
	return nil
}

func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}
