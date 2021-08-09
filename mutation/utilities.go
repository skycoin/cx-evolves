package mutation

import (
	"math/rand"

	cxast "github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
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

// This function writes those bytes to cxprogram.Data
func WriteLiteralArg(cxprogram *cxast.CXProgram, typ int, byts []byte, isGlobal bool) []*cxast.CXExpression {
	pkg, err := cxprogram.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	arg := cxast.MakeArgument("", "", 0)
	arg.AddType(constants.TypeNames[typ])
	arg.ArgDetails.Package = pkg

	var size = len(byts)

	arg.Size = constants.GetArgSize(typ)
	arg.TotalSize = size
	arg.Offset = cxprogram.DataSegmentSize + cxprogram.DataSegmentStartsAt

	if arg.Type == constants.TYPE_STR || arg.Type == constants.TYPE_AFF {
		arg.PassBy = constants.PASSBY_REFERENCE
		arg.Size = constants.TYPE_POINTER_SIZE
		arg.TotalSize = constants.TYPE_POINTER_SIZE
	}

	// A CX program allocates min(INIT_HEAP_SIZE, MAX_HEAP_SIZE) bytes
	// after the stack segment. These bytes are used to allocate the data segment
	// at compile time. If the data segment is bigger than min(INIT_HEAP_SIZE, MAX_HEAP_SIZE),
	// we'll start appending the bytes to AST.Memory.
	// After compilation, we calculate how many bytes we need to add to have a heap segment
	// equal to `minHeapSize()` that is allocated after the data segment.
	if (size + cxprogram.DataSegmentSize + cxprogram.DataSegmentStartsAt) > len(cxprogram.Memory) {
		var i int
		// First we need to fill the remaining free bytes in
		// the current `AST.Memory` slice.
		for i = 0; i < len(cxprogram.Memory)-cxprogram.DataSegmentSize+cxprogram.DataSegmentStartsAt; i++ {
			cxprogram.Memory[cxprogram.DataSegmentSize+cxprogram.DataSegmentStartsAt+i] = byts[i]
		}
		// Then we append the bytes that didn't fit.
		cxprogram.Memory = append(cxprogram.Memory, byts[i:]...)
	} else {
		for i, byt := range byts {
			cxprogram.Memory[cxprogram.DataSegmentSize+cxprogram.DataSegmentStartsAt+i] = byt
		}
	}
	cxprogram.DataSegmentSize += size

	expr := cxast.MakeExpression(nil, "", 0)
	expr.Package = pkg
	expr.Outputs = append(expr.Outputs, arg)
	return []*cxast.CXExpression{expr}
}
