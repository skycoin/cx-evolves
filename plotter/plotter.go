package plotter

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

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

// Points returns plotter x, y points.
func Points(values []float64) plotter.XYs {
	pts := make(plotter.XYs, len(values))
	for i := range pts {
		pts[i].X = float64(i)
		pts[i].Y = values[i]
	}
	return pts
}

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
