package evolve

import (
	"fmt"
	"math"

	cxplotter "github.com/skycoin/cx-evolves/plotter"
)

const (
	averageXLabel = "Generation Number"
	averageYLabel = "Ave Fitness"
	fittestXLabel = "Generation Number"
	fittestYLabel = "Fitness"
)

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

func saveGraphs(aveFitnessValues, fittestValues, histoValues []float64, saveDirectory, benchmarkName string) {
	averageGraphTitle := fmt.Sprintf("Average Fitness Of Individuals (%v)", benchmarkName)
	cxplotter.PointsPlot(cxplotter.PointsPlotCfg{
		Values:       aveFitnessValues,
		Xlabel:       averageXLabel,
		Ylabel:       averageYLabel,
		Title:        averageGraphTitle,
		SaveLocation: saveDirectory + "AverageFitness.png",
	})

	fittestGraphTitle := fmt.Sprintf("Fittest Per Generation N (%v)", benchmarkName)
	cxplotter.PointsPlot(cxplotter.PointsPlotCfg{
		Values:       fittestValues,
		Xlabel:       fittestXLabel,
		Ylabel:       fittestYLabel,
		Title:        fittestGraphTitle,
		SaveLocation: saveDirectory + "FittestPerGeneration.png",
	})

	cxplotter.HistogramPlot(cxplotter.HistoPlotCfg{
		Values:       histoValues,
		Title:        "Fitness Distribution of all programs across all generations",
		SaveLocation: saveDirectory + "Histogram.png",
	})
}
