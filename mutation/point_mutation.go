package mutation

import (
	"fmt"

	cxast "github.com/skycoin/cx/cx/ast"
)

type MutationHandler func(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument)

var (
	PointMutationOperators map[int]MutationHandler
	MutationOpNames        map[int]string
	MutationOpCodes        map[string]int
)

const (
	MOP_INSERT_RAND_ONE_BYTE_I32_LIT = iota + 1
	MOP_INSERT_RAND_TWO_BYTES_I32_LIT
	MOP_INSERT_RAND_FOUR_BYTES_I32_LIT
	MOP_HALF_I32_LIT
	MOP_DOUBLE_I32_LIT
	MOP_ZERO_I32_LIT
	MOP_ADD_ONE_I32_LIT
	MOP_SUB_ONE_I32_LIT
)

// RegisterMutationOperator
func RegisterMutationOperator(code int, name string, handler MutationHandler) {
	// Check if duplicate
	if PointMutationOperators[code] != nil {
		panic(fmt.Sprintf("duplicate opcode %d : '%s' width '%s'.\n", code, name, MutationOpNames[code]))
	}

	PointMutationOperators[code] = handler
	MutationOpNames[code] = name
	MutationOpCodes[name] = code
}
