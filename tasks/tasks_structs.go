package tasks

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
	InpsSize  []int
	OutOffset int
	OutSize   int
}
