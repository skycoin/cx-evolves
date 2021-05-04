package mutation

func RegisterMutationOperators() {
	RegisterMutationOperator(MOP_INSERT_RAND_I8_AS__I32_LIT, "i32.insertLit1Byte", InsertI32Literal_1Byte)
	RegisterMutationOperator(MOP_INSERT_RAND_I16_AS_I32_LIT, "i32.insertLit2Bytes", InsertI32Literal_2Bytes)
	RegisterMutationOperator(MOP_INSERT_RAND__I32_LIT, "i32.insertLit4Bytes", InsertI32Literal_4Bytes)
	RegisterMutationOperator(MOP_HALF_I32_LIT, "i32.halfLit", HalfI32Literal)
	RegisterMutationOperator(MOP_DOUBLE_I32_LIT, "i32.doubleLit", DoubleI32Literal)
	RegisterMutationOperator(MOP_ZERO_I32_LIT, "i32.zeroLit", ZeroI32Literal)
	RegisterMutationOperator(MOP_ADD_ONE_I32_LIT, "i32.addOneLit", AddOneI32Literal)
	RegisterMutationOperator(MOP_ADD_RAND_I32_LIT, "i32.addRandLit", AddRand1ByteI32Literal)
	RegisterMutationOperator(MOP_SUB_ONE_I32_LIT, "i32.subOneLit", SubOneI32Literal)
	RegisterMutationOperator(MOP_SUB_RAND_I32_LIT, "i32.subRandLit", SubRand1ByteI32Literal)
	RegisterMutationOperator(MOP_BIT_OR_I32_LIT, "i32.bitOrLit", BitOrI32Literal)
	RegisterMutationOperator(MOP_BIT_AND_I32_LIT, "i32.bitAndLit", BitAndI32Literal)
	RegisterMutationOperator(MOP_BIT_XOR_I32_LIT, "i32.bitXorLit", BitXorI32Literal)
	RegisterMutationOperator(MOP_OR_I32_LIT, "i32.orLit", OrI32Literal)
	RegisterMutationOperator(MOP_AND_I32_LIT, "i32.andLit", AndI32Literal)
	RegisterMutationOperator(MOP_XOR_I32_LIT, "i32.xorLit", XorI32Literal)
	RegisterMutationOperator(MOP_BIT_ROTATE_LEFT_I32_LIT, "i32.rotateLeftLit", BitRotateLeftI32Literal)
	RegisterMutationOperator(MOP_BIT_ROTATE_RIGHT_I32_LIT, "i32.rotateRightLit", BitRotateRightI32Literal)
	RegisterMutationOperator(MOP_SHIFT_BIT_LEFT_I32_LIT, "i32.shiftLeftLit", ShiftOneBitLeftI32Literal)
	RegisterMutationOperator(MOP_SHIFT_BIT_RIGHT_I32_LIT, "i32.shiftRightLit", ShiftOneBitRightI32Literal)
}
