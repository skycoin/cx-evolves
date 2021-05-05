package mutation_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	cxevolves "github.com/skycoin/cx-evolves/evolve"
	cxmutation "github.com/skycoin/cx-evolves/mutation"
	cxast "github.com/skycoin/cx/cx/ast"
	"github.com/skycoin/cx/cx/astapi"
	cxconstants "github.com/skycoin/cx/cx/constants"
	cxparsingcompletor "github.com/skycoin/cx/cxparser/cxparsingcompletor"
)

func TestPointMutationOperator_InsertI32Literal_1Byte(t *testing.T) {
	prgrm := generateRandomProgram(t, false)

	mainPkg, err := astapi.FindPackage(prgrm, "main")
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario string
		program  *cxast.CXProgram
		pkg      *cxast.CXPackage
		arg      *cxast.CXArgument
	}{
		{
			scenario: "insert literal on expression index 1 input index 0",
			program:  prgrm,
			pkg:      mainPkg,
			arg:      mainPkg.Functions[0].Expressions[1].Inputs[0],
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.InsertI32Literal_1Byte(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			// size will always be 4 bytes because data type is int 32.
			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}

			valueAt2ndByte := tc.program.Memory[tc.program.DataSegmentStartsAt+1]
			if valueAt2ndByte != 0 {
				t.Errorf("want 2nd byte 0, got %v", valueAt2ndByte)
			}

			valueAt3rdByte := tc.program.Memory[tc.program.DataSegmentStartsAt+2]
			if valueAt3rdByte != 0 {
				t.Errorf("want 3rd byte 0, got %v", valueAt3rdByte)
			}

			valueAt4thByte := tc.program.Memory[tc.program.DataSegmentStartsAt+3]
			if valueAt4thByte != 0 {
				t.Errorf("want 4th byte 0, got %v", valueAt4thByte)
			}
		})
	}
}

func TestPointMutationOperator_InsertI32Literal_2Bytes(t *testing.T) {
	prgrm := generateRandomProgram(t, false)

	mainPkg, err := astapi.FindPackage(prgrm, "main")
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario string
		program  *cxast.CXProgram
		pkg      *cxast.CXPackage
		arg      *cxast.CXArgument
	}{
		{
			scenario: "insert literal on expression index 1 input index 0",
			program:  prgrm,
			pkg:      mainPkg,
			arg:      mainPkg.Functions[0].Expressions[1].Inputs[0],
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.InsertI32Literal_2Bytes(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			// size will always be 4 bytes because data type is int 32.
			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}

			valueAt3rdByte := tc.program.Memory[tc.program.DataSegmentStartsAt+2]
			if valueAt3rdByte != 0 {
				t.Errorf("want 3rd byte 0, got %v", valueAt3rdByte)
			}

			valueAt4thByte := tc.program.Memory[tc.program.DataSegmentStartsAt+3]
			if valueAt4thByte != 0 {
				t.Errorf("want 4th byte 0, got %v", valueAt4thByte)
			}
		})
	}
}

func TestPointMutationOperator_InsertI32Literal_4Bytes(t *testing.T) {
	prgrm := generateRandomProgram(t, false)

	mainPkg, err := astapi.FindPackage(prgrm, "main")
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario string
		program  *cxast.CXProgram
		pkg      *cxast.CXPackage
		arg      *cxast.CXArgument
	}{
		{
			scenario: "insert literal on expression index 1 input index 0",
			program:  prgrm,
			pkg:      mainPkg,
			arg:      mainPkg.Functions[0].Expressions[1].Inputs[0],
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.InsertI32Literal_4Bytes(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()
			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize

			// size will always be 4 bytes because data type is int 32.
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}
		})
	}
}

func TestPointMutationOperator_HalfI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario      string
		program       *cxast.CXProgram
		pkg           *cxast.CXPackage
		arg           *cxast.CXArgument
		wantDataValue int
	}{
		{
			scenario:      "half literal on expression index 2 input index 1",
			program:       prgrm,
			pkg:           mainPkg,
			arg:           mainPkg.Functions[0].Expressions[2].Inputs[1],
			wantDataValue: 3, // 6/2=3
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.HalfI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}

			dataVal := binary.LittleEndian.Uint32(dataValue)
			if int(dataVal) != tc.wantDataValue {
				t.Errorf("want data value %v, got %v", tc.wantDataValue, int(dataVal))
			}
		})
	}
}

