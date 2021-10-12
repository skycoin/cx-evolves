package tasks

import (
	"encoding/binary"
	"math"
	"math/rand"
	"time"

	cxast "github.com/skycoin/cx/cx/ast"
	cxexecute "github.com/skycoin/cx/cx/execute"
	"github.com/skycoin/cx/cx/types"
)

// Constants_V1 evolves with constants, 1 i32 input, 1 i32 output
func Constants_V1(ind *cxast.CXProgram, solPrototype EvolveSolProto, cfg TaskConfig) (float64, error) {
	var total int32 = 0
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

		// Give random input for first round
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

		target := round
		consTarg := cfg.ConstantsTarget
		if consTarg != -1 {
			target = consTarg
		}
		score := calculateConstantsScore(data, target, cfg)

		// Check if overflowed
		if total+score < total {
			total = math.MaxInt32
		} else {
			total = total + score
		}
	}

	cxast.PROGRAM = tmp
	return float64(total), nil
}

func calculateConstantsScore(data, target int, cfg TaskConfig) int32 {
	// squared error (output-target)^2
	score := int32(math.Pow(float64(data-target), 2))

	// For 0-256 range benchmark
	if (cfg.ConstantsTarget != -1) && target == 256 && data < target && data > 0 {
		score = 0
	}

	return score
}
