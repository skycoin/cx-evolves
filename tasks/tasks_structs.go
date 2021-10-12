package tasks

import "github.com/skycoin/cx/cx/types"

type TaskConfig struct {
	NumberOfRounds int

	ConstantsTarget int

	UpperRange int
	LowerRange int

	RandSeed       int64
	MazeWidth      int
	MazeHeight     int
	RandomMazeSize bool
}

// Solution Prototype info.
type EvolveSolProto struct {
	InpsSize  []types.Pointer
	OutOffset types.Pointer
	OutSize   types.Pointer
}
