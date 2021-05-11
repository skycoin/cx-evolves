package client

import (
	"time"

	"github.com/henrylee2cn/erpc/v6"
	"github.com/skycoin/cx-evolves/cxexecutes/worker"
	"github.com/skycoin/cx-evolves/tasks"
	cxast "github.com/skycoin/cx/cx/ast"
)

type CallWorkerConfig struct {
	Task    string
	Version int
	Program *cxast.CXProgram
	SolProt *cxast.CXFunction
	TaskCfg tasks.TaskConfig
}

func CallWorker(cWorker CallWorkerConfig, workerAddr string, result *worker.Result) {
	erpc.SetLoggerLevel("OFF")()
	cli := erpc.NewPeer(erpc.PeerConfig{RedialTimes: -1, RedialInterval: time.Second})
	defer cli.Close()
	cli.SetTLSConfig(erpc.GenerateTLSConfigForClient())
	cli.RoutePush(new(Push))

	sess, stat := cli.Dial(workerAddr)
	if !stat.OK() {
		erpc.Fatalf("%v", stat)
	}
	defer sess.Close()

	// Extract solution prototype info
	solProto := tasks.EvolveSolProto{
		OutOffset: cWorker.SolProt.Outputs[0].Offset,
		OutSize:   cWorker.SolProt.Outputs[0].TotalSize,
	}
	solProto.InpsSize = make([]int, len(cWorker.SolProt.Inputs))
	for i := 0; i < len(cWorker.SolProt.Inputs); i++ {
		solProto.InpsSize[i] = cWorker.SolProt.Inputs[i].TotalSize
	}

	// Set worker args
	args := &worker.Args{
		Task:    cWorker.Task,
		Version: cWorker.Version,
		Program: cxast.SerializeCXProgramV2(cWorker.Program, true, false),
		SolProt: solProto,
		Cfg:     cWorker.TaskCfg,
	}

	stat = sess.Call(
		worker.RunProgram,
		args,
		&result,
	).Status()

	if !stat.OK() {
		erpc.Fatalf("%v", stat)
	}
}

// Push push handler
type Push struct {
	erpc.PushCtx
}

// Push handles '/push/status' message
func (p *Push) Status(arg *string) *erpc.Status {
	erpc.Printf("%s", *arg)
	return nil
}
