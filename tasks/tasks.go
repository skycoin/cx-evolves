package tasks

import (
	cxast "github.com/skycoin/cx/cx/ast"
)

type TaskEvaluator func(ind *cxast.CXProgram, solPrototype EvolveSolProto, cfg TaskConfig) (float64, error)

func GetTaskEvaluator(name string, version int) TaskEvaluator {
	switch name {
	case "maze":
		switch version {
		case 1:
			return Maze_V1
		}
	case "constants":
		switch version {
		case 1:
			return Constants_V1
		}
	case "composites":
		switch version {
		case 1:
			return Composites_V1
		}
	case "evens":
		switch version {
		case 1:
			return Evens_V1
		}
	case "network_simulator":
		switch version {
		case 1:
			return NetworkSim_V1
		}
	case "odds":
		switch version {
		case 1:
			return Odds_V1
		}
	case "primes":
		switch version {
		case 1:
			return Primes_V1
		}
	case "range":
		switch version {
		case 1:
			return Range_V1
		}
	default:
		panic("task does not exist: check the spelling of task input")
	}

	return nil
}
