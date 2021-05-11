package evolve

type EvolveConfig struct {
	MazeBenchmark       bool
	ConstantsBenchmark  bool
	EvensBenchmark      bool
	OddsBenchmark       bool
	PrimesBenchmark     bool
	CompositesBenchmark bool
	RangeBenchmark      bool
	NetworkSimBenchmark bool

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
