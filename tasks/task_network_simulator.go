package tasks

import (
	"math/rand"
	"time"

	cxast "github.com/skycoin/cx/cx/ast"
	cxexecute "github.com/skycoin/cx/cx/execute"
	"github.com/skycoin/cx/cx/types"
)

// NetworkSim_V1 for evolve with network sim, 1 i32 input, 1 i32 output
func NetworkSim_V1(ind *cxast.CXProgram, solPrototype EvolveSolProto, cfg TaskConfig) (float64, error) {
	var score int = 0
	for rounds := 0; rounds < cfg.NumberOfRounds; rounds++ {
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

	return float64(score), nil
}

// perByteEvaluation for evolve with network sim transmitter, 1 i32 input, 1 i32 output
func perByteEvaluation_NetworkSim_Transmitter(ind *cxast.CXProgram, solPrototype EvolveSolProto, input []byte) []byte {
	var tmp *cxast.CXProgram = cxast.PROGRAM
	cxast.PROGRAM = ind

	var inpFullByteSize types.Pointer = 0
	for c := 0; c < len(solPrototype.InpsSize); c++ {
		inpFullByteSize += solPrototype.InpsSize[c]
	}

	// We'll store the `i`th inputs on `inps`.
	inps := make([]byte, inpFullByteSize)

	inp := input

	// Copying the input `b`ytes.
	for b := 0; b < len(inp); b++ {
		inps[b] = inp[b]
	}

	injectMainInputs(ind, inps)
	err := cxexecute.RunCompiled(ind, 0, nil)
	if err != nil {
		panic(err)
	}

	byteOut := ind.Memory[solPrototype.OutOffset : solPrototype.OutOffset+solPrototype.OutSize]
	data := byteOut

	cxast.PROGRAM = tmp
	return data
}

// perByteEvaluation for evolve with network sim receiver, 1 i32 input, 1 i32 output
func perByteEvaluation_NetworkSim_Receiver(ind *cxast.CXProgram, solPrototype EvolveSolProto, input []byte) []byte {
	var tmp *cxast.CXProgram = cxast.PROGRAM
	cxast.PROGRAM = ind

	var inpFullByteSize types.Pointer = 0
	for c := 0; c < len(solPrototype.InpsSize); c++ {
		inpFullByteSize += solPrototype.InpsSize[c]
	}

	// We'll store the `i`th inputs on `inps`.
	inps := make([]byte, inpFullByteSize)

	inp := input

	// Copying the input `b`ytes.
	for b := 0; b < len(inp); b++ {
		inps[b] = inp[b]
	}

	injectMainInputs(ind, inps)
	err := cxexecute.RunCompiled(ind, 0, nil)
	if err != nil {
		panic(err)
	}

	byteOut := ind.Memory[solPrototype.OutOffset : solPrototype.OutOffset+solPrototype.OutSize]
	data := byteOut

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
