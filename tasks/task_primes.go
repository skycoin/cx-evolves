package tasks

import (
	"encoding/binary"
	"math/big"
	"math/rand"
	"time"

	cxast "github.com/skycoin/cx/cx/ast"
	cxexecute "github.com/skycoin/cx/cx/execute"
	"github.com/skycoin/cx/cx/types"
)

// Primes_V1 for evolve with prime numbers, 1 i32 input, 1 i32 output
func Primes_V1(ind *cxast.CXProgram, solPrototype EvolveSolProto, cfg TaskConfig) (float64, error) {
	var points int64 = 0
	var tmp *cxast.CXProgram = cxast.PROGRAM
	cxast.PROGRAM = ind

	var inpFullByteSize types.Pointer = 0
	for c := 0; c < len(solPrototype.InpsSize); c++ {
		inpFullByteSize += solPrototype.InpsSize[c]
	}

	// We'll store the `i`th inputs on `inps`.
	inps := make([]byte, inpFullByteSize)
	for round := 0; round < cfg.NumberOfRounds; round++ {
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

		injectMainInputs(ind, inps)
		err := cxexecute.RunCompiled(ind, 0, nil)
		if err != nil {
			panic(err)
		}

		byteOut := ind.Memory[solPrototype.OutOffset : solPrototype.OutOffset+solPrototype.OutSize]
		data := int(binary.BigEndian.Uint32(byteOut))
		// If not a prime, add 1 to total points
		if !big.NewInt(int64(data)).ProbablyPrime(4) {
			points++
		}
	}

	cxast.PROGRAM = tmp
	return float64(points), nil
}
