package evolve

import (
	"encoding/binary"
	"fmt"

	"github.com/skycoin/cx-evolves/cmd/maze"
	cxcore "github.com/skycoin/cx/cx"
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

// perByteEvaluation for evolve with maze, 13 i2 input, 1 i32 output
func perByteEvaluation(ind *cxcore.CXProgram, solPrototype *cxcore.CXFunction, inputs [][]byte, outputs [][]byte) int {
	var move int
	var tmp *cxcore.CXProgram
	tmp = cxcore.PROGRAM
	cxcore.PROGRAM = ind

	inpFullByteSize := 0
	for c := 0; c < len(solPrototype.Inputs); c++ {
		inpFullByteSize += solPrototype.Inputs[c].TotalSize
	}

	// We'll store the `i`th inputs on `inps`.
	inps := make([]byte, inpFullByteSize)
	// `inpsOff` helps us keep track of what byte in `inps` we can write to.
	inpsOff := 0

	for c := 0; c < len(inputs); c++ {
		// The size of the input.
		inpSize := solPrototype.Inputs[c].TotalSize
		// The bytes representing the input.
		inp := inputs[c]

		// Copying the input `b`ytes.
		for b := 0; b < len(inp); b++ {
			inps[inpsOff+b] = inp[b]
		}

		// Updating offset.
		inpsOff += inpSize
	}

	// Injecting the input bytes `inps` to program `ind`.
	injectMainInputs(ind, inps)

	// Running program `ind`.
	ind.RunCompiled(0, nil)

	// Extracting outputs processed by `solPrototype`.
	simOuts := extractMainOutputs(ind, solPrototype)

	data := binary.BigEndian.Uint32(simOuts[0])
	move = int(data)

	cxcore.PROGRAM = tmp
	return move
}

// Evaluate Program as the Maze Player
func mazeMovesEvaluation(ind *cxcore.CXProgram, solPrototype *cxcore.CXFunction, mazeGame maze.Game) float64 {
	player := func(gameMove *maze.GameMove) maze.AgentInput {
		agentInput := maze.AgentInput{
			PassMazeData:             true,
			WallDistanceInputEnabled: true,
			AgentPositionEnabled:     true,
		}
		options := []int{maze.Up, maze.Down, maze.Left, maze.Right}

		move := perByteEvaluation(ind, solPrototype, EncodeParam(gameMove), nil)
		input := options[int(move)%len(options)]
		agentInput.Move = input

		return agentInput
	}

	moves := mazeGame.MazeGame(1, player)
	wg.Done()
	return float64(moves)
}

func EncodeParam(param *maze.GameMove) [][]byte {
	paramCount := 13
	inputs := make([][]byte, paramCount)

	inputs[0] = []byte(fmt.Sprint(int32(param.MoveCount)))
	inputs[1] = []byte(fmt.Sprint(int32(param.ErrorCode)))
	inputs[2] = []byte(param.ErrorMsg)
	inputs[3] = []byte(fmt.Sprint(int32(param.MazeData.Goal.X)))
	inputs[4] = []byte(fmt.Sprint(int32(param.MazeData.Goal.Y)))
	inputs[5] = []byte(fmt.Sprint(int32(param.AgentPosition.X)))
	inputs[6] = []byte(fmt.Sprint(int32(param.AgentPosition.Y)))
	inputs[7] = []byte(fmt.Sprint(int32(param.NumberOfSquaresLeftNorth)))
	inputs[8] = []byte(fmt.Sprint(int32(param.NumberOfSquaresLeftSouth)))
	inputs[9] = []byte(fmt.Sprint(int32(param.NumberOfSquaresLeftEast)))
	inputs[10] = []byte(fmt.Sprint(int32(param.NumberOfSquaresLeftWest)))
	inputs[11] = []byte(fmt.Sprint(int32(param.MazeData.Height)))
	inputs[12] = []byte(fmt.Sprint(int32(param.MazeData.Width)))

	return inputs
}
