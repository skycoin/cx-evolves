package mutation

func RegisterMutationOperators() {
	RegisterMutationOperator(MOP_INSERT_RAND_ONE_BYTE_I32_LIT, "i32.insertLit1Byte", InsertI32Literal_1Byte)
	RegisterMutationOperator(MOP_INSERT_RAND_TWO_BYTES_I32_LIT, "i32.insertLit2Bytes", InsertI32Literal_2Bytes)
	RegisterMutationOperator(MOP_INSERT_RAND_FOUR_BYTES_I32_LIT, "i32.insertLit4Bytes", InsertI32Literal_4Bytes)
	RegisterMutationOperator(MOP_HALF_I32_LIT, "i32.halfLit", HalfI32Literal)
	RegisterMutationOperator(MOP_DOUBLE_I32_LIT, "i32.doubleLit", DoubleI32Literal)
	RegisterMutationOperator(MOP_ZERO_I32_LIT, "i32.zeroLit", ZeroI32Literal)
	RegisterMutationOperator(MOP_ADD_ONE_I32_LIT, "i32.addOneLit", AddOneI32Literal)
	RegisterMutationOperator(MOP_SUB_ONE_I32_LIT, "i32.subOneLit", SubOneI32Literal)
}
