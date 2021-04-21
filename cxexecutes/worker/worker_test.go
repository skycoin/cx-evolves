package worker_test

import (
	"testing"

	"github.com/skycoin/cx-evolves/cxexecutes/worker"
)

func TestGetAvailableWorkers(t *testing.T) {
	tests := []struct {
		scenario             string
		numberOfAvailWorkers int
		wantWorkersAddr      []int
	}{
		{
			scenario:             "1 worker",
			numberOfAvailWorkers: 1,
			wantWorkersAddr:      []int{worker.BasePortNumber},
		},
		{
			scenario:             "2 workers",
			numberOfAvailWorkers: 2,
			wantWorkersAddr: []int{
				worker.BasePortNumber,
				worker.BasePortNumber + 1,
			},
		},
		{
			scenario:             "3 workers",
			numberOfAvailWorkers: 3,
			wantWorkersAddr: []int{
				worker.BasePortNumber,
				worker.BasePortNumber + 1,
				worker.BasePortNumber + 2,
			},
		},
		{
			scenario:             "4 workers",
			numberOfAvailWorkers: 4,
			wantWorkersAddr: []int{
				worker.BasePortNumber,
				worker.BasePortNumber + 1,
				worker.BasePortNumber + 2,
				worker.BasePortNumber + 3,
			},
		},
		{
			scenario:             "5 workers",
			numberOfAvailWorkers: 5,
			wantWorkersAddr: []int{
				worker.BasePortNumber,
				worker.BasePortNumber + 1,
				worker.BasePortNumber + 2,
				worker.BasePortNumber + 3,
				worker.BasePortNumber + 4,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			gotWorkersAddr := worker.GetAvailableWorkers(tc.numberOfAvailWorkers)
			if len(tc.wantWorkersAddr) != len(gotWorkersAddr) {
				t.Errorf("want %v, got %v", tc.wantWorkersAddr, gotWorkersAddr)
			}
		})
	}
}
