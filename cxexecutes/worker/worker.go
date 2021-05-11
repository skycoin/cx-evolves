package worker

import (
	"fmt"

	"github.com/henrylee2cn/erpc/v6"
	cxtasks "github.com/skycoin/cx-evolves/tasks"
	cxast "github.com/skycoin/cx/cx/ast"
)

const (
	RunProgram     = "/program_worker/run_task_evaluator"
	BasePortNumber = 9090
)

type Args struct {
	Task    string
	Version int
	Program []byte
	Cfg     cxtasks.TaskConfig
	SolProt cxtasks.EvolveSolProto
}

type Result struct {
	Output float64
}
type ProgramWorker struct {
	erpc.CallCtx
}

func (pw *ProgramWorker) RunTaskEvaluator(args *Args) (Result, *erpc.Status) {
	prgrmInBytes := args.Program
	prgrm := cxast.DeserializeCXProgramV2(prgrmInBytes, false)
	prgrm.Memory = cxast.MakeProgram().Memory

	evaluate := cxtasks.GetTaskEvaluator(args.Task, args.Version)
	output, err := evaluate(prgrm, args.SolProt, args.Cfg)
	if err != nil {
		return Result{}, erpc.NewStatus(1, fmt.Sprintf("%v", err))
	}
	res := Result{
		Output: output,
	}

	return res, nil
}

func GetAvailableWorkers(numberOfAvailableWorkers int) []int {
	var workersAddr []int
	for i := 0; i < numberOfAvailableWorkers; i++ {
		workersAddr = append(workersAddr, BasePortNumber+i)
	}
	return workersAddr
}
