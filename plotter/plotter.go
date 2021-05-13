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
