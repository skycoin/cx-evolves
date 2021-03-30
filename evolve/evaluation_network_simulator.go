package evolve

import (
	"math/rand"
	"time"

	cxast "github.com/skycoin/cx/cx/ast"
	cxexecute "github.com/skycoin/cx/cx/execute"
)

// perByteEvaluation for evolve with network sim, 1 i32 input, 1 i32 output
func perByteEvaluation_NetworkSim(ind *cxast.CXProgram, solPrototype *cxast.CXFunction, numberOfRounds int) int64 {
	var score int = 0
	for rounds := 0; rounds < numberOfRounds; rounds++ {
		// Generate random Input
		rand.Seed(time.Now().Unix())
		input := toByteArray(int32(rand.Int()))

		// Get output from transmitter
		transmitterOutput := perByteEvaluation_NetworkSim_Transmitter(ind, solPrototype, input)

		// Input noise here

		// Get output from receiver
		receiverOutput := perByteEvaluation_NetworkSim_Receiver(ind, solPrototype, transmitterOutput)

		// Get score by counting number of diff bits between generated input and receiverOutput
		score = score + countDifferentBits(input, receiverOutput)
	}

	return int64(score)
}

// perByteEvaluation for evolve with network sim transmitter, 1 i32 input, 1 i32 output
func perByteEvaluation_NetworkSim_Transmitter(ind *cxast.CXProgram, solPrototype *cxast.CXFunction, input []byte) []byte {
	var tmp *cxast.CXProgram = cxast.PROGRAM
	cxast.PROGRAM = ind

	inpFullByteSize := 0
	for c := 0; c < len(solPrototype.Inputs); c++ {
		inpFullByteSize += solPrototype.Inputs[c].TotalSize
	}

	// We'll store the `i`th inputs on `inps`.
	inps := make([]byte, inpFullByteSize)

	inp := input

	// Copying the input `b`ytes.
	for b := 0; b < len(inp); b++ {
		inps[b] = inp[b]
	}

	// Injecting the input bytes `inps` to program `ind`.
	injectMainInputs(ind, inps)

	// Running program `ind`.
	cxexecute.RunCompiled(ind, 0, nil)

	// Extracting outputs processed by `solPrototype`.
	simOuts := extractMainOutputs(ind, solPrototype)
	data := simOuts[0]

	cxast.PROGRAM = tmp
	return data
}

// perByteEvaluation for evolve with network sim receiver, 1 i32 input, 1 i32 output
func perByteEvaluation_NetworkSim_Receiver(ind *cxast.CXProgram, solPrototype *cxast.CXFunction, input []byte) []byte {
	var tmp *cxast.CXProgram = cxast.PROGRAM
	cxast.PROGRAM = ind

	inpFullByteSize := 0
	for c := 0; c < len(solPrototype.Inputs); c++ {
		inpFullByteSize += solPrototype.Inputs[c].TotalSize
	}

	// We'll store the `i`th inputs on `inps`.
	inps := make([]byte, inpFullByteSize)

	inp := input

	// Copying the input `b`ytes.
	for b := 0; b < len(inp); b++ {
		inps[b] = inp[b]
	}

	// Injecting the input bytes `inps` to program `ind`.
	injectMainInputs(ind, inps)

	// Running program `ind`.
	cxexecute.RunCompiled(ind, 0, nil)

	// Extracting outputs processed by `solPrototype`.
	simOuts := extractMainOutputs(ind, solPrototype)

	data := simOuts[0]

	cxast.PROGRAM = tmp
	return data
}

func countDifferentBits(a []byte, b []byte) int {
	var count int

	for i, val := range a {
		bitPos := 1
		for z := 0; z < 8; z++ {
			if int(val)&bitPos != int(b[i])&bitPos {
				count++
			}
			bitPos = bitPos << 1
		}
	}
	return count
}
