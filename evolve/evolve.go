package evolve

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/skycoin/cx-evolves/cxexecutes/worker"
	workerclient "github.com/skycoin/cx-evolves/cxexecutes/worker/client"
	cxplotter "github.com/skycoin/cx-evolves/plotter"
	cxprobability "github.com/skycoin/cx-evolves/probability"
	cxtasks "github.com/skycoin/cx-evolves/tasks"
	cxast "github.com/skycoin/cx/cx/ast"
	cxconstants "github.com/skycoin/cx/cx/constants"
)

// Used for concurrent output evaluation
var wg = sync.WaitGroup{}

func (pop *Population) Evolve(cfg EvolveConfig) {
	var availPorts []int
	var saveDirectory string
	var plotData cxplotter.PlotData
	plotData.Title = getBenchmarkName(&cfg)

	output := make([]float64, pop.PopulationSize)
	numIter := pop.Iterations
	solProt := pop.FunctionToEvolve
	fnToEvolveName := solProt.Name
	sPrgrm := cxast.SerializeCXProgramV2(pop.Individuals[0], true, true)

	setEpochLength(&cfg)
	saveDirectory = makeDirectory(&cfg)

	logF, err := setupLogger(fmt.Sprintf("%v-log.txt", time.Now().Format(time.RFC3339)), saveDirectory)
	if err != nil {
		panic(err)
	}
	defer logF.Close()
	log.SetOutput(logF)

	log.Printf("Benchmark config: %+v\n", cfg)
	log.Printf("Generations: %v\n", pop.Iterations)
	log.Printf("Population Size: %v\n", pop.PopulationSize)
	log.Printf("Expressions count: %v\n", pop.ExpressionsCount)

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
	for c := 1; c <= numIter; c++ {
		startTimeGeneration := time.Now()

		// Selection process.
		pop1Idx, pop2Idx := tournamentSelection(output, 0.5, true)
		dead1Idx, dead2Idx := tournamentSelection(output, 0.5, false)

		pop1MainPkg, err := pop.Individuals[pop1Idx].GetPackage(cxconstants.MAIN_PKG)
		if err != nil {
			log.Printf("error get package: %v", err)
			panic(err)
		}

		parent1, err := pop1MainPkg.GetFunction(fnToEvolveName)
		if err != nil {
			log.Printf("error get function: %v", err)
			panic(err)
		}

		pop2MainPkg, err := pop.Individuals[pop2Idx].GetPackage(cxconstants.MAIN_PKG)
		if err != nil {
			log.Printf("error get package: %v", err)
			panic(err)
		}

		parent2, err := pop2MainPkg.GetFunction(fnToEvolveName)
		if err != nil {
			log.Printf("error get function: %v", err)
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

		// if random-search is true, there will be no mutation.
		if cfg.RandomSearch {
			// Replace tournament loser with new generated program.
			ReplaceIndividualWithRandomExpressions(pop.Individuals[dead1Idx], pop.FunctionToEvolve, pop.FunctionSet, pop.ExpressionsCount)
			ReplaceIndividualWithRandomExpressions(pop.Individuals[dead2Idx], pop.FunctionToEvolve, pop.FunctionSet, pop.ExpressionsCount)
		} else {
			mutationOption := cxprobability.GetRandIndex(cfg.MutationCrossoverCDF)
			switch mutationOption {
			case 1:
				ReplaceRandomIndividualWithRandom(pop, sPrgrm)
			case 2:
				pointMutation(pop, cfg.PointMutationOperatorCDF)
			case 3:
				// Replace tournament losers with children of tournament winners.
				replaceSolution(pop.Individuals[dead1Idx], fnToEvolveName, child1)
				replaceSolution(pop.Individuals[dead2Idx], fnToEvolveName, child2)
			}
		}

		if cxtasks.IsMazeTask(cfg.TaskName) {
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
					log.Printf("error run benchmark: %v", err)
					output[j] = float64(math.MaxInt32)
				}

				// Append back the worker port number used so that
				// it can be used by another go routine.
				availPortsCh <- currPortNum

				// fmt.Printf("output of program[%v]:%v\n", j, output[j])
			}(i, c, cfg)
		}
		wg.Wait()

		// Update data points values
		err = cxplotter.UpdateDataPoints(&plotData, c, output, saveDirectory)
		if err != nil {
			log.Printf("error updating data points: %v", err)
			panic(err)
		}

		if cfg.SaveAST || c == numIter-1 {
			err := SaveAST(pop.Individuals[getFittestIndex(output)], saveDirectory, c)
			if err != nil {
				log.Printf("error saving ast: %v", err)
				panic(err)
			}
		}

		fmt.Printf("Time to finish generation[%v]=%v\n", c, time.Since(startTimeGeneration))
	}
}

func RunBenchmark(cxprogram *cxast.CXProgram, solProt *cxast.CXFunction, cfg EvolveConfig) (output float64, err error) {
	var result worker.Result

	taskCfg := setTaskParams(cfg)
	workerAddr := fmt.Sprintf(":%v", cfg.WorkerPortNum)
	workerclient.CallWorker(
		workerclient.CallWorkerConfig{
			Task:    cfg.TaskName,
			Version: cfg.Version,
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
