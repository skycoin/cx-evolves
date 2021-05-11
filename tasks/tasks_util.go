package tasks

import (
	"encoding/binary"

	cxast "github.com/skycoin/cx/cx/ast"
)

func toByteArray(i int32) []byte {
	arr := make([]byte, 4)
	binary.BigEndian.PutUint32(arr, uint32(i))
	return arr
}

// injectMainInputs injects `inps` at the beginning of `prgrm`'s memory,
// which should always represent the memory sent to the first expression contained
// in `prgrm`'s `main`'s function.
func injectMainInputs(prgrm *cxast.CXProgram, inps []byte) {
	for i := 0; i < len(inps); i++ {
		prgrm.Memory[i] = inps[i]
	}
}
