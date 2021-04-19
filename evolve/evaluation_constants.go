package evolve

import (
	"encoding/binary"
	"math"
	"math/rand"
	"time"

	cxast "github.com/skycoin/cx/cx/ast"
	cxexecute "github.com/skycoin/cx/cx/execute"
)

// perByteEvaluation for evolve with constants, 1 i32 input, 1 i32 output
func perByteEvaluation_Constants(ind *cxast.CXProgram, solPrototype *cxast.CXFunction, numberOfRounds int) float64 {
	var total float64 = 0
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

		// Give rnadom input for first round
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

		// squared error (output-target)^2
		score := math.Pow(float64(int(data)-round), 2)

		// Check if overflowed
		if total+score < total {
			total = math.MaxFloat64
		} else {
			total = total + score
		}
	}

	cxast.PROGRAM = tmp
	return total
}
