package plotter

import (
	"gonum.org/v1/plot/plotter"
)

// Points returns plotter x, y points.
func Points(values []float64) plotter.XYs {
	pts := make(plotter.XYs, len(values))
	for i := range pts {
		pts[i].X = float64(i)
		pts[i].Y = values[i]
	}
	return pts
}

func UpdateDataPoints(plotData *PlotData, generation int, output []float64, saveDirectory string) error {
	plotDataPoints := PlotDataPoints{
		Generation: generation,
	}
	plotDataPoints.Output = append(plotDataPoints.Output, output...)
	plotData.Data = append(plotData.Data, plotDataPoints)

	err := SavePlotGraphDataToJSON(*plotData, saveDirectory+"data_points")
	if err != nil {
		return err
	}
	return nil
}
