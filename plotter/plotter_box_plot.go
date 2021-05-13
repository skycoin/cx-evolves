package plotter

import (
	"fmt"
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

const (
	boxWidth = 5
)

///////////////////////////////////////////////////
// Source: gonum/plot/plotter/boxplot.go
//
// The fence values are 1.5x the interquartile before
// the first quartile and after the third quartile.  Any
// value that is outside of the fences are drawn as
// Outside points.  The adjacent values (to which the
// whiskers stretch) are the minimum and maximum
// values that are not outside the fences.
///////////////////////////////////////////////////

func NewBoxPlot(title, xLabel, yLabel string) *plot.Plot {
	boxPlot := plot.New()
	boxPlot.Title.Text = title
	boxPlot.Y.Label.Text = yLabel
	boxPlot.X.Label.Text = xLabel

	return boxPlot
}

func AddDataToBoxPlot(boxPlot *plot.Plot, values []float64, generation int) {
	plotValues := make(plotter.ValueLabels, len(values))
	for i, val := range values {
		fmt.Printf("values[%v]=%v\n", i, val)
		plotValues[i].Value = val
		plotValues[i].Label = fmt.Sprintf("%4.4f", val)
	}

	// Make boxes for our data and add them to the plot.
	plotValuesBox, err := plotter.NewBoxPlot(vg.Points(boxWidth), float64(generation), plotValues)
	if err != nil {
		panic(err)
	}
	plotValuesBox.FillColor = color.RGBA{127, 188, 165, 1}
	plotValuesBox.MedianStyle = draw.LineStyle{
		Color:    color.Black,
		Width:    vg.Points(5),
		Dashes:   []vg.Length{},
		DashOffs: 0,
	}
	// Make a vertical box plot.
	plotValuesLabels, err := plotValuesBox.OutsideLabels(plotValues)
	if err != nil {
		panic(err)
	}
	boxPlot.Add(plotValuesBox, plotValuesLabels)
}

func SaveBoxPlot(boxPlot *plot.Plot, SaveLocation string) {
	err := boxPlot.Save(500, 500, SaveLocation)
	if err != nil {
		panic(err)
	}
}

func ResetBoxPlot(boxPlot *plot.Plot) *plot.Plot {
	title := boxPlot.Title.Text
	xLabel := boxPlot.X.Label.Text
	yLabel := boxPlot.Y.Label.Text

	return NewBoxPlot(title, xLabel, yLabel)
}
