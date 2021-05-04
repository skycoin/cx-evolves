package mutation_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"testing"

	"github.com/skycoin/cx-evolves/mutation"
	cxast "github.com/skycoin/cx/cx/ast"
	cxastapi "github.com/skycoin/cx/cx/astapi"
	cxconstants "github.com/skycoin/cx/cx/constants"
	cxparsingcompletor "github.com/skycoin/cx/cxparser/cxparsingcompletor"
)

func TestGetCompatiblePositionForOperator(t *testing.T) {
	tests := []struct {
		scenario     string
		program      *cxast.CXProgram
		functionName string
		operatorName string
		wantLines    []int
		wantErr      error
	}{
		{
			scenario:     "valid function name and operator name - i16.add",
			program:      generateSampleStaticProgram(t, "main", "TestFunction", false),
			functionName: "TestFunction",
			operatorName: "i16.add",
			wantLines:    []int{0, 1},
			wantErr:      nil,
		},
		{
			scenario:     "valid function name and operator name - i32.add",
			program:      generateSampleStaticProgram(t, "main", "TestFunction", false),
			functionName: "TestFunction",
			operatorName: "i32.add",
			wantLines:    []int{2},
			wantErr:      nil,
		},
		{
			scenario:     "valid function name and operator name - jmp",
			program:      generateSampleStaticProgram(t, "main", "TestFunction", false),
			functionName: "TestFunction",
			operatorName: "jmp",
			wantLines:    []int{0, 1, 2},
			wantErr:      nil,
		},
		{
			scenario:     "valid function name but invalid operator name",
			program:      generateSampleStaticProgram(t, "main", "TestFunction", false),
			functionName: "TestFunction",
			operatorName: "i256.div",
			wantLines:    []int{},
			wantErr:      errors.New("standard library function not found"),
		},
		{
			scenario:     "valid operator name but invalid function name",
			program:      generateSampleStaticProgram(t, "main", "TestFunction", false),
			functionName: "Unknown",
			operatorName: "i32.sub",
			wantLines:    []int{},
			wantErr:      errors.New("function not found"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			gotLines, err := mutation.GetCompatiblePositionForOperator(tc.program, tc.functionName, tc.operatorName)
			gotErr := err != nil
			wantErr := tc.wantErr != nil
			if gotErr != wantErr {
				t.Fatalf("want err %v, got %v", tc.wantErr, err)
			}

			if !reflect.DeepEqual(gotLines, tc.wantLines) {
				t.Errorf("want lines %v, got %v", tc.wantLines, gotLines)
			}
		})
	}
}

func TestReplaceArgInput(t *testing.T) {
	pkgName := "main"
	fnName := "TestFunction"
	cxProgram := generateSampleStaticProgram(t, pkgName, fnName, false)
	fn, err := cxastapi.FindFunction(cxProgram, fnName)
	if err != nil {
		t.Fatalf("error in finding function")
	}

	argsAvailable, _ := cxastapi.GetAccessibleArgsForFunctionByType(cxProgram, pkgName, fnName, fn.Expressions[1].Inputs[0].Type)
	tests := []struct {
		scenario string
		cxExpr   *cxast.CXExpression
		argIndex int
		argToPut *cxast.CXArgument
		wantErr  error
	}{
		{
			scenario: "valid expression and argument",
			cxExpr:   fn.Expressions[1],
			argIndex: 1,
			argToPut: argsAvailable[0],
			wantErr:  nil,
		},
		{
			scenario: "invalid arg index",
			cxExpr:   fn.Expressions[1],
			argIndex: 4,
			argToPut: argsAvailable[0],
			wantErr:  errors.New("invalid arg index"),
		},
		{
			scenario: "invalid arg type",
			cxExpr:   fn.Expressions[2],
			argIndex: 0,
			argToPut: argsAvailable[0],
			wantErr:  errors.New("arg types are not the same"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			err := mutation.ReplaceArgInput(tc.cxExpr, tc.argIndex, tc.argToPut)
			gotErr := err != nil
			wantErr := tc.wantErr != nil
			if gotErr != wantErr {
				t.Fatalf("got err %v, want %v", err, tc.wantErr)
			}

			if !wantErr && tc.cxExpr.Inputs[tc.argIndex].ArgDetails.Name != tc.argToPut.ArgDetails.Name {
				t.Errorf("want arg %v, got %v", tc.argToPut.ArgDetails.Name, tc.cxExpr.Inputs[tc.argIndex].ArgDetails.Name)
			}
		})
	}
}

// Output of this generator is:
// Program
// 0.- Package: main
// Functions
// 		0.- Function: TestFunction (inputOne i32) (outputOne i16)
// 				0.- Expression: z i16 = add(x i16, y i16)
// 				1.- Expression: z i16 = sub(x i16, y i16)
// 				2.- Expression: z i32 = mul(x i32, y i32)
func generateSampleStaticProgram(t *testing.T, pkgName, fnName string, withLiteral bool) *cxast.CXProgram {
	var cxProgram *cxast.CXProgram

	// Needed for AddNativeExpressionToFunction
	// because of dependency on cxast.OpNames
	cxparsingcompletor.InitCXCore()
	cxProgram = cxast.MakeProgram()

	err := cxastapi.AddEmptyPackage(cxProgram, pkgName)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = cxastapi.AddEmptyFunctionToPackage(cxProgram, pkgName, fnName)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = cxastapi.AddNativeInputToFunction(cxProgram, pkgName, fnName, "inputOne", cxconstants.TYPE_I32)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = cxastapi.AddNativeOutputToFunction(cxProgram, pkgName, fnName, "outputOne", cxconstants.TYPE_I16)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = cxastapi.AddNativeExpressionToFunction(cxProgram, fnName, cxconstants.OP_ADD)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = cxastapi.AddNativeInputToExpression(cxProgram, pkgName, fnName, "x", cxconstants.TYPE_I16, 0)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = cxastapi.AddNativeInputToExpression(cxProgram, pkgName, fnName, "y", cxconstants.TYPE_I16, 0)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = cxastapi.AddNativeOutputToExpression(cxProgram, pkgName, fnName, "z", cxconstants.TYPE_I16, 0)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = cxastapi.AddNativeExpressionToFunction(cxProgram, fnName, cxconstants.OP_SUB)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = cxastapi.AddNativeInputToExpression(cxProgram, pkgName, fnName, "x", cxconstants.TYPE_I16, 1)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = cxastapi.AddNativeInputToExpression(cxProgram, pkgName, fnName, "y", cxconstants.TYPE_I16, 1)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = cxastapi.AddNativeOutputToExpression(cxProgram, pkgName, fnName, "z", cxconstants.TYPE_I16, 1)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = cxastapi.AddNativeExpressionToFunction(cxProgram, fnName, cxconstants.OP_MUL)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	err = cxastapi.AddNativeInputToExpression(cxProgram, pkgName, fnName, "x", cxconstants.TYPE_I32, 2)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	if withLiteral {
		buf := new(bytes.Buffer)
		var num int32 = 6
		binary.Write(buf, binary.LittleEndian, num)
		err = cxastapi.AddLiteralInputToExpression(cxProgram, "main", "TestFunction", buf.Bytes(), cxconstants.TYPE_I32, 2)
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
	}

	err = cxastapi.AddNativeOutputToExpression(cxProgram, pkgName, fnName, "z", cxconstants.TYPE_I32, 2)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}

	cxProgram.PrintProgram()
	return cxProgram
}
