package evolve

import (
	"fmt"

	cxplotter "github.com/skycoin/cx-evolves/plotter"
)

const (
	averageXLabel = "Generation Number"
	averageYLabel = "Ave Fitness"
	fittestXLabel = "Generation Number"
	fittestYLabel = "Fitness"
)

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
