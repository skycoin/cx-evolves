package tasks

import (
	cxast "github.com/skycoin/cx/cx/ast"
)

type TaskEvaluator func(ind *cxast.CXProgram, solPrototype EvolveSolProto, cfg TaskConfig) (float64, error)

const (
	_ = iota
	Task_Maze
	Task_Constants
	Task_Composites
	Task_Evens
	Task_Network_Simulator
	Task_Odds
	Task_Primes
	Task_Range
)

var (
	TaskName = map[int]string{
		Task_Maze:              "maze",
		Task_Constants:         "constants",
		Task_Composites:        "composites",
		Task_Evens:             "evens",
		Task_Network_Simulator: "network_simulator",
		Task_Odds:              "odds",
		Task_Primes:            "primes",
		Task_Range:             "range",
	}
)

func GetTaskEvaluator(name string, version int) TaskEvaluator {
	switch name {
	case TaskName[Task_Maze]:
		switch version {
		case 1:
			return Maze_V1
		}
	case TaskName[Task_Constants]:
		switch version {
		case 1:
			return Constants_V1
		}
	case TaskName[Task_Composites]:
		switch version {
		case 1:
			return Composites_V1
		}
	case TaskName[Task_Evens]:
		switch version {
		case 1:
			return Evens_V1
		case 2:
			return Evens_V2
		}
	case TaskName[Task_Network_Simulator]:
		switch version {
		case 1:
			return NetworkSim_V1
		}
	case TaskName[Task_Odds]:
		switch version {
		case 1:
			return Odds_V1
		}
	case TaskName[Task_Primes]:
		switch version {
		case 1:
			return Primes_V1
		}
	case TaskName[Task_Range]:
		switch version {
		case 1:
			return Range_V1
		}
	default:
		panic("task does not exist: check the spelling of task input")
	}

	return nil
}

func IsMazeTask(taskName string) bool {
	return taskName == TaskName[Task_Maze]
}

func IsConstantsTask(taskName string) bool {
	return taskName == TaskName[Task_Constants]
}

func IsCompositesTask(taskName string) bool {
	return taskName == TaskName[Task_Composites]
}

func IsEvensTask(taskName string) bool {
	return taskName == TaskName[Task_Evens]
}

func IsNetworkSimulatorTask(taskName string) bool {
	return taskName == TaskName[Task_Network_Simulator]
}

func IsOddsTask(taskName string) bool {
	return taskName == TaskName[Task_Odds]
}

func IsPrimesTask(taskName string) bool {
	return taskName == TaskName[Task_Primes]
}

func IsRangeTask(taskName string) bool {
	return taskName == TaskName[Task_Range]
}
