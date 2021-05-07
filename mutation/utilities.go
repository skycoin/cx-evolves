package mutation

import (
	"math/rand"
)

// GenerateRandomBytes to generate random bytes with uniform length []byte of 4 output.
// Expecting to generate only with sizes 1, 2, and 4.
func GenerateRandomBytes(size int) (blk []byte, err error) {
	blk = make([]byte, size)
	_, err = rand.Read(blk)

	switch size {
	case 1:
		blk = append(blk, []byte{0, 0, 0}...) // Fill extra 3 zeroes
	case 2:
		blk = append(blk, []byte{0, 0}...) // Fill extra 2 zeroes
	}
	return
}

func GetMutationOperatorFunctionSet(fnNames []string) (fns []MutationHandler) {
	for _, fnName := range fnNames {
		fn := PointMutationOperators[MutationOpCodes[fnName]]
		if fn == nil {
			panic("mutation operator function not found.")
		}

		fns = append(fns, fn)
	}
	return fns
}

func GetAllMutationOperatorFunctionSet() (fns []MutationHandler) {
	for _, fn := range PointMutationOperators {
		fns = append(fns, fn)
	}
	return fns
}
