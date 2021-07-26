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

	RandomSearch bool

	RandSeed int64
}

type GraphCfg struct {
	Output        []float64
	FittestIndex  *int
	HistoValues   *[]float64
	MostFit       *[]float64
	AverageValues *[]float64
	EvolveCfg     *EvolveConfig
	PopuSize      int
}
