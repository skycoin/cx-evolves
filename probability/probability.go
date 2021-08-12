package probability

import (
	"math/rand"
)

func GetRandIndex(cdf []float32) int {
	r := rand.Float32()

	idx := 0
	for r > cdf[idx] {
		idx++
	}
	return idx
}

func NewProbability(numberOfOptions int) []float32 {
	// Set density equally.
	dist := float32(1) / float32(numberOfOptions)

	pdf := []float32{} // probability density function
	cdf := []float32{} // cummulative distribution function
	for i := 0; i < numberOfOptions; i++ {
		pdf = append(pdf, dist)
		cdf = append(cdf, 0.00)
	}

	// Get cdf
	cdf[0] = pdf[0]
	for i := 1; i < numberOfOptions; i++ {
		cdf[i] = cdf[i-1] + pdf[i]
	}
	// fmt.Println(cdf)
	return cdf
}