func TestPointMutationOperator_DoubleI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario      string
		program       *cxast.CXProgram
		pkg           *cxast.CXPackage
		arg           *cxast.CXArgument
		wantDataValue int
	}{
		{
			scenario:      "doubles literal on expression index 2 input index 1",
			program:       prgrm,
			pkg:           mainPkg,
			arg:           mainPkg.Functions[0].Expressions[2].Inputs[1],
			wantDataValue: 12, // 6*2=12
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.DoubleI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}

			dataVal := binary.LittleEndian.Uint32(dataValue)
			if int(dataVal) != tc.wantDataValue {
				t.Errorf("want data value %v, got %v", tc.wantDataValue, int(dataVal))
			}
		})
	}
}

func TestPointMutationOperator_ZeroI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario      string
		program       *cxast.CXProgram
		pkg           *cxast.CXPackage
		arg           *cxast.CXArgument
		wantDataValue int
	}{
		{
			scenario:      "make literal on expression index 2 input index 1 become zero",
			program:       prgrm,
			pkg:           mainPkg,
			arg:           mainPkg.Functions[0].Expressions[2].Inputs[1],
			wantDataValue: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.ZeroI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}

			dataVal := binary.LittleEndian.Uint32(dataValue)
			if int(dataVal) != tc.wantDataValue {
				t.Errorf("want data value %v, got %v", tc.wantDataValue, int(dataVal))
			}
		})
	}
}

func TestPointMutationOperator_AddOneI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario      string
		program       *cxast.CXProgram
		pkg           *cxast.CXPackage
		arg           *cxast.CXArgument
		wantDataValue int
	}{
		{
			scenario:      "adds 1 on literal on expression index 2 input index 1",
			program:       prgrm,
			pkg:           mainPkg,
			arg:           mainPkg.Functions[0].Expressions[2].Inputs[1],
			wantDataValue: 7, // 6+1=7
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.AddOneI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}

			dataVal := binary.LittleEndian.Uint32(dataValue)
			if int(dataVal) != tc.wantDataValue {
				t.Errorf("want data value %v, got %v", tc.wantDataValue, int(dataVal))
			}
		})
	}
}

func TestPointMutationOperator_AddRand1ByteI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario     string
		program      *cxast.CXProgram
		pkg          *cxast.CXPackage
		arg          *cxast.CXArgument
		currentValue int
	}{
		{
			scenario:     "adds random int32 on literal on expression index 2 input index 1",
			program:      prgrm,
			pkg:          mainPkg,
			arg:          mainPkg.Functions[0].Expressions[2].Inputs[1],
			currentValue: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.AddRand1ByteI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}

			dataVal := binary.LittleEndian.Uint32(dataValue)
			if int(dataVal) <= tc.currentValue {
				t.Errorf("want data value greater than %v, got %v", tc.currentValue, int(dataVal))
			}
		})
	}
}

func TestPointMutationOperator_SubOneI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario      string
		program       *cxast.CXProgram
		pkg           *cxast.CXPackage
		arg           *cxast.CXArgument
		wantDataValue int
	}{
		{
			scenario:      "subtracts 1 on literal on expression index 2 input index 1",
			program:       prgrm,
			pkg:           mainPkg,
			arg:           mainPkg.Functions[0].Expressions[2].Inputs[1],
			wantDataValue: 5, // 6-1=5
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.SubOneI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}

			dataVal := binary.LittleEndian.Uint32(dataValue)
			if int(dataVal) != tc.wantDataValue {
				t.Errorf("want data value %v, got %v", tc.wantDataValue, int(dataVal))
			}
		})
	}
}

