package mutation

import (
	"math/rand"

	cxast "github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/constants"
	"github.com/skycoin/cx/cx/types"
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
func WriteLiteralArg(cxprogram *cxast.CXProgram, typ types.Code, byts []byte, isGlobal bool) []*cxast.CXExpression {
	pkg, err := cxprogram.GetCurrentPackage()
	if err != nil {
		panic(err)
	}

	arg := cxast.MakeArgument("", "", 0)
	arg.AddType(typ)
	arg.ArgDetails.Package = pkg

	var size = len(byts)

	arg.Size = types.Code(typ).Size()
	arg.TotalSize = types.Pointer(size)
	arg.Offset = cxprogram.Data.Size + cxprogram.Data.StartsAt

	if arg.Type == types.STR || arg.Type == types.AFF {
		arg.PassBy = constants.PASSBY_REFERENCE
		arg.Size = types.POINTER_SIZE
		arg.TotalSize = types.POINTER_SIZE
	}

	// A CX program allocates min(INIT_HEAP_SIZE, MAX_HEAP_SIZE) bytes
	// after the stack segment. These bytes are used to allocate the data segment
	// at compile time. If the data segment is bigger than min(INIT_HEAP_SIZE, MAX_HEAP_SIZE),
	// we'll start appending the bytes to AST.Memory.
	// After compilation, we calculate how many bytes we need to add to have a heap segment
	// equal to `minHeapSize()` that is allocated after the data segment.
	if (size + int(cxprogram.Data.Size) + int(cxprogram.Data.StartsAt)) > len(cxprogram.Memory) {
		var i int
		// First we need to fill the remaining free bytes in
		// the current `AST.Memory` slice.
		for i = 0; i < len(cxprogram.Memory)-int(cxprogram.Data.Size)+int(cxprogram.Data.StartsAt); i++ {
			cxprogram.Memory[int(cxprogram.Data.Size)+int(cxprogram.Data.StartsAt)+i] = byts[i]
		}
		// Then we append the bytes that didn't fit.
		cxprogram.Memory = append(cxprogram.Memory, byts[i:]...)
	} else {
		for i, byt := range byts {
			cxprogram.Memory[int(cxprogram.Data.Size)+int(cxprogram.Data.StartsAt)+i] = byt
		}
	}
	cxprogram.Data.Size += types.Pointer(size)

	expr := cxast.MakeExpression(nil, "", 0)
	expr.Package = pkg
	expr.Outputs = append(expr.Outputs, arg)
	return []*cxast.CXExpression{expr}
}
