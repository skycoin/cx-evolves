package evolve

import (
	"fmt"
	"math"
)

const (
// averageXLabel = "Generation Number"
// averageYLabel = "Ave Fitness"
// fittestXLabel = "Generation Number"
// fittestYLabel = "Fitness"

// averageTitle = "Average Fitness Of Individuals"
// fittestTitle = "Fittest Per Generation N"
// histoTitle   = "Fitness Distribution of all programs across all generations"

// averageFileExtension = "AverageFitness.png"
// fittestFileExtension = "FittestPerGeneration.png"
// histoFileExtension   = "Histogram.png"
// boxPlotExtension     = "BoxPlot.png"
)

func UpdateGraphValues(cfg GraphCfg) error {
	var total float64 = 0
	var fittest float64 = cfg.Output[0]
	for z := 0; z < len(cfg.Output); z++ {
		fitness := cfg.Output[z]
		total = total + fitness

		// Get Best fitness per generation
		if fitness < fittest {
			fittest = fitness
			*cfg.FittestIndex = z
		}

		// Add fitness for histogram
		*cfg.HistoValues = append(*cfg.HistoValues, float64(fitness))
	}

	ave := total / float64(cfg.PopuSize)
	fmt.Printf("Average score=%v\n", ave)
	if cfg.EvolveCfg.UseAntiLog2 {
		ave = math.Pow(2, ave)
		fittest = math.Pow(2, fittest)
	}

	// Add average values for Average fitness graph
	*cfg.AverageValues = append(*cfg.AverageValues, ave)

	// Add fittest values for Fittest per generation graph
	*cfg.MostFit = append(*cfg.MostFit, fittest)
	return nil
}

// func saveGraphs(aveFitnessValues, fittestValues, histoValues []float64, saveDirectory, benchmarkName string) {
// 	averageGraphTitle := fmt.Sprintf(averageTitle+" (%v)", benchmarkName)
// 	cxplotter.PointsPlot(cxplotter.PointsPlotCfg{
// 		Values:       aveFitnessValues,
// 		Xlabel:       averageXLabel,
// 		Ylabel:       averageYLabel,
// 		Title:        averageGraphTitle,
// 		SaveLocation: saveDirectory + averageFileExtension,
// 	})

// 	fittestGraphTitle := fmt.Sprintf(fittestTitle+" (%v)", benchmarkName)
// 	cxplotter.PointsPlot(cxplotter.PointsPlotCfg{
// 		Values:       fittestValues,
// 		Xlabel:       fittestXLabel,
// 		Ylabel:       fittestYLabel,
// 		Title:        fittestGraphTitle,
// 		SaveLocation: saveDirectory + fittestFileExtension,
// 	})

// 	cxplotter.HistogramPlot(cxplotter.HistoPlotCfg{
// 		Values:       histoValues,
// 		Title:        histoTitle,
// 		SaveLocation: saveDirectory + histoFileExtension,
// 	})
// }
