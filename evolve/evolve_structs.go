package evolve

type EvolveConfig struct {
	TaskName string

	MazeHeight     int
	MazeWidth      int
	RandomMazeSize bool

	NumberOfRounds int

	ConstantsTarget int

	UpperRange int
	LowerRange int

	EpochLength int
	PlotFitness bool
	SaveAST     bool
	UseAntiLog2 bool

	WorkerPortNum    int
	WorkersAvailable int

	RandSeed int64
}
