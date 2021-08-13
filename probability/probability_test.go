package probability_test

import (
	"testing"

	cxprobability "github.com/skycoin/cx-evolves/probability"
)

func TestProbability(t *testing.T) {
	tests := []struct {
		scenario           string
		percentProbability []float32
		numberOfSamples    int
	}{
		{
			scenario:           "10 options, 100 samples",
			percentProbability: []float32{10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
			numberOfSamples:    100,
		},
		{
			scenario:           "15 options, 200 samples",
			percentProbability: []float32{20, 20, 20, 20, 10, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			numberOfSamples:    200,
		},
		{
			scenario:           "9 options, 201 samples",
			percentProbability: []float32{10, 10, 10, 10, 10, 25, 5, 10, 10},
			numberOfSamples:    201,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cdf := cxprobability.NewProbability(tc.percentProbability)

			samples := []float32{}
			for i := 0; i < len(tc.percentProbability); i++ {
				samples = append(samples, 0.00)
			}

			for i := 0; i < tc.numberOfSamples; i++ {
				samples[cxprobability.GetRandIndex(cdf)]++
			}

			// Total
			var total float32 = 0
			for i := 0; i < len(tc.percentProbability); i++ {
				total += samples[i]
			}

			if (total / float32(tc.numberOfSamples)) != 1.00 {
				t.Errorf("want total 1, got %v", total)
			}
		})
	}
}
