package mutation_test

import (
	"testing"

	cxmutation "github.com/skycoin/cx-evolves/mutation"
)

func TestProbability(t *testing.T) {
	tests := []struct {
		scenario        string
		numberOfOptions int
		numberOfSamples int
	}{
		{
			scenario:        "10 options, 100 samples",
			numberOfOptions: 10,
			numberOfSamples: 100,
		},
		{
			scenario:        "15 options, 200 samples",
			numberOfOptions: 15,
			numberOfSamples: 200,
		},
		{
			scenario:        "9 options, 201 samples",
			numberOfOptions: 9,
			numberOfSamples: 201,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.NewProbability(tc.numberOfOptions)

			samples := []float32{}
			for i := 0; i < tc.numberOfOptions; i++ {
				samples = append(samples, 0.00)
			}

			for i := 0; i < tc.numberOfSamples; i++ {
				samples[cxmutation.GetRandIndex()]++
			}

			// Total
			var total float32 = 0
			for i := 0; i < tc.numberOfOptions; i++ {
				total += samples[i]
			}

			if (total / float32(tc.numberOfSamples)) != 1.00 {
				t.Errorf("want total 1, got %v", total)
			}
		})
	}
}
