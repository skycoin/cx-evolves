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
	cxtasks "github.com/skycoin/cx-evolves/tasks"
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

	// Create box plot
	boxPlotTitle := fmt.Sprintf("Box Plot"+" (%v)", getBenchmarkName(&cfg))
	evolveBoxPlot := cxplotter.NewBoxPlot(boxPlotTitle, fittestXLabel, fittestYLabel)

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
	for c := 0; c < numIter; c++ {
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

		// if random-search is true, there will be no mutation.
		if cfg.RandomSearch {
			// Replace tournament loser with new generated program.
			GenerateNewIndividualWithRandomExpressions(pop.Individuals[dead1Idx], pop.FunctionToEvolve, pop.FunctionSet, pop.ExpressionsCount)
			GenerateNewIndividualWithRandomExpressions(pop.Individuals[dead2Idx], pop.FunctionToEvolve, pop.FunctionSet, pop.ExpressionsCount)
		} else {
			randomMutation(pop, sPrgrm)

			// Point Mutation
			pointMutation(pop)

			// Replacing individuals in population.
			replaceSolution(pop.Individuals[dead1Idx], fnToEvolveName, child1)
			replaceSolution(pop.Individuals[dead2Idx], fnToEvolveName, child2)
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
		err = UpdateGraphValues(GraphCfg{
			Output:        output,
			FittestIndex:  &fittestIndex,
			HistoValues:   &histoValues,
			MostFit:       &mostFit,
			AverageValues: &averageValues,
			EvolveCfg:     &cfg,
			PopuSize:      pop.PopulationSize,
		})
		if err != nil {
			panic(err)
		}

		cxplotter.AddDataToBoxPlot(evolveBoxPlot, output, c)
		// For now only latest 10 generations to show on the graph.
		if (c+10)%cfg.EpochLength == 0 {
			// Reset Box Plot
			evolveBoxPlot = cxplotter.ResetBoxPlot(evolveBoxPlot)
		}
		if cfg.SaveAST || c == numIter-1 {
			err := SaveAST(pop.Individuals[fittestIndex], saveDirectory, c)
			if err != nil {
				panic(err)
			}
		}

		if (cxtasks.IsMazeTask(cfg.TaskName) && c != 0 && c%cfg.EpochLength == 0) || (!cxtasks.IsMazeTask(cfg.TaskName) && c != 0 && c%100 == 0) {
			graphTitle := fmt.Sprintf(averageTitle+" (%v)", getBenchmarkName(&cfg))
			cxplotter.PointsPlot(cxplotter.PointsPlotCfg{
				Values:       averageValues,
				Xlabel:       averageXLabel,
				Ylabel:       averageYLabel,
				Title:        graphTitle,
				SaveLocation: saveDirectory + fmt.Sprintf("Generation_%v_", c) + averageFileExtension,
			})

			// Save Box Plot
			cxplotter.SaveBoxPlot(evolveBoxPlot, saveDirectory+fmt.Sprintf("Generation_%v_", c)+boxPlotExtension)
		}

		fmt.Printf("Time to finish generation[%v]=%v\n", c, time.Since(startTimeGeneration))
	}

	if cfg.PlotFitness {
		saveGraphs(averageValues, mostFit, histoValues, saveDirectory, getBenchmarkName(&cfg))
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
