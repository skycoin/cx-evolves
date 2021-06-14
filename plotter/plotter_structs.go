package plotter

import "gonum.org/v1/plot/plotter"

type PointsPlotCfg struct {
	Values       []float64
	Xlabel       string
	Ylabel       string
	Title        string
	SaveLocation string
}

type HistoPlotCfg struct {
	Values       plotter.Values
	Title        string
	SaveLocation string
}
