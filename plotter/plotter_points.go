package plotter

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func PointsPlot(cfg PointsPlotCfg) {
	p := plot.New()

	p.Title.Text = cfg.Title
	p.X.Label.Text = cfg.Xlabel
	p.Y.Label.Text = cfg.Ylabel

	err := plotutil.AddLinePoints(p,
		"line", Points(cfg.Values))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, cfg.SaveLocation); err != nil {
		panic(err)
	}
}
