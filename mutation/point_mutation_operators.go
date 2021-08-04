package mutation

import (
	"bytes"
	"encoding/binary"
	"math/rand"

	cxast "github.com/skycoin/cx/cx/ast"
	cxconstants "github.com/skycoin/cx/cx/constants"
)

// InsertI32Literal inserts a random n byte i32 literal as replacement for cxarg.
func InsertI32Literal(n int, cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	_ = cxarg
	byts, err := GenerateRandomBytes(n)
	if err != nil {
		panic(err)
	}

	// Initialize prereq for WritePrimary()
	cxprogram.CurrentPackage = pkg

	// create literal arg
	litArg := WriteLiteralArg(cxprogram, cxconstants.TYPE_I32, byts, false)
	arg := litArg[0].Outputs[0]
	arg.ArgDetails.Package = pkg

	// replace cxarg with new literal arg
	cxarg = arg
}

// InsertI32Literal_1Byte inserts a random 1 byte i32 literal as replacement for cxarg.
func InsertI32Literal_1Byte(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	InsertI32Literal(1, cxprogram, pkg, cxarg)
}

// InsertI32Literal_2Bytes inserts a random 2 bytes i32 literal as replacement for cxarg.
func InsertI32Literal_2Bytes(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	InsertI32Literal(2, cxprogram, pkg, cxarg)
}

// InsertI32Literal_4Bytes inserts a random 4 bytes i32 literal as replacement for cxarg.
func InsertI32Literal_4Bytes(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	InsertI32Literal(4, cxprogram, pkg, cxarg)
}

// HalfI32Literal divides an i32 literal value by two.
func HalfI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]
	// Divide the value by 2
	newValue := (binary.LittleEndian.Uint32(currentByts)) / 2

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// DoubleI32Literal multiplies an i32 literal value by two.
func DoubleI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]

	// Multiply the value by 2
	newValue := (binary.LittleEndian.Uint32(currentByts)) * 2

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// ZeroI32Literal sets an i32 literal value to zero.
func ZeroI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]
	newByts := []byte{}

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// AddOneiteral adds an i32 literal value with 1.
func AddOneI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]

	// Add 1 to the value
	newValue := (binary.LittleEndian.Uint32(currentByts)) + 1

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// AddRand1ByteLiteral adds a random 1 byte to the existing literal value.
func AddRand1ByteI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]
	gen1Byt, err := GenerateRandomBytes(1)
	if err != nil {
		panic(err)
	}
	randomInt8 := binary.LittleEndian.Uint32(gen1Byt)

	newValue := (binary.LittleEndian.Uint32(currentByts)) + randomInt8

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// SubOneiteral substracts an i32 literal value with 1.
func SubOneI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]

	// Subtract 1 to the value
	newValue := (binary.LittleEndian.Uint32(currentByts)) - 1

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// SubRand1ByteLiteral subtracts a random 1 byte to the existing literal value.
func SubRand1ByteI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]
	gen1Byt, err := GenerateRandomBytes(1)
	if err != nil {
		panic(err)
	}
	randomInt8 := binary.LittleEndian.Uint32(gen1Byt)

	newValue := (binary.LittleEndian.Uint32(currentByts)) - randomInt8

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// BitOrI32Literal sets random nth bit to 1.
func BitOrI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]
	currVal := binary.LittleEndian.Uint32(currentByts)

	gen1Byt, err := GenerateRandomBytes(1)
	if err != nil {
		panic(err)
	}
	randomInt8 := binary.LittleEndian.Uint32(gen1Byt)

	// Sets random nth bit to 1.
	newValue := currVal | (1 << (randomInt8 % 32))

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// BitAndI32Literal sets all bits to zero except for random nth bit.
func BitAndI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]
	currVal := binary.LittleEndian.Uint32(currentByts)

	gen1Byt, err := GenerateRandomBytes(1)
	if err != nil {
		panic(err)
	}
	randomInt8 := binary.LittleEndian.Uint32(gen1Byt)

	// Sets rall bits to 0 except for random nth bit.
	newValue := currVal & (1 << (randomInt8 % 32))

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// BitXorI32Literal flips random nth bit.
func BitXorI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]
	currVal := binary.LittleEndian.Uint32(currentByts)

	gen1Byt, err := GenerateRandomBytes(1)
	if err != nil {
		panic(err)
	}
	randomInt8 := binary.LittleEndian.Uint32(gen1Byt)

	// Sets random nth bit to 1.
	newValue := currVal ^ (1 << (randomInt8 % 32))

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// BitRotateLeftI32Literal rotates whole 32 byte to the left.
func BitRotateLeftI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]
	currVal := binary.LittleEndian.Uint32(currentByts)

	gen1Byt, err := GenerateRandomBytes(1)
	if err != nil {
		panic(err)
	}
	randomInt8 := binary.LittleEndian.Uint32(gen1Byt)

	// Rotate whole 32bit with rand
	randRotate := randomInt8 % 32
	newValue := (currVal << randRotate) | (currVal >> (32 - randRotate))

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// BitRotateRightI32Literal rotates whole 32 byte to the right.
func BitRotateRightI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]
	currVal := binary.LittleEndian.Uint32(currentByts)

	gen1Byt, err := GenerateRandomBytes(1)
	if err != nil {
		panic(err)
	}
	randomInt8 := binary.LittleEndian.Uint32(gen1Byt)

	// Rotate whole 32bit with rand
	randRotate := randomInt8 % 32
	newValue := (currVal >> randRotate) | (currVal << (32 - randRotate))

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// OrI32Literal does OR operation on the literal value and random uint32 value.
func OrI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]
	currVal := binary.LittleEndian.Uint32(currentByts)
	newValue := currVal | rand.Uint32()

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// AndI32Literal does AND operation on the literal value and random uint32 value.
func AndI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]
	currVal := binary.LittleEndian.Uint32(currentByts)
	newValue := currVal & rand.Uint32()

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// XorI32Literal does XOR operation on the literal value and random uint32 value.
func XorI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]
	currVal := binary.LittleEndian.Uint32(currentByts)
	newValue := currVal ^ rand.Uint32()

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// ShiftOneBitLeft shifts whole ltieral value one bit to left.
func ShiftOneBitLeftI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]
	currVal := binary.LittleEndian.Uint32(currentByts)
	newValue := currVal << 1

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}

// ShiftOneBitRight shifts whole ltieral value one bit to right.
func ShiftOneBitRightI32Literal(cxprogram *cxast.CXProgram, pkg *cxast.CXPackage, cxarg *cxast.CXArgument) {
	argOffset := cxarg.Offset

	// Check if has value in data segment memory
	if argOffset == 0 {
		return
	}

	currentByts := cxprogram.Memory[argOffset : argOffset+cxarg.TotalSize]
	currVal := binary.LittleEndian.Uint32(currentByts)
	newValue := currVal >> 1

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, newValue)

	newByts := buf.Bytes()

	lenCurrByts := len(currentByts)
	lenNewByts := len(newByts)
	if lenCurrByts > lenNewByts {
		diff := lenCurrByts - lenNewByts
		for i := 0; i < diff; i++ {
			newByts = append(newByts, byte(0))
		}
	}

	// Overwrite current value with new value
	x := 0
	for i := argOffset; i < argOffset+cxarg.TotalSize; i++ {
		cxprogram.Memory[i] = newByts[x]
		x++
	}
}
