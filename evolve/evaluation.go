package evolve

import (
	"log"
	"math"

	cxast "github.com/skycoin/cx/cx/ast"
)

// Codes associated to each of the mutation functions.
const (
	EvaluationPerByte = iota
)

// getCrossoverFn returns the crossover function associated to `mutationCode`.
// func (pop *Population) getEvaluationFn(evaluationCode int) func(*cxcore.CXFunction, *cxcore.CXFunction) (*cxcore.CXFunction, *cxcore.CXFunction) {
// 	switch evaluationCode {
// 	case perByteEvaluation:
// 		return randomMutation
// 	}
// }

// Original perByteEvaluation
// perByteEvaluation ...
// func perByteEvaluation(ind *cxcore.CXProgram, solPrototype *cxcore.CXFunction, inputs [][]byte, outputs [][]byte) float64 {
// 	var tmp *cxcore.CXProgram
// 	tmp = cxcore.PROGRAM
// 	cxcore.PROGRAM = ind

// 	// TODO: We're calculating the error in here.
// 	/// Migrate to functions when we have other fitness functions.

// 	inpFullByteSize := 0
// 	for c := 0; c < len(solPrototype.Inputs); c++ {
// 		inpFullByteSize += solPrototype.Inputs[c].TotalSize
// 	}

// 	var sum float64

// 	// `numElts` represents the number of elements per input array calculated by the inputs function.
// 	// All the inputs represent arrays of the same size, regardless of element type
// 	// (for example, 10 `i32`s and 10 `f64`s). So it is safe to assume that
// 	// looping over `inputs[0]` will make us loop over all `inputs` from 1 to N.
// 	numElts := len(inputs[0]) / solPrototype.Inputs[0].TotalSize

// 	for i := 0; i < numElts; {
// 		// Now we'll loop over each of the `inputs`.
// 		/// We want to extract the `i`th element from each of the `inputs`.
// 		/// For example, if we are sending two arrays (inputs), a [10]i32 and a [10]f64,
// 		/// we want to extract the `i`th i32 and the `i`th f64 and send those two inputs to the solution.

// 		// We'll store the `i`th inputs on `inps`.
// 		inps := make([]byte, inpFullByteSize)
// 		// `inpsOff` helps us keep track of what byte in `inps` we can write to.
// 		inpsOff := 0

// 		for c := 0; c < len(inputs); c++ {
// 			// The size of the input.
// 			inpSize := solPrototype.Inputs[c].TotalSize
// 			// The bytes representing the input.
// 			inp := inputs[c][inpSize*i:inpSize*(i+1)]

// 			// Copying the input `b`ytes.
// 			for b := 0; b < len(inp); b++ {
// 				inps[inpsOff+b] = inp[b]
// 			}

// 			// Updating offset.
// 			inpsOff += inpSize
// 		}

// 		// Updating how many `b`ytes we read from `inputs[0]`.
// 		// b += solPrototype.Inputs[0].TotalSize

// 		// Injecting the input bytes `inps` to program `ind`.
// 		injectMainInputs(ind, inps)

// 		// Running program `ind`.
// 		ind.RunCompiled(0, nil)

// 		// Extracting outputs processed by `solPrototype`.
// 		simOuts := extractMainOutputs(ind, solPrototype)

// 		// Comparing real vs simulated outputs (error).
// 		for o := 0; o < len(solPrototype.Outputs); o++ {
// 			outSize := solPrototype.Outputs[o].TotalSize
// 			for b := 0; b < len(simOuts[o]); b++ {
// 				// Comparing byte by byte.
// 				sum += math.Abs(float64(outputs[o][i*outSize+b] - simOuts[o][b]))
// 			}
// 		}
// 		i++
// 	}

// 	cxcore.PROGRAM = tmp
// 	return sum
// }

func SelectRankCutoff(parent1Fitness, parent2Fitness float64, solProt, child1, child2, parent1, parent2 *cxast.CXFunction, cfg EvolveConfig, currPortNum int, sPrgrm []byte) error {
	// Create cx program for child1
	child1Ind := cxast.DeserializeCXProgramV2(sPrgrm, true)
	child2Ind := cxast.DeserializeCXProgramV2(sPrgrm, true)

	replaceSolution(child1Ind, solProt.Name, child1)
	replaceSolution(child2Ind, solProt.Name, child2)

	cfg.WorkerPortNum = currPortNum

	child1Fitness, err := RunBenchmark(child1Ind, solProt, cfg)
	if err != nil {
		log.Printf("error run benchmark: %v", err)
		child1Fitness = float64(math.MaxInt32)
	}

	child2Fitness, err := RunBenchmark(child1Ind, solProt, cfg)
	if err != nil {
		log.Printf("error run benchmark: %v", err)
		child2Fitness = float64(math.MaxInt32)
	}

	if child1Fitness < parent1Fitness || child1Fitness < parent2Fitness {
		// Child is pruned, not added to population
		child1 = parent1
	}

	if child2Fitness < parent1Fitness || child2Fitness < parent2Fitness {
		// Child is pruned, not added to population
		child2 = parent2
	}

	_ = child1
	_ = child2

	return nil
}
