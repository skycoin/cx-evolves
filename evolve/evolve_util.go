package evolve

import (
	"fmt"
	"os"
	"time"

	crypto_rand "crypto/rand"
	"encoding/binary"

	cxtasks "github.com/skycoin/cx-evolves/tasks"
)

func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

func makeDirectory(cfg *EvolveConfig) string {
	var dir string

	if cfg.PlotFitness || cfg.SaveAST {
		name := getBenchmarkName(cfg)
		dir = fmt.Sprintf("./Results/%v-%v/", time.Now().Unix(), name)

		// create directory
		_ = os.Mkdir(dir, 0700)
		_ = os.Mkdir(dir+"AST/", 0700)

	}
	return dir
}

func getBenchmarkName(cfg *EvolveConfig) string {
	var name string
	if cxtasks.IsMazeTask(cfg.TaskName) {
		// Maze-2x2
		mazeSize := fmt.Sprintf("%vx%v", cfg.MazeWidth, cfg.MazeHeight)
		if cfg.RandomMazeSize {
			mazeSize = "random"
		}

		name = fmt.Sprintf("%v-%v", "Maze", mazeSize)
	}

	if cxtasks.IsConstantsTask(cfg.TaskName) {
		name = "Constants"
		if cfg.ConstantsTarget != -1 {
			name = fmt.Sprintf("%v-%v", name, cfg.ConstantsTarget)
		}
	}

	if cxtasks.IsEvensTask(cfg.TaskName) {
		name = "Evens"
	}

	if cxtasks.IsOddsTask(cfg.TaskName) {
		name = "Odds"
	}

	if cxtasks.IsPrimesTask(cfg.TaskName) {
		name = "Primes"
	}

	if cxtasks.IsCompositesTask(cfg.TaskName) {
		name = "Composites"
	}

	if cxtasks.IsRangeTask(cfg.TaskName) {
		name = fmt.Sprintf("%v-%v-%v", "Range", cfg.LowerRange, cfg.UpperRange)
	}

	if cxtasks.IsNetworkSimulatorTask(cfg.TaskName) {
		name = "NetworkSim"
	}

	if cfg.Version != 1 {
		name = name + fmt.Sprintf("-%v", cfg.Version)
	}
	return name
}

func setEpochLength(cfg *EvolveConfig) {
	if cfg.EpochLength == 0 {
		cfg.EpochLength = 1
	}
}

func generateNewSeed(generationCount int, cfg EvolveConfig) int64 {
	if generationCount%cfg.EpochLength == 0 || generationCount == 0 {
		var b [8]byte
		_, err := crypto_rand.Read(b[:])
		if err != nil {
			panic("cannot seed math/rand package with cryptographically secure random number generator")
		}
		return int64(binary.LittleEndian.Uint64(b[:]))
	}
	return cfg.RandSeed
}

func setTaskParams(cfg EvolveConfig) cxtasks.TaskConfig {
	// Set Task Cfg
	taskCfg := cxtasks.TaskConfig{
		NumberOfRounds:  cfg.NumberOfRounds,
		ConstantsTarget: cfg.ConstantsTarget,
		UpperRange:      cfg.UpperRange,
		LowerRange:      cfg.LowerRange,
		RandSeed:        cfg.RandSeed,
		MazeWidth:       cfg.MazeWidth,
		MazeHeight:      cfg.MazeHeight,
		RandomMazeSize:  cfg.RandomMazeSize,
	}

	return taskCfg
}