func TestPointMutationOperator_SubRand1ByteI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario     string
		program      *cxast.CXProgram
		pkg          *cxast.CXPackage
		arg          *cxast.CXArgument
		currentValue int
	}{
		{
			scenario:     "subtracts random int32 on literal on expression index 2 input index 1",
			program:      prgrm,
			pkg:          mainPkg,
			arg:          mainPkg.Functions[0].Expressions[2].Inputs[1],
			currentValue: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.SubRand1ByteI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}

			dataVal := binary.LittleEndian.Uint32(dataValue)
			if int(dataVal) == tc.currentValue {
				t.Errorf("want data value not %v, got %v", tc.currentValue, dataVal)
			}
		})
	}
}

func TestPointMutationOperator_BitOrI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario     string
		program      *cxast.CXProgram
		pkg          *cxast.CXPackage
		arg          *cxast.CXArgument
		currentValue int
	}{
		{
			scenario:     "literal on expression index 2 input index 1",
			program:      prgrm,
			pkg:          mainPkg,
			arg:          mainPkg.Functions[0].Expressions[2].Inputs[1],
			currentValue: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.BitOrI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}
		})
	}
}

func TestPointMutationOperator_BitAndI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario     string
		program      *cxast.CXProgram
		pkg          *cxast.CXPackage
		arg          *cxast.CXArgument
		currentValue int
	}{
		{
			scenario:     "literal on expression index 2 input index 1",
			program:      prgrm,
			pkg:          mainPkg,
			arg:          mainPkg.Functions[0].Expressions[2].Inputs[1],
			currentValue: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.BitAndI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}
		})
	}
}

func TestPointMutationOperator_BitXorI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario     string
		program      *cxast.CXProgram
		pkg          *cxast.CXPackage
		arg          *cxast.CXArgument
		currentValue int
	}{
		{
			scenario:     "literal on expression index 2 input index 1",
			program:      prgrm,
			pkg:          mainPkg,
			arg:          mainPkg.Functions[0].Expressions[2].Inputs[1],
			currentValue: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.BitXorI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}
		})
	}
}

func TestPointMutationOperator_BitRotateLeftI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario     string
		program      *cxast.CXProgram
		pkg          *cxast.CXPackage
		arg          *cxast.CXArgument
		currentValue int
	}{
		{
			scenario:     "literal on expression index 2 input index 1",
			program:      prgrm,
			pkg:          mainPkg,
			arg:          mainPkg.Functions[0].Expressions[2].Inputs[1],
			currentValue: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.BitRotateLeftI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}
		})
	}
}

func TestPointMutationOperator_BitRotateRightI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario     string
		program      *cxast.CXProgram
		pkg          *cxast.CXPackage
		arg          *cxast.CXArgument
		currentValue int
	}{
		{
			scenario:     "literal on expression index 2 input index 1",
			program:      prgrm,
			pkg:          mainPkg,
			arg:          mainPkg.Functions[0].Expressions[2].Inputs[1],
			currentValue: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.BitRotateRightI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}
		})
	}
}

func TestPointMutationOperator_OrI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario     string
		program      *cxast.CXProgram
		pkg          *cxast.CXPackage
		arg          *cxast.CXArgument
		currentValue int
	}{
		{
			scenario:     "literal on expression index 2 input index 1",
			program:      prgrm,
			pkg:          mainPkg,
			arg:          mainPkg.Functions[0].Expressions[2].Inputs[1],
			currentValue: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.OrI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}
		})
	}
}

func TestPointMutationOperator_AndI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario     string
		program      *cxast.CXProgram
		pkg          *cxast.CXPackage
		arg          *cxast.CXArgument
		currentValue int
	}{
		{
			scenario:     "literal on expression index 2 input index 1",
			program:      prgrm,
			pkg:          mainPkg,
			arg:          mainPkg.Functions[0].Expressions[2].Inputs[1],
			currentValue: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.AndI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}
		})
	}
}

func TestPointMutationOperator_XorI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario     string
		program      *cxast.CXProgram
		pkg          *cxast.CXPackage
		arg          *cxast.CXArgument
		currentValue int
	}{
		{
			scenario:     "literal on expression index 2 input index 1",
			program:      prgrm,
			pkg:          mainPkg,
			arg:          mainPkg.Functions[0].Expressions[2].Inputs[1],
			currentValue: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.XorI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}
		})
	}
}

