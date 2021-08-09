package plotter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gonum.org/v1/plot/plotter"
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

type PlotData struct {
	Title string           `json:"title"`
	Data  []PlotDataPoints `json:"data"`
}

type PlotDataPoints struct {
	Generation int       `json:"generation"`
	Output     []float64 `json:"output"`
}

func SavePlotGraphDataToJSON(data PlotData, filename string) error {
	file, _ := json.MarshalIndent(data, "", " ")

	os.Remove(fmt.Sprintf("%s.json", filename))
	_ = ioutil.WriteFile(fmt.Sprintf("%s.json", filename), file, 0777)
	return nil
}
