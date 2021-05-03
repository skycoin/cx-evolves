package evolve

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
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
		name = "Range"
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

func toByteArray(i int32) []byte {
	arr := make([]byte, 4)
	binary.BigEndian.PutUint32(arr, uint32(i))
	return arr
}
