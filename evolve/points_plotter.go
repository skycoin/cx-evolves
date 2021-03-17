package evolve

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func pointsPlot(values []float64, Xlabel, Ylabel, title, saveLocation string) {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = Xlabel
	p.Y.Label.Text = Ylabel

	err := plotutil.AddLinePoints(p,
		"line", Points(values))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, saveLocation); err != nil {
		panic(err)
	}
}

// Points returns some random x, y points.
func Points(values []float64) plotter.XYs {
	pts := make(plotter.XYs, len(values))
	for i := range pts {
		pts[i].X = float64(i)
		pts[i].Y = values[i]
	}
	return pts
}

func histogramPlot(values plotter.Values, title, saveLocation string) {
	p := plot.New()
	p.Title.Text = title

	hist, err := plotter.NewHist(values, 500)
	if err != nil {
		panic(err)
	}
	p.Add(hist)

	if err := p.Save(8*vg.Inch, 8*vg.Inch, saveLocation); err != nil {
		panic(err)
	}
}
