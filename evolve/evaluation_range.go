package evolve

import (
	"encoding/binary"
	"math/rand"
	"time"

	cxast "github.com/skycoin/cx/cx/ast"
	cxexecute "github.com/skycoin/cx/cx/execute"
)

// perByteEvaluation for evolve with range, 1 i32 input, 1 i32 output
func perByteEvaluation_Range(ind *cxast.CXProgram, solPrototype *cxast.CXFunction, numberOfRounds, upperRange, lowerRange int) float64 {
	var points int64 = 0
	var tmp *cxast.CXProgram = cxast.PROGRAM
	cxast.PROGRAM = ind

	inpFullByteSize := 0
	for c := 0; c < len(solPrototype.Inputs); c++ {
		inpFullByteSize += solPrototype.Inputs[c].TotalSize
	}

	// We'll store the `i`th inputs on `inps`.
	inps := make([]byte, inpFullByteSize)
	for round := 0; round < numberOfRounds; round++ {
		rand.Seed(time.Now().Unix())
		in := round
		for in == 0 {
			in = rand.Int()
		}

		inp := toByteArray(int32(in))

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

		data := binary.BigEndian.Uint32(simOuts[0])

		// if not within range, add 1 to total points
		if data > uint32(upperRange) || data < uint32(lowerRange) {
			points++
		}
	}

	cxast.PROGRAM = tmp
	return float64(points)
}
