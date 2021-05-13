package plotter

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func HistogramPlot(cfg HistoPlotCfg) {
	p := plot.New()
	p.Title.Text = cfg.Title

	hist, err := plotter.NewHist(cfg.Values, 500)
	if err != nil {
		panic(err)
	}
	p.Add(hist)

	if err := p.Save(8*vg.Inch, 8*vg.Inch, cfg.SaveLocation); err != nil {
		panic(err)
	}
}
