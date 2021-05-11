package tasks

import (
	"encoding/binary"
	"math/rand"
	"time"

	cxast "github.com/skycoin/cx/cx/ast"
	cxexecute "github.com/skycoin/cx/cx/execute"
)

// Odds_V1 for evolve with odd numbers, 1 i32 input, 1 i32 output
func Odds_V1(ind *cxast.CXProgram, solPrototype EvolveSolProto, cfg TaskConfig) (float64, error) {
	var points int64 = 0
	var tmp *cxast.CXProgram = cxast.PROGRAM
	cxast.PROGRAM = ind

	inpFullByteSize := 0
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
		// If not odd, add 1 to total points
		if int64(data)%2 == 0 {
			points++
		}
	}

	cxast.PROGRAM = tmp
	return float64(points), nil
}