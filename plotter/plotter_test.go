package plotter_test

import (
	"testing"

	cxplotter "github.com/skycoin/cx-evolves/plotter"
)

func TestSaveJSON(t *testing.T) {
	tests := []struct {
		scenario string
		filename string
		data     []cxplotter.PlotData
		wantErr  error
	}{
		{
			scenario: "One",
			filename: "test.json",
			data: []cxplotter.PlotData{
				{
					Title: "Maze",
					Data: []cxplotter.PlotDataPoints{
						{
							Generation: 1,
							Output:     []float64{1, 2, 3},
						},
						{
							Generation: 2,
							Output:     []float64{4, 5, 6},
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			err := cxplotter.AppendToFile("", "./"+tc.filename)
			if err != tc.wantErr {
				t.Errorf("want err %v, got %v", tc.wantErr, err)
			}

			err = cxplotter.AddTitleToJSON(tc.data[0].Title, "./"+tc.filename)
			if err != tc.wantErr {
				t.Errorf("want err %v, got %v", tc.wantErr, err)
			}

			err = cxplotter.AddDataToJSON(tc.data[0].Data[0], "./"+tc.filename)
			if err != tc.wantErr {
				t.Errorf("want err %v, got %v", tc.wantErr, err)
			}

			err = cxplotter.AddCommaToJSON("./" + tc.filename)
			if err != tc.wantErr {
				t.Errorf("want err %v, got %v", tc.wantErr, err)
			}

			err = cxplotter.AddDataToJSON(tc.data[0].Data[1], "./"+tc.filename)
			if err != tc.wantErr {
				t.Errorf("want err %v, got %v", tc.wantErr, err)
			}

			err = cxplotter.AddClosingToJSON("./" + tc.filename)
			if err != tc.wantErr {
				t.Errorf("want err %v, got %v", tc.wantErr, err)
			}
		})
	}
}
