package evolve

import (
	"fmt"
	"os"
	"time"

	crypto_rand "crypto/rand"
	"encoding/binary"

	"github.com/skycoin/cx-evolves/tasks"
)

func makeDirectory(cfg *EvolveConfig) string {
	var dir string

	if cfg.PlotFitness || cfg.SaveAST {
		name := getBenchmarkName(cfg)
		dir = fmt.Sprintf("./Results/%v-%v/", time.Now().Unix(), name)

		// create directory
		_ = os.Mkdir(dir, 0700)

		if cfg.SaveAST {
			_ = os.Mkdir(dir+"AST/", 0700)
		}
	}
	return dir
}

func getBenchmarkName(cfg *EvolveConfig) string {
	var name string
	if cfg.MazeBenchmark {
		// Maze-2x2
		mazeSize := fmt.Sprintf("%vx%v", cfg.MazeWidth, cfg.MazeHeight)
		if cfg.RandomMazeSize {
			mazeSize = "random"
		}

		name = fmt.Sprintf("%v-%v", "Maze", mazeSize)
	}

	if cfg.ConstantsBenchmark {
		name = "Constants"
		if cfg.ConstantsTarget != -1 {
			name = fmt.Sprintf("%v-%v", name, cfg.ConstantsTarget)
		}
	}

	if cfg.EvensBenchmark {
		name = "Evens"
	}

	if cfg.OddsBenchmark {
		name = "Odds"
	}

	if cfg.PrimesBenchmark {
		name = "Primes"
	}

	if cfg.CompositesBenchmark {
		name = "Composites"
	}

	if cfg.RangeBenchmark {
		name = fmt.Sprintf("%v-%v-%v", "Range", cfg.LowerRange, cfg.UpperRange)
	}

	if cfg.NetworkSimBenchmark {
		name = "NetworkSim"
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

func setTaskParams(cfg EvolveConfig) tasks.TaskConfig {
	// Set Task Cfg
	taskCfg := tasks.TaskConfig{
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
