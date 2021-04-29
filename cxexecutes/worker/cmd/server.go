package main

import (
	"flag"
	"sync"

	"github.com/henrylee2cn/erpc/v6"

	"github.com/skycoin/cx-evolves/cxexecutes/worker"
	cxopcodes "github.com/skycoin/cx/cx/opcodes"
)

var workers int

func init() {
	flag.IntVar(&workers, "workers", 1, "number of workers")
}

func main() {
	// runtime.GOMAXPROCS(1)
	flag.Parse()
	cxopcodes.RegisterOpcodes()
	deployWorker(workers)
}

func deployWorker(workers int) {
	// graceful
	go erpc.GraceSignal()

	wg := new(sync.WaitGroup)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		portNumber := worker.BasePortNumber + i
		go func() {
			// server peer
			srv := erpc.NewPeer(erpc.PeerConfig{
				CountTime:   true,
				ListenPort:  uint16(portNumber),
				PrintDetail: false,
			})
			srv.SetTLSConfig(erpc.GenerateTLSConfigForServer())

			// router
			srv.RouteCall(new(worker.ProgramWorker))

			// listen and serve
			err := srv.ListenAndServe()
			if err != nil {
				panic(err)
			}
			wg.Done()
		}()

		erpc.GetLogger().Printf("listen and serve: %v", portNumber)
		erpc.FlushLogger()
	}
	wg.Wait()
}
