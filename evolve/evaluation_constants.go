package evolve

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/skycoin/cx-evolves/cxexecutes/worker"
	workerclient "github.com/skycoin/cx-evolves/cxexecutes/worker/client"
	cxast "github.com/skycoin/cx/cx/ast"
)

// perByteEvaluation_Constants evolves with constants, 1 i32 input, 1 i32 output
func perByteEvaluation_Constants(ind *cxast.CXProgram, solPrototype *cxast.CXFunction, cfg *EvolveConfig) float64 {
	var total int32 = 0
	var tmp *cxast.CXProgram = cxast.PROGRAM
	cxast.PROGRAM = ind

	inpFullByteSize := 0
	for c := 0; c < len(solPrototype.Inputs); c++ {
		inpFullByteSize += solPrototype.Inputs[c].TotalSize
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

		var result worker.Result
		workerAddr := fmt.Sprintf(":%v", cfg.WorkerPortNum)
		workerclient.CallWorker(
			workerclient.CallWorkerConfig{
				Program:   ind,
				Input:     inps,
				OutputArg: solPrototype.Outputs[0],
			},
			workerAddr,
			&result,
		)

		data := int(binary.BigEndian.Uint32(result.Output))

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
	return float64(total)
}

func calculateConstantsScore(data, target int, cfg *EvolveConfig) int32 {
	// squared error (output-target)^2
	score := int32(math.Pow(float64(data-target), 2))

	// For 0-256 range benchmark
	if (cfg.ConstantsTarget != -1) && target == 256 && data < target && data > 0 {
		score = 0
	}

	return score
}
