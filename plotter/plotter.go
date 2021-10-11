package plotter

import (
	"encoding/json"
	"os"

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

func AddTitleToJSON(title, saveDirectory string) error {
	var plotData PlotData
	plotData.Title = title

	data, err := json.Marshal(plotData)
	if err != nil {
		return err
	}

	data = append(data[0:len(data)-5], []byte("[")...)
	err = AppendToFile(string(data), saveDirectory)
	if err != nil {
		return err
	}

	return nil
}

func AddDataToJSON(dataPoints PlotDataPoints, saveDirectory string) error {
	data, err := json.Marshal(dataPoints)
	if err != nil {
		return err
	}

	dataStr := "{" + string(data[1:len(data)-1]) + "}"
	err = AppendToFile(string(dataStr), saveDirectory)
	if err != nil {
		return err
	}

	return nil
}

func AddCommaToJSON(saveDirectory string) error {
	data := ","
	err := AppendToFile(string(data), saveDirectory)
	if err != nil {
		return err
	}
	return nil
}

func AddClosingToJSON(saveDirectory string) error {
	data := "]}"
	err := AppendToFile(string(data), saveDirectory)
	if err != nil {
		return err
	}
	return nil
}

func AppendToFile(data string, saveDir string) error {
	f, err := os.OpenFile(saveDir, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(data); err != nil {
		return err
	}

	return nil
}
