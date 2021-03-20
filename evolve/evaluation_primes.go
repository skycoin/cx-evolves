package evolve

import (
	"encoding/binary"
	"math/big"
	"math/rand"
	"time"

	cxcore "github.com/skycoin/cx/cx"
)

// perByteEvaluation for evolve with prime numbers, 1 i32 input, 1 i32 output
func perByteEvaluation_Primes(ind *cxcore.CXProgram, solPrototype *cxcore.CXFunction, numberOfRounds int) int64 {
	var points int64 = 0
	var tmp *cxcore.CXProgram = cxcore.PROGRAM
	cxcore.PROGRAM = ind

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
		ind.RunCompiled(0, nil)

		// Extracting outputs processed by `solPrototype`.
		simOuts := extractMainOutputs(ind, solPrototype)

		data := binary.BigEndian.Uint32(simOuts[0])

		// If not a prime, add 1 to total points
		if !big.NewInt(int64(data)).ProbablyPrime(4) {
			points++
		}
	}

	cxcore.PROGRAM = tmp
	return points
}