func TestPointMutationOperator_ShiftOneBitLeftI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario     string
		program      *cxast.CXProgram
		pkg          *cxast.CXPackage
		arg          *cxast.CXArgument
		currentValue int
	}{
		{
			scenario:     "literal on expression index 2 input index 1",
			program:      prgrm,
			pkg:          mainPkg,
			arg:          mainPkg.Functions[0].Expressions[2].Inputs[1],
			currentValue: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.ShiftOneBitLeftI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}

			dataVal := binary.LittleEndian.Uint32(dataValue)
			if int(dataVal) == tc.currentValue {
				t.Errorf("want data value not %v, got %v", tc.currentValue, dataVal)
			}
		})
	}
}

func TestPointMutationOperator_ShiftOneBitRightI32Literal(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	// Generates a literal on expression
	// index 2 input index 1 with value of 6.
	prgrm := generateSampleStaticProgram(t, pkgName, fnName, true)

	mainPkg, err := astapi.FindPackage(prgrm, pkgName)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		scenario     string
		program      *cxast.CXProgram
		pkg          *cxast.CXPackage
		arg          *cxast.CXArgument
		currentValue int
	}{
		{
			scenario:     "literal on expression index 2 input index 1",
			program:      prgrm,
			pkg:          mainPkg,
			arg:          mainPkg.Functions[0].Expressions[2].Inputs[1],
			currentValue: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			cxmutation.ShiftOneBitRightI32Literal(tc.program, tc.pkg, tc.arg)
			tc.program.PrintProgram()

			dataValue := tc.program.Memory[tc.program.DataSegmentStartsAt : tc.program.DataSegmentStartsAt+tc.program.DataSegmentSize]
			printDataInfo(t, tc.program.DataSegmentSize, dataValue)

			dataSegSize := tc.program.DataSegmentSize
			if dataSegSize != 4 {
				t.Errorf("want data segment size 4, got %v", dataSegSize)
			}

			dataVal := binary.LittleEndian.Uint32(dataValue)
			if int(dataVal) == tc.currentValue {
				t.Errorf("want data value not %v, got %v", tc.currentValue, dataVal)
			}
		})
	}
}

func generateRandomProgram(t *testing.T, withLiteral bool) *cxast.CXProgram {
	var cxProgram *cxast.CXProgram

	// Needed for AddNativeExpressionToFunction
	// because of dependency on cxast.OpNames
	cxparsingcompletor.InitCXCore()
	cxProgram = cxast.MakeProgram()

	err := astapi.AddEmptyPackage(cxProgram, "main")
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddEmptyFunctionToPackage(cxProgram, "main", "TestFunction")
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddNativeInputToFunction(cxProgram, "main", "TestFunction", "inputOne", cxconstants.TYPE_I32)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = astapi.AddNativeOutputToFunction(cxProgram, "main", "TestFunction", "outputOne", cxconstants.TYPE_I32)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
	functionSetNames := []string{"i32.add", "i32.mul", "i32.sub", "i32.eq", "i32.uneq", "i32.gt", "i32.gteq", "i32.lt", "i32.lteq", "bool.not", "bool.or", "bool.and", "bool.uneq", "bool.eq", "i32.neg", "i32.abs", "i32.bitand", "i32.bitor", "i32.bitxor", "i32.bitclear", "i32.bitshl", "i32.bitshr", "i32.max", "i32.min", "i32.rand"}
	fns := cxevolves.GetFunctionSet(functionSetNames)

	fn, _ := cxProgram.GetFunction("TestFunction", "main")
	pkg, _ := cxProgram.GetPackage("main")
	cxevolves.GenerateRandomExpressions(fn, pkg, fns, 30)

	if withLiteral {
		buf := new(bytes.Buffer)
		var num int32 = 5
		binary.Write(buf, binary.LittleEndian, num)
		err = astapi.AddLiteralInputToExpression(cxProgram, "main", "TestFunction", buf.Bytes(), cxconstants.TYPE_I32, 2)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	}

	cxProgram.PrintProgram()
	return cxProgram
}

func printDataInfo(t *testing.T, length int, value []byte) {
	t.Logf("Data segment length=%v\n", length)
	t.Logf("Data segment value=%v\n", value)
}
